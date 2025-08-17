<template>
  <v-app>
    <!-- Sidebar -->
    <Sidebar :mini="mini" @toggle="toggleMini" />

    <!-- Main Area -->
    <div
      class="main-area"
      :style="{ marginLeft: `${mini ? miniWidth : drawerWidth}px` }"
    >
      <!-- Header -->
      <Header :mini="mini" @toggle="toggleMini" />

      <!-- Page Content -->
      <v-main>
        <v-container class="mt-4">
          <router-view />
        </v-container>
      </v-main>
    </div>
  </v-app>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Sidebar from './Sidebar.vue'
import Header from './Header.vue'

const mini = ref(false)
const drawerWidth = 240
const miniWidth = 64

const toggleMini = () => {
  mini.value = !mini.value
}

const auth = useAuthStore()
const router = useRouter()

onMounted(() => {
  if (!auth.user) {
    router.push({ name: 'Login' })
  }
})
</script>

<style scoped>
.main-area {
  transition: margin-left 0.3s ease;
}
</style>