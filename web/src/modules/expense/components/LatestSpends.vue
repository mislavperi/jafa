<script setup lang="ts">
import { ref} from "vue"

import Panel from "primevue/panel"
import DataTable from "primevue/datatable"
import Column from "primevue/column"
import ColumnGroup from "primevue/columngroup"

type MockSpend = {
    name: string
    value: string
}

const mockSpends = ref<MockSpend[]>([
    {
    name: "Groceries",
    value: "70"
    }
])
</script>

<template>
<Panel header="I am a latest expense page, see what you spent on last">
    <DataTable :value="mockSpends" table-style="min-width: 50rem" >
        <Column field="name" header="Name"></Column>
        <Column field="value" header="Value"></Column>
    </DataTable>
</Panel>
</template>
