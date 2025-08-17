import { ref } from "vue";

const toast = ref({
  show: false,
  message: "",
  color: "success",
});

export function useToast() {
  const showToast = (message, color = "success") => {
    toast.value.message = message;
    toast.value.color = color;
    toast.value.show = true;
  };

  return {
    toast,
    showToast,
  };
}