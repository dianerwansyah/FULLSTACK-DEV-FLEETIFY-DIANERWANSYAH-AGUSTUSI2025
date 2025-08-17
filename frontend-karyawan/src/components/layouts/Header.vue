<template>
  <v-app-bar
    color="primary"
    dark
    flat
    :style="{ top: '0px', left: '0px', position: 'relative' }"
  >
    <v-toolbar-title>KaryawanApp</v-toolbar-title>

    <v-spacer />

    <v-btn icon @click="logout">
      <v-icon>mdi-logout</v-icon>
    </v-btn>
  </v-app-bar>
</template>

<script setup>
import { useRouter } from "vue-router"
import { useAuthStore } from '@/stores/auth';

const router = useRouter()

async function logout() {
  const auth = useAuthStore();
  try {
    await api.post('/api/auth/logout');
  } catch (e) {
    console.warn('Logout failed:', e);
  }
  auth.clearUser();
  router.push('/login');
}
</script>

<style scoped>
.v-app-bar {
  width: 100% !important;
}
</style>