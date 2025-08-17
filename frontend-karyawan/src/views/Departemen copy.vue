<template>
  <v-container>
    <!-- Judul -->
    <h1 class="text-h5 font-weight-bold mb-4">Departemen</h1>

    <!-- Toggle Filter -->
    <v-btn
      variant="text"
      size="small"
      class="mb-2"
      @click="showFilter = !showFilter"
    >
      <v-icon left>{{ showFilter ? "mdi-eye-off" : "mdi-eye" }}</v-icon>
      {{ showFilter ? "Hide Filter" : "Show Filter" }}
    </v-btn>

    <!-- Filter -->
    <BaseFilter v-if="showFilter" :fields="filterFields" class="mb-6" />

    <!-- Add Button -->
    <v-row justify="end" class="mb-6">
      <v-col cols="12" md="auto">
        <v-btn color="primary" @click="openAddDialog">
          <v-icon left>mdi-plus</v-icon>
          Add Departemen
        </v-btn>
      </v-col>
    </v-row>
    <!-- Table -->
    <BaseTable
      :headers="headers"
      :items="materials"
      :loading="loading"
      :total-items="totalItems"
      :page="page"
      :items-per-page="itemsPerPage"
      :sort-by="sortBy"
      :sort-desc="sortDesc"
      @update:page="(val) => (page = val)"
      @update:itemsPerPage="(val) => (itemsPerPage = val)"
      @update:sortBy="setSortBy"
      @update:sortDesc="setSortDesc"
    >
      <template #item.image="{ item }">
        <v-img
          :src="item.image"
          max-height="80"
          max-width="80"
          cover
          class="rounded border"
        />
      </template>

      <template #item.actions="{ item }">
        <div class="d-flex align-center" style="gap: 8px">
          <v-btn icon size="small" color="blue" @click="editMaterial(item)">
            <v-icon>mdi-pencil</v-icon>
          </v-btn>
          <v-btn icon size="small" color="red" @click="deleteMaterial(item)">
            <v-icon>mdi-delete</v-icon>
          </v-btn>
          <v-btn icon size="small" color="grey" @click="viewAudit(item)">
            <v-icon>mdi-history</v-icon>
          </v-btn>
        </div>
      </template>
    </BaseTable>

    <!-- Dialog Form -->
    <BaseFormDialog
      v-model="showDialog"
      :title="editedMaterial ? 'Edit Material' : 'Add Material'"
      :fields="formFields"
      :previewImage="previewImage"
      @save="saveMaterial"
    />
  </v-container>

  <ConfirmDeleteDialog
    v-model="showDeleteDialog"
    :itemName="materialToDelete?.name"
    @confirm="confirmDelete"
    @cancel="showDeleteDialog = false"
  />

  <!-- Snackbar feedback -->
  <v-snackbar
    v-model="toast.show"
    :color="toast.color"
    timeout="3000"
    location="bottom right"
    elevation="4"
    style="margin-bottom: 20px"
  >
    {{ toast.message }}
  </v-snackbar>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import api from "@/axios";
import BaseTable from "@/components/BaseTable.vue";
import BaseFilter from "@/components/BaseFilter.vue";
import BaseFormDialog from "@/components/BaseFormDialog.vue";
import { usePaginatedTable } from "@/components/usePaginatedTable";
import ConfirmDeleteDialog from "@/components/ConfirmDeleteDialog.vue";
import { useToast } from "@/components/useToast";


const { toast, showToast } = useToast();
const showFilter = ref(false)

/* ===== Rules ===== */
const rules = {
  required: (v) => !!v || "This field is required",
  number: (v) => v === null || v === "" || !isNaN(v) || "Must be a number",
  maxSize: (value) => {
    if (!value) return true;
    const file = Array.isArray(value) ? value[0] : value;
    return file.size <= 5 * 1024 * 1024 || "File must be less than 5MB";
  },
};

/* ===== Schema ===== */
const materialSchema = [
  {
    key: "name",
    label: "Material Name",
    type: "text",
    default: "",
    rules: [rules.required],
  },
  {
    key: "category_id",
    label: "Category",
    type: "select",
    default: "",
    rules: [rules.required],
  },
  {
    key: "stock",
    label: "Stock",
    type: "number",
    default: 0,
    rules: [rules.required, rules.number],
  },
  {
    key: "price",
    label: "Price",
    type: "number",
    default: 0,
    rules: [rules.required, rules.number],
  },
  {
    key: "unit",
    label: "Unit",
    type: "text",
    default: "",
    rules: [rules.required],
  },
  {
    key: "description",
    label: "Description",
    type: "text",
    default: "",
    rules: [rules.required],
  },
  {
    key: "is_active",
    label: "Active",
    type: "text",
    default: "",
    rules: [rules.required],
  },
  {
    key: "image",
    label: "Material Image",
    type: "file",
    default: null,
    rules: [rules.required, rules.maxSize],
    accept: "image/*",
    showSize: true,
    prependIcon: "mdi-image",
  },
];

/* ===== Filters ===== */
const filters = {
  name: ref(""),
  stock: ref(null),
  price: ref(null),
  category: ref(""),
};

