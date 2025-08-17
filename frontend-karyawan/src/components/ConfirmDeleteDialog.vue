<template>
  <v-dialog
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    max-width="400"
  >
    <v-card>
      <v-card-title class="text-h6">Confirm Delete</v-card-title>
      <v-card-text>
        <div class="text-body-2">
          Are you sure you want to delete <strong>{{ itemName }}</strong
          >?
        </div>
      </v-card-text>
      <v-card-actions>
        <v-spacer />
        <v-btn text @click="$emit('cancel')">Cancel</v-btn>
        <v-btn color="red" @click="$emit('confirm')">Delete</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
const props = defineProps({
  modelValue: Boolean,
  itemName: {
    type: String,
    default: "this item",
  },
});

defineEmits(["confirm", "cancel", "update:modelValue"]);
</script>
