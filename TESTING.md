# Запуск тестов

## Все сервисы

```bash
# User service (8 unit-тестов)
cd services/user && go test ./tests/... -v

# Session service (6 unit-тестов)
cd services/session && go test ./tests/... -v

# Parking service (8 unit-тестов) — новые
cd services/parking && go test ./tests/... -v

# Notification service (8 unit-тестов) — новые
cd services/notification/tests && go test . -v

# Gateway — auth + middleware (15 unit-тестов) — новые
cd services/gateway && go test ./tests/... -v
```

## Все разом

```bash
cd parking-system
for svc in user session parking gateway; do
  echo "=== $svc ==="
  (cd services/$svc && go test ./tests/... -v 2>&1)
done
echo "=== notification ==="
(cd services/notification/tests && go test . -v 2>&1)
```

## Покрытие

```bash
cd services/user    && go test ./tests/... -cover
cd services/session && go test ./tests/... -cover
cd services/parking && go test ./tests/... -cover
cd services/gateway && go test ./tests/... -cover
```

---

## Фронтенд

```bash
cd frontend
npm install
npm run dev          # http://localhost:3000
npm run build        # production build
```

Или через Docker:
```bash
docker compose up --build
# frontend → http://localhost:3000
# gateway  → http://localhost:8080
```
