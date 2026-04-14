<template>
  <!-- Admin custom home override -->
  <div v-if="homeContent" class="min-h-screen">
    <iframe v-if="isHomeContentUrl" :src="homeContent.trim()" class="h-screen w-full border-0" allowfullscreen />
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Proxy IP Landing Page -->
  <div
    v-else
    class="relative flex min-h-screen flex-col overflow-hidden bg-gradient-to-br from-slate-50 via-blue-50/30 to-slate-100 dark:from-dark-950 dark:via-dark-900 dark:to-dark-950"
  >
    <!-- Background blobs -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <div class="absolute -right-40 -top-40 h-96 w-96 rounded-full bg-blue-400/20 blur-3xl"></div>
      <div class="absolute -bottom-40 -left-40 h-96 w-96 rounded-full bg-blue-500/15 blur-3xl"></div>
      <div class="absolute inset-0 bg-[linear-gradient(rgba(59,130,246,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(59,130,246,0.03)_1px,transparent_1px)] bg-[size:64px_64px]"></div>
    </div>

    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <div class="flex items-center gap-2.5">
          <div class="h-9 w-9 overflow-hidden rounded-xl shadow-md">
            <img :src="siteLogo || '/logo.png'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="font-bold text-gray-900 dark:text-white text-lg">{{ siteName }}</span>
        </div>

        <div class="flex items-center gap-3">
          <LocaleSwitcher />
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

          <template v-if="isAuthenticated">
            <router-link
              to="/proxy/marketplace"
              class="inline-flex items-center gap-1.5 rounded-full bg-blue-600 px-4 py-1.5 text-sm font-medium text-white transition-colors hover:bg-blue-700"
            >
              购买代理
            </router-link>
            <router-link
              :to="dashboardPath"
              class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-3 transition-colors hover:bg-gray-800 dark:bg-gray-800"
            >
              <span class="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-br from-blue-400 to-blue-600 text-[10px] font-semibold text-white">
                {{ userInitial }}
              </span>
              <span class="text-xs font-medium text-white">控制台</span>
            </router-link>
          </template>
          <template v-else>
            <router-link
              to="/login"
              class="inline-flex items-center rounded-full bg-gray-900 px-4 py-1.5 text-sm font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800"
            >
              登录
            </router-link>
          </template>
        </div>
      </nav>
    </header>

    <!-- Hero -->
    <main class="relative z-10 flex-1 px-6 pt-16 pb-24">
      <div class="mx-auto max-w-6xl">
        <div class="mb-16 flex flex-col items-center justify-between gap-12 lg:flex-row lg:gap-16">
          <!-- Left: Text -->
          <div class="flex-1 text-center lg:text-left">
            <div class="mb-4 inline-flex items-center gap-2 rounded-full bg-blue-100 dark:bg-blue-900/40 px-3 py-1 text-xs font-medium text-blue-700 dark:text-blue-300">
              <span class="h-1.5 w-1.5 rounded-full bg-blue-500 animate-pulse"></span>
              真实住宅 IP · 即买即用
            </div>
            <h1 class="mb-5 text-4xl font-bold text-gray-900 dark:text-white md:text-5xl lg:text-6xl leading-tight">
              家庭 IP 代理<br />
              <span class="text-blue-600 dark:text-blue-400">按需租用</span>
            </h1>
            <p class="mb-8 text-lg text-gray-600 dark:text-dark-300 md:text-xl max-w-xl">
              真实住宅 IP，支持 HTTP 代理与 VLESS 双协议，按天/月计费，流量封顶保障成本，独享 IP 不共享。
            </p>

            <div class="flex flex-wrap gap-3 justify-center lg:justify-start">
              <router-link
                :to="isAuthenticated ? '/proxy/marketplace' : '/login'"
                class="inline-flex items-center gap-2 rounded-xl bg-blue-600 hover:bg-blue-700 px-7 py-3 text-base font-semibold text-white shadow-lg shadow-blue-500/30 transition-colors"
              >
                立即购买
                <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25" />
                </svg>
              </router-link>
              <router-link
                v-if="isAuthenticated"
                to="/proxy/rentals"
                class="inline-flex items-center gap-2 rounded-xl border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 px-7 py-3 text-base font-semibold text-gray-700 dark:text-gray-200 hover:border-blue-400 transition-colors"
              >
                我的代理
              </router-link>
            </div>
          </div>

          <!-- Right: IP Node visualization -->
          <div class="flex flex-1 justify-center lg:justify-end">
            <div class="relative w-80">
              <!-- Globe card -->
              <div class="rounded-2xl bg-gradient-to-br from-blue-600 to-blue-800 p-6 text-white shadow-2xl shadow-blue-500/30">
                <div class="mb-4 flex items-center gap-2">
                  <div class="h-2 w-2 rounded-full bg-green-400 animate-pulse"></div>
                  <span class="text-sm font-medium opacity-80">节点在线</span>
                </div>
                <div class="text-4xl mb-4">🌍</div>
                <div class="space-y-2">
                  <div
                    v-for="(node, i) in previewNodes"
                    :key="i"
                    class="flex items-center justify-between rounded-lg bg-white/10 px-3 py-2 text-sm"
                  >
                    <span>{{ node.flag }} {{ node.name }}</span>
                    <span class="text-xs text-green-300 font-medium">可用</span>
                  </div>
                </div>
                <div class="mt-4 pt-3 border-t border-white/20 text-xs text-white/60 text-center">
                  支持 HTTP · VLESS · Shadowrocket
                </div>
              </div>
              <!-- Floating badge -->
              <div class="absolute -bottom-4 -right-4 rounded-xl bg-white dark:bg-gray-800 shadow-lg px-4 py-2 text-sm font-semibold text-gray-900 dark:text-white border border-gray-100 dark:border-gray-700">
                独享 IP · 不共享
              </div>
            </div>
          </div>
        </div>

        <!-- Feature tags -->
        <div class="mb-12 flex flex-wrap justify-center gap-4">
          <div
            v-for="tag in featureTags"
            :key="tag.text"
            class="inline-flex items-center gap-2 rounded-full border border-gray-200/50 bg-white/80 dark:border-dark-700/50 dark:bg-dark-800/80 px-5 py-2.5 shadow-sm backdrop-blur-sm"
          >
            <span class="text-base">{{ tag.icon }}</span>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ tag.text }}</span>
          </div>
        </div>

        <!-- Features Grid -->
        <div class="grid gap-6 md:grid-cols-3">
          <div
            v-for="f in features"
            :key="f.title"
            class="group rounded-2xl border border-gray-200/50 bg-white/60 dark:border-dark-700/50 dark:bg-dark-800/60 p-6 backdrop-blur-sm transition-all hover:shadow-xl hover:shadow-blue-500/10"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl shadow-lg transition-transform group-hover:scale-110 text-2xl"
              :class="f.bg"
            >{{ f.icon }}</div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">{{ f.title }}</h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 dark:border-dark-800/50 px-6 py-6">
      <p class="text-center text-sm text-gray-500 dark:text-dark-400">
        &copy; {{ currentYear }} {{ siteName }}. 保留所有权利。
      </p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'

