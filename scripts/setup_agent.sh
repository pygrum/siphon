#!/bin/bash

set -e

BASEDIR=$(dirname "$0")
BINFILE="$1"
HOMEDIR="${HOME}/.siphon"
NAME="agent"

if [ $# -ne 1 ]; then
  echo "invalid usage"
  echo "usage: ${0} <path/to/agent/binary>"
  exit 1
fi

OLDDIR=${PWD}
cd "${BASEDIR}/.."

mkdir "agent_files"
cp "$BINFILE" "agent_files"
cd "agent_files"

echo "generating agent keys..."
bash ../scripts/keygen.sh "$NAME"

echo "done"

AGENT_FOLDER="$(md5sum "$BINFILE" | cut -f1 -d' ')"
mkdir "${HOMEDIR}/${AGENT_FOLDER}"
echo "created agent folder: ${HOMEDIR}/${AGENT_FOLDER} - saving agent certificate to folder..."
cp "${NAME}.crt" "${HOMEDIR}/${AGENT_FOLDER}"
echo "done. configure Siphon to use this certificate for this agent (agents set <id> --cert-file=${HOMEDIR}/${AGENT_FOLDER}/${NAME}.crt)"

echo "generating default configuration file..."
echo 'cache: true' >> config.yaml
echo "cert_file: /root/${NAME}.crt" >> config.yaml
echo "key_file: /root/${NAME}.key" >> config.yaml
echo "monitor_folders:" >> config.yaml
echo '  - path: "/tmp"' >> config.yaml
echo '    recursive: true' >> config.yaml
echo "saved as config.yaml"

echo "agent files saved to directory 'agent_files'"
echo "follow the guide at https://github.com/pygrum/siphon/blob/main/docs/DOCS.md to finish configuring the agent."
cd ${OLDDIR}
