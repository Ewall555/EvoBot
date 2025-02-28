#!/usr/bin/env bash
#
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
cd "${BASE_DIR}" || exit 1
. "${BASE_DIR}/utils.sh"
yq="${BASE_DIR}/../statics/bin/yq"
chmod +x "${yq}"

function set_docker_subnet() {
  ENV_FILE="${BASE_DIR}/../compose/.env" 
  read -p "是否使用默认Docker网段(192.168.250.0/24)? (y/n): " use_default
  if [ -z "$use_default" ] || [ "$use_default" == "y" ]; then
    DOCKER_SUBNET="192.168.250.0/24"
  else
    read -p "Enter DOCKER SUBNET: " DOCKER_SUBNET
    echo
  fi

  echo "DOCKER_SUBNET=$DOCKER_SUBNET" >> "$ENV_FILE"

  echo "DOCKER_SUBNET information updated"
}


function set_db() {
  ENV_FILE="${BASE_DIR}/../compose/.env" 
  read -p "Do you want to use default values for DB connection? (y/n): " use_default
  if [ -z "$use_default" ] || [ "$use_default" == "y" ]; then
    DB_HOST="postgres"
    DB_PORT=5432
    DB_PASSWORD=$(openssl rand -base64 12)
    DB_USERNAME="postgres"

    DB_DATABASE="evobot"
  else
    read -p "Enter DB Host: " DB_HOST
    echo
    read -p "Enter DB User: " DB_USERNAME
    echo
    read -p "Enter DB PORT: " DB_PORT
    echo
    read -s -p "Enter DB Password: " DB_PASSWORD
    echo
    read -p "Enter DB Database: " DB_DATABASE
  fi
  echo "DB_HOST=$DB_HOST" >> "$ENV_FILE"
  echo "DB_PASSWORD=$DB_PASSWORD" >> "$ENV_FILE"
  echo "DB_USERNAME=$DB_USERNAME" >> "$ENV_FILE"
  echo "DB_DATABASE=$DB_DATABASE" >> "$ENV_FILE"
  echo "DB_PORT=$DB_PORT" >> "$ENV_FILE"

  echo "Database connection information updated in $ENV_FILE"
}

function set_redis() {
  ENV_FILE="${BASE_DIR}/../compose/.env" 
  read -p "Do you want to use default values for Redis connection? (y/n): " use_default
  if [ -z "$use_default" ] || [ "$use_default" == "y" ]; then
    REDIS_HOST="redis"
    REDIS_PASSWORD=$(openssl rand -base64 12)
    REDIS_PORT=6379
  else
    read -p "Enter Redis Host: " REDIS_HOST
    echo
    read -s -p "Enter Redis Password: " REDIS_PASSWORD
    echo
    read -p "Enter Redis Port: " REDIS_PORT
  fi

  echo "REDIS_HOST=$REDIS_HOST" >> "$ENV_FILE"
  echo "REDIS_PORT=$REDIS_PORT" >> "$ENV_FILE"
  echo "REDIS_PASSWORD=$REDIS_PASSWORD" >> "$ENV_FILE"
  mkdir -p "${VOLUME_DIR}/config/redis"
  REDIS_CONF="${VOLUME_DIR}/config/redis/redis.conf"
  echo "bind 0.0.0.0" > "$REDIS_CONF"
  echo "port $REDIS_PORT" >> "$REDIS_CONF"
  echo "requirepass $REDIS_PASSWORD" >> "$REDIS_CONF"
  echo "daemonize yes" >> "$REDIS_CONF"
  echo "Redis connection information updated"
  echo "Redis configuration file generated at: $REDIS_CONF"

}

function set_port() {
  ENV_FILE="${BASE_DIR}/../compose/.env" 
  read -p "Do you want to use default port 24916 (y/n): " use_default
  if [ -z "$use_default" ] || [ "$use_default" == "y" ]; then
    EVO_PORT=24916
  else
    read -p "Enter PORT: " EVO_PORT
    echo
  fi

  echo "EVO_PORT=$EVO_PORT" >> "$ENV_FILE"
  echo "Port information updated in $ENV_FILE"
}

function set_vol() {
  ENV_FILE="${BASE_DIR}/../compose/.env" 
  read -p "Do you want to use default values for Volume path? (y/n): " use_default
  if [ -z "$use_default" ] || [ "$use_default" == "y" ]; then
    VOLUME_DIR="/data/evobot"
  else
    read -p "Enter vol path: " VOLUME_DIR
    echo
  fi
  if [ ! -d "$VOLUME_DIR" ]; then
    mkdir -p "$VOLUME_DIR" 
    chmod 755 "$VOLUME_DIR"  
  fi
  echo "VOLUME_DIR=$VOLUME_DIR" > "$ENV_FILE"
  mkdir -p "${VOLUME_DIR}/config/"
  cp -p "${BASE_DIR}/config-example.yaml" "${VOLUME_DIR}/config/config.yaml"
  echo "Volume path information updated"
  export VOLUME_DIR=$VOLUME_DIR
}

function set_config_file() {
  source "${BASE_DIR}/../compose/.env"
  config_file="${VOLUME_DIR}/config/config.yaml"
  $yq -i ".db.host = \"${DB_HOST}\" | .db.username = \"${DB_USERNAME}\" | .db.password = \"${DB_PASSWORD}\" | .db.database = \"${DB_DATABASE}\" | .db.port = ${DB_PORT} | .redis.host = \"${REDIS_HOST}\" | .redis.password = \"${REDIS_PASSWORD}\" | .redis.port = ${REDIS_PORT}" $config_file
}

function main() {
  if set_vol; then
    echo_done
  fi  
  if set_port; then
    echo_done
  fi  
  if set_db; then
    echo_done
  fi
  if set_redis; then
    echo_done
  fi
  if set_config_file; then
    echo_done
  fi 
  if set_docker_subnet; then
    echo_done
  fi
}

if [[ "$0" == "${BASH_SOURCE[0]}" ]]; then
  main
fi