<template>
  <v-app>
    <!-- App Bar -->
    <v-app-bar app color="white" elevation="2">
      <v-toolbar-title>KaryawanApp</v-toolbar-title>
      <v-spacer></v-spacer>
      <v-btn text to="/login" class="mr-4">Login</v-btn>
      <!-- <v-btn text to="/register">Register</v-btn> -->
    </v-app-bar>

    <!-- Main Content -->
    <v-main class="flex flex-col justify-center items-center text-center px-6">
      <h2 class="text-4xl font-extrabold mt-4 mb-4">Welcome to Karyawan Management</h2>
      <p class="mb-8 text-gray-600">
        Manage your employees, attendance, and roles with secure access and full audit trail.
      </p>
      <v-btn color="primary" class="mt-4">Get Started</v-btn>
    </v-main>

    <!-- Footer -->
    <v-footer app padless>
      <v-col class="text-center" cols="12">
        &copy; 2025 KaryawanApp. All rights reserved.
      </v-col>
    </v-footer>
  </v-app>
</template>