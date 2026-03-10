<script setup lang="ts">
import { ref } from 'vue'
import { Form, FormField } from '@primevue/forms'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import MultiSelect from 'primevue/multiselect'
import Message from 'primevue/message'
import Button from 'primevue/button'
import { useCreateExpense } from '../composables/useExpenses'
import { useAllTags, useCreateTag, useAddTagToExpense } from '../composables/useTags'
import type { Tag } from '../models/expense'
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

function resolver({ values }: { values: Record<string, any> }) {
  const errors: Record<string, { message: string }[]> = {}
  if (!values.name?.trim()) errors.name = [{ message: 'Name is required' }]
  if (values.amount == null) errors.amount = [{ message: 'Amount is required' }]
  if (values.cost == null) errors.cost = [{ message: 'Cost is required' }]
  return { errors }
}

async function onSubmit(event: FormSubmitEvent) {
  if (!event.valid) return
  submitError.value = null
  const { name, amount, cost, tags } = event.values as {
    name: string
    amount: number
    cost: number
    tags: Tag[]
  }
  try {
    const expense = await createExpense({ name, amount, cost })
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
            :options="allTags ?? []"
            option-label="name"
            placeholder="Select tags..."
            display="chip"
            class="w-full"
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
