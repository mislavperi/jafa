<script setup lang="ts">
import { ref, computed } from 'vue'
import { Form, FormField } from '@primevue/forms'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import MultiSelect from 'primevue/multiselect'
import Message from 'primevue/message'
import Button from 'primevue/button'
import Select from 'primevue/select'
import Checkbox from 'primevue/checkbox'
import DatePicker from 'primevue/datepicker'
import { useCreateExpense } from '../composables/useExpenses'
import { useAllTags, useCreateTag, useAddTagToExpense } from '../composables/useTags'
import type { Tag, RecurringSchedule, RecurrenceInterval } from '../models/expense'
import type { FormInstance, FormSubmitEvent } from '@primevue/forms'

const props = defineProps<{ visible: boolean }>()
const emit = defineEmits<{ 'update:visible': [value: boolean] }>()

const formRef = ref<FormInstance>()
const newTagName = ref('')
const submitError = ref<string | null>(null)

const { mutateAsync: createExpense, isPending } = useCreateExpense()
const { data: allTags } = useAllTags()
const { mutateAsync: createTag } = useCreateTag()
const { mutateAsync: addTag } = useAddTagToExpense()

const tagsWithSelectAll = computed(() => {
  const tags = allTags.value ?? []
  return [
    { id: -1, name: 'Select All', isSelectAll: true },
    ...tags,
  ]
})

function handleTagsChange(event: { value: Tag[] }) {
  const selected = event.value
  if (selected.some((t) => t.id === -1)) {
    formRef.value?.setFieldValue('tags', [...tagsWithSelectAll.value.slice(1)])
  } else {
    formRef.value?.setFieldValue('tags', selected)
  }
}

interface FormValues {
  name?: string
  amount?: number
  cost?: number
  isRecurring?: boolean
  recurrenceInterval?: RecurrenceInterval
  recurrenceDay?: number
  [key: string]: unknown
}

function resolver({ values }: { values: FormValues }) {
  const errors: Record<string, { message: string }[]> = {}
  if (!values.name?.trim()) errors.name = [{ message: 'Name is required' }]
  if (values.amount == null) errors.amount = [{ message: 'Amount is required' }]
  if (values.cost == null) errors.cost = [{ message: 'Cost is required' }]
  if (values.isRecurring && !values.recurrenceInterval) {
    errors.recurrenceInterval = [{ message: 'Recurrence interval is required' }]
  }
  if (values.isRecurring && !values.recurrenceDay) {
    errors.recurrenceDay = [{ message: 'Day is required' }]
  }
  return { errors }
}

async function onSubmit(event: FormSubmitEvent) {
  if (!event.valid) return
  submitError.value = null
  const { name, amount, cost, tags, isRecurring, recurrenceInterval, recurrenceDay, recurrenceStartDate } = event.values as {
    name: string
    amount: number
    cost: number
    tags: Tag[]
    isRecurring: boolean
    recurrenceInterval: RecurrenceInterval | null
    recurrenceDay: number | null
    recurrenceStartDate: Date | null
  }
  try {
    let recurringSchedule: RecurringSchedule | undefined
    if (isRecurring && recurrenceInterval && recurrenceDay) {
      const dateStr = recurrenceStartDate
        ? recurrenceStartDate.toISOString().split('T')[0]
        : new Date().toISOString().split('T')[0]
      recurringSchedule = {
        interval: recurrenceInterval,
        dayOfMonth: recurrenceDay,
        startDate: dateStr ?? '',
      }
    }
    const expense = await createExpense({ name, amount, cost, recurringSchedule })
    await Promise.all((tags ?? []).map((tag) => addTag({ expenseId: expense.id, tagId: tag.id })))
    formRef.value?.reset()
    newTagName.value = ''
    emit('update:visible', false)
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Failed to save expense'
  }
}

async function handleCreateTag() {
  const name = newTagName.value.trim()
  if (!name) return
  const tag = await createTag(name)
  const current: Tag[] = formRef.value?.getFieldState('tags')?.value ?? []
  formRef.value?.setFieldValue('tags', [...current, tag])
  newTagName.value = ''
}

const intervalOptions = [
  { label: 'Monthly', value: 'monthly' },
  { label: 'Yearly', value: 'yearly' },
]

const dayOptions = Array.from({ length: 28 }, (_, i) => ({
  label: String(i + 1),
  value: i + 1,
}))
</script>

