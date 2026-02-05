<script setup lang="ts">
import Panel from "primevue/panel"
import DataTable from "primevue/datatable"
import Column from "primevue/column"

import { useExpenses } from "../composables/useExpenses"

const { data: expenses, isLoading, isError } = useExpenses()
</script>

<template>
    <div class="grid grid-cols-2 gap-4">
        <Panel header="I am a latest expense page, see what you spent on last">
            <DataTable :value="expenses" :loading="isLoading" paginator :rows="5">
                <Column field="name" header="Name"></Column>
                <Column field="amount" header="Amount"></Column>
            </DataTable>
        </Panel>
        <Panel header="I am a latest expense page, see what you spent on last">
            <p>Some other content</p>
        </Panel>
    </div>
</template>
