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
      :items="departments"
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
      <template #item.actions="{ item }">
        <div class="d-flex align-center" style="gap: 8px">
          <v-btn icon size="small" color="blue" @click="editDepartment(item)">
            <v-icon>mdi-pencil</v-icon>
          </v-btn>
          <v-btn icon size="small" color="red" @click="deleteDepartment(item)">
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </div>
      </template>
    </BaseTable>

    <!-- Dialog Form -->
    <BaseFormDialog
      v-model="showDialog"
      :title="editedDepartment ? 'Edit Department' : 'Add Department'"
      :fields="formFields"
      @save="saveDepartment"
    />
  </v-container>

  <!-- Confirm Delete -->
  <ConfirmDeleteDialog
    v-model="showDeleteDialog"
    :itemName="departmentToDelete?.name"
    @confirm="confirmDelete"
    @cancel="showDeleteDialog = false"
  />

  <!-- Snackbar -->
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
import { ref, computed, onMounted } from "vue";
import api from "@/axios";
import BaseTable from "@/components/BaseTable.vue";
import BaseFilter from "@/components/BaseFilter.vue";
import BaseFormDialog from "@/components/BaseFormDialog.vue";
import ConfirmDeleteDialog from "@/components/ConfirmDeleteDialog.vue";
import { usePaginatedTable } from "@/components/usePaginatedTable";
import { useToast } from "@/components/useToast";
import moment from "moment";

const showFilter = ref(false);
const showDialog = ref(false);
const showDeleteDialog = ref(false);
const editedDepartment = ref(null);
const departmentToDelete = ref(null);
const { toast, showToast } = useToast();

/* ===== Rules ===== */
const rules = {
  required: (v) => !!v || "This field is required",
};

/* ===== Schema ===== */
const departmentSchema = [
  {
    key: "departementName",
    label: "Department Name",
    type: "text",
    default: "",
    rules: [rules.required],
  },
  {
    key: "maxClockInTime",
    label: "Max Clock In Time",
    type: "time",
    default: "",
    rules: [rules.required],
  },
  {
    key: "maxClockOutTime",
    label: "Max Clock Out Time",
    type: "time",
    default: "",
    rules: [rules.required],
  },
];

const filters = {
  departementName: ref(""),
};

const filterFields = computed(() => [
  {
    label: "Filter by Departement Name",
    type: "text",
    model: filters.departementName,
  },
]);

const normalizeDepartment = (item) => ({
  id: item.id,
  departementName: item.departementName || "Unnamed",
  maxClockInTime: moment.utc(item.maxClockInTime).format("HH:mm"),
  maxClockOutTime: moment.utc(item.maxClockOutTime).format("HH:mm"),
});

const {
  items: departments,
  loading,
  totalItems,
  page,
  itemsPerPage,
  sortBy,
  sortDesc,
  fetchData,
} = usePaginatedTable({
  endpoint: "api/departement/GetData",
  normalizeFn: normalizeDepartment,
  filters,
});

const setSortBy = (val) => {
  sortBy.value = Array.isArray(val) ? val : [val];
};
const setSortDesc = (val) => {
  sortDesc.value = Array.isArray(val) ? val : [val];
};

const headers = [
  { title: "Department Name", key: "departementName" },
  { title: "Max Clock In Time", key: "maxClockInTime" },
  { title: "Max Clock Out Time", key: "maxClockOutTime" },
  { title: "Actions", key: "actions", sortable: false, width: "140px" },
];

const getDefaultDepartment = () =>
  departmentSchema.reduce((acc, f) => ({ ...acc, [f.key]: f.default }), {});

const newDepartment = ref(getDefaultDepartment());

const formFields = departmentSchema.map((f) => ({
  ...f,
  model: ref(f.default),
}));

const fillFormFields = () => {
  formFields.forEach((field) => {
    field.model.value = newDepartment.value[field.key];
  });
};

const openAddDialog = () => {
  editedDepartment.value = null;
  newDepartment.value = getDefaultDepartment();
  fillFormFields();
  showDialog.value = true;
};

const editDepartment = (item) => {
  editedDepartment.value = { ...item };
  newDepartment.value = { ...getDefaultDepartment(), ...item };
  fillFormFields();
  showDialog.value = true;
};

const saveDepartment = async () => {
  const isEdit = !!editedDepartment.value?.id;
  const payload = {};

  formFields.forEach((field) => {
    const key = field.key;
    const newVal = field.model.value;
    const oldVal = editedDepartment.value?.[key];

    if (isEdit) {
      if (newVal !== oldVal) {
        payload[key] = newVal;
      }
    } else {
      payload[key] = newVal;
    }
  });


  try {
    if (isEdit) {
      await api.put(`api/departement/${editedDepartment.value.id}`, payload);
      showToast("Department updated successfully", "success");
    } else {
      await api.post("api/departement", payload);
      showToast("Department created successfully", "success");
    }

    showDialog.value = false;
    fetchData();
  } catch (err) {
    console.error("Gagal simpan:", err);
    showToast("Failed to save department", "error");
  }
};

const deleteDepartment = (item) => {
  departmentToDelete.value = item;
  showDeleteDialog.value = true;
};

const confirmDelete = async () => {
  try {
    await api.delete(`api/departement/${departmentToDelete.value.id}`);
    showToast("Department deleted successfully", "success");
    fetchData();
  } catch (err) {
    console.error("Gagal hapus:", err);
    showToast("Failed to delete department", "error");
  } finally {
    showDeleteDialog.value = false;
    departmentToDelete.value = null;
  }
};

/* ===== Init ===== */
onMounted(() => {
  fetchData();
});
</script>
