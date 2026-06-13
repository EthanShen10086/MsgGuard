# mTLS Certificate Rotation

## Scope
- Caddy / nginx ingress client-auth CA
- Optional gateway admin path enforcement

## Rotation steps

1. Generate new CA + certs (`deploy/mtls/gen-certs.sh` or cert-manager)
2. Deploy new CA secret to ingress **alongside** old CA (dual trust)
3. Re-issue client certificates to all admin operators
4. Remove old CA from trust store after 7-day overlap
5. Verify: `curl --cert client.crt --key client.key --cacert ca.crt https://api.../health`

## Gateway header mode

When using Caddy `X-Client-Cert-Subject`:
- Rotation is at the **edge** only; gateway env unchanged
- Restart gateway after updating `MTLS_CLIENT_HEADER` if renamed

## Rollback

Set `MTLS_ADMIN_REQUIRED=false` in ConfigMap / compose env and redeploy.