const categoryOptions = ref([]);
const fetchCategories = async () => {
  try {
    const res = await api.get("api/material-categories");
    categoryOptions.value = res.data.map((cat) => ({
      title: cat.name,
      value: cat.uuid,
    }));
  } catch (err) {
    console.error("Gagal ambil kategori:", err);
  }
};

const filterFields = computed(() => [
  { label: "Filter by Name", type: "text", model: filters.name },
  { label: "Stock", type: "number", model: filters.stock },
  { label: "Price", type: "number", model: filters.price },
  {
    label: "Category",
    type: "select",
    model: filters.category,
    items: categoryOptions.value.length
      ? categoryOptions.value
      : [{ title: "Loading...", value: null }],
    disabled: !categoryOptions.value.length,
  },
]);

/* ===== Table Data ===== */
const normalizeMaterial = (item) => {
  const category = categoryOptions.value.find(
    (cat) => cat.value === item.category_id
  );
  return {
    uuid: item.uuid,
    category_id: item.category_id,
    category_name: category ? category.title : "-",
    name: item.name || "Unnamed",
    stock: item.stock ?? 0,
    price: item.price ?? 0,
    unit: item.unit,
    description: item.description,
    is_active: item.is_active,
    image: item.image ?? null,
  };
};

const {
  items: materials,
  loading,
  totalItems,
  page,
  itemsPerPage,
  sortBy,
  sortDesc,
  fetchData,
} = usePaginatedTable({
  endpoint: "api/materials",
  normalizeFn: normalizeMaterial,
  filters,
});

const setSortBy = (val) => {
  sortBy.value = Array.isArray(val) ? val : [val];
};

const setSortDesc = (val) => {
  sortDesc.value = Array.isArray(val) ? val : [val];
};

const showDeleteDialog = ref(false);
const materialToDelete = ref(null);

const deleteMaterial = (item) => {
  materialToDelete.value = item;
  showDeleteDialog.value = true;
};

/* ===== Headers ===== */
const headers = [
  { title: "Material Name", key: "name" },
  { title: "Category", key: "category_name" },
  { title: "Stock", key: "stock" },
  { title: "Price", key: "price" },
  { title: "Unit", key: "unit" },
  { title: "Active", key: "is_active" },
  { title: "Actions", key: "actions", sortable: false, width: "140px" },
];

/* ===== Form State ===== */
const showDialog = ref(false);
const editedMaterial = ref(null);

const getDefaultMaterial = () =>
  materialSchema.reduce((acc, f) => ({ ...acc, [f.key]: f.default }), {});

const newMaterial = ref(getDefaultMaterial());

const formFields = materialSchema.map((f) => ({
  ...f,
  model: ref(f.default),
}));

const previewImage = computed(() => {
  const img = newMaterial.value.image;
  if (!img) return null;
  return typeof img === "string" ? img : URL.createObjectURL(img);
});

const fillFormFields = () => {
  formFields.forEach((field) => {
    field.model.value = newMaterial.value[field.key];
  });
};

const injectCategoryItems = () => {
  const categoryField = formFields.find((f) => f.key === "category_id");
  if (categoryField) {
    categoryField.items = categoryOptions.value;
    categoryField.disabled = !categoryOptions.value.length;
  }
};

/* ===== Actions ===== */
const openAddDialog = () => {
  editedMaterial.value = null;
  newMaterial.value = getDefaultMaterial();
  injectCategoryItems();
  fillFormFields();
  showDialog.value = true;
};

const editMaterial = (item) => {
  editedMaterial.value = { ...item };
  newMaterial.value = { ...getDefaultMaterial(), ...item };
  injectCategoryItems();
  fillFormFields();
  showDialog.value = true;
};

const saveMaterial = async () => {
  const formData = new FormData();

  if (editedMaterial.value?.uuid) {
    formData.append("_method", "PUT");
  }

  formFields.forEach((field) => {
    const val = field.model.value;
    if (field.type === "file" && val) {
      formData.append(field.key, Array.isArray(val) ? val[0] : val);
    } else {
      formData.append(field.key, val);
    }
  });

  try {
    if (editedMaterial.value?.uuid) {
      await api.post(`api/materials/${editedMaterial.value.uuid}`, formData);
      showToast("Material updated successfully", "success");
    } else {
      await api.post("api/materials", formData);
      showToast("Material created successfully", "success");
    }
    showDialog.value = false;
    fetchData();
  } catch (err) {
    console.error("Gagal simpan:", err);
    showToast("Failed to save material", "error");
  }
};

const confirmDelete = async () => {
  try {
    await api.delete(`api/materials/${materialToDelete.value.uuid}`);
    showToast("Material deleted successfully", "success");
    fetchData();
  } catch (err) {
    console.error("Gagal hapus:", err);
    showToast("Failed to delete material", "error");
  } finally {
    showDeleteDialog.value = false;
    materialToDelete.value = null;
  }
};

/* ===== Audit Trail ===== */
const showAuditDialog = ref(false);
const selectedMaterialId = ref(null);
const viewAudit = (item) => {
  selectedMaterialId.value = item.uuid;
  showAuditDialog.value = true;
};

/* ===== Init ===== */
onMounted(async () => {
  await fetchCategories();
  fetchData();
});

watch(categoryOptions, () => {
  injectCategoryItems();
});
</script>
