#!/bin/bash
TINYB_FILE="${PWD}/data/tinybird/.tinyb"
if ! command -v jq &> /dev/null; then
    echo "🚨 jq no está instalado. Instálalo para continuar."
    exit 1
fi

if [[ -f "$TINYB_FILE" ]]; then
    TB_TOKEN=$(jq -r '.token' "$TINYB_FILE")
else
    echo "🚨 El archivo .tinyb no se encuentra en el directorio actual (${PWD}/data/tinybird)."
    exit 1
fi

echo '** 🚨🚨🚨 FORCE 🚨🚨🚨 starting upload process'
docker run -v ${PWD}/data/tinybird:/mnt/data tinybirdco/tinybird-cli-docker bash -c "
tb auth --token $TB_TOKEN && \
tb push /mnt/data/datasources/*.datasource --fixtures --force
"
