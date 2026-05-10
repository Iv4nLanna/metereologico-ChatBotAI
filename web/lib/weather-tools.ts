import { tool } from "ai"
import { z } from "zod"

const API_URL = process.env.API_URL ?? "http://localhost:8080"

export const weatherTools = {
  getCurrentWeather: tool({
    description:
      "Get current weather conditions for a city. Use when the user asks about current weather, temperature, or conditions right now.",
    inputSchema: z.object({
      city: z.string().describe("City name, e.g. 'São Paulo', 'Curitiba', 'Lisbon'"),
    }),
    execute: async ({ city }) => {
      const res = await fetch(
        `${API_URL}/weather?city=${encodeURIComponent(city)}`
      )
      if (!res.ok) {
        throw new Error(`Weather API returned ${res.status} for city: ${city}`)
      }
      return res.json()
    },
  }),

  getForecast: tool({
    description:
      "Get weather forecast for a city for the next N days. Use when the user asks about future weather, tomorrow, or a specific number of days.",
    inputSchema: z.object({
      city: z.string().describe("City name"),
      days: z
        .number()
        .min(1)
        .max(7)
        .describe("Number of forecast days between 1 and 7"),
    }),
    execute: async ({ city, days }) => {
      const res = await fetch(
        `${API_URL}/forecast?city=${encodeURIComponent(city)}&days=${days}`
      )
      if (!res.ok) {
        throw new Error(`Forecast API returned ${res.status} for city: ${city}`)
      }
      return res.json()
    },
  }),
}
