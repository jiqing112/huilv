// Small rAF number tween for the converter result. Reduced motion -> instant.
export function tween(
  from: number,
  to: number,
  duration: number,
  onTick: (v: number) => void,
): void {
  if (matchMedia('(prefers-reduced-motion: reduce)').matches || from === to) {
    onTick(to)
    return
  }
  const t0 = performance.now()
  const step = (t: number) => {
    const p = Math.min(1, (t - t0) / duration)
    const e = 1 - Math.pow(1 - p, 3)
    onTick(from + (to - from) * e)
    if (p < 1) requestAnimationFrame(step)
  }
  requestAnimationFrame(step)
}
