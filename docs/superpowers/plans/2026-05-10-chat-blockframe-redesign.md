# Chat Block-Frame Redesign Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Aplicar estilo neobrutalist block-frame ao chat meteorológico — fundo cream com dot-grid, coluna full-height centralizada com borda preta e shadow offset, tipografia bold Space Grotesk + Inter, bubbles com borda/shadow, botão pink com hover lift.

**Architecture:** Redesign puramente visual — hooks, estado e chamadas de API permanecem intactos. As mudanças se concentram em layout.tsx (fontes), globals.css (variáveis CSS + utilitários block-frame) e nos 4 componentes de UI (Chat, Message, SuggestionChips, WeatherBadge).

**Tech Stack:** Next.js 16, TypeScript, Tailwind CSS v4, next/font/google (Space Grotesk + Inter)

---

## Mapeamento de Arquivos

| Arquivo | Ação | Responsabilidade |
|---|---|---|
| `web/app/layout.tsx` | Modificar | Carregar Space Grotesk + Inter via next/font/google |
| `web/app/globals.css` | Modificar | Variáveis CSS block-frame, dot-grid no body, `.bf-btn` e `.bf-chip` |
| `web/components/Chat.tsx` | Modificar | Wrapper da coluna, header, loading, footer/form |
| `web/components/Message.tsx` | Modificar | Bubble pink (usuário) e branco (bot) com borda/shadow |
| `web/components/SuggestionChips.tsx` | Modificar | Chips quadrados com shadow + welcome state |
| `web/components/WeatherBadge.tsx` | Modificar | Badge de cidade com estilo block-frame |

---

### Task 1: Adicionar fontes Google em layout.tsx

**Files:**
- Modify: `web/app/layout.tsx`

- [ ] **Step 1: Substituir o conteúdo completo de `web/app/layout.tsx`**

```tsx
import type { Metadata } from "next"
import { Space_Grotesk, Inter } from "next/font/google"
import "./globals.css"

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-space-grotesk",
  weight: ["400", "500", "600", "700"],
})

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
  weight: ["400", "500", "600", "700", "800", "900"],
})

export const metadata: Metadata = {
  title: "⛅ WeatherBot",
  description: "Consulte o clima de qualquer cidade com IA",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="pt-BR">
      <body className={`${spaceGrotesk.variable} ${inter.variable}`}>
        {children}
      </body>
    </html>
  )
}
```

- [ ] **Step 2: Verificar tipos**

```bash
cd /home/ivan/metereologico-ChatBotAI/web && npx tsc --noEmit
```

Esperado: sem erros de tipo.

- [ ] **Step 3: Commit**

```bash
cd /home/ivan/metereologico-ChatBotAI && git add web/app/layout.tsx && git commit -m "feat: load Space Grotesk and Inter fonts via next/font/google"
```

---

### Task 2: Configurar globals.css com tema block-frame

**Files:**
- Modify: `web/app/globals.css`

- [ ] **Step 1: Substituir o conteúdo completo de `web/app/globals.css`**

```css
@import "tailwindcss";

:root {
  --bf-cream: #FFDC8B;
  --bf-pink: #FE90E8;
  --bf-yellow: #F7CB46;
  --bf-black: #000000;
  --bf-white: #FFFFFF;
  --bf-offwhite: #FFFDF5;
  --bf-shadow-sm: 4px 4px 0px #000000;
  --bf-shadow: 8px 8px 0px #000000;
}

@theme inline {
  --font-sans: var(--font-inter);
  --font-display: var(--font-space-grotesk);
}

body {
  background-color: var(--bf-cream);
  background-image: radial-gradient(circle, #000 1.2px, transparent 1.2px);
  background-size: 24px 24px;
  font-family: var(--font-inter), Inter, sans-serif;
  color: var(--bf-black);
}

/* Neobrutalist button — offset shadow + hover lift, skip when disabled */
.bf-btn {
  box-shadow: var(--bf-shadow-sm);
  transition: transform 0.1s ease, box-shadow 0.1s ease;
}
.bf-btn:hover:not(:disabled) {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px #000;
}
.bf-btn:active:not(:disabled) {
  transform: translate(2px, 2px);
  box-shadow: none;
}

/* Neobrutalist chip — same hover lift */
.bf-chip {
  box-shadow: var(--bf-shadow-sm);
  transition: transform 0.1s ease, box-shadow 0.1s ease;
}
.bf-chip:hover {
  transform: translate(-2px, -2px);
  box-shadow: 6px 6px 0px #000;
}
.bf-chip:active {
  transform: translate(2px, 2px);
  box-shadow: none;
}
```

- [ ] **Step 2: Commit**

