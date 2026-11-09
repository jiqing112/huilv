<script lang="ts">
  import { app, hist, ensurePair, saveApp } from '$lib/state.svelte'
  import { viewport } from '$lib/media.svelte'
  import { pctChange, fmtRate } from '$lib/format'
  import { ArrowDownRight, ArrowUpRight } from 'lucide-svelte'
  import RateChart from './RateChart.svelte'

  const RANGES = [
    { days: 30, label: '1月' },
    { days: 90, label: '3月' },
    { days: 180, label: '6月' },
    { days: 365, label: '1年' },
    { days: 1825, label: '5年' },
  ]

  const points = $derived(hist.pair[`${app.from}|${app.to}|${app.range}`] ?? [])

  const stats = $derived.by(() => {
    if (points.length < 2) return null
    const rs = points.map((p) => p.r)
    return {
      high: Math.max(...rs),
      low: Math.min(...rs),
      chg: pctChange(points),
    }
  })

  $effect(() => {
    const f = app.from
    const t = app.to
    ensurePair(f, t, app.range)
  })

  function setRange(days: number) {
    app.range = days
    saveApp()
  }
</script>

<div class="animate-[fx-rise_0.5s_0.06s_ease-out_backwards] flex flex-col rounded-2xl border bg-card p-4 shadow-[0_1px_2px_rgba(0,0,0,0.04)] sm:p-5">
  <div class="flex flex-wrap items-center justify-between gap-2">
    <div
      class="inline-flex rounded-lg border bg-muted/60 p-0.5"
      role="tablist"
      aria-label="时间范围"
    >
      {#each RANGES as r (r.days)}
        <button
          type="button"
          role="tab"
          aria-selected={app.range === r.days}
          onclick={() => setRange(r.days)}
          class="h-7 rounded-md px-2.5 text-xs font-medium transition-all {app.range === r.days
            ? 'bg-background text-foreground shadow-xs'
            : 'text-muted-foreground hover:text-foreground'}"
        >
          {r.label}
        </button>
      {/each}
    </div>

    {#if stats}
      <div class="flex items-center gap-3.5 text-xs tabular-nums text-muted-foreground">
        <span class="inline-flex items-center gap-0.5 text-up">
          <ArrowUpRight class="size-3" />{fmtRate(stats.high)}
        </span>
        <span class="inline-flex items-center gap-0.5 text-down">
          <ArrowDownRight class="size-3" />{fmtRate(stats.low)}
        </span>
        <span class={stats.chg >= 0 ? 'text-up' : 'text-down'}>
          {stats.chg > 0 ? '+' : ''}{stats.chg.toFixed(2)}%
        </span>
      </div>
    {/if}
  </div>

  <div class="mt-3 min-h-0 flex-1 landscape-short:flex-1">
    <RateChart {points} fill={viewport.wide} height={viewport.landscapeShort && !viewport.wide ? 180 : 260} />
  </div>
</div>
