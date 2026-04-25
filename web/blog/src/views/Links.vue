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
              <el-icon class="link-icon"><Link /></el-icon>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-empty v-if="links.length === 0" description="暂无友情链接" />

      <div class="apply-section">
        <el-divider />
        <h3>{{ applyTitle || '申请友链' }}</h3>
        <p v-if="applyDescription">{{ applyDescription }}</p>
        <p v-if="applyEmail">📧 Email: {{ applyEmail }}</p>
        <p v-else-if="applyDescription" style="color: var(--theme-text-secondary);">📧 Email: 暂未配置</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Link } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '@/utils/api'

const links = ref([])
const applyTitle = ref('')
const applyDescription = ref('')
const applyEmail = ref('')

const loadLinks = async () => {
  try {
    const response = await api.get('/links', { params: { status: 1 } })
    links.value = response.data || []
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const loadSettings = async () => {
  try {
    const response = await api.get('/settings/public')
    applyTitle.value = response.data?.link_apply_title || ''
    applyDescription.value = response.data?.link_apply_description || ''
    applyEmail.value = response.data?.link_apply_email || ''
  } catch (error) {
    console.error('加载设置失败:', error)
  }
}

const openLink = (url) => {
  window.open(url, '_blank')
}

onMounted(() => {
  loadLinks()
  loadSettings()
})
</script>

<style scoped>
.links-page {
  padding: var(--spacing-xl) 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.links-page h1 {
  text-align: center;
  margin-bottom: var(--spacing-sm);
  font-size: var(--font-size-2xl);
  color: var(--theme-text-primary);
}

.subtitle {
  text-align: center;
  color: var(--theme-text-secondary);
  margin-bottom: var(--spacing-xl);
  font-size: var(--font-size-base);
  line-height: var(--line-height-base);
}

.link-card {
  cursor: pointer;
  margin-bottom: var(--spacing-lg);
  transition: transform var(--transition-slow), box-shadow var(--transition-base);
  height: 100%;
  box-shadow: var(--shadow-sm);
  border-radius: var(--radius-lg);
}

.link-card:hover {
  transform: translateY(-5px);
  box-shadow: var(--shadow-md);
}

.link-content {
  text-align: center;
  position: relative;
}

.link-content h3 {
  margin: var(--spacing-md) 0 var(--spacing-sm);
  font-size: var(--font-size-lg);
  color: var(--theme-text-primary);
}

.link-description {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  line-height: var(--line-height-base);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  min-height: 40px;
}

.link-icon {
  position: absolute;
  top: var(--spacing-xs);
  right: var(--spacing-xs);
  color: var(--theme-primary);
}

.apply-section {
  margin-top: calc(var(--spacing-xl) + var(--spacing-lg));
  text-align: center;
  color: var(--theme-text-secondary);
}

.apply-section h3 {
  margin-bottom: var(--spacing-md);
  color: var(--theme-text-primary);
  font-size: var(--font-size-lg);
}

.apply-section p {
  margin: var(--spacing-xs) 0;
  line-height: var(--line-height-base);
}
</style>

