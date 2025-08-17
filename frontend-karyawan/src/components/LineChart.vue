<template>
  <div>
    <Line :data="chartData" :options="chartOptions" />
  </div>
</template>

<script setup>
import { Line } from "vue-chartjs";
import {
  Chart as ChartJS,
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
} from "chart.js";

ChartJS.register(
  Title,
  Tooltip,
  Legend,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale
);

// Props
const props = defineProps({
  labels: { type: Array, required: true },
  data: { type: Array, required: true },
});

// Chart Data
const chartData = {
  labels: props.labels,
  datasets: [
    {
      label: "Jumlah Kehadiran",
      data: props.data,
      fill: false,
      borderColor: "#1976D2",
      backgroundColor: "#1976D2",
      tension: 0.3,
      pointRadius: 4,
    },
  ],
};

// Chart Options
const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: { display: false },
    tooltip: {
      callbacks: {
        label: (ctx) => ` ${ctx.parsed.y} hadir`,
      },
    },
  },
  scales: {
    x: {
      ticks: { color: "#555" },
      grid: { display: false },
    },
    y: {
      beginAtZero: true,
      ticks: { stepSize: 5, color: "#555" },
      grid: { color: "#eee" },
    },
  },
};
</script>

<style scoped>
div {
  height: 300px;
}
</style>
