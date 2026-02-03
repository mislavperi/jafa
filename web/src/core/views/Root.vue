<script setup lang="ts">
import Button from 'primevue/button';
import Avatar from 'primevue/avatar';
import OverlayBadge from 'primevue/overlaybadge';

</script>

<template>
    <div class="flex border-b-1">
        <Button icon="pi pi-moon" />
        <OverlayBadge value="4" severity="danger" class="inline-flex">
            <Avatar label="A" size="large" />
        </OverlayBadge>
    </div>
    <div>
        <slot></slot>
    </div>
</template>