import { UIMessage } from "ai"
import WeatherBadge from "./WeatherBadge"

interface Props {
  message: UIMessage
}

function extractCity(message: UIMessage): string | null {
  for (const part of message.parts) {
    // Handle both static tool parts (type: "tool-<name>") and dynamic tool parts
    if (
      (part.type.startsWith("tool-") || part.type === "dynamic-tool") &&
      "state" in part &&
      part.state === "output-available" &&
      "output" in part &&
      part.output !== null &&
      typeof part.output === "object"
    ) {
      const output = part.output as Record<string, unknown>
      if (typeof output.city === "string") {
        return output.city
      }
    }
  }
  return null
}

export default function Message({ message }: Props) {
  const isUser = message.role === "user"

  const textContent = message.parts
    .filter((p) => p.type === "text")
    .map((p) => {
      if (p.type === "text") return p.text
      return ""
    })
    .join("")

  const city = isUser ? null : extractCity(message)

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
}
