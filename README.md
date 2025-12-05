![CI/CD Pipeline](https://github.com/Atmosfr/go-hello-prod/actions/workflows/ci.yml/badge.svg)
![Docker Image Size](https://ghcr-badge.egpl.dev/atmosfr/go-hello-prod/size?tag=main)
![Go Version](https://img.shields.io/badge/Go-1.25-blue)

# Go Hello Prod — Production-Ready Service 2025

Clean, secure, modern Go service built from scratch to full deployment.

**27 MB image · No `:latest` · Non-root · Full CI/CD · Docker Compose**

### Features (Senior level 2025–2026)

- Go 1.25 + structured `slog` logging
- Middleware: request logging + panic recovery
- Graceful shutdown
- Configuration via `.env`
- Multi-stage Dockerfile → final image **~8 MB**
- `wolfi-base` (maximum security + Windows compatibility)
- One-command launch with `docker compose up`
- GitHub Actions: lint → test → build → push to GitHub Container Registry
- Zero `:latest` anywhere

### Quick start

```bash
git clone https://github.com/Atmosfr/go-hello-prod.git
cd go-hello-prod
docker compose up --build
# → http://localhost:8080/hello
```

### Endpoints

| Method | Path     | Description                                            |
| ------ | -------- | ------------------------------------------------------ |
| GET    | `/hello` | Returns hello message + timestamp                      |
| GET    | `/panic` | Triggers panic → tests recovery (service stays alive!) |

### Deploy

Image is automatically built and published on every push to `main`:

→ [GitHub Container Registry](https://github.com/Atmosfr/go-hello-prod/pkgs/container/go-hello-prod)
