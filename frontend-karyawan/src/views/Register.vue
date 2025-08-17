<template>
  <v-container class="fill-height" fluid>
    <v-row justify="center" align="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card elevation="2" class="pa-6">
          <!-- Back Button -->
          <v-btn icon @click="router.push('/')">
            <v-icon>mdi-arrow-left</v-icon>
          </v-btn>

          <v-card-title class="text-h5 text-center">
            Create an Account
          </v-card-title>

          <v-form @submit.prevent="submitForm" ref="formRef" v-model="valid">
            <v-text-field
              v-model="form.username"
              label="Username"
              :rules="[(v) => !!v || 'Username is required']"
              autocomplete="off"
            />
            <v-text-field
              v-model="form.name"
              label="Full Name"
              :rules="[(v) => !!v || 'Full Name is required']"
              autocomplete="off"
            />
            <v-text-field
              v-model="form.email"
              label="Email"
              type="email"
              :rules="[(v) => !!v || 'Email is required']"
              autocomplete="email"
            />
            <v-text-field
              v-model="form.password"
              label="Password"
              type="password"
              :rules="[(v) => !!v || 'Password is required']"
              autocomplete="new-password"
            />
            <v-text-field
              v-model="form.password_confirmation"
              label="Confirm Password"
              type="password"
              :rules="[(v) => !!v || 'Confirm Password is required']"
              autocomplete="new-password"
            />

            <v-btn
              :disabled="!valid || loading"
              type="submit"
              :loading="loading"
              color="primary"
              class="mt-4"
              block
            >
              Register
            </v-btn>
          </v-form>

          <div class="text-center mt-4">
            <router-link to="/login" class="text-blue-600 hover:text-blue-800">
              Already have an account? Login here
            </router-link>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { reactive, ref } from "vue";
import api from "@/axios";
import { useRouter } from "vue-router";

const router = useRouter();
const formRef = ref(null);

const form = reactive({
  username: "",
  name: "",
  email: "",
  password: "",
  password_confirmation: "",
});

const loading = ref(false);

const submitForm = async () => {
  loading.value = true;

  try {
    await api.post("api/register", form);
    alert("Registration successful! Please login.");
    router.push("/login");
  } catch (error) {
    if (error.response?.data?.errors) {
      Object.entries(error.response.data.errors).forEach(([key, messages]) => {
        errors[key] = messages.join(" ");
      });
    } else {
      alert("An error occurred, please try again.");
    }
  } finally {
    loading.value = false;
  }
};
</script>
