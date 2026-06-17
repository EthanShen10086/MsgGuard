# Security Policy

## Reporting a vulnerability

**Do not** open public GitHub issues for security vulnerabilities.

Report privately via **[GitHub Private Vulnerability Reporting](https://docs.github.com/en/code-security/security-advisories/working-with-repository-security-advisories/configuring-private-vulnerability-reporting-for-a-repository)**:

1. Open this repository on GitHub
2. Go to **Security** → **Report a vulnerability**

## What to include

- Affected endpoint (gateway, iOS extension, admin SPA, model CDN)
- Steps to reproduce
- Data impact (SMS body, feedback samples, admin tokens)

## Scope

- Gateway API (`/api/v1/*`)
- Admin OIDC and device tokens
- Model download URLs
- iOS App Group on-device data

## Response SLA

| Stage | Target |
|-------|--------|
| Acknowledgment | 72 hours |
| Initial triage | 7 days |
| Fix or mitigation plan | 30 days (severity-dependent) |

See also [docs/security/THREAT_MODEL.md](docs/security/THREAT_MODEL.md).
