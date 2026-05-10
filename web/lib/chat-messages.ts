import { convertToModelMessages, type ModelMessage, type UIMessage } from "ai"

function flattenTextOnlyContent(message: ModelMessage): ModelMessage {
  if (message.role !== "user" && message.role !== "assistant") {
    return message
  }

  if (!("content" in message) || !Array.isArray(message.content)) {
    return message
  }

  if (!message.content.every((part) => part.type === "text")) {
    return message
  }

  return {
    ...message,
    content: message.content.map((part) => part.text).join(""),
  }
}

export async function normalizeChatMessages(
  messages: UIMessage[]
): Promise<ModelMessage[]> {
  const modelMessages = await convertToModelMessages(messages)

  return modelMessages.map(flattenTextOnlyContent)
}
