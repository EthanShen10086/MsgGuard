# Data Governance

## Data Sources & Licenses
| Source | License | Citation |
|--------|---------|----------|
| UCI SMS Spam | Public domain | UCI ML Repository |
| FBS SMS (GitHub) | Research | CCS'20 paper |
| HF Multilingual | See HF card | dataset card |
| Seed data (repo) | MIT (project) | `ml/datasets/seed/` |

## PII Handling
- `ml/pipeline/clean.py` redacts phone numbers and ID numbers
- User feedback stores hash+label by default; full text only with opt-in
- Gateway middleware applies PII redaction on logs

## Retention
- Feedback samples: 90 days rolling (configurable)
- Crash reports: opt-in, 30 days
- Model artifacts: versioned in S3/filesystem indefinitely

## Cost
Run `python ml/costs/estimate.py` for monthly estimate.