```bash
cd /home/ivan/metereologico-ChatBotAI && git add web/app/globals.css && git commit -m "feat: add block-frame CSS variables, dot-grid body and bf utility classes"
```

---

### Task 3: Restyle Chat.tsx

**Files:**
- Modify: `web/components/Chat.tsx`

- [ ] **Step 1: Substituir o bloco `return (...)` de Chat.tsx**

Manter todo o código de hooks e handlers acima do `return`. Substituir apenas o JSX retornado:

```tsx
  return (
    <div
      className="flex flex-col h-screen w-full max-w-[680px] mx-auto bg-white"
      style={{ borderLeft: "4px solid #000", borderRight: "4px solid #000", boxShadow: "var(--bf-shadow)" }}
    >
      <header
        className="flex-none p-4 bg-white"
        style={{ borderBottom: "4px solid #000" }}
      >
        <div className="flex flex-col items-center gap-1">
          <span
            className="text-xs font-bold uppercase tracking-widest px-2 py-0.5"
            style={{
              fontFamily: "var(--font-space-grotesk)",
              background: "var(--bf-yellow)",
              border: "2px solid #000",
            }}
          >
            Assistente Climático
          </span>
          <span
            className="text-2xl font-black uppercase tracking-tight"
            style={{ fontFamily: "var(--font-space-grotesk)" }}
          >
            WeatherBot
          </span>
        </div>
      </header>

      <main className="flex-1 overflow-y-auto p-4 space-y-4 bg-white">
        {messages.length === 0 && (
          <SuggestionChips onSelect={handleSuggestion} />
        )}
        {messages
          .filter((m) => m.role !== "system")
          .map((m) => (
            <Message key={m.id} message={m} />
          ))}
        {isLoading && (
          <div className="flex justify-start">
            <div
              className="px-4 py-3 bg-white"
              style={{ border: "2px solid #000", boxShadow: "var(--bf-shadow-sm)" }}
            >
              <div className="flex gap-1.5 items-center h-4">
                <span className="w-2 h-2 bg-black animate-bounce [animation-delay:-0.3s]" />
                <span className="w-2 h-2 bg-black animate-bounce [animation-delay:-0.15s]" />
                <span className="w-2 h-2 bg-black animate-bounce" />
              </div>
            </div>
          </div>
        )}
        <div ref={bottomRef} />
      </main>

      <footer
        className="flex-none p-4 bg-white"
        style={{ borderTop: "4px solid #000" }}
      >
        <form onSubmit={handleSubmit} className="flex gap-2">
          <input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Ex: Vai chover em São Paulo amanhã?"
            className="flex-1 px-4 py-2 text-sm outline-none"
            style={{
              background: "var(--bf-offwhite)",
              border: "4px solid #000",
              fontFamily: "var(--font-inter)",
            }}
            disabled={isLoading}
          />
          <button
            type="submit"
            disabled={isLoading || !input.trim()}
            className="bf-btn px-5 py-2 text-sm font-black uppercase tracking-wide disabled:opacity-50 disabled:cursor-not-allowed"
            style={{
              background: "var(--bf-pink)",
              border: "3px solid #000",
              fontFamily: "var(--font-space-grotesk)",
            }}
          >
            Enviar
          </button>
        </form>
      </footer>
    </div>
  )
```

- [ ] **Step 2: Verificar tipos**

```bash
cd /home/ivan/metereologico-ChatBotAI/web && npx tsc --noEmit
```

Esperado: sem erros de tipo.

- [ ] **Step 3: Commit**

```bash
cd /home/ivan/metereologico-ChatBotAI && git add web/components/Chat.tsx && git commit -m "feat: apply block-frame style to Chat wrapper, header, loading indicator and footer"
```

---

### Task 4: Restyle Message.tsx

**Files:**
- Modify: `web/components/Message.tsx`

- [ ] **Step 1: Substituir o bloco `return (...)` de Message.tsx**

Manter todo o código acima do `return` intacto (`extractCity`, `textContent`, `city`, etc.). Substituir apenas o JSX:

```tsx
  return (
    <div className={`flex ${isUser ? "justify-end" : "justify-start"} mb-3`}>
      <div
        className="max-w-xs lg:max-w-md px-4 py-2"
        style={{
          background: isUser ? "var(--bf-pink)" : "var(--bf-white)",
          border: "3px solid #000",
          boxShadow: "var(--bf-shadow-sm)",
          fontFamily: "var(--font-inter)",
        }}
      >
        <p className="text-sm whitespace-pre-wrap">{textContent}</p>
        {city && <WeatherBadge city={city} />}
      </div>
    </div>
  )
```

- [ ] **Step 2: Verificar tipos**

