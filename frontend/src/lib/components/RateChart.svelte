<script lang="ts">
  import type { Point } from '$lib/api'
  import { fmtRate } from '$lib/format'

  let {
    points = [],
    height = 260,
    fill = false,
  }: { points?: Point[]; height?: number; fill?: boolean } = $props()

  let width = $state(0)
  let measuredH = $state(0)
  let hover = $state<number | null>(null)
  const uid = `fx-${Math.random().toString(36).slice(2, 8)}`

  const H = $derived(fill ? measuredH || 320 : height)

  const PAD = { t: 16, r: 54, b: 26, l: 12 }

  type XY = { x: number; y: number; p: Point }

  let geo = $derived.by(() => {
    const w = width || 640
    const h = H
    const plotW = Math.max(10, w - PAD.l - PAD.r)
    const plotH = Math.max(10, h - PAD.t - PAD.b)
    if (points.length < 2) return { w, plotW, plotH, xy: [] as XY[], min: 0, max: 1 }
    const rs = points.map((p) => p.r)
    let min = Math.min(...rs)
    let max = Math.max(...rs)
    if (min === max) {
      const d = Math.abs(min) || 1
      min -= d * 0.1
      max += d * 0.1
    } else {
      const pad = (max - min) * 0.08
      min -= pad
      max += pad
    }
    const span = max - min
    const xy = points.map((p, i) => ({
      x: PAD.l + (i / (points.length - 1)) * plotW,
      y: PAD.t + (1 - (p.r - min) / span) * plotH,
      p,
    }))
    return { w, plotW, plotH, xy, min, max }
  })

  function niceTicks(min: number, max: number, count: number): number[] {
    const span = max - min
    if (span <= 0) return [min]
    const step0 = span / count
    const mag = Math.pow(10, Math.floor(Math.log10(step0)))
    const step =
      [1, 2, 5, 10].map((m) => m * mag).find((s) => span / s <= count) ?? mag * 10
    const out: number[] = []
    for (let v = Math.ceil(min / step) * step; v <= max + 1e-9; v += step) out.push(v)
    return out
  }

  let ticks = $derived(niceTicks(geo.min, geo.max, 4))
  let yFor = (v: number) =>
    PAD.t + (1 - (v - geo.min) / (geo.max - geo.min || 1)) * geo.plotH

  let geoH = $derived(H)

  let linePath = $derived.by(() => {
    if (geo.xy.length < 2) return ''
    return (
      'M' + geo.xy.map((q) => `${q.x.toFixed(1)} ${q.y.toFixed(1)}`).join('L')
    )
  })

  let areaPath = $derived.by(() => {
    if (!linePath) return ''
    const bottom = PAD.t + geo.plotH
    const first = geo.xy[0]
    const last = geo.xy[geo.xy.length - 1]
    return `${linePath}L${last.x.toFixed(1)} ${bottom}L${first.x.toFixed(1)} ${bottom}Z`
  })

  const X_LABELS = 4
  let xLabels = $derived.by(() => {
    if (geo.xy.length < 2) return []
    const n = geo.xy.length
    const idxs = Array.from({ length: X_LABELS }, (_, i) =>
      Math.round((i * (n - 1)) / (X_LABELS - 1)),
    )
    return [...new Set(idxs)].map((i) => {
      const q = geo.xy[i]
      const date = new Date(q.p.d + 'T00:00:00')
      const long = points.length > 200
      return {
        x: q.x,
        label: date.toLocaleDateString('zh-CN', {
          year: long ? 'numeric' : undefined,
          month: 'short',
          day: 'numeric',
        }),
      }
    })
  })

  let hoverXY = $derived(hover !== null ? geo.xy[hover] : null)

  function onMove(e: PointerEvent) {
    const rect = (e.currentTarget as SVGElement).getBoundingClientRect()
    const x = e.clientX - rect.left
    if (geo.xy.length < 2) return
    const plotW = geo.plotW
    const t = Math.min(1, Math.max(0, (x - PAD.l) / plotW))
    hover = Math.round(t * (geo.xy.length - 1))
  }

  function fmtDate(d: string, withYear: boolean) {
    const date = new Date(d + 'T00:00:00')
    return date.toLocaleDateString('zh-CN', {
      year: withYear ? 'numeric' : undefined,
      month: 'short',
      day: 'numeric',
    })
  }

  let withYear = $derived(points.length > 200)
