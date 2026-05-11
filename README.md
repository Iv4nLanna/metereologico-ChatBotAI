# Weather Platform ⛅

Plataforma de previsão do tempo com **API Go** e **AI Agent em Next.js** que responde perguntas em linguagem natural sobre o clima.

---

## Demo ao vivo

| | URL |
|---|---|
| Chat | [metereologico-chat-bot-ai.vercel.app](https://metereologico-chat-bot-ai.vercel.app) |
| API | [chatbot-production-a38f.up.railway.app/health](https://chatbot-production-a38f.up.railway.app/health) |
| Dashboard de métricas | [chatbot-production-a38f.up.railway.app/debug](https://chatbot-production-a38f.up.railway.app/debug) |

---

## Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                        Usuário                               │
└─────────────────────┬───────────────────────────────────────┘
                      │ pergunta em linguagem natural
                      ▼
┌─────────────────────────────────────────────────────────────┐
│              Next.js Frontend  (Vercel)                      │
│                                                             │
│  ┌──────────────┐    ┌──────────────────────────────────┐   │
│  │  Chat UI     │    │  AI Agent  /api/chat             │   │
│  │  (useChat)   │───▶│  Vercel AI SDK + Groq LLM        │   │
│  └──────────────┘    │  llama-3.3-70b-versatile         │   │
│                      └──────────────┬───────────────────┘   │
└─────────────────────────────────────┼───────────────────────┘
                                      │ tool calls (server-side)
                                      ▼
┌─────────────────────────────────────────────────────────────┐
│              Go API  (Railway)                               │
│                                                             │
│  ┌──────────┐  ┌───────────┐  ┌────────────────────────┐   │
│  │ /weather │  │ /forecast │  │ /metrics  /debug        │   │
│  └────┬─────┘  └─────┬─────┘  └────────────────────────┘   │
│       │              │                                       │
│  ┌────▼──────────────▼────┐                                 │
│  │   Cache in-memory      │  TTL 5 min                      │
│  └────────────┬───────────┘                                 │
└───────────────┼─────────────────────────────────────────────┘
                │ HTTP
                ▼
┌─────────────────────────────┐
│  Open-Meteo API (gratuita)  │  dados meteorológicos globais
└─────────────────────────────┘
```

### Fluxo de uma pergunta

1. Usuário digita *"Vai chover em Lisboa amanhã?"* no chat
2. `useChat` envia a mensagem para `/api/chat` (Next.js)
3. O agent normaliza as mensagens (UIMessage → ModelMessage) e chama o Groq
4. O LLM decide chamar a tool `getForecast` com `city="Lisboa"` e `days=1`
5. O agent faz GET na API Go: `/forecast?city=Lisboa&days=1`
6. A API consulta o cache; se miss, chama a Open-Meteo e armazena
7. O LLM recebe os dados e gera a resposta em streaming para o usuário

---

## Endpoints da API

| Método | Rota | Descrição |
|---|---|---|
| GET | `/health` | Status do serviço com uptime |
| GET | `/weather?city={city}` | Condições atuais |
| GET | `/forecast?city={city}&days={n}` | Previsão de 1–7 dias |
| GET | `/metrics` | Snapshot JSON de métricas ao vivo |
| GET | `/debug` | Dashboard visual de métricas (auto-refresh 5s) |

```bash
curl "https://chatbot-production-a38f.up.railway.app/weather?city=São Paulo"
curl "https://chatbot-production-a38f.up.railway.app/forecast?city=Lisboa&days=5"
curl "https://chatbot-production-a38f.up.railway.app/metrics"
```

| Código | Situação |
|---|---|
| 200 | Sucesso |
| 400 | Parâmetro ausente ou inválido |
| 404 | Cidade não encontrada |
| 502 | Falha ao buscar dados na Open-Meteo |

---

## Observabilidade

Dashboard embutido na API — sem Prometheus, sem Grafana, sem configuração extra.

**Acesse ao vivo:** [chatbot-production-a38f.up.railway.app/debug](https://chatbot-production-a38f.up.railway.app/debug)

```
┌──────────────┬──────────────┬───────────────────────┐
│ Total Req.   │ Error Rate   │     Latency P95        │
│   1,247      │   0.8%  ●ok  │       142ms            │
├──────────────┴──────────────┴───────────────────────┤
│ Cache Hit Rate   ████████░░   78.6%                  │
│                  980 hits / 1247 total lookups        │
├─────────────────────────────────────────────────────┤
│ ● live — refreshing every 5s              [pause]    │
└─────────────────────────────────────────────────────┘
```

| Métrica | Descrição |
|---|---|
| `total_requests` | Total de requests desde o startup |
| `total_errors` | Requests com status 4xx ou 5xx |
| `error_rate_pct` | Taxa de erro em % |
| `p95_latency_ms` | Latência no percentil 95 (histograma: <10ms / <50ms / <100ms / <500ms / ≥500ms) |
| `cache_hits` | Lookups servidos do cache |
| `cache_misses` | Lookups que foram à Open-Meteo |
| `cache_hit_rate_pct` | Taxa de cache hits em % |

---

## Decisões técnicas

### Por que Vercel AI SDK?

Integração nativa com Next.js, streaming de respostas sem boilerplate, e suporte a tool calling com tipagem TypeScript. Alternativas consideradas: LangChain.js (mais pesado, overkill para um agent single-purpose) e chamadas diretas à API Groq (sem streaming nem tool loop gerenciado).

### Por que Groq + llama-3.3-70b?

Velocidade de inferência ~10× maior que providers tradicionais, tier gratuito generoso para desenvolvimento, e o llama-3.3-70b tem excelente desempenho em tool calling e respostas multilíngues.

### Por que Open-Meteo?

Zero configuração — não exige API key, é gratuita, tem cobertura global e dados precisos. Elimina uma dependência de cadastro no setup do projeto.

### Por que cache in-memory com TTL de 5 minutos?

Dados meteorológicos mudam lentamente. Um TTL de 5 minutos é suficiente para evitar chamadas repetidas à Open-Meteo sem servir dados desatualizados. Para deploy single-instance (caso de uso atual), in-memory é mais simples e performático que Redis.

### Por que Gin (Go)?

Framework HTTP mais adotado no ecossistema Go, middleware ecosystem maduro, performance adequada para uma API de previsão do tempo.

### Por que `log/slog` com JSON?

Stdlib do Go 1.21+, sem dependências externas. Saída em JSON é diretamente ingestível por qualquer stack de log aggregation (Loki, CloudWatch, Datadog).

### Observabilidade embutida

Em vez de adicionar Prometheus + Grafana ao `docker-compose.yml`, a API expõe `/metrics` (JSON) e `/debug` (dashboard HTML com `go:embed`). Zero dependências externas, zero configuração extra, demonstra o mesmo conceito de forma acessível.

> **Nota sobre segurança:** os endpoints `/metrics` e `/debug` estão intencionalmente públicos neste projeto para facilitar a avaliação — qualquer pessoa com a URL pode inspecionar as métricas ao vivo. Em um ambiente profissional, esses endpoints seriam protegidos por autenticação (ex: token Bearer via variável de ambiente), pois expõem informações operacionais internas como taxa de erros, latência e uptime do serviço.

### Por que Vercel + Railway?

O frontend (Vercel) chama a API (Railway) server-side — `GROQ_API_KEY` e `API_URL` nunca chegam ao browser. Vercel tem integração nativa com Next.js e CI/CD automático via GitHub. Railway faz deploy direto a partir do `Dockerfile` existente, injeta `PORT` automaticamente e expõe HTTPS sem configuração de servidor.

---

## Rodando localmente

**Pré-requisitos:** Docker Desktop + chave gratuita do [Groq](https://console.groq.com)

```bash
git clone https://github.com/Iv4nLanna/medereologico-ChatBotAI.git
cd medereologico-ChatBotAI
cp .env.example .env          # edite .env e insira sua GROQ_API_KEY
docker compose up --build
```

| Serviço | URL |
|---|---|
| Chat | http://localhost:3000 |
| API Go | http://localhost:8080 |
| Dashboard de métricas | http://localhost:8080/debug |

---

## Testes

```bash
# Unitários (Go)
cd api && go test ./...

# Integração (requer a API rodando)
cd api && go test -tags=integration ./...

# Frontend (TypeScript)
cd web && node --test lib/chat-messages.test.mts
```

---

## Considerações para evolução

| Área | Recomendação |
|---|---|
| Cache | Substituir in-memory por **Redis** para suportar múltiplas instâncias |
| Secrets | Usar secrets manager (AWS Secrets Manager, HashiCorp Vault) para `GROQ_API_KEY` |
| Rate limiting | Middleware de rate limit por IP na API Go |
| Resiliência | Retry com exponential backoff nas chamadas à Open-Meteo |
| Observabilidade | Exportar métricas para Prometheus + Grafana; traces com OpenTelemetry |
