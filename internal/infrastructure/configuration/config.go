package configuration

import (
	"os"

	"go-vents/internal/domain/types"

	"github.com/spf13/viper"
)

//[.'.]:> 📦 SISTEMA DE CONFIGURACIÓN
//[.'.]:> ========================
//[.'.]:> Este módulo maneja la configuración de nuestra aplicación siguiendo este flujo:
//[.'.]:>
//[.'.]:> 1️⃣ ARCHIVO CONFIG: Busca primero un archivo app.env en el directorio especificado
//[.'.]:>    Si lo encuentra, carga todas sus variables como configuración base
//[.'.]:>
//[.'.]:> 2️⃣ VARIABLES DE ENTORNO: Después de cargar el archivo (o si no existe):
//[.'.]:>    - Comprueba si existen variables de entorno con los mismos nombres
//[.'.]:>    - Las variables de entorno tienen PRIORIDAD y sobrescriben los valores del archivo
//[.'.]:>
//[.'.]:> 3️⃣ VALORES POR DEFECTO: Como último recurso, si algún valor sigue vacío,
//[.'.]:>    se aplican valores predeterminados para garantizar que la app pueda funcionar

// [.'.]:> 🧩 Estructura que contiene toda la configuración de la aplicación
// [.'.]:> Cada campo se mapea a una variable de entorno o valor en app.env del mismo nombre
type EventConfiguration struct {
	DatabaseDriver    string `mapstructure:"DATABASE_DRIVER"`
	DatabaseDsn       string `mapstructure:"DATABASE_DSN"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
}

// [.'.]:> 🔄 Función principal que carga toda la configuración
// [.'.]:> Recibe la ruta donde buscar el archivo app.env y devuelve la configuración completa
// [.'.]:> Si hay algún error durante la carga, lo devuelve para que la aplicación pueda manejarlo
func LoadConfiguration(configFilesPath string, logger types.Logger) (config EventConfiguration, err error) {
	viper.New()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configFilesPath)
	_ = viper.BindEnv("DATABASE_DRIVER")
	_ = viper.BindEnv("DATABASE_DSN")
	_ = viper.BindEnv("GRPC_SERVER_ADDRESS")
	viper.AutomaticEnv()

	if err := loadConfigFile(logger); err != nil {
		return config, err
	}

	overrideWithEnvVars()

	err = viper.Unmarshal(&config)
	if err == nil {
		logger.Info("['.']:> ✨ Configuración cargada correctamente ✨")
		logger.Info("['.']:> ----------------------------------------")
		logger.Info("['.']:> 🚀 gRPC Server.: " + config.GRPCServerAddress + "\n")
		logger.Info("['.']:> 📁 Database DSN: " + config.DatabaseDsn + "\n")
		logger.Info("['.']:> ----------------------------------------")
		applyFallbackEnvVars(&config)
	}

	return config, err
}

// [.'.]:> 📄 Intenta cargar el archivo de configuración app.env
// [.'.]:> Si el archivo no existe, lo maneja elegantemente y permite continuar
// [.'.]:> usando solo variables de entorno
func loadConfigFile(logger types.Logger) error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
		logger.Info("['.']:> 📋 Archivo de configuración no encontrado, usando solo variables de entorno")
	} else {
		logger.Info("['.']:> 📋 Archivo de configuración cargado correctamente")
	}

	return nil
}

// [.'.]:> 🔀 Sobrescribe valores del archivo con variables de entorno
// [.'.]:> Esta es la clave para que las variables de entorno tengan prioridad
// [.'.]:> sobre el archivo de configuración
func overrideWithEnvVars() {
	driverEnv := os.Getenv("DATABASE_DRIVER")
	if driverEnv != "" && viper.GetString("DATABASE_DRIVER") == "" {
		viper.Set("DATABASE_DRIVER", driverEnv)
	}

	dsnEnv := os.Getenv("DATABASE_DSN")
	if dsnEnv != "" && viper.GetString("DATABASE_DSN") == "" {
		viper.Set("DATABASE_DSN", dsnEnv)
	}

	grpcEnv := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcEnv != "" && viper.GetString("GRPC_SERVER_ADDRESS") == "" {
		viper.Set("GRPC_SERVER_ADDRESS", grpcEnv)
	}
}

// [.'.]:> 🔒 Aplica valores de respaldo directamente desde variables de entorno
// [.'.]:> como última red de seguridad para los campos que aún estén vacíos
// [.'.]:> después de procesar el archivo y las variables a través de viper
func applyFallbackEnvVars(config *EventConfiguration) {
	if config.DatabaseDriver == "" {
		config.DatabaseDriver = os.Getenv("DATABASE_DRIVER")
	}
	if config.DatabaseDsn == "" {
		config.DatabaseDsn = os.Getenv("DATABASE_DSN")
	}
	if config.GRPCServerAddress == "" {
		config.GRPCServerAddress = os.Getenv("GRPC_SERVER_ADDRESS")
	}
}
