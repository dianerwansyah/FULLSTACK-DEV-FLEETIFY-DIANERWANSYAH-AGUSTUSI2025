<template>
  <v-card class="mt-4">
    <v-data-table-server
      :headers="headers"
      :items="items"
      :items-length="totalItems"
      :loading="loading"
      :page="page"
      :items-per-page="itemsPerPage"
      :sort-by="sortBy"
      :sort-desc="sortDesc"
      :multi-sort="multiSort"
      class="elevation-1"
      :footer-props="{
        'items-per-page-options': [5, 10, 20, 25],
        'show-current-page': true,
        'show-first-last-page': true,
      }"
      @update:page="(val) => emit('update:page', val)"
      @update:items-per-page="(val) => emit('update:itemsPerPage', val)"
      @update:sort-by="(val) => emit('update:sortBy', val)"
      @update:sort-desc="(val) => emit('update:sortDesc', val)"
    >
      <!-- Slot forwarding -->
      <template v-for="(_, slotName) in $slots" #[slotName]="slotProps">
        <slot :name="slotName" v-bind="slotProps" />
      </template>
    </v-data-table-server>
  </v-card>
</template>

<script setup>
const props = defineProps({
  headers: {
    type: Array,
    required: true,
  },
  items: {
    type: Array,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  totalItems: {
    type: Number,
    default: 0,
  },
  page: {
    type: Number,
    default: 1,
  },
  itemsPerPage: {
    type: Number,
    default: 10,
  },
  sortBy: {
    type: Array,
    default: () => [],
  },
  sortDesc: {
    type: Array,
    default: () => [],
  },
  multiSort: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits([
  "update:page",
  "update:itemsPerPage",
  "update:sortBy",
  "update:sortDesc",
]);
</script>