</script>

<div
  class={fill ? 'relative h-full w-full select-none' : 'relative w-full select-none'}
  style={fill ? undefined : `height: ${height}px`}
  bind:clientWidth={width}
  bind:clientHeight={measuredH}
>
  {#if geo.xy.length >= 2}
    <svg
      class="absolute inset-0 h-full w-full overflow-visible"
      viewBox="0 0 {geo.w} {geoH}"
      onpointermove={onMove}
      onpointerleave={() => (hover = null)}
      role="img"
      aria-label="汇率走势图"
    >
      <defs>
        <linearGradient id="{uid}-area" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stop-color="var(--color-fx)" stop-opacity="0.16" />
          <stop offset="100%" stop-color="var(--color-fx)" stop-opacity="0" />
        </linearGradient>
      </defs>

      <!-- gridlines + y labels -->
      {#each ticks as t (t)}
        <line
          x1={PAD.l}
          x2={PAD.l + geo.plotW}
          y1={yFor(t)}
          y2={yFor(t)}
          class="stroke-border"
          stroke-width="1"
        />
        <text
          x={PAD.l + geo.plotW + 8}
          y={yFor(t)}
          dominant-baseline="middle"
          class="fill-muted-foreground text-[10px] tabular-nums"
        >
          {fmtRate(t)}
        </text>
      {/each}

      <!-- x labels -->
      {#each xLabels as l (l.x)}
        <text
          x={l.x}
          y={geoH - 8}
          text-anchor={l.x < 40 ? 'start' : l.x > geo.w - 60 ? 'end' : 'middle'}
          class="fill-muted-foreground text-[10px] tabular-nums"
        >
          {l.label}
        </text>
      {/each}

      <!-- area + line -->
      {#key points}
        <path d={areaPath} fill="url(#{uid}-area)" class="fx-area" />
        <path
          d={linePath}
          fill="none"
          stroke="var(--color-fx)"
          stroke-width="1.8"
          stroke-linecap="round"
          stroke-linejoin="round"
          pathLength="1"
          class="fx-line"
        />
      {/key}

      <!-- crosshair -->
      {#if hoverXY}
        <line
          x1={hoverXY.x}
          x2={hoverXY.x}
          y1={PAD.t}
          y2={PAD.t + geo.plotH}
          class="stroke-muted-foreground/40"
          stroke-width="1"
          stroke-dasharray="3 3"
        />
        <circle
          cx={hoverXY.x}
          cy={hoverXY.y}
          r="4"
          class="fill-fx stroke-background"
          stroke-width="2"
        />
      {/if}
    </svg>

    {#if hoverXY}
      <div
        class="pointer-events-none absolute z-10 -translate-x-1/2 rounded-md border bg-popover px-2 py-1 text-center shadow-md"
        style="left: {Math.min(Math.max(hoverXY.x, 44), geo.w - 44)}px; top: {hoverXY.y - 14}px; transform: translate(-50%, -100%);"
      >
        <div class="text-[10px] text-muted-foreground">
          {fmtDate(hoverXY.p.d, withYear)}
        </div>
        <div class="text-xs font-semibold tabular-nums">{fmtRate(hoverXY.p.r)}</div>
      </div>
    {/if}
  {:else}
    <div class="absolute inset-0 animate-pulse rounded-xl bg-muted/50"></div>
  {/if}
</div>

<style>
  .fx-line {
    stroke-dasharray: 1;
    stroke-dashoffset: 1;
    animation: fx-draw 0.9s cubic-bezier(0.3, 0, 0.2, 1) forwards;
  }
  .fx-area {
    opacity: 0;
    animation: fx-fade 0.6s 0.35s ease-out forwards;
  }
  @keyframes fx-fade {
    to {
      opacity: 1;
    }
  }
</style>
