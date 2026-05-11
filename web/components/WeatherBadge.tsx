interface Props {
  city: string
}

export default function WeatherBadge({ city }: Props) {
  return (
    <span
      className="inline-flex items-center gap-1 text-xs font-bold uppercase mt-2 px-1.5 py-0.5"
      style={{
        border: "1.5px solid #000",
        background: "var(--bf-yellow)",
        fontFamily: "var(--font-space-grotesk)",
      }}
    >
      📍 {city}
    </span>
  )
}
