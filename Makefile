.DEFAULT_GOAL := help
.PHONY: help start test test-e2e proto db-start db-stop db-create db-drop db-seed spanner-start spanner-stop spanner-create spanner-drop format lint support-install-linter support-uninstall-linter build appsec-install appsec-uninstall appsec-test tidy

help:
	@echo ""
	@echo "['.']:> =================================================="
	@echo "['.']:> 🚀 MARKITOS-IT GO VENTS COMMANDS"
	@echo "['.']:> =================================================="
	@printf "  \033[36m%-24s\033[0m %s\n" "help" "Muestra este menú de ayuda interactivo"
	@printf "  \033[36m%-24s\033[0m %s\n" "build" "Genera el binario final de la aplicación en la carpeta dist"
	@printf "  \033[36m%-24s\033[0m %s\n" "start" "Inicia la aplicación localmente"
	@printf "  \033[36m%-24s\033[0m %s\n" "test" "Ejecuta la suite de pruebas de la aplicación"
	@printf "  \033[36m%-24s\033[0m %s\n" "test-e2e" "Ejecuta las pruebas End-to-End mediante gRPC"
	@printf "  \033[36m%-24s\033[0m %s\n" "proto" "Genera los archivos de código a partir de los de gRPC .proto"
	@printf "  \033[36m%-24s\033[0m %s\n" "lint" "Analiza el código Go con golangci-lint"
	@printf "  \033[36m%-24s\033[0m %s\n" "lint-fix" "Formatea el código Go automáticamente (gofmt, goimports)"
	@printf "  \033[36m%-24s\033[0m %s\n" "tidy" "Limpia y actualiza las dependencias de Go (go mod tidy)"
	@printf "  \033[36m%-24s\033[0m %s\n" "db-create" "Crea la base de datos en PostgreSQL"
	@printf "  \033[36m%-24s\033[0m %s\n" "db-drop" "Elimina (drop) la base de datos por completo"
	@printf "  \033[36m%-24s\033[0m %s\n" "db-start" "Inicia el entorno de la base de datos"
	@printf "  \033[36m%-24s\033[0m %s\n" "db-stop" "Detiene el entorno de la base de datos"
	@printf "  \033[36m%-24s\033[0m %s\n" "db-seed" "Puebla la base de datos con datos de prueba a través de gRPC"
	@printf "  \033[36m%-24s\033[0m %s\n" "spanner-start" "Inicia el entorno del emulador de Spanner"
	@printf "  \033[36m%-24s\033[0m %s\n" "spanner-stop" "Detiene el entorno del emulador de Spanner"
	@printf "  \033[36m%-24s\033[0m %s\n" "spanner-create" "Crea la base de datos en Spanner"
	@printf "  \033[36m%-24s\033[0m %s\n" "spanner-drop" "Elimina (drop) la base de datos en Spanner por completo"
	@printf "  \033[36m%-24s\033[0m %s\n" "support-install-linter" "Instala la herramienta golangci-lint"
	@printf "  \033[36m%-24s\033[0m %s\n" "support-uninstall-linter" "Desinstala la herramienta golangci-lint"
	@printf "  \033[36m%-24s\033[0m %s\n" "appsec-install" "Instala herramientas de seguridad (Snyk, Gitleaks)"
	@printf "  \033[36m%-24s\033[0m %s\n" "appsec-uninstall" "Desinstala las herramientas de seguridad"
	@printf "  \033[36m%-24s\033[0m %s\n" "appsec-test" "Ejecuta las pruebas de seguridad (Snyk, Gitleaks)"
	@echo "['.']:> =================================================="
	@echo ""


start:
	bash bin/app/start.sh

test:
	bash bin/app/test.sh

test-e2e:
	bash bin/app/test_e2e_grpc.sh

test-verbose:
	bash bin/app/test-verbose.sh

proto:
	bash bin/app/proto.sh

db-start:
	bash bin/database/postgres/start.sh

db-stop:
	bash bin/database/postgres/stop.sh

db-create:
	bash bin/database/postgres/create.sh

db-drop:
	bash bin/database/postgres/drop.sh

db-seed:
	bash bin/database/postgres/seed_data.sh

spanner-start:
	bash bin/database/spanner/start.sh

spanner-stop:
	bash bin/database/spanner/stop.sh

spanner-create:
	bash bin/database/spanner/create.sh

spanner-drop:
	bash bin/database/spanner/drop.sh

lint:
	bash bin/code/lint.sh

lint-fix:
	bash bin/code/lint-fix.sh

support-install-linter:
	bash bin/support/install-linter.sh

support-uninstall-linter:
	bash bin/support/uninstall-linter.sh

build:
	bash bin/app/build.sh

clean:
	bash bin/app/clean.sh

appsec-test:
	bash bin/appsec/test.sh

appsec-install:
	bash bin/appsec/install.sh

appsec-uninstall:
	bash bin/appsec/uninstall.sh

appsec-pre-commit:
	bash bin/appsec/pre-commit.sh

appsec-pre-push:
	bash bin/appsec/pre-push.sh

tidy:
	go mod tidy