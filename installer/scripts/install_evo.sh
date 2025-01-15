#!/usr/bin/env bash
#
BASE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
cd "${BASE_DIR}" || exit 1
. "${BASE_DIR}/utils.sh"

function load_images() {
   images_dir="${BASE_DIR}/../statics/images"
   if [ -d "$images_dir" ] && [ -n "$(ls -A "$images_dir")" ]; then
       for file in "$images_dir"/*.tar; do
           if [ -f "$file" ]; then
               docker load -i "$file"
           fi
       done
   else
       echo_red "No Images"
   fi
}

function main() {
  echo_logo
  if ! bash "${BASE_DIR}/install_docker.sh"; then
    exit 1
  fi
  echo_green ">>> $(gettext 'complete')"
  echo_yellow "3. $(gettext 'Load Images')"
  if ! load_images ; then
      exit 1
  fi
  echo_green ">>> $(gettext 'complete')"
  echo_yellow "4. $(gettext 'Setup Config')"
  if ! bash "${BASE_DIR}/config.sh"; then
    exit 1
  fi
  echo_green ">>> $(gettext 'complete')"  
}

if [[ "$0" == "${BASH_SOURCE[0]}" ]]; then
  main
fi
