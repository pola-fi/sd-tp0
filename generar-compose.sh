#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "Uso: $0 <archivo-salida> <cantidad-clientes>" >&2
  exit 1
fi

output_file="$1"
client_count="$2"

if ! [[ "$client_count" =~ ^[0-9]+$ ]]; then
  echo "Error: la cantidad de clientes debe ser un entero no negativo" >&2
  exit 1
fi

NOMBRES=("Juan" "Maria" "Carlos" "Ana" "Luis" "Laura" "Pedro" "Sofia" "Diego" "Valentina")
APELLIDOS=("Garcia" "Rodriguez" "Lopez" "Martinez" "Gonzalez" "Perez" "Sanchez" "Ramirez" "Torres" "Diaz")

generar_dni() {
  echo $((10000000 + RANDOM % 90000000))
}

generar_fecha() {
  year=$((1960 + RANDOM % 41))
  month=$((1 + RANDOM % 12))
  day=$((1 + RANDOM % 28))
  printf "%04d-%02d-%02d" $year $month $day
}

generar_numero() {
  echo $((1000 + RANDOM % 9000))
}

# ej6: cada cliente con AGENCY_ID=N lee .data/agency-N.csv en el container (volumen ./.data).
# Si falta el CSV, crear un stub mínimo (no sobrescribe archivos ya existentes).
ensure_agency_csv() {
  local id="$1"
  mkdir -p .data
  local f=".data/agency-${id}.csv"
  if [[ -f "$f" ]]; then
    return 0
  fi
  cat >"$f" <<'EOF'
A,B,00000000,2000-01-01,1000
A,B,00000001,2000-01-01,1001
A,B,00000002,2000-01-01,1002
EOF
}

for ((j = 1; j <= client_count; j++)); do
  ensure_agency_csv "$j"
done

{
  cat <<'YAML'
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /app/main.py
    environment:
      - PYTHONUNBUFFERED=1
    volumes:
      - ./server/config.ini:/app/config.ini
    networks:
      - testing_net

YAML

  for ((i = 1; i <= client_count; i++)); do
    agency_id=$i
    nombre_idx=$((RANDOM % ${#NOMBRES[@]}))
    apellido_idx=$((RANDOM % ${#APELLIDOS[@]}))
    nombre=${NOMBRES[$nombre_idx]}
    apellido=${APELLIDOS[$apellido_idx]}
    documento=$(generar_dni)
    nacimiento=$(generar_fecha)
    numero=$(generar_numero)

    cat <<YAML
  client${i}:
    container_name: client${i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=${i}
      - AGENCY_ID=${agency_id}
      - NOMBRE=${nombre}
      - APELLIDO=${apellido}
      - DOCUMENTO=${documento}
      - NACIMIENTO=${nacimiento}
      - NUMERO=${numero}
    volumes:
      - ./client/config.yaml:/app/config.yaml
      - ./.data:/app/.data
    networks:
      - testing_net
    depends_on:
      - server

YAML
  done

  cat <<'YAML'
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
YAML
} > "$output_file"

echo "Docker Compose generado en $output_file con $client_count cliente(s)."

