package testdb

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"go-vents/internal/domain/types"
	"go-vents/internal/infrastructure/configuration"
	"go-vents/internal/infrastructure/database"
	slogcolored "go-vents/internal/infrastructure/logging/slog-colored"

	spannergorm "github.com/googleapis/go-gorm-spanner"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance   *gorm.DB
	dbOnce       sync.Once
	repoInstance types.EventRepository
	repoOnce     sync.Once
)

// [.'.]:> 🔄 Obtiene la conexión a la base de datos para pruebas
// [.'.]:> Usa la misma configuración que el código de producción
// [.'.]:> Si hay variables de entorno, tienen prioridad sobre el archivo de configuración
func GetDB() *gorm.DB {
	if os.Getenv("DATABASE_DRIVER") == "spanner" {
		return nil
	}

	dbOnce.Do(func() {
		var dsn string
		var configSource string
		driver := os.Getenv("DATABASE_DRIVER")
		if driver == "" {
			driver = "postgres"
		}

		cleanDriver := strings.ReplaceAll(strings.ReplaceAll(driver, "\n", ""), "\r", "")
		log.Println("['.']:> ==============================================")
		log.Println("['.']:> 🧪 INICIALIZANDO ENTORNO DE PRUEBAS (" + strings.ToUpper(cleanDriver) + ") 🧪")
		log.Println("['.']:> ==============================================")

		if envDSN := os.Getenv("DATABASE_DSN"); envDSN != "" {
			dsn = envDSN
			configSource = "ENV VARS"
			log.Println("['.']:> 🌟 ORIGEN DE CONFIGURACIÓN: VARIABLES DE ENTORNO")
		} else {
			logger := slogcolored.NewColoredSLogger()
			config, err := configuration.LoadConfiguration("../../..", logger)
			if err != nil {
				log.Printf("['.']:> ⚠️ No se pudo cargar la configuración: %v", err)
				if driver == "mariadb" {
					dsn = "root:admin@tcp(localhost:3306)/markitos-it-svc-event?charset=utf8mb4&parseTime=True&loc=Local"
				} else {
					dsn = "host=localhost user=admin password=admin dbname=markitos-it-svc-event sslmode=disable"
				}
				configSource = "HARDCODED DEFAULTS"
				log.Println("['.']:> 🌟 ORIGEN DE CONFIGURACIÓN: VALORES PREDETERMINADOS INTERNOS")
			} else {
				dsn = config.DatabaseDsn
				configSource = "CONFIG FILE"
				log.Println("['.']:> 🌟 ORIGEN DE CONFIGURACIÓN: ARCHIVO DE CONFIGURACIÓN")
			}
		}

		maskedDSN := maskPassword(dsn)
		sanitizedDSN := strings.NewReplacer("\n", "", "\r", "").Replace(maskedDSN)

		log.Println("['.']:> ----------------------------------------------")
		log.Println("['.']:> 🔍 Modo de configuración:", configSource)
		fmt.Println("['.']:> 🔌 Conectando a base de datos:", sanitizedDSN)
		log.Println("['.']:> ----------------------------------------------")

		var db *gorm.DB
		var err error

		if driver == "mariadb" {
			db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		} else {
			db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		}

		if err != nil {
			log.Println("['.']:> ❌ ERROR DE CONEXIÓN A BASE DE DATOS ❌")
			log.Println("['.']:> ==============================================")
			log.Fatalf("['.']:> Error: %v", err)
		}

		log.Println("['.']:> ✅ CONEXIÓN EXITOSA A BASE DE DATOS")
		log.Println("['.']:> ==============================================")

		dbInstance = db
		_ = dbInstance.AutoMigrate(&types.Event{}, &types.Queue{}, &types.Subscription{})
	})

	return dbInstance
}

func GetRepository() types.EventRepository {
	repoOnce.Do(func() {
		driver := os.Getenv("DATABASE_DRIVER")

		switch driver {
		case "spanner":
			log.Println("['.']:> ==============================================")
			log.Println("['.']:> 🧪 INICIALIZANDO ENTORNO DE PRUEBAS (SPANNER) 🧪")
			log.Println("['.']:> ==============================================")
			dsn := os.Getenv("DATABASE_DSN")
			db, err := gorm.Open(spannergorm.New(spannergorm.Config{
				DriverName: "spanner",
				DSN:        dsn,
			}), &gorm.Config{})
			if err != nil {
				log.Fatalf("['.']:> Error conectando a Spanner: %v", err)
			}
			repoInstance = database.NewEventSpannerRepository(db)
			log.Printf("['.']:> 📦 Repositorio de prueba inicializado (Spanner)")
		case "mariadb":
			db := GetDB()
			repo := database.NewEventMariaDBRepository(db)
			repoInstance = repo
			log.Printf("['.']:> 📦 Repositorio de prueba inicializado (MariaDB)")
		default:
			db := GetDB()
			repo := database.NewEventPostgresRepository(db)
			repoInstance = repo
			log.Printf("['.']:> 📦 Repositorio de prueba inicializado (Postgres)")
		}
	})

	return repoInstance
}

// [.'.]:> 🔒 Oculta la contraseña en la cadena de conexión
// [.'.]:> para no exponer información sensible en los logs
func maskPassword(dsn string) string {
	if dsn == "" {
		return "¡No configurada!"
	}

	if strings.Contains(dsn, "password=") {
		parts := strings.Split(dsn, " ")
		for i, part := range parts {
			if strings.HasPrefix(part, "password=") {
				parts[i] = "password=******"
			}
		}
		return strings.Join(parts, " ")
	}

	// Simple mask for mariadb dsn format user:pass@tcp(...)
	if strings.Contains(dsn, "@tcp") && strings.Contains(dsn, ":") {
		parts := strings.Split(dsn, "@")
		if len(parts) > 1 {
			credParts := strings.Split(parts[0], ":")
			if len(credParts) == 2 {
				return fmt.Sprintf("%s:******@%s", credParts[0], parts[1])
			}
		}
	}

	return dsn
}
