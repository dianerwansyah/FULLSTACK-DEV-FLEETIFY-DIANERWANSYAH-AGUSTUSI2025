<template>
  <v-container>
    <h1 class="text-h5 font-weight-bold mb-4">Data Absensi Karyawan</h1>

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

    <!-- Filter with Expand Transition -->
    <v-expand-transition>
      <div v-show="showFilter">
        <BaseFilter :fields="filterFields" class="mb-6" />
      </div>
    </v-expand-transition>

    <!-- Table -->
    <BaseTable
      :headers="headers"
      :items="logs"
      :loading="loading"
      :total-items="totalItems"
      :page="page"
      :items-per-page="itemsPerPage"
      :sort-by="sortBy"
      :sort-desc="sortDesc"
      @update:page="(val) => (page = val)"
      @update:itemsPerPage="(val) => (itemsPerPage = val)"
      @update:sortBy="(val) => (sortBy = Array.isArray(val) ? val : [val])"
      @update:sortDesc="(val) => (sortDesc = Array.isArray(val) ? val : [val])"
    >
      <template #item.attendanceType="{ item }">
        <v-chip
          :color="item.attendanceType === 'in' ? 'green' : 'red'"
          size="small"
          variant="elevated"
        >
          <v-icon start>
            {{ item.attendanceType === "in" ? "mdi-clock-in" : "mdi-clock-out" }}
          </v-icon>
          {{ item.attendanceType === "in" ? "Clock In" : "Clock Out" }}
        </v-chip>
      </template>

      <template #item.status="{ item }">
        <v-chip
          :color="getStatusColor(item.status)"
          size="small"
          variant="elevated"
        >
          {{ item.status }}
        </v-chip>
      </template>

      <template #item.description="{ item }">
        <v-tooltip location="top">
          <template #activator="{ props }">
            <span v-bind="props">
              {{
                item.description.length > 40
                  ? item.description.slice(0, 40) + "..."
                  : item.description
              }}
            </span>
          </template>
          <span style="white-space: pre-line; max-width: 300px">
            {{ item.description }}
          </span>
        </v-tooltip>
      </template>
    </BaseTable>
  </v-container>
</template>

<script setup>
import { ref, computed, onMounted, watch } from "vue";
import moment from "moment";
import api from "@/axios";
import BaseTable from "@/components/BaseTable.vue";
import BaseFilter from "@/components/BaseFilter.vue";
import { usePaginatedTable } from "@/components/usePaginatedTable";
import { useToast } from "@/components/useToast";

const { toast, showToast } = useToast();

/* ===== Toggle Filter ===== */
const showFilter = ref(true);

/* ===== Departemen Options ===== */
const departementOptions = ref([]);
const loadingDepartement = ref(false);

const fetchDepartement = async () => {
  loadingDepartement.value = true;
  try {
    const res = await api.post("api/departement/GetData", {});
    const raw = res.data.data ?? [];
    departementOptions.value = raw.map((d) => ({
      title: d.departementName ?? "-",
      value: d.id ?? null,
    }));
  } catch (err) {
    console.error("Gagal fetch departement:", err);
  } finally {
    loadingDepartement.value = false;
  }
};

/* ===== Filter Models ===== */
const filters = {
  "date_attendance.gte": ref(moment().startOf("month").format("YYYY-MM-DD")),
  "date_attendance.lte": ref(moment().endOf("month").format("YYYY-MM-DD")),
  departementID: ref(""),
};

/* ===== Filter Fields for BaseFilter ===== */
const filterFields = computed(() => [
  {
    label: "Dari Tanggal",
    type: "date",
    model: filters["date_attendance.gte"],
    max: filters["date_attendance.lte"].value,
    menu: ref(false),
  },
  {
    label: "Sampai Tanggal",
    type: "date",
    model: filters["date_attendance.lte"],
    min: filters["date_attendance.gte"].value,
    menu: ref(false),
  },
  {
    label: "Departemen",
    type: "select",
    model: filters["departementID"],
    items: departementOptions.value,
    itemTitle: "title",
    itemValue: "value",
    clearable: true,
    loading: loadingDepartement.value,
  },
]);

/* ===== Reset Menu on Show Filter ===== */
watch(showFilter, (val) => {
  if (val) {
    filterFields.value.forEach((field) => {
      if (
        field.type === "date" &&
        field.menu &&
        typeof field.menu.value !== "undefined"
      ) {
        field.menu.value = false;
      }
    });
  }
});

/* ===== Status Color Helper ===== */
const getStatusColor = (status) => {
  if (status === "Terlambat") return "orange";
  if (status === "Pulang Cepat") return "deep-orange";
  if (status === "Tepat") return "green";
  return "grey";
};

/* ===== Normalizer ===== */
const normalizeLog = (item) => ({
  id: item.id,
  employeeID: item.employeeID,
  employeeName: item.employeeName,
  departementName: item.departementName || "-",
  dateAttendance: moment(item.dateAttendance).format("dddd, DD MMMM YYYY"),
  clock: item.clock ? moment(item.clock).format("HH:mm:ss") : "-",
  maxClock: item.maxClock || "-",
  status: item.status || "-",
  attendanceType: item.attendanceType,
  description: item.description || "-",
});

/* ===== Table State ===== */
const {
  items: logs,
  loading,
  totalItems,
  page,
  itemsPerPage,
  sortBy,
  sortDesc,
  fetchData,
} = usePaginatedTable({
  endpoint: "api/attendance/GetData",
  normalizeFn: normalizeLog,
  filters,
});

/* ===== Headers ===== */
const headers = [
  { title: "Tanggal", key: "dateAttendance" },
  { title: "Nama", key: "employeeName" },
  { title: "Departemen", key: "departementName" },
  { title: "Tipe", key: "attendanceType" },
  { title: "Jam", key: "clock" },
  { title: "Batas Jam", key: "maxClock" },
  { title: "Status", key: "status" },
  { title: "Deskripsi", key: "description" },
];

/* ===== Init ===== */
onMounted(async () => {
  fetchData();
  await fetchDepartement();
});
</script>
