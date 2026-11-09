<script lang="ts">
  import { app, market, hist, ensurePair, amountValue, crossRate, saveApp } from '$lib/state.svelte'
  import { fmtAmount, fmtRate, pctChange } from '$lib/format'
  import { tween } from '$lib/tween'
  import { ArrowUpDown, Calculator } from 'lucide-svelte'
  import CurrencyPicker from './CurrencyPicker.svelte'
  import Sparkline from './Sparkline.svelte'
  import ChangeChip from './ChangeChip.svelte'

  let side = $state<'from' | 'to'>('from')
  let swapAngle = $state(0)
  let passiveDisp = $state('')
  let prevNum = 0
  let fromInputEl = $state<HTMLInputElement>()

  // Discoverability for arithmetic in the amount field: a clickable example
  // on first use, then live "= result" feedback while an expression is typed.
  let calcHint = $state(localStorage.getItem('fx-calc-seen') !== '1')
  const isExpr = $derived(/[+\-*/()×÷]/.test(app.amount))
  const exprVal = $derived(isExpr ? amountValue() : null)

  function markCalcSeen() {
    if (!calcHint) return
    calcHint = false
    localStorage.setItem('fx-calc-seen', '1')
  }

  function demo() {
    side = 'from'
    app.amount = '250*4'
    markCalcSeen()
    fromInputEl?.focus()
  }

  const rate = $derived(crossRate(market.usd, app.from, app.to))

  const fromVal = $derived.by(() => {
    if (!market.ok) return null
    const n = amountValue()
    if (n === null) return null
    return side === 'from' ? n : rate > 0 ? n / rate : null
  })

  const toVal = $derived.by(() => {
    if (!market.ok) return null
    const n = amountValue()
    if (n === null) return null
    return side === 'from' ? (rate > 0 ? n * rate : null) : n
  })

  const decimals = $derived({
    from: market.currencies.find((c) => c.code === app.from)?.decimals ?? 2,
    to: market.currencies.find((c) => c.code === app.to)?.decimals ?? 2,
  })

  // Passive row shows the computed conversion, number-morphed on change.
  $effect(() => {
    const target = side === 'from' ? toVal : fromVal
    if (target === null) {
      passiveDisp = ''
      prevNum = 0
      return
    }
    const dec = side === 'from' ? decimals.to : decimals.from
    const prev = prevNum
    prevNum = target
    tween(prev, target, 220, (v) => (passiveDisp = fmtAmount(v, dec)))
  })

  function focusSide(s: 'from' | 'to') {
    if (side === s) return
    const v = s === 'from' ? fromVal : toVal
    side = s
    app.amount =
      v === null
        ? ''
        : String(Math.round(v * 1e6) / 1e6)
    saveApp()
  }

  function onInput(v: string) {
    app.amount = v
    if (/[+*/×÷]/.test(v)) markCalcSeen()
    saveApp()
  }

  function swap() {
    const f = app.from
    app.from = app.to
    app.to = f
    swapAngle += 180
    saveApp()
  }

  const spark = $derived(hist.pair[`${app.from}|${app.to}|30`] ?? [])
  const chg = $derived(pctChange(spark))

  $effect(() => {
    void market.ok
    ensurePair(app.from, app.to, 30)
  })
</script>

<div class="animate-[fx-rise_0.5s_ease-out_backwards]">
  <div
    class="relative rounded-2xl border bg-card shadow-[0_1px_2px_rgba(0,0,0,0.04)]"
  >
    <!-- from row -->
    <div class="flex items-center gap-3 py-3 pr-4 pl-5">
      <input
        bind:this={fromInputEl}
        class="w-0 min-w-0 flex-1 bg-transparent text-3xl font-medium tracking-tight tabular-nums outline-none placeholder:text-muted-foreground/50 sm:text-4xl"
        style="caret-color: var(--fx)"
        value={side === 'from' ? app.amount : passiveDisp}
        oninput={(e) => onInput(e.currentTarget.value)}
        onfocus={() => focusSide('from')}
        inputmode="decimal"
        autocomplete="off"
        spellcheck="false"
        placeholder="金额或算式"
        aria-label="源金额"
      />
      {#if side === 'from' && exprVal !== null}
        <span
          class="shrink-0 rounded-md bg-muted px-2 py-1 text-sm font-medium tabular-nums text-muted-foreground"
        >
          = {fmtAmount(exprVal, decimals.from)}
        </span>
      {/if}
      {#if calcHint && side === 'from' && !isExpr}
        <button
          type="button"
          onclick={demo}
          class="hidden shrink-0 items-center gap-1.5 rounded-md px-2 py-1 text-xs text-muted-foreground/70 transition-colors hover:bg-accent hover:text-foreground sm:inline-flex"
          title="金额栏支持加减乘除"
        >
          <Calculator class="size-3.5" />
          试试 250*4
        </button>
      {/if}
      <CurrencyPicker bind:value={app.from} label="源货币" />
    </div>

    <div class="border-t"></div>

    <!-- to row -->
    <div class="flex items-center gap-3 py-3 pr-4 pl-5">
      <input
        class="w-0 min-w-0 flex-1 bg-transparent text-3xl font-medium tracking-tight tabular-nums outline-none placeholder:text-muted-foreground/50 sm:text-4xl"
        style="caret-color: var(--fx)"
        value={side === 'to' ? app.amount : passiveDisp}
        oninput={(e) => onInput(e.currentTarget.value)}
        onfocus={() => focusSide('to')}
        inputmode="decimal"
        autocomplete="off"
        spellcheck="false"
        placeholder="金额或算式"
        aria-label="目标金额"
      />
      {#if side === 'to' && exprVal !== null}
        <span
          class="shrink-0 rounded-md bg-muted px-2 py-1 text-sm font-medium tabular-nums text-muted-foreground"
        >
          = {fmtAmount(exprVal, decimals.to)}
        </span>
      {/if}
      <CurrencyPicker bind:value={app.to} label="目标货币" />
    </div>

    <button
      type="button"
      onclick={swap}
      aria-label="交换货币"
      class="absolute top-1/2 right-4 z-10 flex size-8 items-center justify-center rounded-full border bg-background text-muted-foreground shadow-xs transition-[transform,color,background-color] duration-200 hover:text-foreground hover:shadow-sm active:scale-95"
      style="transform: translateY(-50%) rotate({swapAngle}deg)"
    >
      <ArrowUpDown class="size-3.5" />
    </button>
  </div>

  <!-- rate meta -->
  <div class="mt-3 flex items-center justify-between gap-4 px-1">
    <div class="min-w-0 space-y-0.5 text-sm tabular-nums text-muted-foreground">
      {#if market.ok && rate > 0}
        <div class="truncate">
          1 {app.from} =
          <span class="font-medium text-foreground">{fmtRate(rate)}</span>
          {app.to}
        </div>
        <div class="truncate">
          1 {app.to} =
          <span class="font-medium text-foreground">{fmtRate(1 / rate)}</span>
          {app.from}
        </div>
      {:else}
        <div class="inline-block h-4 w-40 animate-pulse rounded bg-muted"></div>
      {/if}
    </div>
    {#if spark.length > 1}
      <div class="flex shrink-0 items-center gap-2.5">
        <Sparkline points={spark} class="hidden h-5 w-24 sm:block" strokeWidth={2} stroke="var(--color-fx)" />
        <ChangeChip pct={chg} suffix="30日" class="text-xs" />
      </div>
    {/if}
  </div>
</div>
