<script lang="ts">
  import { fly } from 'svelte/transition'
  import { Check, ChevronDown, Search } from 'lucide-svelte'
  import { market, saveApp } from '$lib/state.svelte'
  import { cn } from '$lib/utils'
  import CurrencyBadge from './CurrencyBadge.svelte'

  let {
    value = $bindable(),
    label = '选择货币',
    compact = false,
    align = 'right',
  }: {
    value: string
    label?: string
    compact?: boolean
    align?: 'left' | 'right'
  } = $props()

  let open = $state(false)
  let dropUp = $state(false)
  let q = $state('')
  let active = $state(0)
  let btn = $state<HTMLButtonElement>()
  let pop = $state<HTMLDivElement>()
  let inputEl = $state<HTMLInputElement>()

  const list = $derived.by(() => {
    const s = q.trim().toLowerCase()
    const all = market.currencies
    if (!s) {
      return [...all.filter((c) => c.popular), ...all.filter((c) => !c.popular)]
    }
    return all.filter(
      (c) =>
        c.code.toLowerCase().includes(s) ||
        c.zh.includes(q.trim()) ||
        c.name.toLowerCase().includes(s) ||
        c.symbol.toLowerCase().includes(s),
    )
  })

  const popularCount = $derived(q.trim() ? 0 : market.currencies.filter((c) => c.popular).length)

  function pick(code: string) {
    value = code
    saveApp()
    close()
  }

  function openPop() {
    // Flip upward when there isn't room below (e.g. the rates list at page bottom).
    const r = btn?.getBoundingClientRect()
    dropUp = !!r && r.bottom + 400 > window.innerHeight && r.top > 420
    q = ''
    active = Math.max(0, list.findIndex((c) => c.code === value))
    open = true
    setTimeout(() => inputEl?.focus(), 0)
  }

  function close() {
    open = false
    btn?.focus()
  }

  function onkeydown(e: KeyboardEvent) {
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      active = Math.min(active + 1, list.length - 1)
    } else if (e.key === 'ArrowUp') {
      e.preventDefault()
      active = Math.max(active - 1, 0)
    } else if (e.key === 'Enter') {
      e.preventDefault()
      if (list[active]) pick(list[active].code)
    } else if (e.key === 'Escape') {
      e.preventDefault()
      close()
    }
  }

  function onDocPointer(e: PointerEvent) {
    if (open && !btn?.contains(e.target as Node) && !pop?.contains(e.target as Node)) {
      open = false
    }
  }

  $effect(() => {
    if (open && pop) {
      const el = pop.querySelector('[data-active="true"]')
      el?.scrollIntoView({ block: 'nearest' })
    }
  })
</script>

<svelte:document onpointerdown={onDocPointer} />

<div class="relative">
  <button
    bind:this={btn}
    type="button"
    aria-haspopup="listbox"
    aria-expanded={open}
    aria-label={label}
    onclick={openPop}
    class={cn(
      'group flex items-center gap-2 rounded-lg transition-colors hover:bg-accent',
      compact ? 'px-1.5 py-1' : 'px-2 py-1.5',
    )}
  >
    <CurrencyBadge code={value} class={compact ? 'h-6 w-6 text-[9px]' : 'h-7 w-9 text-[10px]'} />
    <span class="text-sm font-semibold tracking-wide">{value}</span>
    <ChevronDown
      class={cn(
        'size-3.5 text-muted-foreground transition-transform duration-200 group-hover:text-foreground',
        open && 'rotate-180',
      )}
    />
  </button>

  {#if open}
    <div
      bind:this={pop}
      transition:fly={{ y: 4, duration: 140 }}
      class={cn(
        'absolute z-50 w-72 rounded-xl border bg-popover shadow-lg',
        'bg-card',
        dropUp ? 'bottom-full mb-2' : 'top-full mt-2',
        align === 'right' ? 'right-0' : 'left-0',
      )}
      role="dialog"
      aria-label={label}
    >
      <div class="flex items-center gap-2 border-b px-3">
        <Search class="size-3.5 shrink-0 text-muted-foreground" />
        <input
          bind:this={inputEl}
          bind:value={q}
          onkeydown={onkeydown}
          class="h-10 w-full bg-transparent text-sm outline-none placeholder:text-muted-foreground"
          placeholder="搜索"
          aria-label={label}
        />
      </div>
      <div class="max-h-72 overflow-y-auto p-1" role="listbox">
        {#each list as c, i (c.code)}
          {#if i === popularCount && popularCount > 0 && popularCount < list.length}
            <div class="mx-2 my-1 border-t"></div>
          {/if}
          <button
            type="button"
            role="option"
            aria-selected={c.code === value}
            data-active={i === active}
            onpointerenter={() => (active = i)}
            onclick={() => pick(c.code)}
            class="flex w-full items-center gap-2.5 rounded-lg px-2 py-1.5 text-left {i === active
              ? 'bg-accent'
              : ''}"
          >
            <CurrencyBadge code={c.code} class="h-6 w-6 text-[9px]" />
            <span class="text-sm font-semibold tracking-wide">{c.code}</span>
            <span class="truncate text-xs text-muted-foreground">{c.zh}</span>
            <span class="ml-auto flex items-center gap-1 text-xs text-muted-foreground">
              {c.symbol}
              {#if c.code === value}<Check class="size-3.5 text-foreground" />{/if}
            </span>
          </button>
        {/each}
        {#if list.length === 0}
          <div class="px-3 py-8 text-center text-sm text-muted-foreground">无结果</div>
        {/if}
      </div>
    </div>
  {/if}
</div>
