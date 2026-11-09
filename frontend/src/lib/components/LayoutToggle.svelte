<script lang="ts">
  import { scale } from 'svelte/transition'
  import { RectangleHorizontal, RectangleVertical } from 'lucide-svelte'
  import { viewport, toggleLayout } from '$lib/media.svelte'

  const forced = $derived(viewport.forced)
</script>

<button
  type="button"
  onclick={toggleLayout}
  aria-label={forced ? '恢复自动布局' : '切换为横排布局'}
  aria-pressed={forced}
  title={forced ? '恢复自动布局' : '切换为横排布局'}
  class="flex size-8 items-center justify-center rounded-lg transition-colors {forced
    ? 'bg-accent text-foreground'
    : 'text-muted-foreground hover:bg-accent hover:text-foreground'}"
>
  {#if forced}
    <span in:scale={{ start: 0.6, duration: 160 }} class="contents">
      <RectangleHorizontal class="size-4" />
    </span>
  {:else}
    <span in:scale={{ start: 0.6, duration: 160 }} class="contents">
      <RectangleVertical class="size-4" />
    </span>
  {/if}
</button>
