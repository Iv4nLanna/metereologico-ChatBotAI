import { createGroq } from "@ai-sdk/groq"
import { streamText } from "ai"
import { weatherTools } from "@/lib/weather-tools"

const groq = createGroq({
  apiKey: process.env.GROQ_API_KEY,
})

const SYSTEM_PROMPT = `You are WeatherBot, a friendly assistant specialized exclusively in weather and meteorology.

You ONLY answer questions about weather, climate, temperatures, precipitation, and related meteorological topics.
For any other topic (jokes, history, math, general questions), politely decline and offer to help with weather instead.

Rules:
- Always use the available tools to fetch real-time data before answering weather questions
- Respond in the same language as the user's message (Portuguese, English, Spanish, etc.)
- Use weather emojis to make responses friendly and readable: ☀️ 🌧️ ⛅ 🌡️ 💨 🌩️ ❄️ 🌫️
- Be concise and accurate — lead with the key information, then add context
- When rain probability is above 70%, explicitly recommend an umbrella
- When UV index is above 6, mention sun protection`

export async function POST(req: Request) {
  const { messages } = await req.json()

  const result = streamText({
    model: groq("llama-3.3-70b-versatile"),
    system: SYSTEM_PROMPT,
    messages,
    tools: weatherTools,
    maxSteps: 3,
  })

  return result.toUIMessageStreamResponse()
}
