<script lang="ts">
  import { Search } from 'lucide-svelte'
  import { app, market, hist, crossRate, ensureBoard, saveApp } from '$lib/state.svelte'
  import { fmtRate, pctChange } from '$lib/format'
  import CurrencyPicker from './CurrencyPicker.svelte'
  import CurrencyBadge from './CurrencyBadge.svelte'
  import Sparkline from './Sparkline.svelte'
  import ChangeChip from './ChangeChip.svelte'

  let q = $state('')
  let searchEl = $state<HTMLInputElement>()

  const rows = $derived.by(() => {
    const s = q.trim().toLowerCase()
    return market.currencies
      .filter((c) => c.code !== app.from)
      .filter(
        (c) =>
          !s ||
          c.code.toLowerCase().includes(s) ||
          c.zh.includes(q.trim()) ||
          c.name.toLowerCase().includes(s),
      )
      .sort((a, b) => b.popular - a.popular || a.code.localeCompare(b.code))
  })

  const board = $derived(hist.board[app.from] ?? {})

  function rateOf(code: string): number {
    return crossRate(market.usd, app.from, code)
  }

  function onPick(code: string) {
    app.to = code
    saveApp()
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  function onKeydown(e: KeyboardEvent) {
    const tag = (e.target as HTMLElement)?.tagName
    if (tag === 'INPUT' || tag === 'TEXTAREA') return
    if (e.key === '/') {
      e.preventDefault()
      searchEl?.focus()
    }
  }

  $effect(() => {
    const f = app.from
    if (f) ensureBoard(f)
  })
</script>

<svelte:window onkeydown={onKeydown} />

<div class="animate-[fx-rise_0.5s_0.12s_ease-out_backwards] rounded-2xl border bg-card shadow-[0_1px_2px_rgba(0,0,0,0.04)]">
  <div class="flex items-center gap-2 border-b px-3 py-1.5">
    <Search class="size-4 shrink-0 text-muted-foreground" />
    <input
      bind:this={searchEl}
      bind:value={q}
      class="h-8 w-full min-w-0 flex-1 bg-transparent text-sm outline-none placeholder:text-muted-foreground"
      placeholder="搜索"
      aria-label="搜索货币"
    />
    <div class="h-4 w-px shrink-0 bg-border"></div>
    <CurrencyPicker bind:value={app.from} label="基准货币" compact />
  </div>

  <div class="fx-board-rows p-1">
    {#if market.ok}
      {#each rows as c (c.code)}
        <button
          type="button"
          onclick={() => onPick(c.code)}
          title="设为目标货币"
          class="group flex h-11 w-full items-center gap-3 rounded-lg px-2.5 text-left transition-colors hover:bg-accent"
        >
          <CurrencyBadge code={c.code} class="h-6 w-6 text-[9px]" />
          <span class="text-sm font-semibold tracking-wide">{c.code}</span>
          <span class="hidden min-w-0 truncate text-xs text-muted-foreground sm:inline landscape-short:hidden">{c.zh}</span>
          <span class="ml-auto flex items-center gap-3">
            {#if board[c.code]?.length > 1}
              <Sparkline points={board[c.code]} class="hidden h-5 w-16 sm:block" />
            {/if}
            <span class="w-20 text-right text-sm font-medium tabular-nums sm:w-24">
              {fmtRate(rateOf(c.code))}
            </span>
            <ChangeChip pct={pctChange(board[c.code])} class="w-16 justify-end text-xs sm:w-20" />
          </span>
        </button>
      {/each}
      {#if rows.length === 0}
        <div class="px-3 py-10 text-center text-sm text-muted-foreground">无结果</div>
      {/if}
    {:else}
      {#each Array(8) as _, i (i)}
        <div class="flex h-11 items-center gap-3 rounded-lg px-2.5">
          <div class="h-6 w-6 animate-pulse rounded-md bg-muted"></div>
          <div class="h-3 w-10 animate-pulse rounded bg-muted"></div>
          <div class="ml-auto h-3 w-20 animate-pulse rounded bg-muted"></div>
        </div>
      {/each}
    {/if}
  </div>
</div>
