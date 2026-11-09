<script lang="ts">
  let { class: klass = '' }: { class?: string } = $props()
</script>

<span
  class="flex size-7 items-center justify-center rounded-[8px] bg-foreground text-background {klass}"
>
  <svg
    class="size-4"
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    stroke-width="2.2"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
  >
    <path d="M4 8h13M13.5 4.5 17 8l-3.5 3.5" />
    <path d="M20 16H7m3.5 3.5L7 16l3.5-3.5" />
  </svg>
</span>
