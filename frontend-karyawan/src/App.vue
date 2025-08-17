<template>
  <v-app>
    <router-view />
  </v-app>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();

const showDialog = ref(false);
const showAuditDialog = ref(false);
const snackbar = ref({
  show: false,
  color: "",
  text: "",
});

onMounted(() => {
  window.addEventListener("unauthenticated", () => {
    showDialog.value = false;
    showAuditDialog.value = false;

    router.push({ name: "login" });

    snackbar.value = {
      show: true,
      color: "error",
      text: "Sesi kamu sudah berakhir. Silakan login ulang.",
    };
  });
});
</script>
