# CDN Deployment

Host static assets on Cloudflare or Aliyun CDN:

- `cdn.msgguard.app/rules/{version}/rules.json` — rule bundles with ETag
- `cdn.msgguard.app/models/{version}/spam_classifier.mlmodel` — Core ML models

Configure cache TTL: rules 1h, models 24h.

Origin: S3/R2 bucket synced from `ml/output/`.
