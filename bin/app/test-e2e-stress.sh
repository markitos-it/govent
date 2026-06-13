#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'
ENVIRONMENT_FILE="bin/shared/environment.sh"
source $ENVIRONMENT_FILE

SERVER="localhost:30000"
SERVICE="event.Eventservice"
SOURCE="EventSource"

# Constantes para evitar números mágicos
NUM_USERS=10
HOW_MANY_EVENTS_PER_USER=20
TOTAL_EVENTS=$((NUM_USERS * HOW_MANY_EVENTS_PER_USER))

echo "🧹 Cleaning up database..."
docker exec -it goevents-postgres psql -U admin -d goevents -c "TRUNCATE TABLE queue, events, subscriptions RESTART IDENTITY CASCADE;" > /dev/null

echo "🚀 Starting Isolated E2E gRPC Test..."
EVENT_IDS=()

# 1. Crear suscripciones y eventos por cada usuario
for i in $(seq 1 $NUM_USERS); do
    USER_SLUG="EventTest_user_$i"
    SUB_PAYLOAD=$(jq -n --arg name "user_$i" --arg slug "$USER_SLUG" --arg src "$SOURCE" \
        '{subscriber_name: $name, event_name: $slug, source: $src}')
    grpcurl -plaintext -d "$SUB_PAYLOAD" $SERVER $SERVICE/CreateSubscription > /dev/null 2>&1
    
    # 2. Crear eventos específicos para este SLUG
    for e in $(seq 1 $HOW_MANY_EVENTS_PER_USER); do
        CREATE_PAYLOAD=$(jq -n --arg slug "$USER_SLUG" --arg src "$SOURCE" --arg pld "{\"msg\": \"event_$e\"}" \
            '{slug: $slug, source: $src, payload: $pld}')
        
        RESP=$(grpcurl -plaintext -d "$CREATE_PAYLOAD" $SERVER $SERVICE/CreateEvent)
        EVENT_IDS+=("$(echo "$RESP" | jq -r '.id')")
    done
    echo -ne "🛠  Provisioning user_$i with $HOW_MANY_EVENTS_PER_USER events...\r"
done
echo -e "\n✅ $TOTAL_EVENTS events created across $NUM_USERS unique slugs."

# 3. PullMessages por cada usuario y validación
for i in $(seq 1 $NUM_USERS); do
    USER_SLUG="EventTest_user_$i"
    PULL_PAYLOAD=$(jq -n --arg slug "$USER_SLUG" --arg src "$SOURCE" '{event_name: $slug, source: $src}')
    
    PULL_RESP=$(grpcurl -plaintext -d "$PULL_PAYLOAD" $SERVER $SERVICE/PullMessages)
    COUNT=$(echo "$PULL_RESP" | jq '.messages | length // 0')
    
    if [ "$COUNT" -ne "$HOW_MANY_EVENTS_PER_USER" ]; then
        echo "❌ Failure: Expected $HOW_MANY_EVENTS_PER_USER messages for $USER_SLUG, found $COUNT."
        exit 1
    fi

    QUEUE_IDS=$(echo "$PULL_RESP" | jq -c '[.messages[].id]')
    ACK_PAYLOAD=$(jq -n --argjson ids "$QUEUE_IDS" '{queue_ids: $ids}')
    grpcurl -plaintext -d "$ACK_PAYLOAD" $SERVER $SERVICE/AckMessages > /dev/null
done
echo "✅ All $NUM_USERS isolated queues validated and acknowledged."

# 4. Limpieza de eventos
echo "🗑️  Deleting $TOTAL_EVENTS events..."
for i in "${!EVENT_IDS[@]}"; do
    DELETE_PAYLOAD=$(jq -n --arg id "${EVENT_IDS[$i]}" '{id: $id}')
    grpcurl -plaintext -d "$DELETE_PAYLOAD" $SERVER $SERVICE/DeleteEvent > /dev/null
    echo -ne "🔥 Deleting: $((i + 1))/$TOTAL_EVENTS\r"
done
echo -e "\n✅ All events deleted."
