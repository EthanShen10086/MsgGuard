# Changelog

All notable changes to MsgGuard. Format based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

### Added
- Organization-standard governance (LICENSE, SECURITY, CODEOWNERS, dependabot)
- Go unit tests for auth RBAC and OIDC admin allowlist
- CI verify job and blocking security scans

## [0.1.0] - 2026-06-17

### Added
- Hybrid SMS filter: L0 heuristic, L1 Bayes, L2 CoreML, optional cloud LLM
- Go gateway with RBAC, device tokens, OIDC admin SSO
- Redis token revocation, App Store webhook skeleton
- ML flywheel with benchmark CI gate
- iOS Message Filter extension + macOS Mail host
- Admin web SPA, Helm charts, staging/release CI workflows
- SLO, runbooks, threat model, 33-chapter cookbook

[Unreleased]: https://github.com/EthanShen10086/MsgGuard/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/EthanShen10086/MsgGuard/releases/tag/v0.1.0
