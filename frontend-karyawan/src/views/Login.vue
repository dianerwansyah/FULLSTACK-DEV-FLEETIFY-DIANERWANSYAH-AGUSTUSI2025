<template>
  <v-container class="fill-height d-flex justify-center align-center">
    <v-card width="400" class="pa-6 rounded-xl">
      <!-- Tombol Back -->
      <v-btn icon variant="text" @click="$router.back()" class="mb-4">
        <v-icon>mdi-arrow-left</v-icon>
      </v-btn>

      <h2 class="text-h5 mb-6 font-weight-bold text-center">Login</h2>

      <v-form v-model="valid" @submit.prevent="handleLogin">
        <v-text-field
          label="Employee ID"
          v-model="state.employeeID"
          :rules="[(v) => !!v || 'Employee is required']"
          outlined
          prepend-inner-icon="mdi-account"
          autocomplete="off"
          class="mb-4"
        />

        <v-text-field
          label="Password"
          v-model="state.password"
          :rules="[(v) => !!v || 'Password is required']"
          type="password"
          outlined
          prepend-inner-icon="mdi-lock"
          autocomplete="new-password"
          class="mb-6"
        />

        <v-btn
          :disabled="!valid || state.loading"
          type="submit"
          color="primary"
          size="large"
          block
        >
          <template v-if="state.loading">
            <v-progress-circular indeterminate color="white" size="20" />
          </template>
          <template v-else> Login </template>
        </v-btn>
      </v-form>

      <v-divider class="my-6"></v-divider>

      <!-- <div class="text-center">
        <span class="text-caption">Don't have an account?</span>
        <v-btn
          variant="text"
          color="primary"
          class="text-caption"
          @click="$router.push('/register')"
        >
          Register
        </v-btn>
      </div> -->
    </v-card>
  </v-container>
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
</template>

<script setup>
import { ref, reactive } from "vue";
import api from "@/axios";
import { useRouter } from "vue-router";
import { useToast } from "@/components/useToast";
import { useAuthStore } from '@/stores/auth';

const router = useRouter();
const valid = ref(false);
const { toast, showToast } = useToast();
const auth = useAuthStore();

const state = reactive({
  employeeID: "",
  password: "",
  loading: false,
});

const handleLogin = async () => {
  state.loading = true;

  try {
    const response = await api.post("api/auth/login", {
      employeeID: state.employeeID,
      password: state.password,
    });

    const me = await api.get("/api/auth/me");
    auth.setUser(me.data)

    showToast("Login successful", "success");

    setTimeout(() => {
      router.push("/dashboard");
    }, 1000);
  } catch (error) {
    let message = "An error occurred, please try again.";

    if (error.response) {
      message =
        error.response.data?.error ||
        error.response.data?.message ||
        JSON.stringify(error.response.data) ||
        message;
    } else if (error.request) {
      message = "No response from server. Please try again later.";
    } else {
      message = error.error;
    }

    showToast(message, "error");
  } finally {
    state.loading = false;
  }
};
</script>
