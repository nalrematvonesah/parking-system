# notification-service

NATS consumer that turns domain events into emails. Owner: shared.

## Subscribed subjects

| Subject | Producer | Action |
|---|---|---|
| `user.registered` | user-service | welcome email |
| `parking.started` | session-service | parking-started notification |
| `payment.completed` | session-service | payment receipt |

## SMTP modes

If `SMTP_USER`, `SMTP_PASS` and `SMTP_FROM` are all set, the service sends
real emails via the configured SMTP server.

If any of them is empty, the service runs in **dev mode** and just logs each
email to stdout. This is the default in `docker-compose.yml` so the demo
works without real credentials.

Real SMTP example (Gmail App Password):

```env
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=you@gmail.com
SMTP_PASS=xxxxxxxxxxxxxxxx
SMTP_FROM=you@gmail.com
```
