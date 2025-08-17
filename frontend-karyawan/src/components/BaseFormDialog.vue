<template>
  <v-dialog
    :model-value="modelValue"
    @update:model-value="(val) => emit('update:modelValue', val)"
    max-width="500px"
  >
    <v-card>
      <v-card-title class="headline">
        {{ title }}
      </v-card-title>

      <v-card-text>
        <v-form ref="formRef" v-model="formValid">
          <v-row>
            <v-col
              v-for="(field, index) in fields"
              :key="index"
              :cols="field.cols || 12"
            >
              <component
                :is="getComponent(field)"
                v-model="field.model.value"
                v-bind="getProps(field)"
                @change="handleFileChange(field)"
              />
            </v-col>
          </v-row>

          <v-img
            v-if="previewImage"
            :src="previewImage"
            max-height="150"
            class="mt-4 rounded border"
            cover
          />
        </v-form>
      </v-card-text>

      <v-card-actions>
        <v-spacer />
        <v-btn text @click="onCancel">Cancel</v-btn>
        <v-btn color="primary" @click="onSave">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup>
import { ref, watch } from "vue";
import TimePickerWrapper from "@/components/TimePickerWrapper.vue";

const props = defineProps({
  modelValue: Boolean,
  title: String,
  fields: Array,
  previewImage: String,
});

const emit = defineEmits(["update:modelValue", "save", "update:previewImage"]);

const formRef = ref(null);
const formValid = ref(false);

const getComponent = (field) => {
  switch (field.type) {
    case "number":
      return "v-text-field";
    case "select":
      return "v-select";
    case "file":
      return "v-file-input";
    case "time":
      return TimePickerWrapper;
    default:
      return "v-text-field";
  }
};

const getProps = (field) => {
  const props = {
    label: field.label,
    clearable: true,
    rules: field.rules || [],
    prependIcon: field.prependIcon,
    disabled: field.disabled || false,
  };

  if (field.type === "number") {
    props.type = "number";
  } else if (field.type === "file") {
    props.accept = field.accept || "image/*";
    props.showSize = field.showSize || false;
  } else if (field.type === "select") {
    props.items = Array.isArray(field.options) ? field.options : [];
    props.itemTitle = field.itemTitle || "title";
    props.itemValue = field.itemValue || "value";
    props.returnObject = field.returnObject || false;
  }

  return props;
};

const handleFileChange = (field) => {
  if (field.type !== "file") return;

  const file = Array.isArray(field.model.value)
    ? field.model.value[0]
    : field.model.value;

  if (!file) {
    emit("update:previewImage", null);
    return;
  }

  const reader = new FileReader();
  reader.onload = (e) => {
    emit("update:previewImage", e.target.result);
  };
  reader.readAsDataURL(file);
};

const getDefaultValue = (type) => {
  switch (type) {
    case "number":
      return 0;
    case "file":
      return null;
    case "time":
      return "";
    default:
      return "";
  }
};

const resetFormFields = () => {
  props.fields.forEach((field) => {
    if (field.model) {
      field.model.value = field.default ?? getDefaultValue(field.type);
    }
  });
  emit("update:previewImage", null);
};

const onCancel = () => {
  resetFormFields();
  emit("update:modelValue", false);
};

const onSave = () => {
  if (!formRef.value?.validate()) return;
  emit("save");
};

watch(
  () => props.modelValue,
  (val) => {
    if (!val) resetFormFields();
  }
);
</script>