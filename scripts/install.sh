#!/bin/bash

set -e

read -p "both Go & Make must be installed before proceeding. Have you installed both? y/N " yn

case $yn in
y | Y)
  echo "proceeding with installation"
  ;;
*)
  echo "please install golang and make first, as they are requirements."
  exit 1
  ;;
esac

BASEDIR=$(dirname "$0")
NAME="siphon"
HOMEPATH=${HOME}/.siphon

if [ ! -d "${HOMEPATH}" ]; then
  echo "creating home folder..."
  mkdir "${HOMEPATH}"
  echo "created home folder at ${HOMEPATH}"
else
  echo "home folder ${HOMEPATH} already exists - exiting"
  exit 1
fi

echo "generating x509 key pair..."
bash "${BASEDIR}/keygen.sh" "$NAME"

OLDPATH=${PWD}
SRCPATH=$(cd "${BASEDIR}/.." && pwd)

echo "moving key pair to home folder..."
mv "${NAME}.crt" "${HOMEPATH}/${NAME}.crt"
mv "${NAME}.key" "${HOMEPATH}/${NAME}.key"
echo "done"

echo "copying source to siphon home folder..."
cp -r "${SRCPATH}" "${HOMEPATH}"
echo "copied"
echo "generating configuration file..."

CONFIG="${HOMEPATH}/.siphon.yaml"
echo "refreshrate: 1" >> "$CONFIG"
echo "cert_file: ${HOMEPATH}/${NAME}.crt" >> "$CONFIG"
echo "key_file: ${HOMEPATH}/${NAME}.key" >> "$CONFIG"

echo "building the binary - please wait..."
make build
mv siphon ${HOME}/.local/bin
if ! command -v siphon &> /dev/null; then
  mkdir -p "${HOME}/.local/bin" 2>/dev/null
  echo 'export PATH=$PATH:$HOME/.local/bin' >> "${HOME}/.profile"
fi
echo "done"
echo "installation complete"
echo "logout and log back in for changes to take effect"
cd "${OLDPATH}"
