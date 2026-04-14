<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">代理套餐</h2>
          <button @click="openCreate" class="btn btn-primary btn-sm">添加套餐</button>
        </div>
      </template>

      <template #default>
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" />
        </div>

        <div v-else-if="products.length === 0" class="py-12 text-center text-gray-500 dark:text-slate-400">暂无套餐</div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-slate-700">
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ID</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">套餐名称</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">价格</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">时长</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">流量</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">排序</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">上架</th>
                <th class="px-3 py-3 text-left font-medium text-gray-500 dark:text-slate-400">操作</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="product in products" :key="product.id" class="border-b border-gray-100 dark:border-slate-800">
                <td class="px-3 py-3 text-gray-900 dark:text-slate-100">{{ product.id }}</td>
                <td class="px-3 py-3 font-medium text-gray-900 dark:text-slate-100">{{ product.name }}</td>
                <td class="px-3 py-3 text-gray-900 dark:text-slate-100">¥{{ product.price }}</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ product.duration_days }} 天</td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">
                  {{ product.traffic_limit_gb === 0 ? '不限' : `${product.traffic_limit_gb} GB` }}
                </td>
                <td class="px-3 py-3 text-gray-600 dark:text-slate-400">{{ product.sort_order }}</td>
                <td class="px-3 py-3">
                  <span :class="product.is_active ? 'text-green-600 dark:text-green-400' : 'text-gray-400'">
                    {{ product.is_active ? '✓' : '✗' }}
                  </span>
                </td>
                <td class="px-3 py-3">
                  <div class="flex gap-3">
                    <button @click="openEdit(product)" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">编辑</button>
                    <button @click="handleDelete(product.id)" class="text-sm text-red-600 hover:text-red-700 dark:text-red-400">删除</button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>
    </TablePageLayout>

    <!-- Create/Edit Dialog -->
    <BaseDialog :show="dialogOpen" @close="closeDialog" :title="editingProduct ? '编辑套餐' : '添加套餐'">
      <div class="space-y-4">
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">套餐名称 *</label>
          <input v-model="form.name" class="input w-full" placeholder="月度基础版" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">价格（元） *</label>
            <input v-model="form.price" class="input w-full" placeholder="29.00" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">时长（天） *</label>
            <input v-model.number="form.duration_days" type="number" min="1" class="input w-full" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">流量限制（GB，0=不限）</label>
            <input v-model.number="form.traffic_limit_gb" type="number" min="0" class="input w-full" />
          </div>
          <div>
            <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">排序</label>
            <input v-model.number="form.sort_order" type="number" class="input w-full" />
          </div>
        </div>
        <div>
          <label class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">描述</label>
          <input v-model="form.description" class="input w-full" />
        </div>
        <div class="flex items-center gap-2">
          <input id="is-active" v-model="form.is_active" type="checkbox" class="rounded" />
          <label for="is-active" class="text-sm text-gray-700 dark:text-slate-300">上架销售</label>
        </div>
      </div>
      <template #footer>
        <button @click="closeDialog" class="btn btn-secondary">取消</button>
        <button @click="save" :disabled="saving" class="btn btn-primary">{{ saving ? '保存中...' : '保存' }}</button>
      </template>
    </BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import { proxyAdminAPI } from '@/api/admin/proxy'
import type { AdminProxyProduct } from '@/api/admin/proxy'

const products = ref<AdminProxyProduct[]>([])
const loading = ref(true)
const dialogOpen = ref(false)
const editingProduct = ref<AdminProxyProduct | null>(null)
const saving = ref(false)

const defaultForm = (): Omit<AdminProxyProduct, 'id'> => ({
  name: '',
  description: '',
  duration_days: 30,
  traffic_limit_gb: 0,
  price: '',
  sort_order: 0,
  is_active: true,
})

const form = ref(defaultForm())

async function reload() {
  loading.value = true
  try {
    products.value = await proxyAdminAPI.listProducts()
  } finally {
    loading.value = false
  }
}

function openCreate() {
  editingProduct.value = null
  form.value = defaultForm()
  dialogOpen.value = true
}

function openEdit(product: AdminProxyProduct) {
  editingProduct.value = product
  form.value = { ...product }
  dialogOpen.value = true
}

function closeDialog() {
  dialogOpen.value = false
}

async function save() {
  saving.value = true
  try {
    if (editingProduct.value) {
      await proxyAdminAPI.updateProduct(editingProduct.value.id, form.value)
    } else {
      await proxyAdminAPI.createProduct(form.value)
    }
    closeDialog()
    await reload()
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '保存失败')
  } finally {
    saving.value = false
  }
}

async function handleDelete(id: number) {
  if (!confirm('确认删除此套餐？')) return
  try {
    await proxyAdminAPI.deleteProduct(id)
    await reload()
  } catch (e: unknown) {
    alert(e instanceof Error ? e.message : '删除失败')
  }
}

onMounted(reload)
</script>
