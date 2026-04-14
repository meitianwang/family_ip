<template>
  <AppLayout>
  <div class="max-w-5xl mx-auto px-4 py-8">
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white">家庭 IP 代理</h1>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">选择一个可用 IP 节点和套餐，购买后获得专属代理访问权限</p>
    </div>

    <!-- Filters -->
    <div class="flex gap-3 mb-6 flex-wrap">
      <select
        v-model="filterCountry"
        class="rounded-lg border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-sm px-3 py-2 text-gray-700 dark:text-gray-300"
      >
        <option value="">全部国家/地区</option>
        <option v-for="c in countries" :key="c.code" :value="c.code">{{ c.name }}</option>
      </select>
    </div>

    <!-- Node list -->
    <div v-if="loading" class="text-center py-12 text-gray-400">加载中...</div>
    <div v-else-if="nodes.length === 0" class="text-center py-12 text-gray-400">暂无可用节点</div>
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="node in filteredNodes"
        :key="node.id"
        class="border border-gray-200 dark:border-gray-700 rounded-xl p-4 bg-white dark:bg-gray-800 cursor-pointer hover:border-blue-400 transition-colors"
        :class="{ 'border-blue-500 ring-2 ring-blue-300': selectedNode?.id === node.id }"
        @click="selectNode(node)"
      >
        <div class="flex items-start justify-between">
          <div>
            <div class="flex items-center gap-2 mb-1">
              <span class="text-lg">{{ countryFlag(node.country_code) }}</span>
              <span class="font-semibold text-gray-900 dark:text-white">{{ node.city || node.country }}</span>
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400">{{ node.isp }}</div>
            <div class="text-xs text-gray-400 dark:text-gray-500 mt-1 font-mono">{{ node.ip_address }}</div>
          </div>
          <span class="text-xs px-2 py-1 rounded-full bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300">
            可用
          </span>
        </div>
        <div v-if="node.tags.length > 0" class="mt-3 flex flex-wrap gap-1">
          <span
            v-for="tag in node.tags"
            :key="tag"
            class="text-xs px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300"
          >{{ tag }}</span>
        </div>
        <div class="mt-3 text-xs text-gray-400">
          HTTP :{{ node.http_port }} &nbsp;·&nbsp; VLESS :{{ node.vless_port }}
        </div>
      </div>
    </div>

    <!-- Product selection modal overlay -->
    <div
      v-if="selectedNode"
      class="fixed inset-0 bg-black/50 z-50 flex items-end sm:items-center justify-center p-4"
      @click.self="selectedNode = null"
    >
      <div class="bg-white dark:bg-gray-900 rounded-2xl w-full max-w-md p-6 shadow-xl">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
            {{ selectedNode.city || selectedNode.country }} — 选择套餐
          </h2>
          <button @click="selectedNode = null" class="text-gray-400 hover:text-gray-600 text-xl leading-none">&times;</button>
        </div>

        <div v-if="loadingProducts" class="text-center py-6 text-gray-400">加载中...</div>
        <div v-else class="space-y-3">
          <div
            v-for="product in products"
            :key="product.id"
            class="border rounded-xl p-4 cursor-pointer transition-colors"
            :class="selectedProduct?.id === product.id
              ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/30'
              : 'border-gray-200 dark:border-gray-700 hover:border-blue-300'"
            @click="selectedProduct = product"
          >
            <div class="flex items-center justify-between">
              <span class="font-medium text-gray-900 dark:text-white">{{ product.name }}</span>
              <span class="font-bold text-blue-600 dark:text-blue-400">¥{{ product.price }}</span>
            </div>
            <div class="text-sm text-gray-500 dark:text-gray-400 mt-1">
              {{ product.duration_days }} 天 &nbsp;·&nbsp;
              {{ product.traffic_limit_gb === 0 ? '不限流量' : product.traffic_limit_gb + ' GB' }}
            </div>
            <p v-if="product.description" class="text-xs text-gray-400 mt-1">{{ product.description }}</p>
          </div>
        </div>

        <!-- Payment type -->
        <div v-if="selectedProduct" class="mt-4">
          <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">支付方式</label>
          <div class="flex gap-2 flex-wrap">
            <button
              v-for="pt in payTypes"
              :key="pt.value"
              class="px-4 py-2 rounded-lg border text-sm transition-colors"
              :class="payType === pt.value
                ? 'border-blue-500 bg-blue-50 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300'
                : 'border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:border-blue-300'"
              @click="payType = pt.value"
            >{{ pt.label }}</button>
          </div>
        </div>

        <button
          v-if="selectedProduct"
          class="mt-5 w-full py-3 rounded-xl bg-blue-600 hover:bg-blue-700 text-white font-semibold transition-colors disabled:opacity-60"
          :disabled="submitting"
          @click="purchase"
        >
          {{ submitting ? '处理中...' : `确认购买 ¥${selectedProduct.price}` }}
        </button>
      </div>
    </div>
  </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import { listNodes, listProducts, createRental } from '@/api/proxy'
import type { ProxyNode, ProxyProduct } from '@/api/proxy'

const router = useRouter()

const nodes = ref<ProxyNode[]>([])
const products = ref<ProxyProduct[]>([])
const loading = ref(true)
const loadingProducts = ref(false)
const selectedNode = ref<ProxyNode | null>(null)
const selectedProduct = ref<ProxyProduct | null>(null)
const payType = ref('alipay')
const submitting = ref(false)
const filterCountry = ref('')

const payTypes = [
  { value: 'alipay', label: '支付宝' },
  { value: 'wxpay', label: '微信支付' },
  { value: 'stripe', label: 'Stripe' }
]

const countries = computed(() => {
  const seen = new Set<string>()
  const result: { code: string; name: string }[] = []
  for (const n of nodes.value) {
    if (n.country_code && !seen.has(n.country_code)) {
      seen.add(n.country_code)
      result.push({ code: n.country_code, name: n.country || n.country_code })
    }
  }
  return result
})

const filteredNodes = computed(() => {
  if (!filterCountry.value) return nodes.value
  return nodes.value.filter(n => n.country_code === filterCountry.value)
})

function countryFlag(code: string): string {
  if (!code || code.length !== 2) return '🌐'
  return code.toUpperCase().replace(/./g, c =>
    String.fromCodePoint(c.charCodeAt(0) + 0x1F1A5)
  )
}

async function selectNode(node: ProxyNode) {
  selectedNode.value = node
  selectedProduct.value = null
  if (products.value.length === 0) {
    loadingProducts.value = true
    try {
      products.value = await listProducts()
    } finally {
      loadingProducts.value = false
    }
  }
}

async function purchase() {
  if (!selectedNode.value || !selectedProduct.value) return
  submitting.value = true
  try {
    const result = await createRental(selectedNode.value.id, selectedProduct.value.id, payType.value)
    // Navigate to payment result page
    if (result.pay_url) {
      window.location.href = result.pay_url
    } else {
      router.push({ name: 'PaymentResult', query: { order_id: result.order_id } })
    }
  } catch (e: unknown) {
    const msg = e instanceof Error ? e.message : '购买失败，请重试'
    alert(msg)
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  try {
    nodes.value = await listNodes()
  } finally {
    loading.value = false
  }
})
</script>
