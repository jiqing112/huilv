<script lang="ts">
  import { bootstrap, market } from '$lib/state.svelte'
  import Logo from '$lib/components/Logo.svelte'
  import ThemeToggle from '$lib/components/ThemeToggle.svelte'
  import LayoutToggle from '$lib/components/LayoutToggle.svelte'
  import Converter from '$lib/components/Converter.svelte'
  import ChartView from '$lib/components/ChartView.svelte'
  import RateBoard from '$lib/components/RateBoard.svelte'

  bootstrap()
</script>

  <div class="relative min-h-dvh">
    <div class="fx-glow pointer-events-none absolute inset-x-0 top-0 h-[420px]"></div>

    <div class="fx-shell relative mx-auto w-full max-w-2xl px-4 pb-16 landscape-short:max-w-5xl sm:px-6">
      <header class="flex items-center justify-between py-5 landscape-short:py-3">
        <a href="/" class="flex items-center gap-2.5 rounded-lg" aria-label="汇率">
          <Logo />
          <span class="text-[15px] font-semibold tracking-tight">汇率</span>
        </a>
        <div class="flex items-center gap-1">
          <LayoutToggle />
          <ThemeToggle />
        </div>
      </header>

      <main
        class="flex flex-col gap-12 landscape-short:grid landscape-short:grid-cols-[minmax(0,1.05fr)_minmax(0,1fr)] landscape-short:items-start landscape-short:gap-x-8 landscape-short:gap-y-6"
      >
        <div
          data-panel="converter"
          class="landscape-short:sticky landscape-short:top-3 landscape-short:min-w-0 landscape-short:col-start-1 landscape-short:row-start-1 landscape-short:row-span-2 landscape-short:self-start"
        >
          <Converter />
        </div>
        <div data-panel="chart" class="landscape-short:min-w-0 landscape-short:col-start-2 landscape-short:row-start-1">
          <ChartView />
        </div>
        <div data-panel="board" class="landscape-short:min-w-0 landscape-short:col-start-2 landscape-short:row-start-2">
          <RateBoard />
        </div>
      </main>

      <footer class="pt-14 text-center text-xs text-muted-foreground/70 tabular-nums landscape-short:pt-8">
        {#if market.date}
          欧洲央行参考汇率 · {market.date}
        {/if}
      </footer>
    </div>
  </div>
