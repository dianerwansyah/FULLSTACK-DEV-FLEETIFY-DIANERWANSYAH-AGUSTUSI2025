import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import "./style.css";
import "vuetify/styles";
import "@mdi/font/css/materialdesignicons.css";
import { createVuetify } from "vuetify";
import * as components from "vuetify/components";
import * as directives from "vuetify/directives";
import { createPinia } from "pinia";


const vuetify = createVuetify({
  components,
  directives,
  icons: {
    defaultSet: "mdi",
  },
});

const app = createApp(App);
const pinia = createPinia();

app.use(pinia); 
app.use(vuetify);
app.use(router);
app.mount("#app");