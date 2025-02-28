#!/usr/bin/env bash
#
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
echo $BASE_DIR
DOCKER_CONFIG="/etc/docker/daemon.json"
. "${BASE_DIR}/utils.sh"
cd "${BASE_DIR}" || exit 1

function install_docker() {
  if [[ ! -d "/usr/local/bin" ]]; then
    mkdir -p /usr/local/bin
  fi
  tar --no-same-owner --strip-components=1 -xf "${BASE_DIR}/../statics/docker/docker.tar.gz" -C /usr/local/bin/ || {
    echo_red "$(gettext 'Failed to extract docker.tar.gz')"
    exit 1
  }
  \cp -f "${BASE_DIR}/../statics/docker/docker.service" /etc/systemd/system/
}

function install_compose() {
  if [[ ! -f "/usr/local/libexec/docker/cli-plugins/docker-compose" ]]; then
    if [[ ! -d "/usr/local/libexec/docker/cli-plugins" ]]; then
      mkdir -p /usr/local/libexec/docker/cli-plugins
    fi
    \cp -f "${BASE_DIR}/../statics/docker/docker-compose" /usr/local/libexec/docker/cli-plugins/
    sudo ln -s /usr/local/libexec/docker/cli-plugins/docker-compose /usr/local/bin/docker-compose
  fi
}

function start_docker() {
  if command -v systemctl&>/dev/null; then
    systemctl daemon-reload
    systemctl enable docker
    systemctl start docker
  fi
  if ! docker ps &>/dev/null; then
    echo_failed
    exit 1
  fi
}

function check_docker_start() {
  if ! docker ps &>/dev/null; then
    start_docker
  fi
}

function check_docker_compose() {
  if ! docker compose version &>/dev/null; then
    echo_red "Failed"
    echo_yellow "retry docker-compose"
  else
    echo_yellow "Success"
  fi
}

function main() {
  echo_yellow "1. $(gettext 'Install Docker')"
  if command -v docker &>/dev/null && command -v docker-compose &>/dev/null; then
    echo_green "Docker 和 Docker Compose 已安装"
  else
    install_docker
    install_compose
  fi
  echo_yellow "\n2. $(gettext 'Start Docker')"
  check_docker_start
  check_docker_compose
}

if [[ "$0" == "${BASH_SOURCE[0]}" ]]; then
  main
fi