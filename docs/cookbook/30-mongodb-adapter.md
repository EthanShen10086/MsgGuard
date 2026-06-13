# Cookbook 30 — MongoDB Port Adapter

## 切换步骤

1. 启动 MongoDB（Tier 1 compose 或本地）：
   ```bash
   docker run -d --name msgguard-mongo -p 27017:27017 \
     -e MONGO_INITDB_ROOT_USERNAME=msgguard \
     -e MONGO_INITDB_ROOT_PASSWORD=msgguard mongo:7
   ```

2. 合并配置：
   ```bash
   # database.driver=mongodb, dsn 见 deploy/config.mongodb.yaml
   export DATABASE_DSN="mongodb://msgguard:msgguard@localhost:27017/?authSource=admin"
   export DATABASE_DRIVER=mongodb
   ```

3. 启动 Gateway（Container 自动 wire `pkg/adapters/mongodb`）：
   ```bash
   cd services/gateway && CONFIG_PATH=../../deploy/config.yaml go run ./cmd/server
   ```

## 验收

```bash
# 写入 feedback
curl -X POST localhost:8080/api/v1/feedback \
  -H 'Content-Type: application/json' \
  -d '{"body":"mongo test","label":"ham"}'

# 需 token 读取
TOKEN=$(curl -s -X POST localhost:8080/api/v1/auth/token \
  -H 'Content-Type: application/json' \
  -d '{"roles":["admin"]}' | jq -r .access_token)
curl -H "Authorization: Bearer $TOKEN" localhost:8080/api/v1/feedback
```

Mongo 不可达时 Container **自动回退** memory store，无需改业务代码。
