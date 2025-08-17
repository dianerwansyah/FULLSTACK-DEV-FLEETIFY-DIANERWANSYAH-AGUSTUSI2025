<template>
  <v-navigation-drawer
    v-model="drawer"
    :mini-variant="mini"
    permanent
    class="custom-drawer"
    :style="{ width: mini ? `${miniWidth}px` : `${drawerWidth}px` }"
  >
    <!-- Toggle Button -->
    <v-list-item
      @click="$emit('toggle')"
      class="justify-center align-center"
      :style="{ height: '64px' }"
    >
      <v-icon>{{ mini ? "mdi-menu-open" : "mdi-menu" }}</v-icon>
    </v-list-item>

    <v-divider />

    <!-- Menu Items -->
    <v-list dense>
      <v-list-item
        v-for="item in filteredMenu"
        :key="item.path"
        @click="navigate(item.path)"
        link
      >
        <div class="d-flex align-center">
          <v-list-item-icon class="mr-2">
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>
          <v-list-item-title v-if="!mini">
            {{ item.meta.title }}
          </v-list-item-title>
        </div>
      </v-list-item>
    </v-list>
  </v-navigation-drawer>
</template>

<script setup>
import { ref, computed } from "vue"
import { useRouter } from "vue-router"
import menu from "@/router/menu.json"

const props = defineProps({ mini: Boolean })
const drawer = ref(true)
const router = useRouter()

const drawerWidth = 240
const miniWidth = 64

const filteredMenu = computed(() =>
  menu.filter((item) => item.meta.requiresAuth === true)
)

const navigate = (path) => {
  router.push(path)
}
</script>

<style scoped>
.custom-drawer {
  position: fixed;
  top: 0;
  left: 0;
  height: 100vh;
  z-index: 20;
  transition: width 0.3s ease;
}
</style>