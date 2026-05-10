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
        className={`max-w-xs lg:max-w-md px-4 py-2 rounded-2xl ${
          isUser
            ? "bg-blue-500 text-white rounded-br-sm"
            : "bg-gray-100 text-gray-800 rounded-bl-sm"
        }`}
      >
        <p className="text-sm whitespace-pre-wrap">{textContent}</p>
        {city && <WeatherBadge city={city} />}
      </div>
    </div>
  )
}
