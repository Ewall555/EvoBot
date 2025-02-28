#!/usr/bin/env bash
#

BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

function check_root() {
  [[ "$(id -u)" == 0 ]]
}

echo_yellow() {
    YELLOW='\033[1;33m'
    NC='\033[0m'        
    echo -e "${YELLOW}$1${NC}" 
}

echo_red() {
    RED='\033[0;31m'  
    NC='\033[0m'      

    echo -e "${RED}$1${NC}" 
}

echo_green() {
    GREEN='\033[0;32m'  
    NC='\033[0m'        

    echo -e "${GREEN}$1${NC}" 
}

echo_done() {
  echo_green "done"
}

function prepare_compose_bin() {
  chown -R root:root "${BASE_DIR}/../statics/docker/docker-compose"
  chmod +x "${BASE_DIR}/../statics/docker/docker-compose"
}

function echo_purple() {
    echo -e "\033[0;35m$@\033[0m"
}

function echo_logo() {
  cat <<"EOF"
    ______            ____        __     ___    __      __                   ____   ___ ____ 
   / ____/   ______  / __ )____  / /_   /   |  / /___  / /_  ____ _   _   __/ __ \ <  // __ \ 
  / __/ | | / / __ \/ __  / __ \/ __/  / /| | / / __ \/ __ \/ __ `/  | | / / / / / / // / / / 
 / /___ | |/ / /_/ / /_/ / /_/ / /_   / ___ |/ / /_/ / / / / /_/ /   | |/ / /_/ / / // /_/ /  
/_____/ |___/\____/_____/\____/\__/  /_/  |_/_/ .___/_/ /_/\__,_/    |___/\____(_)_(_)____/   
                                             /_/                                             

EOF

}

function get_config() {
env_file="${BASE_DIR}/../compose/.env"
key=$1
default=${2-''}
value=$(grep "^${key}=" "${env_file}" | awk -F= '{ print $2 }' | awk -F' ' '{ print $1 }' | tail -1)
if [[ -z "$value" ]];then
value="$default"
fi
echo "${value}"
}

function get_docker_compose_cmd_line() {
  db_host=$(get_config DB_HOST)
  redis_host=$(get_config REDIS_HOST)
  cmd="docker compose -f compose/evobot.yml -f compose/evo_net.yml"
  if [[ "${db_host}"="postgres" ]]; then
    cmd+=" -f compose/postgresql.yml"
  fi
  if [[ "${redis_host}"="redis" ]]; then
    cmd+=" -f compose/redis.yml"
  fi
  echo "${cmd}"  
}