<template>
  <Dialog
    :visible="props.visible"
    @update:visible="emit('update:visible', $event)"
    header="Add Expense"
    modal
    :style="{ width: '26rem' }"
  >
    <Form ref="formRef" :resolver="resolver" @submit="onSubmit" class="flex flex-col gap-4 pt-1">
      <FormField v-slot="$field" name="name" initialValue="">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Name</label>
          <InputText v-bind="$field.props" placeholder="e.g. Coffee" class="w-full" />
          <Message v-if="$field.invalid" severity="error" size="small" variant="simple">
            {{ $field.error?.message }}
          </Message>
        </div>
      </FormField>

      <FormField v-slot="$field" name="amount" :initialValue="null">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Amount</label>
          <InputNumber
            v-bind="$field.props"
            placeholder="e.g. 2"
            :min="0"
            :min-fraction-digits="0"
            :max-fraction-digits="3"
            class="w-full"
          />
          <Message v-if="$field.invalid" severity="error" size="small" variant="simple">
            {{ $field.error?.message }}
          </Message>
        </div>
      </FormField>

      <FormField v-slot="$field" name="cost" :initialValue="null">
        <div class="flex flex-col gap-1">
          <label class="text-sm font-medium">Cost ($)</label>
          <InputNumber
            v-bind="$field.props"
            placeholder="e.g. 4.50"
            :min="0"
            :min-fraction-digits="2"
            :max-fraction-digits="3"
            prefix="$"
            class="w-full"
          />
          <Message v-if="$field.invalid" severity="error" size="small" variant="simple">
            {{ $field.error?.message }}
          </Message>
        </div>
      </FormField>

      <FormField v-slot="$field" name="tags" :initialValue="[]">
        <div class="flex flex-col gap-2">
          <label class="text-sm font-medium">Tags</label>
          <MultiSelect
            v-bind="$field.props"
            :options="tagsWithSelectAll"
            option-label="name"
            placeholder="Select tags..."
            display="chip"
            class="w-full"
            @change="handleTagsChange"
          />
          <div class="flex gap-2">
            <InputText
              v-model="newTagName"
              placeholder="New tag name..."
              size="small"
              class="flex-1"
              @keyup.enter="handleCreateTag"
            />
            <Button
              label="Create"
              size="small"
              severity="secondary"
              :disabled="!newTagName.trim()"
              @click="handleCreateTag"
            />
          </div>
        </div>
      </FormField>

      <FormField v-slot="$field" name="isRecurring" :initialValue="false">
        <div class="flex items-center gap-2">
          <Checkbox v-bind="$field.props" :binary="true" input-id="isRecurring" />
          <label for="isRecurring" class="text-sm font-medium cursor-pointer">Recurring expense</label>
        </div>
      </FormField>

      <template v-if="formRef?.getFieldState('isRecurring')?.value">
        <FormField v-slot="$field" name="recurrenceInterval" :initialValue="null">
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium">Frequency</label>
            <Select
              v-bind="$field.props"
              :options="intervalOptions"
              option-label="label"
              option-value="value"
              placeholder="Select frequency..."
              class="w-full"
            />
            <Message v-if="$field.invalid" severity="error" size="small" variant="simple">
              {{ $field.error?.message }}
            </Message>
          </div>
        </FormField>

        <FormField v-slot="$field" name="recurrenceDay" :initialValue="null">
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium">Day of {{ formRef?.getFieldState('recurrenceInterval')?.value === 'yearly' ? 'Month' : 'Month' }}</label>
            <Select
              v-bind="$field.props"
              :options="dayOptions"
              option-label="label"
              option-value="value"
              placeholder="Select day..."
              class="w-full"
            />
            <Message v-if="$field.invalid" severity="error" size="small" variant="simple">
              {{ $field.error?.message }}
            </Message>
          </div>
        </FormField>

        <FormField v-slot="$field" name="recurrenceStartDate" :initialValue="null">
          <div class="flex flex-col gap-1">
            <label class="text-sm font-medium">Start Date</label>
            <DatePicker
              v-bind="$field.props"
              date-format="yy-mm-dd"
              placeholder="Select start date..."
              class="w-full"
            />
          </div>
        </FormField>
      </template>

    </Form>

    <Message v-if="submitError" severity="error" class="mx-1 mb-2">{{ submitError }}</Message>

    <template #footer>
      <Button
        label="Cancel"
        severity="secondary"
        text
        @click="emit('update:visible', false)"
      />
      <Button label="Add" :loading="isPending" @click="formRef?.submit()" />
    </template>
  </Dialog>
</template>
