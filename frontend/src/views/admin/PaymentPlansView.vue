<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('payment.admin.subscriptionPlans') }}</h2>
          <button @click="openCreateDialog" class="btn btn-primary btn-sm">{{ t('common.create') }}</button>
        </div>
      </template>

      <template #default>
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
        </div>

        <div v-else-if="plans.length === 0" class="py-12 text-center text-gray-500 dark:text-slate-400">
          {{ t('payment.admin.noPlans') }}
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-slate-700">
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ID</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.planName') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.price') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.validity') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.groupId') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.forSale') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="plan in plans" :key="plan.id" class="border-b border-gray-100 dark:border-slate-800">
                <td class="px-4 py-3 text-gray-900 dark:text-slate-100">{{ plan.id }}</td>
                <td class="px-4 py-3 font-medium text-gray-900 dark:text-slate-100">{{ plan.name }}</td>
                <td class="px-4 py-3 text-gray-900 dark:text-slate-100">
                  <span v-if="plan.original_price" class="mr-1 text-xs line-through text-gray-400">¥{{ plan.original_price }}</span>
                  ¥{{ plan.price }}
                </td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ plan.validity_days }} {{ plan.validity_unit }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ plan.group_id || '-' }}</td>
                <td class="px-4 py-3">
                  <span :class="plan.for_sale ? 'text-green-600 dark:text-green-400' : 'text-gray-400'" role="img" :aria-label="plan.for_sale ? t('payment.admin.forSale') : t('payment.admin.notForSale')">
                    {{ plan.for_sale ? '✓' : '✗' }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <div class="flex gap-2">
                    <button @click="openEditDialog(plan)" :aria-label="t('common.edit') + ' ' + plan.name" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">
                      {{ t('common.edit') }}
                    </button>
                    <button @click="handleDelete(plan.id)" :aria-label="t('common.delete') + ' ' + plan.name" class="text-sm text-red-600 hover:text-red-700 dark:text-red-400">
                      {{ t('common.delete') }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>
    </TablePageLayout>

    <!-- Create/Edit Dialog -->
    <BaseDialog :show="dialogOpen" @close="dialogOpen = false" :title="editingPlan ? t('common.edit') : t('common.create')">
      <div class="space-y-4">
        <div>
          <label for="plan-name" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.planName') }}</label>
          <input id="plan-name" v-model="form.name" class="input w-full" maxlength="100" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="plan-price" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.price') }}</label>
            <input id="plan-price" v-model="form.price" type="number" step="0.01" min="0" class="input w-full" placeholder="99.00" />
          </div>
          <div>
            <label for="plan-orig-price" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.originalPrice') }}</label>
            <input id="plan-orig-price" v-model="form.original_price" type="number" step="0.01" min="0" class="input w-full" placeholder="129.00" />
          </div>
        </div>
        <div class="grid grid-cols-3 gap-4">
          <div>
            <label for="plan-validity" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.validityDays') }}</label>
            <input id="plan-validity" v-model.number="form.validity_days" type="number" min="1" class="input w-full" />
          </div>
          <div>
            <label for="plan-unit" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.validityUnit') }}</label>
            <select id="plan-unit" v-model="form.validity_unit" class="input w-full">
              <option value="day">Day</option>
              <option value="week">Week</option>
              <option value="month">Month</option>
            </select>
          </div>
          <div>
            <label for="plan-group" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.groupId') }}</label>
            <input id="plan-group" v-model.number="form.group_id" type="number" class="input w-full" />
          </div>
        </div>
        <div>
          <label for="plan-desc" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.description') }}</label>
          <textarea id="plan-desc" v-model="form.description" rows="2" class="input w-full" maxlength="500" />
        </div>
        <div>
          <label for="plan-features" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.channel.features') }}</label>
          <input id="plan-features" v-model="form.features" class="input w-full" placeholder="feature1,feature2" maxlength="500" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="plan-product" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.productName') }}</label>
            <input id="plan-product" v-model="form.product_name" class="input w-full" maxlength="255" />
          </div>
          <div>
            <label for="plan-sort" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.sortOrder') }}</label>
            <input id="plan-sort" v-model.number="form.sort_order" type="number" class="input w-full" />
          </div>
        </div>
        <label for="plan-for-sale" class="flex items-center gap-2 text-sm">
          <input id="plan-for-sale" type="checkbox" v-model="form.for_sale" class="rounded" />
          {{ t('payment.admin.forSale') }}
        </label>
        <div class="flex justify-end gap-3">
          <button @click="dialogOpen = false" class="btn btn-secondary">{{ t('common.cancel') }}</button>
          <button @click="handleSave" class="btn btn-primary" :disabled="formLoading">
            {{ formLoading ? t('payment.processing') : t('common.save') }}
          </button>
        </div>
      </div>
    </BaseDialog>

    <ConfirmDialog
      :show="showDeleteConfirm"
      :title="t('common.delete')"
      :message="t('common.confirmDelete')"
      :danger="true"
      @confirm="confirmDelete"
      @cancel="showDeleteConfirm = false"
    />
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import { adminPayAPI } from '@/api/admin/pay'
import { useAppStore } from '@/stores'
import type { PaymentSubscriptionPlan } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const plans = ref<PaymentSubscriptionPlan[]>([])
const dialogOpen = ref(false)
const editingPlan = ref<PaymentSubscriptionPlan | null>(null)
const formLoading = ref(false)

interface PlanForm {
  name: string
  price: string
  original_price: string
  validity_days: number
  validity_unit: string
  group_id: number | undefined
  description: string
  features: string
  product_name: string
  sort_order: number
  for_sale: boolean
}

const form = reactive<PlanForm>({
  name: '',
  price: '',
  original_price: '',
  validity_days: 30,
  validity_unit: 'day',
  group_id: undefined,
  description: '',
  features: '',
  product_name: '',
  sort_order: 0,
  for_sale: true
})

async function loadPlans() {
  loading.value = true
  try {
    plans.value = await adminPayAPI.listPlans()
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    loading.value = false
  }
}

onMounted(loadPlans)

function resetForm() {
  form.name = ''
  form.price = ''
  form.original_price = ''
  form.validity_days = 30
  form.validity_unit = 'day'
  form.group_id = undefined
  form.description = ''
  form.features = ''
  form.product_name = ''
  form.sort_order = 0
  form.for_sale = true
}

function openCreateDialog() {
  editingPlan.value = null
  resetForm()
  dialogOpen.value = true
}

function openEditDialog(plan: PaymentSubscriptionPlan) {
  editingPlan.value = plan
  form.name = plan.name
  form.price = plan.price
  form.original_price = plan.original_price || ''
  form.validity_days = plan.validity_days
  form.validity_unit = plan.validity_unit
  form.group_id = plan.group_id
  form.description = plan.description || ''
  form.features = plan.features || ''
  form.product_name = plan.product_name || ''
  form.sort_order = plan.sort_order
  form.for_sale = plan.for_sale
  dialogOpen.value = true
}

async function handleSave() {
  // Trim string fields
  form.name = form.name.trim()
  form.description = form.description.trim()
  form.features = form.features.trim()
  form.product_name = form.product_name.trim()

  if (!form.name) {
    appStore.showError(t('common.nameRequired'))
    return
  }
  if (!form.price || parseFloat(form.price) <= 0) {
    appStore.showError(t('payment.admin.price') + ' > 0')
    return
  }

  formLoading.value = true
  try {
    const payload: Omit<PaymentSubscriptionPlan, 'id' | 'created_at' | 'updated_at'> = {
      name: form.name,
      price: form.price,
      original_price: form.original_price || undefined,
      validity_days: form.validity_days,
      validity_unit: form.validity_unit,
      group_id: form.group_id,
      description: form.description || undefined,
      features: form.features || undefined,
      product_name: form.product_name || undefined,
      sort_order: form.sort_order,
      for_sale: form.for_sale
    }
    if (editingPlan.value) {
      await adminPayAPI.updatePlan(editingPlan.value.id, payload)
    } else {
      await adminPayAPI.createPlan(payload)
    }
    dialogOpen.value = false
    loadPlans()
    appStore.showSuccess(t('common.saved'))
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    formLoading.value = false
  }
}

const showDeleteConfirm = ref(false)
const deletingId = ref<number | null>(null)

function handleDelete(id: number) {
  deletingId.value = id
  showDeleteConfirm.value = true
}

async function confirmDelete() {
  showDeleteConfirm.value = false
  if (!deletingId.value) return
  try {
    await adminPayAPI.deletePlan(deletingId.value)
    loadPlans()
    appStore.showSuccess(t('common.deleted'))
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  }
}
</script>
