<script lang="ts">
  import { fade, scale } from 'svelte/transition'
  import { Moon, Sun } from 'lucide-svelte'

  let dark = $state(document.documentElement.classList.contains('dark'))

  function toggle() {
    dark = !dark
    document.documentElement.classList.toggle('dark', dark)
    localStorage.setItem('fx-theme', dark ? 'dark' : 'light')
  }
</script>

<button
  type="button"
  onclick={toggle}
  aria-label={dark ? '切换到浅色模式' : '切换到深色模式'}
  class="flex size-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
>
  {#if dark}
    <span in:scale={{ start: 0.6, duration: 160 }} class="contents">
      <Sun class="size-4" />
    </span>
  {:else}
    <span in:scale={{ start: 0.6, duration: 160 }} class="contents">
      <Moon class="size-4" />
    </span>
  {/if}
</button>
