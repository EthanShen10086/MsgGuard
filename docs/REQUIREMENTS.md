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
- [x] L2 CoreMLClassifier stub in FilterEngine

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
