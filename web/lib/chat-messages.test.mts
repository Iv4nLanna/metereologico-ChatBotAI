import assert from "node:assert/strict"
import test from "node:test"
import type { UIMessage } from "ai"
import { normalizeChatMessages } from "./chat-messages.ts"

test("converts text-only UI messages into plain string model messages for streamText", async () => {
  const messages = [
    {
      id: "1",
      role: "user",
      parts: [{ type: "text", text: "Vai chover em Lisboa hoje?" }],
    },
  ] satisfies UIMessage[]

  const result = await normalizeChatMessages(messages)

  assert.deepEqual(result, [
    {
      role: "user",
      content: "Vai chover em Lisboa hoje?",
    },
  ])
})
