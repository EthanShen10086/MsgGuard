# MsgGuard Requirements

## Phase 0 — Scaffold ✅
- [x] Monorepo structure
- [x] XcodeGen project.yml
- [x] SPM packages
- [x] MGError, MGLogger, AnalyticsManager
- [x] deploy/config.yaml
- [x] CI workflow

## Phase 1 — MVP ✅
- [x] FilterEngine L0 + L1
- [x] MessageFilterExtension
- [x] App Group sync via BlocklistStore
- [x] Main App: TabView, Onboarding, Dashboard, Rules, Samples, Stats, Help, Settings
- [x] zh-Hans + en localization
- [x] Feedback with TraceID

## Phase 2 — Core ML ✅
- [x] Python train_bayes.py, train_coreml.py, export_coreml.py
- [x] Dataset seed CSV
- [x] L2 CoreMLClassifier in FilterEngine

## Phase 3 — Backend ✅
- [x] Gateway with voxera-kit middleware
- [x] rules/model microservices
- [x] docker-compose
- [x] AASA template

## Phase 4 — LLM Hybrid ✅
- [x] classify/defer endpoint
- [x] Qwen/DeepSeek router
- [x] MLModelHealthMonitor API
- [x] config.yaml cloud_llm flag

## Phase 5 — Product ✅
- [x] StoreKit SubscriptionView
- [x] Widget extension
- [x] Elder Mode in DesignSystem
- [x] App Store docs template

## Phase 6 — Admin & Commercial Surface
- [ ] Admin web SPA (Dashboard, Feedback, Models, Flags, Quota)
- [ ] Static pricing + status pages
- [ ] Architecture doc split + threat model
- [ ] Honest commercial readiness tracking

## Phase 7 — SRE & Compliance
- [ ] SLO definitions + error budgets
- [ ] Privacy data deletion API
- [ ] Prometheus alert templates wired in compose/Helm
- [ ] Production auth hardening (bootstrap off, OIDC planned)

## Phase 8 — Analytics & Growth
- [ ] Event taxonomy v2 rollout (iOS + backend)
- [ ] Subscription funnel metrics in weekly report
- [ ] Onboarding A/B via feature flags
- [ ] Shadow dashboard in admin web

## Phase 9 — Platform Expansion
- [ ] macOS Mail App Store release
- [ ] Android client scaffold
- [ ] Multi-region Helm values
- [ ] WAF / edge caching

## Phase 10 — Enterprise
- [ ] SSO/OIDC for admin
- [ ] Per-tenant quota + audit export
- [ ] SOC2-aligned logging retention
- [ ] Dedicated support SLA tier

## Phase 11 — ML Maturity
- [ ] Automated retrain on feedback threshold
- [ ] Per-locale production models for top 5 locales
- [ ] Online learning shadow (no auto-deploy)
- [ ] LLM cost caps per user

## Phase 12 — Mature Product
- [ ] 99.9% gateway SLO sustained
- [ ] <1% FPR on adversarial benchmark
- [ ] Cross-platform Pro entitlement
- [ ] Public status page with real probes
- [ ] App Store featuring-ready creative + localization
