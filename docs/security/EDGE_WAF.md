# Edge WAF & DDoS

MsgGuard API should sit behind an edge proxy with WAF in production. Application code does not replace Cloudflare/AWS WAF.

## Recommended stack

| Layer | Tool | Purpose |
|-------|------|---------|
| DNS + CDN | Cloudflare / AWS CloudFront | TLS termination, caching for static site & model CDN |
| WAF | Cloudflare WAF / AWS WAF | OWASP rules, bot fight, geo block |
| Rate limit | Edge + Gateway `ratelimiter` | Burst at edge; per-IP RPS in gateway |
| mTLS | Ingress + Caddy (Tier 4) | Admin API client certificates |

## Helm ingress annotations (nginx / Cloudflare)

Set in `values-prod.yaml`:

```yaml
ingress:
  waf:
    enabled: true
    annotations:
      nginx.ingress.kubernetes.io/limit-rps: "50"
      nginx.ingress.kubernetes.io/limit-burst-multiplier: "3"
      # Cloudflare: proxy orange-cloud + WAF rules in dashboard
```

## Minimum production rules

1. Block `POST /api/v1/auth/token` at edge when bootstrap disabled.
2. Challenge or rate-limit `POST /api/v1/feedback` and `/api/v1/analytics`.
3. Allowlist Apple webhook IPs for `POST /api/v1/webhooks/appstore` (or verify JWS only).
4. Enable Bot Fight Mode on marketing site; API on separate hostname (`api.msgguard.app`).

## Consumer vs admin

- **Mobile clients** — device token + App Attest (optional); no SSO.
- **Admin console** — OIDC SSO (`OIDC_*` env); WAF still applies but allow corporate IdP redirect URIs.

See [auth-security.md](../architecture/auth-security.md) and Cookbook 27.
