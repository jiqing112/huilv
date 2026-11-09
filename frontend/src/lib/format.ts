export function fmtAmount(n: number, decimals: number): string {
  if (!Number.isFinite(n)) return ''
  return new Intl.NumberFormat('en-US', {
    maximumFractionDigits: decimals,
    minimumFractionDigits: Math.min(decimals, n % 1 === 0 ? 0 : decimals),
  }).format(n)
}

// Smart decimals: fewer for big rates, more precision for tiny ones.
export function fmtRate(r: number): string {
  if (!Number.isFinite(r)) return '—'
  const a = Math.abs(r)
  const d = a === 0 ? 2 : a >= 1000 ? 2 : a >= 100 ? 3 : a >= 1 ? 4 : a >= 0.01 ? 5 : 7
  return new Intl.NumberFormat('en-US', {
    maximumFractionDigits: d,
    minimumFractionDigits: Math.min(2, d),
  }).format(r)
}

export function fmtPct(p: number): string {
  if (!Number.isFinite(p)) return '—'
  return `${p > 0 ? '+' : ''}${p.toFixed(2)}%`
}

export function pctChange(points: { r: number }[] | undefined): number {
  if (!points || points.length < 2) return NaN
  const first = points[0].r
  const last = points[points.length - 1].r
  return first === 0 ? NaN : ((last - first) / first) * 100
}

/**
 * Arithmetic in the amount field: "1200*4+50", "45/3", parentheses, unary
 * minus, thousand separators. Returns null for anything that doesn't parse.
 */
export function evalExpr(src: string): number | null {
  const s = src
    .replace(/[,\s_]/g, '')
    .replace(/[×xX]/g, '*')
    .replace(/[÷]/g, '/')
    .replace(/＋/g, '+')
    .replace(/－/g, '-')
    .replace(/（/g, '(')
    .replace(/）/g, ')')
  if (!s) return null
  let i = 0

  function expr(): number | null {
    let v = term()
    if (v === null) return null
    while (s[i] === '+' || s[i] === '-') {
      const op = s[i++]
      const r = term()
      if (r === null) return null
      v = op === '+' ? v + r : v - r
    }
    return v
  }

  function term(): number | null {
    let v = unary()
    if (v === null) return null
    while (s[i] === '*' || s[i] === '/') {
      const op = s[i++]
      const r = unary()
      if (r === null) return null
      if (op === '/') {
        if (r === 0) return null
        v = v / r
      } else {
        v = v * r
      }
    }
    return v
  }

  function unary(): number | null {
    if (s[i] === '-') {
      i++
      const v = unary()
      return v === null ? null : -v
    }
    if (s[i] === '+') {
      i++
      return unary()
    }
    return atom()
  }

  function atom(): number | null {
    if (s[i] === '(') {
      i++
      const v = expr()
      if (v === null || s[i] !== ')') return null
      i++
      return v
    }
    const start = i
    while (i < s.length && (s[i] >= '0' && s[i] <= '9' || s[i] === '.')) i++
    if (i === start) return null
    const n = Number(s.slice(start, i))
    return Number.isFinite(n) ? n : null
  }

  const v = expr()
  return v !== null && i === s.length && Number.isFinite(v) ? v : null
}
