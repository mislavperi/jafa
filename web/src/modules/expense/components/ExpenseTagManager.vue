<script setup lang="ts">
import { ref, computed, toRef } from 'vue'
import { Form, FormField } from '@primevue/forms'
import type { FormInstance, FormSubmitEvent } from '@primevue/forms'
import Chip from 'primevue/chip'
import Select from 'primevue/select'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import { useExpenseTags, useAddTagToExpense, useRemoveTagFromExpense } from '../composables/useTags'
import { useAllTags, useCreateTag } from '../composables/useTags'

const props = defineProps<{ expenseId: number }>()

const expenseIdRef = toRef(props, 'expenseId')

const { data: expenseTags, isLoading } = useExpenseTags(expenseIdRef)
const { data: allTags } = useAllTags()
const { mutate: addTag, isPending: isAdding } = useAddTagToExpense()
const { mutate: removeTag } = useRemoveTagFromExpense()
const { mutateAsync: createTag, isPending: isCreating } = useCreateTag()

const selectedTag = ref<{ id: number; name: string } | null>(null)
const showNewTag = ref(false)
const newTagFormRef = ref<FormInstance>()

const availableTags = computed(() => {
  if (!allTags.value || !expenseTags.value) return allTags.value ?? []
  const assigned = new Set(expenseTags.value.map((t) => t.id))
  return allTags.value.filter((t) => !assigned.has(t.id))
})

function handleAdd() {
  if (!selectedTag.value) return
  addTag({ expenseId: props.expenseId, tagId: selectedTag.value.id })
  selectedTag.value = null
}

async function handleCreate(event: FormSubmitEvent) {
  if (!event.valid) return
  const tag = await createTag(event.values.tagName)
  addTag({ expenseId: props.expenseId, tagId: tag.id })
  newTagFormRef.value?.reset()
  showNewTag.value = false
}

function handleRemove(tagId: number) {
  removeTag({ expenseId: props.expenseId, tagId })
}
</script>

<template>
  <div class="flex flex-col gap-3 py-2 px-1">
    <div class="flex flex-wrap gap-1.5 min-h-[2rem] items-center">
      <span v-if="isLoading" class="text-sm text-surface-400">Loading...</span>
      <span v-else-if="!expenseTags?.length" class="text-sm text-surface-400">No tags</span>
      <Chip
        v-for="tag in expenseTags"
        :key="tag.id"
        :label="tag.name"
        removable
        @remove="handleRemove(tag.id)"
        class="text-xs"
      />
    </div>

    <div class="flex flex-wrap gap-2 items-center">
      <template v-if="!showNewTag">
        <Select
          v-model="selectedTag"
          :options="availableTags"
          option-label="name"
          placeholder="Pick a tag..."
          size="small"
          class="w-40"
        />
        <Button
          label="Add"
          size="small"
          :disabled="!selectedTag"
          :loading="isAdding"
          @click="handleAdd"
        />
        <Button
          label="New tag"
          size="small"
          severity="secondary"
          text
          @click="showNewTag = true"
        />
      </template>

      <template v-else>
        <Form
          ref="newTagFormRef"
          :resolver="({ values }) => ({ errors: !values.tagName?.trim() ? { tagName: [{ message: 'Required' }] } : {} })"
          @submit="handleCreate"
          class="flex gap-2 items-start"
        >
          <FormField v-slot="$field" name="tagName" initialValue="">
            <InputText
              v-bind="$field.props"
              placeholder="Tag name..."
              size="small"
              class="w-36"
            />
          </FormField>
          <Button
            type="submit"
            label="Create & add"
            size="small"
            :loading="isCreating"
          />
          <Button
            label="Cancel"
            size="small"
            severity="secondary"
            text
            @click="showNewTag = false"
          />
        </Form>
      </template>
    </div>
  </div>
</template>
