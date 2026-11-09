// Shared viewport flag for the landscape layout. CSS handles the grid via the
// `landscape-short` custom variant (media query OR .force-landscape on <html>);
// charts need the flag in JS to resize themselves.
export const LANDSCAPE_SHORT_QUERY =
  '(orientation: landscape) and (max-height: 500px)'

export const viewport = $state({
  landscapeShort: false,
  wide: false,
  forced:
    localStorage.getItem('fx-layout') === 'landscape' ||
    new URLSearchParams(location.search).get('layout') === 'landscape',
})

let mqPhone: MediaQueryList | null = null
let mqWide: MediaQueryList | null = null

function apply() {
  viewport.landscapeShort = viewport.forced || !!mqPhone?.matches
  // Three-column dashboard only when forced AND the screen is wide enough.
  viewport.wide = viewport.landscapeShort && !!mqWide?.matches
}

if (typeof window !== 'undefined' && 'matchMedia' in window) {
  mqPhone = window.matchMedia(LANDSCAPE_SHORT_QUERY)
  mqWide = window.matchMedia('(min-width: 900px)')
  mqPhone.addEventListener('change', apply)
  mqWide.addEventListener('change', apply)
  if (viewport.forced) {
    document.documentElement.classList.add('force-landscape')
  }
  apply()
}

export function toggleLayout() {
  viewport.forced = !viewport.forced
  document.documentElement.classList.toggle('force-landscape', viewport.forced)
  localStorage.setItem('fx-layout', viewport.forced ? 'landscape' : 'auto')
  apply()
}
