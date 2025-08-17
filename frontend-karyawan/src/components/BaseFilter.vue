<template>
  <v-row class="mb-4">
    <v-col
      v-for="(field, index) in fields"
      :key="index"
      :cols="field.cols || 12"
      :md="field.md || 4"
    >
      <!-- Date Picker with Menu -->
      <template v-if="field.type === 'date'">
        <v-menu
          v-model="field.menu.value"
          :close-on-content-click="false"
          transition="scale-transition"
          offset-y
        >
          <template #activator="{ props }">
            <v-text-field
              v-bind="props"
              :label="field.label"
              :model-value="formatDate(field.model.value)"
              readonly
              clearable
              density="compact"
              hide-details
              @click:clear="field.model.value = null"
            />
          </template>
          <v-card>
            <v-sheet max-width="280px" class="mx-auto">
              <v-date-picker
                v-model="field.model.value"
                :min="field.min"
                :max="field.max"
                show-adjacent-months
              />
              <v-card-actions class="justify-end px-4 pb-2">
                <v-btn variant="text" @click="field.menu.value = false"
                  >OK</v-btn
                >
              </v-card-actions>
            </v-sheet>
          </v-card>
        </v-menu>
      </template>

      <!-- Other Field Types -->
      <template v-else>
        <component
          :is="getComponent(field)"
          v-model="field.model.value"
          v-bind="getProps(field)"
        />
      </template>
    </v-col>
  </v-row>
</template>

<script setup>
import { ref } from "vue";
import moment from "moment";

const props = defineProps({
  fields: {
    type: Array,
    required: true,
  },
});

/* ===== Helpers ===== */
const getComponent = (field) => {
  switch (field.type) {
    case "number":
      return "v-text-field";
    case "select":
      return "v-select";
    default:
      return "v-text-field";
  }
};

const getProps = (field) => {
  const baseProps = {
    label: field.label,
    clearable: true,
    density: "compact",
    hideDetails: true,
    type: field.type === "number" ? "number" : "text",
  };

  if (field.type === "select") {
    baseProps.items = field.items || [];
  }

  return baseProps;
};

const formatDate = (val) => {
  return val ? moment(val).format("DD-MM-YYYY") : "";
};
</script>
