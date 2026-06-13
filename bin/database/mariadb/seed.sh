#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

# Intentamos cargar el entorno si existe
ENVIRONMENT_FILE="bin/shared/environment.sh"
if [ -f "$ENVIRONMENT_FILE" ]; then
    source "$ENVIRONMENT_FILE"
fi

# Parámetros con valores por defecto
SERVER=${1:-"localhost:30000"}
NUM_RECORDS=${2:-10}
SERVICE="event.Eventservice"

# Función secuencial para transformar números a letras (1=A, 2=B... 27=AA)
int_to_letters() {
    local num=$1
    local letters=""
    while [ "$num" -gt 0 ]; do
        local remainder=$(( (num - 1) % 26 ))
        local char=$(printf "\\$(printf '%03o' $((65 + remainder)))")
        letters="${char}${letters}"
        num=$(( (num - 1) / 26 ))
    done
    echo "$letters"
}

# Nombres alfabéticos estáticos para los tipos de eventos
EVENT_TYPES=("Alpha" "Bravo" "Charlie")

echo "🚀 Iniciando script de Seed para gRPC en $SERVER..."
echo "📦 Se crearán $NUM_RECORDS suscripciones y eventos (100% letras)."

if ! command -v grpcurl &> /dev/null; then
    echo "❌ grpcurl no está instalado. Ejecuta 'make support-install-grpc-tools' o instálalo manualmente."
    exit 1
fi

for i in $(seq 1 "$NUM_RECORDS"); do
    echo "--------------------------------------------------"
    echo "🔄 Iteración $i de $NUM_RECORDS"
    
    # Generamos un sufijo único hecho únicamente de letras
    LETTERS_ID=$(int_to_letters "$i")
    
    # Construcción de variables usando solo caracteres alfabéticos
    EVENT_NAME="LoadTestEvent${EVENT_TYPES[$((i % 3))]}" 
    SOURCE_NAME="LoadTestSource"
    SUB_NAME="SubscriberWorker${LETTERS_ID}"
    
    echo "1️⃣ Creando Suscripción para $SUB_NAME..."
    SUB_PAYLOAD=$(cat <<EOF
{
  "subscriber_name": "$SUB_NAME",
  "event_name": "$EVENT_NAME",
  "source": "$SOURCE_NAME"
}
EOF
)
    grpcurl -plaintext -d "$SUB_PAYLOAD" "$SERVER" "$SERVICE/CreateSubscription"
    
    echo "2️⃣ Creando Evento..."
    # El payload ahora envía la iteración como texto ("A", "B") y remueve caracteres numéricos
    EVENT_PAYLOAD=$(cat <<EOF
{
  "slug": "$EVENT_NAME",
  "source": "$SOURCE_NAME",
  "payload": "{\"iteration\": \"$LETTERS_ID\", \"message\": \"Mensaje autogenerado en letras\"}"
}
EOF
)
    grpcurl -plaintext -d "$EVENT_PAYLOAD" "$SERVER" "$SERVICE/CreateEvent"
done

echo "--------------------------------------------------"
echo "🎉 ¡Poblado de datos finalizado con éxito!"