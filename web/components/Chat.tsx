"use client"

import { useChat } from "@ai-sdk/react"
import { DefaultChatTransport } from "ai"
import { useEffect, useRef, useState } from "react"
import Message from "./Message"
import SuggestionChips from "./SuggestionChips"

export default function Chat() {
  const { messages, sendMessage, status } = useChat({
    transport: new DefaultChatTransport({ api: "/api/chat" }),
  })
  const [input, setInput] = useState("")
  const isLoading = status === "submitted" || status === "streaming"
  const bottomRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" })
  }, [messages])

  // Update tab title based on last bot message
  useEffect(() => {
    const lastBot = [...messages].reverse().find((m) => m.role === "assistant")
    if (!lastBot) return
    // Extract text from parts (ai@6 UIMessage uses parts array)
    const content = lastBot.parts
      .filter((p) => p.type === "text")
      .map((p) => (p.type === "text" ? p.text : ""))
      .join("")
    const emoji = content.includes("☀️")
      ? "☀️"
      : content.includes("🌧️")
        ? "🌧️"
        : content.includes("❄️")
          ? "❄️"
          : "⛅"
    document.title = `${emoji} WeatherBot`
  }, [messages])

  const handleSuggestion = (text: string) => {
    setInput(text)
  }

  const handleSubmit = (e?: { preventDefault?: () => void }) => {
    e?.preventDefault?.()
    if (!input.trim() || isLoading) return
    sendMessage({ text: input.trim() })
    setInput("")
  }

  return (
    <div
      className="flex flex-col w-full max-w-[680px] mx-auto my-8 bg-white"
      style={{ border: "4px solid #000", boxShadow: "var(--bf-shadow)", height: "calc(100vh - 4rem)" }}
    >
      <header
        className="flex-none p-6 bg-white"
        style={{ borderBottom: "4px solid #000" }}
      >
        <div className="flex flex-col items-center gap-2">
          <span
            className="text-5xl font-black uppercase tracking-tight"
            style={{ fontFamily: "var(--font-space-grotesk)" }}
          >
            WeatherBot
          </span>
          <span
            className="text-xs font-bold uppercase tracking-widest px-2 py-0.5"
            style={{
              fontFamily: "var(--font-space-grotesk)",
              background: "var(--bf-pink)",
              border: "2px solid #000",
            }}
          >
            Assistente Climático
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
}
