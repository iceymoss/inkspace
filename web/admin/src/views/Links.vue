<template>
  <div class="links-page">
    <div class="container">
      <h1>友情链接</h1>
      <p class="subtitle">与优秀的站点互相链接，共同成长</p>

      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8" :lg="6" v-for="link in links" :key="link.id">
          <el-card class="link-card" shadow="hover" @click="openLink(link.url)">
            <div class="link-content">
              <el-avatar :size="60" :src="link.logo" :alt="link.name">
                {{ link.name.charAt(0) }}
              </el-avatar>
              <h3>{{ link.name }}</h3>
              <p class="link-description">{{ link.description }}</p>
              <Link class="link-icon h-4 w-4" />
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-empty v-if="links.length === 0" description="暂无友情链接" />

      <div class="apply-section">
        <el-divider />
        <h3>申请友链</h3>
        <p>如果你也想交换友链，请联系我：</p>
        <p>📧 Email: admin@example.com</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Link } from 'lucide-vue-next'
import { toast } from 'vue-sonner'
import api from '@/utils/api'

const links = ref([])

const loadLinks = async () => {
  try {
    const response = await api.get('/links', { params: { status: 1 } })
    links.value = response.data || []
  } catch (error) {
    toast.error('加载失败')
  }
}

const openLink = (url) => {
  window.open(url, '_blank')
}

onMounted(() => {
  loadLinks()
})
</script>

<style scoped>
.links-page {
  padding: 40px 0;
}

.links-page h1 {
  text-align: center;
  margin-bottom: 10px;
}

.subtitle {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: 40px;
}

.link-card {
  cursor: pointer;
  margin-bottom: 20px;
  transition: transform 0.3s;
  height: 100%;
}

.link-card:hover {
  transform: translateY(-5px);
}

.link-content {
  text-align: center;
  position: relative;
}

.link-content h3 {
  margin: 15px 0 10px;
  font-size: 1.1rem;
}

.link-description {
  color: var(--text-secondary);
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  min-height: 40px;
}

.link-icon {
  position: absolute;
  top: 5px;
  right: 5px;
  color: var(--primary-color);
}

.apply-section {
  margin-top: 60px;
  text-align: center;
  color: var(--text-secondary);
}

.apply-section h3 {
  margin-bottom: 15px;
  color: var(--text-primary);
}

.apply-section p {
  margin: 5px 0;
}
</style>

