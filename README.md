# Session Service Fix — что внутри

Архив с патчем для проекта `parking-system`. Применяется поверх существующих
файлов (заменяет). Распакуй и **скопируй с заменой**.

## Список изменений

### 🆕 services/session/ — полностью переписан

Старый сервис был сломан:
- gRPC сервер создавался пустым (хендлер не регистрировался)
- хардкод `localhost:50051/6379/4222` — не работает в docker-compose
- сигнатуры хендлеров не соответствовали .proto
- не входил в `go.work` и `docker-compose.yml`

Новая версия:
- gRPC методы: `StartSession`, `EndSession`, `GetSession`, `CalculatePrice`
- интеграция с `parking.proto` (через gRPC client) и `session.proto`
- Redis-кэш активных сессий
- NATS events: `parking.started`, `payment.completed` (типизированные)
- транзакционный rollback: если БД упала после assign — слот освобождается
- middleware: recovery + zap logging
- graceful shutdown
- 6 unit-тестов с моками

### ✏️ services/parking/ — обновлены импорты

`parking-proto` → `parking.proto`:
- `go.mod`
- `internal/app/app.go`
- `internal/handler/grpc.go`

### ✏️ services/user/ — обновлены импорты + удалён костыль

- `go.mod` → `user.proto v0.1.0`
- `internal/app/app.go` → `userv1.RegisterUserServiceServer` (из нормального
  сгенерированного кода)
- `internal/handler/grpc.go` → методы под новый `DeleteVehicleRequest`
- **УДАЛИТЬ** файл `internal/pb/user_ext.go` — он больше не нужен,
  Logout и DeleteVehicle теперь в .proto

### ✏️ docker-compose.yml — добавлены nats и session

- сервис `nats` (порты 4222, 8222)
- сервис `session` (порт 50053)
- `postgres` с healthcheck
- инициализационная миграция session
- env-переменные для всех сервисов

### ✏️ go.work — добавлен session

## Как применить

```bash
# 1) из корня проекта parking-system/, заменяем файлы патча:
unzip -o session-service-fix.zip

# 2) удаляем устаревший файл
rm services/user/internal/pb/user_ext.go
rmdir services/user/internal/pb 2>/dev/null || true

# 3) обновляем зависимости
cd services/parking && go mod tidy && cd ../..
cd services/user    && go mod tidy && cd ../..
cd services/session && go mod tidy && cd ../..

# 4) поднимаем
docker compose up --build
```

## Проверка работоспособности

```bash
# 1) создаём пользователя
grpcurl -plaintext -d '{"email":"test@test.com","password":"123456"}' \
  localhost:50052 user.v1.UserService/Register

# 2) добавляем машину
grpcurl -plaintext -d '{"user_id":1,"plate_number":"ABC-123"}' \
  localhost:50052 user.v1.UserService/AddVehicle

# 3) стартуем парковочную сессию
grpcurl -plaintext -d '{"user_id":1,"vehicle_id":1}' \
  localhost:50053 session.v1.SessionService/StartSession

# 4) заканчиваем (заплати!)
grpcurl -plaintext -d '{"session_id":1}' \
  localhost:50053 session.v1.SessionService/EndSession

# 5) сколько свободных мест осталось
grpcurl -plaintext -d '{}' \
  localhost:50051 parking.v1.ParkingService/GetAvailableSlots
```

## Запуск тестов

```bash
cd services/session && go test ./tests/...
cd services/user    && go test ./tests/...
```
