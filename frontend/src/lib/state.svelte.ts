import { api, type Currency, type Point } from './api'
import { evalExpr } from './format'

function load<T>(key: string, fallback: T): T {
  try {
    const v = localStorage.getItem(key)
    return v === null ? fallback : (JSON.parse(v) as T)
  } catch {
    return fallback
  }
}

export const app = $state({
  from: load('fx-from', 'USD'),
  to: load('fx-to', 'CNY'),
  amount: load('fx-amount', '1000'),
  range: load('fx-range', 90),
})

export const market = $state({
  currencies: [] as Currency[],
  usd: {} as Record<string, number>,
  date: '',
  ok: false,
})

// Fetched history, keyed. pair: "FROM|TO|days"; board: FROM -> code -> points.
export const hist = $state({
  pair: {} as Record<string, Point[]>,
  board: {} as Record<string, Record<string, Point[]>>,
})

const pairInflight = new Set<string>()
export function ensurePair(from: string, to: string, days: number): void {
  if (from === to || !from || !to) return
  const key = `${from}|${to}|${days}`
  if (hist.pair[key] || pairInflight.has(key)) return
  pairInflight.add(key)
  api
    .history(from, to, days)
    .then((res) => (hist.pair[key] = res.points))
    .catch(() => {})
    .finally(() => pairInflight.delete(key))
}

const boardInflight = new Set<string>()
export function ensureBoard(base: string): void {
  if (!base || boardInflight.has(base)) return
  boardInflight.add(base)
  api
    .board(base, 30)
    .then((res) => (hist.board[base] = res.series))
    .catch(() => {})
    .finally(() => boardInflight.delete(base))
}

// Pair rate: how many `to` one `from` buys, derived from USD-based rates.
export function crossRate(usd: Record<string, number>, from: string, to: string): number {
  const f = usd[from]
  const t = usd[to]
  return f && t ? t / f : 0
}

/** Amount currently in the converter, respecting arithmetic expressions. */
export function amountValue(): number | null {
  const v = evalExpr(app.amount)
  return v !== null && v >= 0 ? v : null
}

export function saveApp(): void {
  localStorage.setItem('fx-from', JSON.stringify(app.from))
  localStorage.setItem('fx-to', JSON.stringify(app.to))
  localStorage.setItem('fx-amount', JSON.stringify(app.amount))
  localStorage.setItem('fx-range', JSON.stringify(app.range))
}

let bootstrapped = false
export async function bootstrap(): Promise<void> {
  if (bootstrapped) return
  bootstrapped = true
  const pull = async () => {
    try {
      const res = await api.latest('USD')
      market.usd = res.rates
      market.date = res.date
      market.ok = true
    } catch {
      market.ok = market.ok // keep stale data on refresh failure
    }
  }
  api
    .currencies()
    .then((cs) => {
      market.currencies = cs
      // Reset stale picks from an older session (e.g. a delisted currency).
      const codes = new Set(cs.map((c) => c.code))
      if (!codes.has(app.from)) app.from = 'USD'
      if (!codes.has(app.to)) app.to = 'CNY'
    })
    .catch(() => {})
  await pull()
  setInterval(() => {
    if (document.visibilityState === 'visible') pull()
  }, 5 * 60 * 1000)
}
