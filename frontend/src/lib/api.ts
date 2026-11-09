export interface Currency {
  code: string
  name: string
  zh: string
  symbol: string
  decimals: number
  popular: number
}

export interface Point {
  d: string
  r: number
}

export interface BoardResp {
  base: string
  days: number
  series: Record<string, Point[]>
}

async function get<T>(url: string): Promise<T> {
  const res = await fetch(url)
  if (!res.ok) throw new Error(`${url}: ${res.status}`)
  return res.json()
}

export const api = {
  currencies: () => get<Currency[]>('/api/currencies'),
  latest: (base = 'USD') =>
    get<{ date: string; base: string; rates: Record<string, number> }>(
      `/api/latest?base=${base}`,
    ),
  history: (from: string, to: string, days: number) =>
    get<{ from: string; to: string; days: number; points: Point[] }>(
      `/api/history?from=${from}&to=${to}&days=${days}`,
    ),
  board: (base: string, days = 30) =>
    get<BoardResp>(`/api/board?base=${base}&days=${days}`),
}
