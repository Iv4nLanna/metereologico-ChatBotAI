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
