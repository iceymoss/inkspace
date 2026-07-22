import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import App from './App.vue'
import router from './router'
import './assets/main.css'
import './themes/base.css'
import './themes/magazine.css'
import './themes/terminal.css'
import './themes/cozy.css'
import './themes/swiss.css'
import { bootstrapCachedAppearance, useAppearanceStore } from './stores/appearance'

bootstrapCachedAppearance()
const app = createApp(App)
const pinia = createPinia()

// Register Element Plus Icons
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app.use(pinia)
app.use(router)
app.use(ElementPlus)

useAppearanceStore(pinia).initialize()

app.mount('#app')
