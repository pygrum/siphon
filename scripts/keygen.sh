#!/bin/bash

set -e

NAME="$1"

# Set cert names to siphon if name not provided
if [ $# -lt 1 ]; then
  NAME="siphon"
fi

openssl req -newkey rsa:4096 \
            -x509 \
            -sha256 \
            -days 3650 \
            -nodes \
            -out "${NAME}.crt" \
            -keyout "${NAME}.key" \
            -subj "/C=US/ST=New York/L=New York City/O=${NAME}/OU=${NAME}/CN=www.${NAME}.com"

echo "key pair saved at ${NAME}.crt (certificate) and ${NAME}.key (key)"
