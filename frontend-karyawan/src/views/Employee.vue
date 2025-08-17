<template>
  <v-container>
    <!-- Judul -->
    <h1 class="text-h5 font-weight-bold mb-4">Employee</h1>

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
          Add Employee
        </v-btn>
      </v-col>
    </v-row>

    <!-- Table -->
    <BaseTable
      :headers="headers"
      :items="employees"
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
          <v-btn icon size="small" color="blue" @click="editEmployee(item)">
            <v-icon>mdi-pencil</v-icon>
          </v-btn>
          <v-btn icon size="small" color="red" @click="deleteEmployee(item)">
            <v-icon>mdi-delete</v-icon>
          </v-btn>
        </div>
      </template>
    </BaseTable>

    <!-- Dialog Form -->
    <BaseFormDialog
      v-model="showDialog"
      :title="editedEmployee ? 'Edit Employee' : 'Add Employee'"
      :fields="formFields"
      @save="saveEmployee"
    />

    <!-- Confirm Delete -->
    <ConfirmDeleteDialog
      v-model="showDeleteDialog"
      :itemName="employeeToDelete?.employeeName"
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
  </v-container>
</template>

<script setup>
import { ref, reactive, computed, onMounted, toRef } from "vue";
import api from "@/axios";
import BaseTable from "@/components/BaseTable.vue";
import BaseFilter from "@/components/BaseFilter.vue";
import BaseFormDialog from "@/components/BaseFormDialog.vue";
import ConfirmDeleteDialog from "@/components/ConfirmDeleteDialog.vue";
import { usePaginatedTable } from "@/components/usePaginatedTable";
import { useToast } from "@/components/useToast";

const showFilter = ref(false);
const showDialog = ref(false);
const showDeleteDialog = ref(false);
const editedEmployee = ref(null);
const employeeToDelete = ref(null);
const { toast, showToast } = useToast();

/* ===== Rules ===== */
const rules = {
  required: (v) => !!v || "This field is required",
};

/* ===== Schema ===== */
const departmentOptions = ref([]);

const employeeSchema = [
  {
    key: "employeeID",
    label: "Employee ID",
    type: "text",
    default: "",
    rules: [rules.required],
  },
  {
    key: "name",
    label: "Employee Name",
    type: "text",
    default: "",
    rules: [rules.required],
  },
  {
    key: "address",
    label: "Address",
    type: "text",
    default: "",
    rules: [rules.required],
  },
  {
    key: "departmentID",
    label: "Department",
    type: "select",
    options: departmentOptions,
    default: "",
    rules: [rules.required],
  },
];

/* ===== Form Model (Reactive) ===== */
const getDefaultEmployee = () =>
  employeeSchema.reduce((acc, f) => ({ ...acc, [f.key]: f.default }), {});
const formModel = reactive(getDefaultEmployee());

const formFields = computed(() =>
  employeeSchema.map((f) => ({
    ...f,
    model: toRef(formModel, f.key),
    options: f.key === "departmentID" ? departmentOptions.value : f.options,
  }))
);

/* ===== Filter ===== */
const filters = {
  name: ref(""),
  position: ref(""),
};

const filterFields = computed(() => [
  { label: "Filter by Employee", type: "text", model: filters.name },
  { label: "Filter by Position", type: "text", model: filters.position },
]);

/* ===== Normalizer ===== */
const normalizeEmployee = (item) => ({
  id: item.id,
  employeeID: item.employeeID,
  name: item?.name || "Unnamed",
  address: item.address,
  departmentName: item?.departementName || "-",
});

/* ===== Table Composable ===== */
const {
  items: employees,
  loading,
  totalItems,
  page,
  itemsPerPage,
  sortBy,
  sortDesc,
  fetchData,
} = usePaginatedTable({
  endpoint: "api/employee/GetData",
  normalizeFn: normalizeEmployee,
  filters,
});

const setSortBy = (val) => {
  sortBy.value = Array.isArray(val) ? val : [val];
};
const setSortDesc = (val) => {
  sortDesc.value = Array.isArray(val) ? val : [val];
};

/* ===== Headers ===== */
const headers = [
  { title: "Employee ID", key: "employeeID" },
  { title: "Employee Name", key: "name" },
  { title: "Department", key: "departmentName" },
  { title: "Address", key: "address" },
  { title: "Actions", key: "actions", sortable: false, width: "140px" },
];


/* ===== Form Logic ===== */
const openAddDialog = () => {
  editedEmployee.value = null;
  Object.assign(formModel, getDefaultEmployee());
  showDialog.value = true;
};

const editEmployee = async (item) => {
  editedEmployee.value = { ...item };

  if (departmentOptions.value.length === 0) {
    await fetchDepartement();
  }

  const matched = departmentOptions.value.find(
    (opt) => opt.title === item.departmentName
  );
  
  const resolvedDepartmentID = matched?.value ?? "";

  const filled = {
    ...getDefaultEmployee(),
    ...item,
    departmentID: resolvedDepartmentID,
  };

  Object.assign(formModel, filled);
  showDialog.value = true;
};

const getPayloadFromFormFields = () => ({ ...formModel });

const getChangedFields = () => {
  const changes = {};
  for (const key in formModel) {
    if (formModel[key] !== editedEmployee.value?.[key]) {
      changes[key] = formModel[key];
    }
  }
  return changes;
};

const saveEmployee = async () => {
  const payload = editedEmployee.value
    ? getChangedFields()
    : getPayloadFromFormFields();

  try {
    if (editedEmployee.value?.id) {
      if (Object.keys(payload).length === 0) {
        showToast("No changes to update", "info");
        showDialog.value = false;
        return;
      }

      await api.put(`api/employee/${editedEmployee.value.id}`, payload);
      showToast("Employee updated successfully", "success");
    } else {
      await api.post("api/employee", payload);
      showToast("Employee created successfully", "success");
    }

    showDialog.value = false;
    fetchData();
  } catch (err) {
    console.error("Gagal simpan:", err);

    const message =
      err?.response?.data?.error || "Failed to save employee";
    showToast(message, "error");
  }
};

/* ===== Delete Logic ===== */
const deleteEmployee = (item) => {
  employeeToDelete.value = item;
  showDeleteDialog.value = true;
};

const confirmDelete = async () => {
  try {
    await api.delete(`api/employee/${employeeToDelete.value.id}`);
    showToast("Employee deleted successfully", "success");
    fetchData();
  } catch (err) {
    console.error("Gagal hapus:", err);
    showToast("Failed to delete employee", "error");
  } finally {
    showDeleteDialog.value = false;
    employeeToDelete.value = null;
  }
};

/* ===== Fetch Department ===== */
const fetchDepartement = async () => {
  try {
    const res = await api.post("api/departement/GetData", {});
    const raw = res.data.data ?? [];
    departmentOptions.value = raw.map((d) => ({
      title: d.departementName,
      value: d.id,
    }));
  } catch (err) {
    console.error("Gagal fetch departement:", err);
  }
};

/* ===== Init ===== */
onMounted(async () => {
  fetchData();
  await fetchDepartement();
});
</script>