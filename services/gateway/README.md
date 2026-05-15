# API Gateway Demo

Base URL: `http://localhost:8080`

## 1) Register a new user

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@test.com","password":"secret123"}'
```

Response:

```json
{"user_id":1,"token":"eyJhbGciOi..."}
```

Save the token:

```bash
TOKEN="eyJhbGciOi..."
```

## 2) Login (returns a fresh token)

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"demo@test.com","password":"secret123"}'
```

## 3) Add a vehicle (auth required)

```bash
curl -X POST http://localhost:8080/vehicles \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"plate_number":"ABC-123"}'
```

## 4) List vehicles

```bash
curl http://localhost:8080/vehicles \
  -H "Authorization: Bearer $TOKEN"
```

## 5) Available parking slots

```bash
curl http://localhost:8080/slots/available \
  -H "Authorization: Bearer $TOKEN"
```

## 6) Start a parking session

```bash
curl -X POST http://localhost:8080/sessions/start \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"vehicle_id":1}'
```

Response:

```json
{"session_id":1,"slot_id":1,"vehicle_id":1,"start_time_unix":1731600000}
```

## 7) Get current price for the active session

```bash
curl http://localhost:8080/sessions/1/price \
  -H "Authorization: Bearer $TOKEN"
```

## 8) End the session (charges payment)

```bash
curl -X POST http://localhost:8080/sessions/1/end \
  -H "Authorization: Bearer $TOKEN"
```

## 9) Logout

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer $TOKEN"
```
