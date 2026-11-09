<script lang="ts">
  let { code, class: klass = '' }: { code: string; class?: string } = $props()
</script>

<span
  class="inline-flex shrink-0 items-center justify-center rounded-md bg-foreground/[0.06] font-semibold tracking-wide text-foreground/75 select-none dark:bg-foreground/10 dark:text-foreground/80 {klass}"
>
  {code}
</span>
