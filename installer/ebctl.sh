#!/usr/bin/env bash
#
PROJECT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

cd "${PROJECT_DIR}" || exit 1

SCRIPT_DIR="${PROJECT_DIR}/scripts/"
. "${SCRIPT_DIR}/utils.sh"
action=${1-}
target=${2-}
args=("$@")

function help() {
  echo "$(gettext 'EvoBot Deployment Management Script')"
  echo
  echo "Usage: "
  echo "  ./evoctl.sh [COMMAND] [ARGS...]"
  echo "  ./evoctl.sh --help"
  echo
  echo "Installation Commands: "
  echo "  install           $(gettext 'Install EvoBot')"
  echo "  uninstall           $(gettext 'Uninstall EvoBot')"
  echo
  echo "Management Commands: "
  echo "  start             $(gettext 'Start     EvoBot')"
  echo "  stop              $(gettext 'Stop      EvoBot')"
  echo "  restart           $(gettext 'Restart   EvoBot')"
  echo
}

function start() {
  EXE=$(get_docker_compose_cmd_line)
  ${EXE} up -d
}

function stop() {
  EXE=$(get_docker_compose_cmd_line)
  ${EXE} down -v
}

function status() {
  EXE=$(get_docker_compose_cmd_line)
  ${EXE} ps
}

function main() {
  action=$1

  case "$action" in
  install)
    bash "${SCRIPT_DIR}/install_evo.sh"
    ;;
  start)
    start
    ;;
  restart)
    stop && start
    ;;
  stop)
    stop
    ;;
  status)
    status
    ;;
  uninstall)
    bash "${SCRIPT_DIR}/uninstall_evo.sh"
    ;;
  help)
    help
    ;;
  --help)
    help
    ;;
  -h)
    help
    ;;
  *)
    echo "No such command: $action"
    help
    ;;
  esac
}

main "$@"