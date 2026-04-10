<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white">{{ t('payment.admin.channels') }}</h2>
          <button @click="openCreateDialog" class="btn btn-primary btn-sm">{{ t('common.create') }}</button>
        </div>
      </template>

      <template #default>
        <div v-if="loading" class="flex items-center justify-center py-12">
          <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary-500 border-t-transparent" role="status" aria-label="Loading" />
        </div>

        <div v-else-if="channels.length === 0" class="py-12 text-center text-gray-500 dark:text-slate-400">
          {{ t('payment.admin.noChannels') }}
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-sm">
            <thead>
              <tr class="border-b border-gray-200 dark:border-slate-700">
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">ID</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.channelName') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.platform') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.channel.rate') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.groupId') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('payment.admin.enabled') }}</th>
                <th scope="col" class="px-4 py-3 text-left font-medium text-gray-500 dark:text-slate-400">{{ t('common.actions') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="ch in channels" :key="ch.id" class="border-b border-gray-100 dark:border-slate-800">
                <td class="px-4 py-3 text-gray-900 dark:text-slate-100">{{ ch.id }}</td>
                <td class="px-4 py-3 font-medium text-gray-900 dark:text-slate-100">{{ ch.name }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ ch.platform }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ ch.rate_multiplier }}</td>
                <td class="px-4 py-3 text-gray-600 dark:text-slate-400">{{ ch.group_id || '-' }}</td>
                <td class="px-4 py-3">
                  <span :class="ch.enabled ? 'text-green-600 dark:text-green-400' : 'text-gray-400'" role="img" :aria-label="ch.enabled ? t('payment.admin.enabled') : t('payment.admin.disabled')">
                    {{ ch.enabled ? '✓' : '✗' }}
                  </span>
                </td>
                <td class="px-4 py-3">
                  <div class="flex gap-2">
                    <button @click="openEditDialog(ch)" :aria-label="t('common.edit') + ' ' + ch.name" class="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400">
                      {{ t('common.edit') }}
                    </button>
                    <button @click="handleDelete(ch.id)" :aria-label="t('common.delete') + ' ' + ch.name" class="text-sm text-red-600 hover:text-red-700 dark:text-red-400">
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
    <BaseDialog :show="dialogOpen" @close="dialogOpen = false" :title="editingChannel ? t('common.edit') : t('common.create')">
      <div class="space-y-4">
        <div>
          <label for="ch-name" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.channelName') }}</label>
          <input id="ch-name" v-model="form.name" class="input w-full" maxlength="100" />
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="ch-platform" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.platform') }}</label>
            <input id="ch-platform" v-model="form.platform" class="input w-full" placeholder="claude" maxlength="50" />
          </div>
          <div>
            <label for="ch-rate" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.channel.rate') }}</label>
            <input id="ch-rate" v-model="form.rate_multiplier" type="number" step="0.0001" min="0.0001" class="input w-full" placeholder="1.0" />
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label for="ch-group" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.groupId') }}</label>
            <input id="ch-group" v-model.number="form.group_id" type="number" class="input w-full" />
          </div>
          <div>
            <label for="ch-sort" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.sortOrder') }}</label>
            <input id="ch-sort" v-model.number="form.sort_order" type="number" class="input w-full" />
          </div>
        </div>
        <div>
          <label for="ch-desc" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.description') }}</label>
          <textarea id="ch-desc" v-model="form.description" rows="2" class="input w-full" maxlength="500" />
        </div>
        <div>
          <label for="ch-models" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.admin.models') }}</label>
          <input id="ch-models" v-model="form.models" class="input w-full" placeholder="claude-opus-4-6,claude-sonnet-4-6" maxlength="1000" />
        </div>
        <div>
          <label for="ch-features" class="mb-1 block text-sm font-medium text-gray-700 dark:text-slate-300">{{ t('payment.channel.features') }}</label>
          <input id="ch-features" v-model="form.features" class="input w-full" maxlength="500" />
        </div>
        <label for="ch-enabled" class="flex items-center gap-2 text-sm">
          <input id="ch-enabled" type="checkbox" v-model="form.enabled" class="rounded" />
          {{ t('payment.admin.enabled') }}
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
import type { PaymentChannel } from '@/types'

const { t } = useI18n()
const appStore = useAppStore()

const loading = ref(true)
const channels = ref<PaymentChannel[]>([])
const dialogOpen = ref(false)
const editingChannel = ref<PaymentChannel | null>(null)
const formLoading = ref(false)

interface ChannelForm {
  name: string
  platform: string
  rate_multiplier: string
  group_id: number | undefined
  sort_order: number
  description: string
  models: string
  features: string
  enabled: boolean
}

const form = reactive<ChannelForm>({
  name: '',
  platform: 'claude',
  rate_multiplier: '1.0',
  group_id: undefined,
  sort_order: 0,
  description: '',
  models: '',
  features: '',
  enabled: true
})

async function loadChannels() {
  loading.value = true
  try {
    channels.value = await adminPayAPI.listChannels()
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  } finally {
    loading.value = false
  }
}

onMounted(loadChannels)

function resetForm() {
  form.name = ''
  form.platform = 'claude'
  form.rate_multiplier = '1.0'
  form.group_id = undefined
  form.sort_order = 0
  form.description = ''
  form.models = ''
  form.features = ''
  form.enabled = true
}

function openCreateDialog() {
  editingChannel.value = null
  resetForm()
  dialogOpen.value = true
}

function openEditDialog(ch: PaymentChannel) {
  editingChannel.value = ch
  form.name = ch.name
  form.platform = ch.platform
  form.rate_multiplier = ch.rate_multiplier
  form.group_id = ch.group_id
  form.sort_order = ch.sort_order
  form.description = ch.description || ''
  form.models = ch.models || ''
  form.features = ch.features || ''
  form.enabled = ch.enabled
  dialogOpen.value = true
}

async function handleSave() {
  // Trim string fields
  form.name = form.name.trim()
  form.platform = form.platform.trim()
  form.description = form.description.trim()
  form.models = form.models.trim()
  form.features = form.features.trim()

  if (!form.name) {
    appStore.showError(t('common.nameRequired'))
    return
  }
  const rate = parseFloat(form.rate_multiplier)
  if (isNaN(rate) || rate <= 0) {
    appStore.showError(t('payment.channel.rate') + ' > 0')
    return
  }
  if (form.group_id !== undefined && form.group_id !== null && form.group_id < 0) {
    form.group_id = undefined
  }

  formLoading.value = true
  try {
    const payload: Omit<PaymentChannel, 'id' | 'created_at' | 'updated_at'> = {
      name: form.name,
      platform: form.platform,
      rate_multiplier: form.rate_multiplier,
      group_id: form.group_id,
      sort_order: form.sort_order,
      description: form.description || undefined,
      models: form.models || undefined,
      features: form.features || undefined,
      enabled: form.enabled
    }
    if (editingChannel.value) {
      await adminPayAPI.updateChannel(editingChannel.value.id, payload)
    } else {
      await adminPayAPI.createChannel(payload)
    }
    dialogOpen.value = false
    loadChannels()
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
    await adminPayAPI.deleteChannel(deletingId.value)
    loadChannels()
    appStore.showSuccess(t('common.deleted'))
  } catch (err: unknown) {
    appStore.showError(err instanceof Error ? err.message : t('common.error'))
  }
}
</script>
