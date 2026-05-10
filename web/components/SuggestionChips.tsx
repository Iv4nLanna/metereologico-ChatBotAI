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
        <h2 className="text-xl font-semibold text-gray-800">WeatherBot</h2>
        <p className="text-gray-500 text-sm mt-1">
          Pergunte sobre o clima de qualquer cidade do mundo
        </p>
      </div>
      <div className="flex flex-col gap-2 w-full max-w-sm">
        {SUGGESTIONS.map((s) => (
          <button
            key={s}
            onClick={() => onSelect(s)}
            className="border border-gray-200 rounded-xl px-4 py-3 text-sm text-left hover:bg-gray-50 hover:border-gray-300 transition-all"
          >
            {s}
          </button>
        ))}
      </div>
    </div>
  )
}
