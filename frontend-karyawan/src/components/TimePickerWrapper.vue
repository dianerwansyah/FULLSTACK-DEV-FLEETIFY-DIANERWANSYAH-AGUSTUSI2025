<template>
  <v-text-field
    v-model="time"
    :label="label"
    type="time"
    clearable
    @click:clear="clearTime"
  />
</template>

<script setup>
import { ref, watch } from "vue";

const props = defineProps({
  modelValue: String,
  label: String,
});
const emit = defineEmits(["update:modelValue"]);

const time = ref(props.modelValue || "");

watch(() => props.modelValue, (val) => {
  time.value = val;
});

watch(time, (val) => {
  emit("update:modelValue", val);
});

const clearTime = () => {
  time.value = "";
};
</script>