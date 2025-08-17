<template>
  <div>
    <!-- Statistik Ringkas -->
    <v-row class="mb-6" dense>
      <v-col cols="12" md="3" v-for="(label, key) in statLabels" :key="key">
        <v-card elevation="2" class="pa-4">
          <div class="text-subtitle-1 font-weight-medium">{{ label }}</div>
          <div class="text-h5 font-weight-bold">{{ summaryStats[key] }}</div>
        </v-card>
      </v-col>
    </v-row>

    <!-- Grafik Tren Kehadiran -->
    <v-card elevation="2" class="mb-6 pa-4">
      <div class="text-subtitle-1 font-weight-medium mb-2">Tren Kehadiran</div>
      <LineChart
        :labels="attendanceTrend.labels"
        :data="attendanceTrend.data"
      />
    </v-card>

    <!-- Tabel Status Terakhir -->
    <v-card elevation="2" class="pa-4">
      <div class="text-subtitle-1 font-weight-medium mb-2">
        Status Terakhir Karyawan
      </div>
      <v-data-table
        :headers="lastStatusHeaders"
        :items="lastStatusItems"
        :items-per-page="10"
        class="elevation-1"
      >
        <template #item.status="{ item }">
          <v-chip :color="statusColor(item.status)" dark>{{
            item.status
          }}</v-chip>
        </template>
      </v-data-table>
    </v-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import LineChart from "@/components/LineChart.vue";

const summaryStats = ref({
  total: 120,
  hadir: 95,
  terlambat: 15,
  pulangCepat: 10,
});
const statLabels = {
  total: "Total Absensi",
  hadir: "Hadir Tepat",
  terlambat: "Terlambat",
  pulangCepat: "Pulang Cepat",
};

const attendanceTrend = ref({
  labels: ["01 Aug", "02 Aug", "03 Aug", "04 Aug", "05 Aug"],
  data: [20, 22, 18, 25, 35],
});

const lastStatusHeaders = [
  { title: "Nama", key: "employeeName" },
  { title: "Departemen", key: "departementName" },
  { title: "Tanggal", key: "dateAttendance" },
  { title: "Jam", key: "clock" },
  { title: "Status", key: "status" },
];
const lastStatusItems = ref([
  {
    employeeName: "Rina Wijaya",
    departementName: "Finance",
    dateAttendance: "2025-08-17",
    clock: "08:03",
    status: "Terlambat",
  },
  {
    employeeName: "Budi Santoso",
    departementName: "IT",
    dateAttendance: "2025-08-17",
    clock: "07:55",
    status: "Hadir",
  },
  {
    employeeName: "Sari Lestari",
    departementName: "HR",
    dateAttendance: "2025-08-17",
    clock: "16:20",
    status: "Pulang Cepat",
  },
]);

const statusColor = (status) => {
  switch (status) {
    case "Hadir":
      return "green";
    case "Terlambat":
      return "red";
    case "Pulang Cepat":
      return "orange";
    default:
      return "grey";
  }
};
</script>
