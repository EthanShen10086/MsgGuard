#!/usr/bin/env bash
# Generate dev mTLS CA + server + client certs for MsgGuard API.
set -euo pipefail
DIR="$(cd "$(dirname "$0")" && pwd)"
OUT="${DIR}/certs"
mkdir -p "$OUT"
DAYS="${MTLS_CERT_DAYS:-825}"

openssl genrsa -out "$OUT/ca.key" 4096
openssl req -x509 -new -nodes -key "$OUT/ca.key" -sha256 -days "$DAYS" \
  -out "$OUT/ca.crt" -subj "/CN=MsgGuard Dev CA"

openssl genrsa -out "$OUT/server.key" 2048
openssl req -new -key "$OUT/server.key" -out "$OUT/server.csr" \
  -subj "/CN=api.msgguard.app"
openssl x509 -req -in "$OUT/server.csr" -CA "$OUT/ca.crt" -CAkey "$OUT/ca.key" \
  -CAcreateserial -out "$OUT/server.crt" -days "$DAYS" -sha256

openssl genrsa -out "$OUT/client.key" 2048
openssl req -new -key "$OUT/client.key" -out "$OUT/client.csr" \
  -subj "/CN=msgguard-admin"
openssl x509 -req -in "$OUT/client.csr" -CA "$OUT/ca.crt" -CAkey "$OUT/ca.key" \
  -CAcreateserial -out "$OUT/client.crt" -days "$DAYS" -sha256

rm -f "$OUT"/*.csr "$OUT"/ca.srl
echo "Generated certs in $OUT"
echo "  ca.crt / server.crt+key / client.crt+key"
