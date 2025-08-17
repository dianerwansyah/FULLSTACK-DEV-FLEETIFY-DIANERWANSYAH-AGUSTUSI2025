<template>
  <v-container>
    <h1 class="text-h5 font-weight-bold mb-4">Log Absensi Saya</h1>

    <!-- Filter Tanggal -->
    <v-row class="mb-4" dense align="center" justify="start" style="gap: 16px">
      <!-- Dari Tanggal -->
      <v-col cols="12" md="3">
        <v-menu
          v-model="menuFrom"
          :close-on-content-click="false"
          transition="scale-transition"
          offset-y
        >
          <template #activator="{ props }">
            <v-text-field
              :model-value="moment(dateFrom).format('DD-MM-YYYY')"
              label="Dari Tanggal"
              readonly
              density="compact"
              hide-details
              v-bind="props"
            />
          </template>
          <v-card>
            <v-date-picker
              v-model="dateFrom"
              :max="dateTo"
              show-adjacent-months
            />
            <v-card-actions class="justify-end">
              <v-btn variant="text" @click="menuFrom = false">OK</v-btn>
            </v-card-actions>
          </v-card>
        </v-menu>
      </v-col>

      <!-- Sampai Tanggal -->
      <v-col cols="12" md="3">
        <v-menu
          v-model="menuTo"
          :close-on-content-click="false"
          transition="scale-transition"
          offset-y
        >
          <template #activator="{ props }">
            <v-text-field
              :model-value="moment(dateTo).format('DD-MM-YYYY')"
              label="Sampai Tanggal"
              readonly
              density="compact"
              hide-details
              v-bind="props"
            />
          </template>
          <v-card>
            <v-date-picker
              v-model="dateTo"
              :min="dateFrom"
              show-adjacent-months
            />
            <v-card-actions class="justify-end">
              <v-btn variant="text" @click="menuTo = false">OK</v-btn>
            </v-card-actions>
          </v-card>
        </v-menu>
      </v-col>

      <!-- Tombol Filter -->
      <v-col cols="12" md="auto">
        <v-btn color="primary" @click="fetchData" class="mt-1">
          <v-icon start>mdi-filter</v-icon>
          Filter
        </v-btn>
      </v-col>
    </v-row>

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
      @update:page="(val) => (page.value = val)"
      @update:itemsPerPage="(val) => (itemsPerPage.value = val)"
      @update:sortBy="setSortBy"
      @update:sortDesc="setSortDesc"
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
import { ref, computed, onMounted } from "vue";
import moment from "moment";
import BaseTable from "@/components/BaseTable.vue";
import { usePaginatedTable } from "@/components/usePaginatedTable";
import { useToast } from "@/components/useToast";

const { toast, showToast } = useToast();

// Sort
const setSortBy = (val) => {
  sortBy.value = Array.isArray(val) ? val : [val];
};
const setSortDesc = (val) => {
  sortDesc.value = Array.isArray(val) ? val : [val];
};

/* ===== Date Range State ===== */
const dateFrom = ref(moment().startOf("month").toDate());
const dateTo = ref(moment().endOf("month").toDate());
const menuFrom = ref(false);
const menuTo = ref(false);

/* ===== Filters ===== */
const filters = computed(() => ({
  "date_attendance.gte": moment(dateFrom.value).format("YYYY-MM-DD"),
  "date_attendance.lte": moment(dateTo.value).format("YYYY-MM-DD"),
}));

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
  date_attendance: moment(item.dateAttendance).format("dddd, DD MMMM YYYY"),
  departementName: item.departementName || "-",
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
  fetchData: baseFetch,
} = usePaginatedTable({
  endpoint: "api/attendance/logs",
  normalizeFn: normalizeLog,
  filters,
});

/* ===== Headers ===== */
const headers = [
  { title: "Tanggal", key: "date_attendance" },
  { title: "Departemen", key: "departementName" },
  { title: "Tipe", key: "attendanceType" },
  { title: "Jam", key: "clock" },
  { title: "Batas Jam", key: "maxClock" },
  { title: "Status", key: "status" },
  { title: "Deskripsi", key: "description" },
];

/* ===== Fetch with Validation ===== */
const fetchData = async () => {
  if (moment(dateFrom.value).isAfter(dateTo.value)) {
    showToast(
      "Tanggal awal tidak boleh lebih besar dari tanggal akhir",
      "error"
    );
    return;
  }

  try {
    await baseFetch();
  } catch (err) {
    console.error("Gagal ambil data log:", err);
    showToast("Gagal mengambil data absensi", "error");
  }
};

onMounted(() => {
  fetchData();
});
</script>
