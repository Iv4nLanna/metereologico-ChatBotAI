import type { Metadata } from "next"
import "./globals.css"

export const metadata: Metadata = {
  title: "⛅ WeatherBot",
  description: "Consulte o clima de qualquer cidade com IA",
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="pt-BR">
      <body className="bg-white antialiased">{children}</body>
    </html>
  )
}
