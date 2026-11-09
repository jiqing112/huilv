<script lang="ts">
  let {
    points = [],
    class: klass = '',
    stroke = 'var(--color-muted-foreground)',
    strokeWidth = 1.6,
  }: {
    points?: { r: number }[]
    class?: string
    stroke?: string
    strokeWidth?: number
  } = $props()

  const W = 100
  const H = 30
  const P = 2

  let path = $derived.by(() => {
    if (points.length < 2) return ''
    const rs = points.map((p) => p.r)
    const min = Math.min(...rs)
    const max = Math.max(...rs)
    const span = max - min || 1
    const step = (W - P * 2) / (points.length - 1)
    const xy = points.map(
      (p, i) => [P + i * step, P + (1 - (p.r - min) / span) * (H - P * 2)] as const,
    )
    return 'M' + xy.map(([x, y]) => `${x.toFixed(2)} ${y.toFixed(2)}`).join('L')
  })

  let up = $derived(points.length >= 2 && points[points.length - 1].r >= points[0].r)
</script>

<svg
  viewBox="0 0 {W} {H}"
  preserveAspectRatio="none"
  class={klass}
  aria-hidden="true"
>
  {#if path}
    <path
      d={path}
      fill="none"
      stroke={stroke}
      stroke-width={strokeWidth}
      vector-effect="non-scaling-stroke"
      stroke-linecap="round"
      stroke-linejoin="round"
      opacity={up ? 0.9 : 0.65}
    />
  {/if}
</svg>
