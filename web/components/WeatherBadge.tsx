interface Props {
  city: string
}

export default function WeatherBadge({ city }: Props) {
  return (
    <span className="inline-flex items-center gap-1 text-xs text-gray-400 mt-1">
      📍 {city}
    </span>
  )
}
