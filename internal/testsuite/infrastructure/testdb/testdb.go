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
	dbOnce.Do(func() {
		var dsn string
		var configSource string

		log.Println("['.']:> ==============================================")
		log.Println("['.']:> 🧪 INICIALIZANDO ENTORNO DE PRUEBAS 🧪")
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
				dsn = "host=localhost user=admin password=admin dbname=markitos-it-svc-event sslmode=disable"
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

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("['.']:> ❌ ERROR DE CONEXIÓN A BASE DE DATOS ❌")
			log.Println("['.']:> ==============================================")
			log.Fatalf("['.']:> Error: %v", err)
		}

		log.Println("['.']:> ✅ CONEXIÓN EXITOSA A BASE DE DATOS")
		log.Println("['.']:> ==============================================")

		dbInstance = db
		_ = dbInstance.AutoMigrate(&types.Event{})
	})

	return dbInstance
}

// [.'.]:> 🔄 Obtiene el repositorio para pruebas
// [.'.]:> Reutiliza la conexión a la base de datos
func GetRepository() types.EventRepository {
	repoOnce.Do(func() {
		db := GetDB()
		repo := database.NewEventPostgresRepository(db)
		repoInstance = repo
		log.Printf("['.']:> 📦 Repositorio de prueba inicializado")
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

	return dsn
}
