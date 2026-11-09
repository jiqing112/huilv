<script lang="ts">
  import { ArrowUpRight, ArrowDownRight } from 'lucide-svelte'
  import { fmtPct } from '$lib/format'

  let {
    pct,
    suffix = '',
    class: klass = '',
  }: { pct: number; suffix?: string; class?: string } = $props()

  let up = $derived(pct > 0)
  let down = $derived(pct < 0)
</script>

<span
  class="inline-flex items-center gap-0.5 tabular-nums {klass}
    {up ? 'text-up' : down ? 'text-down' : 'text-muted-foreground'}"
>
  {#if up}<ArrowUpRight class="size-3" />{:else if down}<ArrowDownRight class="size-3" />{/if}
  {fmtPct(pct)}{#if suffix}<span class="text-muted-foreground">· {suffix}</span>{/if}
</span>
