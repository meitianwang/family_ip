/**
 * Composable for platform-specific styling (Claude, OpenAI, Gemini, etc.)
 */

import { computed, type Ref } from 'vue'

export interface PlatformStyleResult {
  badgeClass: string
  accentClass: string
  modelTagClass: string
  dotClass: string
  buttonClass: string
}

const STYLES: Record<string, PlatformStyleResult> = {
  claude: {
    badgeClass: 'bg-orange-100 text-orange-700 dark:bg-orange-900/40 dark:text-orange-300',
    accentClass: 'text-orange-600 dark:text-orange-400',
    modelTagClass:
      'border-orange-200 bg-orange-50 text-orange-700 dark:border-orange-800 dark:bg-orange-900/20 dark:text-orange-300',
    dotClass: 'bg-orange-500',
    buttonClass: 'bg-orange-500 hover:bg-orange-600 dark:bg-orange-600 dark:hover:bg-orange-500'
  },
  openai: {
    badgeClass: 'bg-green-100 text-green-700 dark:bg-green-900/40 dark:text-green-300',
    accentClass: 'text-green-600 dark:text-green-400',
    modelTagClass:
      'border-green-200 bg-green-50 text-green-700 dark:border-green-800 dark:bg-green-900/20 dark:text-green-300',
    dotClass: 'bg-green-500',
    buttonClass: 'bg-green-600 hover:bg-green-700 dark:bg-green-600 dark:hover:bg-green-500'
  },
  gemini: {
    badgeClass: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300',
    accentClass: 'text-blue-600 dark:text-blue-400',
    modelTagClass:
      'border-blue-200 bg-blue-50 text-blue-700 dark:border-blue-800 dark:bg-blue-900/20 dark:text-blue-300',
    dotClass: 'bg-blue-500',
    buttonClass: 'bg-blue-600 hover:bg-blue-700 dark:bg-blue-600 dark:hover:bg-blue-500'
  }
}

const DEFAULT_STYLE: PlatformStyleResult = {
  badgeClass: 'bg-slate-100 text-slate-700 dark:bg-slate-700 dark:text-slate-300',
  accentClass: 'text-slate-700 dark:text-slate-300',
  modelTagClass:
    'border-slate-200 bg-slate-50 text-slate-700 dark:border-slate-600 dark:bg-slate-800 dark:text-slate-300',
  dotClass: 'bg-slate-500',
  buttonClass: 'bg-slate-700 hover:bg-slate-800 dark:bg-slate-600 dark:hover:bg-slate-500'
}

/**
 * Returns reactive platform style classes based on a platform string.
 */
export function usePlatformStyle(platform: Ref<string> | (() => string)) {
  const resolve = typeof platform === 'function' ? platform : () => platform.value

  const style = computed(() => {
    const p = (resolve() || '').toLowerCase()
    return STYLES[p] || DEFAULT_STYLE
  })

  return {
    badgeClass: computed(() => style.value.badgeClass),
    accentClass: computed(() => style.value.accentClass),
    modelTagClass: computed(() => style.value.modelTagClass),
    dotClass: computed(() => style.value.dotClass),
    buttonClass: computed(() => style.value.buttonClass)
  }
}

/**
 * Non-reactive version for one-off lookups.
 */
export function getPlatformStyleClasses(platform: string): PlatformStyleResult {
  return STYLES[(platform || '').toLowerCase()] || DEFAULT_STYLE
}
