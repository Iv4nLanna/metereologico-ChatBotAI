# Weather Platform ⛅

Plataforma de previsão do tempo com **API REST em Go** e **AI Agent em Next.js** que responde perguntas em linguagem natural sobre o clima de qualquer cidade do mundo.


---

## Acesso ao projeto (deploy em produção)

> Tudo já está deployado e funcionando. **Não é necessário rodar nada localmente para avaliar.**

### Interface principal

**[metereologico-chat-bot-ai.vercel.app](https://metereologico-chat-bot-ai.vercel.app)**

Acesse o link acima e faça perguntas como:
- *"Vai chover em São Paulo amanhã?"*
- *"Qual a temperatura em Lisboa agora?"*
- *"Previsão para Curitiba nos próximos 5 dias"*
- *"Como está o tempo em Tokyo hoje?"*

### API Go (backend)

| Endpoint | Link |
|---|---|
| Health check | [chatbot-production-a38f.up.railway.app/health](https://chatbot-production-a38f.up.railway.app/health) |
| Dashboard de métricas ao vivo | [chatbot-production-a38f.up.railway.app/debug](https://chatbot-production-a38f.up.railway.app/debug) |
| Snapshot JSON de métricas | [chatbot-production-a38f.up.railway.app/metrics](https://chatbot-production-a38f.up.railway.app/metrics) |

Exemplos diretos via curl:
```bash
curl "https://chatbot-production-a38f.up.railway.app/weather?city=São Paulo"
curl "https://chatbot-production-a38f.up.railway.app/forecast?city=Lisboa&days=5"
```

---

## Stack

| Camada | Tecnologia |
|---|---|
| Backend | Go + Gin |
| Frontend / AI Agent | Next.js 16 + TypeScript |
| Agent Framework | Vercel AI SDK |
| LLM | Groq — llama-3.3-70b-versatile |
| Dados meteorológicos | Open-Meteo API (gratuita, sem API key) |
| Deploy frontend | Vercel |
| Deploy backend | Railway |

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
┌─────────────────────────┐
│  Open-Meteo API         │  dados meteorológicos globais
└─────────────────────────┘
```

**Fluxo de uma pergunta:**

1. Usuário digita *"Vai chover em Lisboa amanhã?"* no chat
2. `useChat` envia histórico completo para `/api/chat` (Next.js server-side)
3. O agent normaliza as mensagens e chama o Groq LLM
4. O LLM decide chamar a tool `getForecast` com `city="Lisboa"` e `days=1`
5. O agent faz GET na API Go: `/forecast?city=Lisboa&days=1`
6. A API consulta o cache (TTL 5min); se miss, chama a Open-Meteo e armazena
7. O LLM recebe os dados meteorológicos e gera a resposta em streaming

---

## API Go — Referência

| Método | Rota | Descrição |
|---|---|---|
| GET | `/health` | Status e uptime do serviço |
| GET | `/weather?city={city}` | Condições atuais |
| GET | `/forecast?city={city}&days={n}` | Previsão de 1–7 dias |
| GET | `/metrics` | Snapshot JSON de métricas |
| GET | `/debug` | Dashboard visual de métricas (auto-refresh 5s) |

**Códigos de resposta:**

| Código | Situação |
|---|---|
| 200 | Sucesso |
| 400 | Parâmetro ausente ou inválido |
| 404 | Cidade não encontrada |
| 502 | Falha ao buscar dados na Open-Meteo |

---

## Observabilidade

Dashboard embutido sem dependências externas — sem Prometheus, sem Grafana.

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
| `p95_latency_ms` | Latência no percentil 95 |
| `cache_hits` / `cache_misses` | Eficiência do cache |
| `cache_hit_rate_pct` | Taxa de cache hits em % |

> **Nota:** `/metrics` e `/debug` estão intencionalmente públicos para facilitar esta avaliação. Em produção seriam protegidos por token Bearer via variável de ambiente.

---

## Decisões técnicas

### Vercel AI SDK (vs LangChain.js)
Integração nativa com Next.js, streaming sem boilerplate e tool calling com tipagem TypeScript completa. LangChain.js foi descartado por ser mais pesado e ter abstrações desnecessárias para um agent single-purpose.

### Groq + llama-3.3-70b-versatile
Velocidade de inferência ~10× maior que providers tradicionais e excelente desempenho em tool calling multilíngue. Tier gratuito generoso para desenvolvimento.

### Open-Meteo
Sem API key, gratuita, cobertura global e dados precisos. Elimina uma dependência de cadastro no setup.

### Cache in-memory com TTL de 5 minutos
Dados meteorológicos mudam lentamente. TTL de 5min evita chamadas repetidas sem servir dados desatualizados. In-memory é suficiente para single-instance — Redis seria o próximo passo para múltiplas instâncias.

### Gin (Go)
Framework HTTP mais adotado no ecossistema Go, middleware ecosystem maduro e performance adequada para a carga esperada.

### `log/slog` com JSON
Stdlib do Go 1.21+, zero dependências externas. Saída JSON ingestível por qualquer stack de observabilidade (Loki, Datadog, CloudWatch).

### Vercel + Railway
Frontend (Vercel) chama a API (Railway) server-side — `GROQ_API_KEY` e `API_URL` nunca chegam ao browser. CI/CD automático via GitHub em ambos.

---

## Setup local

**Pré-requisitos:** Docker Desktop + chave gratuita do [Groq](https://console.groq.com)

```bash
git clone https://github.com/Iv4nLanna/metereologico-ChatBotAI.git
cd metereologico-ChatBotAI
cp .env.example .env          # insira sua GROQ_API_KEY no .env
docker compose up --build
```

| Serviço | URL local |
|---|---|
| Chat | http://localhost:3000 |
| API Go | http://localhost:8080 |
| Dashboard de métricas | http://localhost:8080/debug |

---

## Testes

```bash
# Unitários (Go)
cd api && go test ./...

# Integração — requer a API rodando
cd api && go test -tags=integration ./...

# Frontend (TypeScript)
cd web && node --test lib/chat-messages.test.mts
```

---

## Considerações para produção

| Área | Recomendação |
|---|---|
| Cache | Substituir in-memory por **Redis** para suportar múltiplas instâncias |
| Autenticação | Token Bearer para `/metrics` e `/debug`; auth na `/api/chat` para controle de uso |
| Rate limiting | Middleware por IP na API Go para proteger a cota da Open-Meteo |
| Resiliência | Retry com exponential backoff nas chamadas à Open-Meteo |
| Observabilidade | Exportar métricas para Prometheus + Grafana; traces com OpenTelemetry |
| Secrets | AWS Secrets Manager ou HashiCorp Vault para `GROQ_API_KEY` |
