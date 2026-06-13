# Tier 4 — mTLS at Caddy edge + gateway admin header check

## Generate certificates (dev)

```bash
bash deploy/mtls/gen-certs.sh
```

Outputs under `deploy/mtls/certs/` (gitignored).

## Caddy (recommended)

Use `deploy/caddy/Caddyfile.mtls` — terminates TLS + requires client cert, forwards
`X-Client-Cert-Subject` to gateway.

```bash
export MTLS_ADMIN_REQUIRED=true
export MTLS_CLIENT_HEADER=X-Client-Cert-Subject
docker compose -f deploy/docker-compose.yml -f deploy/docker-compose.mtls.yml up -d
```

## Gateway middleware

When `MTLS_ADMIN_REQUIRED=true`:
- With `MTLS_CLIENT_HEADER` — validates header on `/api/v1/admin/*` (behind Caddy)
- Without header — requires `r.TLS.PeerCertificates` (gateway terminates TLS)

Public routes (`/health`, `/api/v1/classify`, `/api/v1/feedback`) are **not** mTLS-gated.

## Client test

```bash
curl --cert deploy/mtls/certs/client.crt --key deploy/mtls/certs/client.key \
  --cacert deploy/mtls/certs/ca.crt \
  https://localhost:8443/api/v1/admin/metrics/summary \
  -H "Authorization: Bearer $TOKEN"
```

## Production

- Use cert-manager / Vault for rotation — see `docs/runbooks/mtls-rotation.md`
- Helm: `values-mtls.yaml` ingress annotations for nginx mutual-auth
