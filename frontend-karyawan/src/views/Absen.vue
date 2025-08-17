<template>
  <v-container class="pa-4" max-width="500">
    <v-card elevation="4">
      <v-card-title>
        <v-icon class="mr-2">mdi-calendar</v-icon>
        Menu Absen Karyawan
      </v-card-title>

      <v-card-text>
        <div v-if="loading">
          <v-progress-circular indeterminate color="primary" />
        </div>
        <div v-else>
          <!-- Status Absen -->
          <p class="mb-4 d-flex align-center">
            <v-icon
              class="mr-2"
              :color="
                attendance?.clock_in && attendance?.clock_out
                  ? 'green'
                  : attendance?.clock_in
                  ? 'orange'
                  : 'red'
              "
            >
              {{
                attendance?.clock_in && attendance?.clock_out
                  ? "mdi-check-circle"
                  : attendance?.clock_in
                  ? "mdi-clock-outline"
                  : "mdi-alert-circle"
              }}
            </v-icon>
            <strong>Status:</strong>
            <span class="ml-2">
              {{
                attendance?.clock_in && attendance?.clock_out
                  ? "Sudah Clock In & Out"
                  : attendance?.clock_in
                  ? "Sudah Clock In"
                  : "Belum Absen"
              }}
            </span>
          </p>

          <!-- Tombol Absen -->
          <v-btn
            color="green"
            class="mr-2"
            :disabled="!canClockIn"
            @click="handleClockIn"
          >
            <v-icon left>mdi-clock-in</v-icon>
            Clock In
          </v-btn>

          <v-btn color="red" :disabled="!canClockOut" @click="handleClockOut">
            <v-icon left>mdi-clock-out</v-icon>
            Clock Out
          </v-btn>

          <v-divider class="my-4" />

          <!-- Detail Waktu & Lokasi -->
          <p>
            <v-icon class="mr-1">mdi-clock-outline</v-icon>
            Clock In: {{ attendance?.clock_in || "-" }}
          </p>
          <p>
            <v-icon class="mr-1">mdi-clock-outline</v-icon>
            Clock Out: {{ attendance?.clock_out || "-" }}
          </p>
          <p>
            <v-icon class="mr-1">mdi-map-marker</v-icon>
            Lokasi: {{ location || "-" }}
          </p>
        </div>
      </v-card-text>
    </v-card>

    <!-- Dialog Konfirmasi -->
    <v-dialog v-model="confirmType" max-width="420" persistent>
      <v-card elevation="8">
        <v-card-title class="text-h6 d-flex align-center">
          <v-icon color="primary" class="mr-2">mdi-help-circle</v-icon>
          Konfirmasi {{ confirmType === "clock_in" ? "Clock In" : "Clock Out" }}
        </v-card-title>

        <v-card-text class="text-body-1">
          Kamu akan melakukan
          <strong>{{
            confirmType === "clock_in" ? "Clock In" : "Clock Out"
          }}</strong>
          sekarang. Pastikan lokasi dan waktu sudah sesuai.
        </v-card-text>

        <v-card-actions class="justify-end">
          <v-btn
            variant="text"
            color="grey-darken-1"
            @click="confirmType = null"
          >
            <v-icon left>mdi-arrow-left</v-icon>
            Batal
          </v-btn>
          <v-btn color="primary" variant="elevated" @click="confirmAttendance">
            <v-icon left>mdi-check-bold</v-icon>
            Ya, Lanjut
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

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
import { ref, onMounted, computed } from "vue";
import api from "@/axios";
import { useToast } from "@/components/useToast";

const attendance = ref(null);
const original = ref(null);
const loading = ref(false);
const location = ref("Lokasi tidak tersedia");

const confirmType = ref(null);
const { toast, showToast } = useToast();

const getDeviceInfo = () => navigator.userAgent;

const fetchLocation = () => {
  if (!navigator.geolocation) {
    location.value = "Geolocation tidak didukung";
    return;
  }

  navigator.geolocation.getCurrentPosition(
    (pos) => {
      const { latitude, longitude } = pos.coords;
      location.value = `Lat: ${latitude.toFixed(5)}, Lng: ${longitude.toFixed(
        5
      )}`;
    },
    (err) => {
      console.warn("Gagal ambil lokasi:", err);
      location.value = "Gagal ambil lokasi";
    },
    { enableHighAccuracy: true, timeout: 5000 }
  );
};

const fetchAttendance = async () => {
  loading.value = true;
  try {
    const res = await api.get("/api/attendance/today");
    attendance.value = res.data;
    original.value = JSON.parse(JSON.stringify(res.data));
  } catch (err) {
    console.error("Failed to fetch attendance", err);
    attendance.value = null;

    showToast(
      "Gagal mengambil data absen: " +
        (err?.response?.data?.error || err.message || "Unknown error"),
      "error"
    );
  } finally {
    loading.value = false;
  }
};
const postAttendance = async (type) => {
  const description = `Clock-${type} via web at ${
    location.value
  }, device: ${getDeviceInfo()}`;
  try {
    await api.post("/api/attendance", { type, description });
    await fetchAttendance();

    showToast(
      `Berhasil ${type === "clock_in" ? "Clock In" : "Clock Out"}`,
      "success"
    );
  } catch (err) {
    console.error("Failed to post attendance", err);
    showToast(
      "Gagal absen: " + (err?.response?.data?.error || "Unknown error"),
      "error"
    );
  }
};

const handleClockIn = () => (confirmType.value = "clock_in");
const handleClockOut = () => (confirmType.value = "clock_out");

const confirmAttendance = () => {
  const type = confirmType.value;
  confirmType.value = null;
  postAttendance(type);
};

const canClockIn = computed(() => !attendance.value?.clock_in);
const canClockOut = computed(
  () => attendance.value?.clock_in && !attendance.value?.clock_out
);

onMounted(() => {
  fetchLocation();
  fetchAttendance();
});
</script>