const authStore = useAuthStore()
const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '家庭 IP 代理')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const isHomeContentUrl = computed(() => {
  const c = homeContent.value.trim()
  return c.startsWith('http://') || c.startsWith('https://')
})

const isDark = ref(document.documentElement.classList.contains('dark'))
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => authStore.user?.email?.charAt(0).toUpperCase() || '')
const currentYear = computed(() => new Date().getFullYear())

const previewNodes = [
  { flag: '🇺🇸', name: '美国 · 洛杉矶' },
  { flag: '🇯🇵', name: '日本 · 东京' },
  { flag: '🇸🇬', name: '新加坡' },
]

const featureTags = [
  { icon: '🏠', text: '真实住宅 IP' },
  { icon: '🔒', text: '独享不共享' },
  { icon: '📊', text: '流量精确计费' },
  { icon: '⚡', text: '即买即用' },
]

const features = [
  {
    icon: '🌐',
    bg: 'bg-gradient-to-br from-blue-500 to-blue-600',
    title: '多国住宅节点',
    desc: '覆盖美国、日本、新加坡等多个地区，真实家庭宽带 IP，绕过各类检测。',
  },
  {
    icon: '🔌',
    bg: 'bg-gradient-to-br from-violet-500 to-violet-600',
    title: 'HTTP + VLESS 双协议',
    desc: '同时提供 HTTP 代理账号密码和 VLESS 链接，兼容 Shadowrocket 等主流客户端。',
  },
  {
    icon: '💳',
    bg: 'bg-gradient-to-br from-emerald-500 to-emerald-600',
    title: '灵活计费方案',
    desc: '按天或按月租用，设置流量上限，支持支付宝、微信支付、Stripe 多种支付方式。',
  },
]

function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  const saved = localStorage.getItem('theme')
  if (saved === 'dark' || (!saved && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
  authStore.checkAuth()
  if (!appStore.publicSettingsLoaded) appStore.fetchPublicSettings()
})
</script>
