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
    <div className="flex flex-col h-screen max-w-2xl mx-auto">
      <header className="flex-none p-4 border-b bg-white">
        <div className="text-center">
          <span className="font-semibold text-gray-800">⛅ WeatherBot</span>
        </div>
      </header>

      <main className="flex-1 overflow-y-auto p-4 space-y-4">
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
            <div className="bg-gray-100 rounded-2xl rounded-bl-none px-4 py-2">
              <div className="flex gap-1 items-center h-5">
                <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce [animation-delay:-0.3s]" />
                <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce [animation-delay:-0.15s]" />
                <span className="w-2 h-2 bg-gray-400 rounded-full animate-bounce" />
              </div>
            </div>
          </div>
        )}
        <div ref={bottomRef} />
      </main>

      <footer className="flex-none p-4 border-t bg-white">
        <form onSubmit={handleSubmit} className="flex gap-2">
          <input
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder="Ex: Vai chover em São Paulo amanhã?"
            className="flex-1 border border-gray-200 rounded-full px-4 py-2 text-sm outline-none focus:ring-2 focus:ring-blue-400 focus:border-transparent"
            disabled={isLoading}
          />
          <button
            type="submit"
            disabled={isLoading || !input.trim()}
            className="bg-blue-500 hover:bg-blue-600 disabled:opacity-50 text-white rounded-full px-5 py-2 text-sm font-medium transition-colors"
          >
            Enviar
          </button>
        </form>
      </footer>
    </div>
  )
}