```bash
cd /home/ivan/metereologico-ChatBotAI/web && npx tsc --noEmit
```

Esperado: sem erros de tipo.

- [ ] **Step 3: Commit**

```bash
cd /home/ivan/metereologico-ChatBotAI && git add web/components/Message.tsx && git commit -m "feat: apply block-frame style to message bubbles (pink user, white bot)"
```

---

### Task 5: Restyle SuggestionChips.tsx e WeatherBadge.tsx

**Files:**
- Modify: `web/components/SuggestionChips.tsx`
- Modify: `web/components/WeatherBadge.tsx`

- [ ] **Step 1: Substituir o conteúdo completo de `web/components/SuggestionChips.tsx`**

```tsx
const SUGGESTIONS = [
  "☀️ Vai chover em São Paulo hoje?",
  "🌡️ Qual a temperatura em Lisboa amanhã?",
  "💨 Previsão de Curitiba para esta semana?",
]

interface Props {
  onSelect: (text: string) => void
}

export default function SuggestionChips({ onSelect }: Props) {
  return (
    <div className="flex flex-col items-center justify-center py-16 gap-6">
      <div className="text-center">
        <div className="text-5xl mb-3">⛅</div>
        <h2
          className="text-2xl font-black uppercase tracking-tight"
          style={{ fontFamily: "var(--font-space-grotesk)" }}
        >
          WeatherBot
        </h2>
        <p className="text-sm mt-2" style={{ fontFamily: "var(--font-inter)" }}>
          Pergunte sobre o clima de qualquer cidade do mundo
        </p>
      </div>
      <div className="flex flex-col gap-3 w-full max-w-sm">
        {SUGGESTIONS.map((s) => (
          <button
            key={s}
            onClick={() => onSelect(s)}
            className="bf-chip px-4 py-3 text-sm text-left font-medium bg-white"
            style={{ border: "3px solid #000", fontFamily: "var(--font-inter)" }}
          >
            {s}
          </button>
        ))}
      </div>
    </div>
  )
}
```

- [ ] **Step 2: Substituir o conteúdo completo de `web/components/WeatherBadge.tsx`**

```tsx
interface Props {
  city: string
}

export default function WeatherBadge({ city }: Props) {
  return (
    <span
      className="inline-flex items-center gap-1 text-xs font-bold uppercase mt-2 px-1.5 py-0.5"
      style={{
        border: "1.5px solid #000",
        background: "var(--bf-yellow)",
        fontFamily: "var(--font-space-grotesk)",
      }}
    >
      📍 {city}
    </span>
  )
}
```

- [ ] **Step 3: Verificar tipos**

```bash
cd /home/ivan/metereologico-ChatBotAI/web && npx tsc --noEmit
```

Esperado: sem erros de tipo.

- [ ] **Step 4: Commit**

```bash
cd /home/ivan/metereologico-ChatBotAI && git add web/components/SuggestionChips.tsx web/components/WeatherBadge.tsx && git commit -m "feat: apply block-frame style to suggestion chips and weather badge"
```

---

### Task 6: Verificação visual

**Files:** nenhum arquivo alterado

- [ ] **Step 1: Instalar dependências e iniciar dev server**

```bash
cd /home/ivan/metereologico-ChatBotAI/web && npm install && npm run dev
```

Esperado: `▲ Next.js 16.x.x` rodando em `http://localhost:3000`

- [ ] **Step 2: Abrir `http://localhost:3000` e verificar cada item**

- [ ] Fundo cream (#FFDC8B) visível nas laterais da coluna com dot-grid
- [ ] Coluna do chat centralizada, max 680px, fundo branco
- [ ] Shadow offset preto visível à direita/baixo da coluna
- [ ] Header: label chip amarelo "ASSISTENTE CLIMÁTICO" + título "WEATHERBOT" em bold uppercase
- [ ] Suggestion chips: quadrados, borda preta 3px, shadow offset + hover com lift
- [ ] Enviar uma mensagem de teste (ex: "Clima em São Paulo?")
- [ ] Bubble do usuário: fundo pink (#FE90E8), borda preta, shadow offset, sem border-radius
- [ ] Bubble do bot: fundo branco, borda preta, shadow offset, sem border-radius
- [ ] WeatherBadge aparece dentro do bubble do bot com fundo amarelo e borda preta
- [ ] Loading (3 quadrados pretos pulsando) aparece durante a resposta
- [ ] Input: borda preta 4px, fundo offwhite (#FFFDF5), sem border-radius
- [ ] Botão ENVIAR: pink, borda preta, uppercase — hover lift + shadow maior, disabled com opacidade 50%

- [ ] **Step 3: Encerrar dev server**

```
Ctrl+C
```
