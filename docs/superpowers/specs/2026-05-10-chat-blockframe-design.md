# Chat Block-Frame Design

**Data:** 2026-05-10
**Projeto:** metereologico-ChatBotAI — frontend Next.js (`web/`)
**Escopo:** Redesign visual do componente `Chat.tsx` para estilo neobrutalist block-frame

---

## Objetivo

Aplicar o estilo block-frame (neobrutalism) ao chat meteorológico. O chat deve ficar centralizado na tela com largura limitada, fundo cream com dot-grid, bordas grossas pretas com shadow offset, e destaque pink — fiel à paleta do template `beautiful-html-templates/block-frame`.

---

## Paleta de Cores

| Variável        | Valor     | Uso                                   |
|-----------------|-----------|---------------------------------------|
| `--bf-cream`    | `#FFDC8B` | Fundo da página                       |
| `--bf-pink`     | `#FE90E8` | Mensagens do usuário, botão enviar    |
| `--bf-blue`     | `#C0F7FE` | Reservado / futuro uso                |
| `--bf-green`    | `#99E885` | Reservado / futuro uso                |
| `--bf-yellow`   | `#F7CB46` | Label chip do header                  |
| `--bf-black`    | `#000000` | Bordas, sombras, texto                |
| `--bf-white`    | `#FFFFFF` | Fundo do card do chat, bubbles do bot |
| `--bf-offwhite` | `#FFFDF5` | Fundo do input                        |
| `--bf-border`   | `4px solid #000` | Borda padrão                   |
| `--bf-shadow`   | `8px 8px 0px #000` | Sombra offset padrão         |
| `--bf-shadow-sm`| `4px 4px 0px #000` | Sombra offset pequena        |

---

## Tipografia

- **Títulos / Labels:** `Space Grotesk`, weights 700–900, uppercase
- **Corpo / Input / Mensagens:** `Inter`, weights 400–700
- Importados via `next/font/google` em `web/app/layout.tsx`

---

## Layout da Página

- `body`: `min-h-screen`, fundo `--bf-cream`, dot-grid (`radial-gradient` 24px), `display: flex`, `align-items: center`, `justify-content: center`
- **Coluna do chat:** `max-w-[680px] w-full h-screen flex flex-col`, fundo branco, `border: var(--bf-border)`, `box-shadow: var(--bf-shadow)`

---

## Componentes

### Header

- Label chip: texto `ASSISTENTE CLIMÁTICO`, uppercase, font Space Grotesk, fundo `--bf-yellow`, `border: 2px solid #000`, padding `2px 8px`
- Título: `WEATHERBOT`, `font-weight: 900`, uppercase, Inter, ~28px
- Separador inferior: `border-bottom: var(--bf-border)`

### Área de Mensagens

- Fundo branco, `flex-1`, overflow-y scroll
- **Bubble — usuário:** fundo `--bf-pink`, `border: 3px solid #000`, `box-shadow: var(--bf-shadow-sm)`, sem border-radius, alinhamento direita
- **Bubble — bot:** fundo `--bf-white`, `border: 3px solid #000`, `box-shadow: var(--bf-shadow-sm)`, sem border-radius, alinhamento esquerda
- **Loading:** três blocos quadrados (`8px × 8px`) com animação de pulso substituindo os dots atuais
- **Suggestion chips:** `border: 2px solid #000`, `box-shadow: var(--bf-shadow-sm)`, sem border-radius; hover: `translate(-2px, -2px)` + shadow maior

### Footer / Input

- Separador superior: `border-top: var(--bf-border)`
- **Input:** `border: var(--bf-border)`, sem border-radius, fundo `--bf-offwhite`, font Inter
- **Botão ENVIAR:** texto uppercase bold, fundo `--bf-pink`, `border: 3px solid #000`, `box-shadow: var(--bf-shadow-sm)`, sem border-radius
  - Hover: `transform: translate(-2px, -2px)`, shadow `6px 6px 0px #000`
  - Active: `transform: translate(2px, 2px)`, sem shadow
  - Disabled: opacidade reduzida, cursor not-allowed

---

## Abordagem de Implementação

**Abordagem A** — CSS puro + Tailwind no componente existente:

1. Adicionar `Space Grotesk` + `Inter` em `web/app/layout.tsx` via `next/font/google`
2. Adicionar variáveis CSS (`--bf-*`) em `web/app/globals.css`
3. Adicionar dot-grid ao `body` em `globals.css`
4. Reescrever classes Tailwind em `web/components/Chat.tsx` aplicando o novo visual

**O que NÃO muda:**
- Hook `useChat` e toda a lógica de estado
- Auto-scroll com `useRef`
- Atualização dinâmica do `document.title`
- Suggestion chips (apenas restyling)
- Endpoints da API (`/api/chat`)

---

## Arquivos Afetados

| Arquivo                          | Mudança                              |
|----------------------------------|--------------------------------------|
| `web/app/layout.tsx`             | Adicionar fontes Google              |
| `web/app/globals.css`            | Variáveis CSS + dot-grid no body     |
| `web/components/Chat.tsx`        | Reescrever classes de estilo         |
