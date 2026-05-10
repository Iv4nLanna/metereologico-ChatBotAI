# Weather Platform ⛅

Plataforma de previsão do tempo com **API Go** + **AI Agent Next.js** que responde perguntas em linguagem natural.

## Requisitos

- Docker e Docker Compose
- Chave de API do [Groq](https://console.groq.com) (gratuita)

## Quickstart

```bash
# 1. Clone e configure
git clone <repo-url>
cd weather-platform
cp .env.example .env
# Edite .env e adicione: GROQ_API_KEY=sua_chave_aqui

# 2. Suba tudo
docker compose up --build

# 3. Acesse
# Frontend: http://localhost:3000
# API:      http://localhost:8080
```

## Endpoints da API

| Método | Rota | Exemplo |
|--------|------|---------|
| GET | `/health` | `curl localhost:8080/health` |
| GET | `/weather?city={city}` | `curl "localhost:8080/weather?city=Curitiba"` |
| GET | `/forecast?city={city}&days={n}` | `curl "localhost:8080/forecast?city=SP&days=3"` |

## Testes

```bash
# Testes unitários (Go)
cd api && go test ./...

# Testes de integração (requer servidor rodando)
cd api && go test -tags=integration ./...
```

## Decisões Técnicas

| Decisão | Escolha | Motivo |
|---------|---------|--------|
| Agent Framework | Vercel AI SDK | Integração nativa com Next.js, streaming built-in, tool calling simples |
| LLM | Groq (llama-3.3-70b) | Alta velocidade, tier gratuito generoso, ótimo tool calling |
| Weather API | Open-Meteo | Gratuita, sem API key, dados precisos e globais |
| HTTP Framework (Go) | Gin | Amplamente adotado, bom middleware ecosystem |
| Cache | In-memory (TTL 5min) | Zero dependência, adequado para escala single-instance |
| Logs | log/slog (JSON) | Stdlib Go 1.21+, estruturado, pronto para ingestão em log aggregators |

## Considerações para Produção

- **Cache distribuído**: substituir map in-memory por Redis para múltiplas instâncias
- **Secrets**: usar um secrets manager (AWS Secrets Manager, HashiCorp Vault) para `GROQ_API_KEY`
- **TLS**: terminar HTTPS no load balancer (AWS ALB, Cloudflare)
- **Rate limiting**: adicionar middleware de rate limit por IP na API Go
- **Retry com backoff**: nas chamadas à Open-Meteo para resiliência a falhas transitórias
- **Observabilidade**: exportar métricas para Prometheus + Grafana; traces com OpenTelemetry
- **Deploy sugerido**: API Go → Railway ou Fly.io | Next.js → Vercel
