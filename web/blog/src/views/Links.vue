<template>
  <div class="links-page">
    <div class="container">
      <div class="page-kicker">
        LINKS · DIRECTORY
      </div>
      <h1>友情链接</h1>
      <p class="subtitle">
        与优秀的站点互相链接，共同成长
      </p>

      <el-row :gutter="20">
        <el-col
          v-for="link in links"
          :key="link.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <el-card
            class="link-card"
            shadow="hover"
            tabindex="0"
            @click="openLink(link.url)"
            @keyup.enter="openLink(link.url)"
          >
            <div class="link-content">
              <el-avatar
                :size="60"
                :src="link.logo"
                :alt="link.name"
              >
                {{ link.name.charAt(0) }}
              </el-avatar>
              <h3>{{ link.name }}</h3>
              <p class="link-description">
                {{ link.description }}
              </p>
              <el-icon class="link-icon">
                <Link />
              </el-icon>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-empty
        v-if="links.length === 0"
        description="暂无友情链接"
      />

      <div class="apply-section">
        <el-divider />
        <h3>{{ applyTitle || '申请友链' }}</h3>
        <p v-if="applyDescription">
          {{ applyDescription }}
        </p>
        <p v-if="applyEmail">
          📧 Email: {{ applyEmail }}
        </p>
        <p
          v-else-if="applyDescription"
          style="color: var(--text-secondary);"
        >
          📧 Email: 暂未配置
        </p>
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
  padding: 40px 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
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
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
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

/* Magazine adaptation */
.links-page { padding: 62px 0 80px; background: var(--theme-bg-primary); }
.links-page .container { max-width: 1060px; padding: 0 32px; }
.page-kicker { margin-bottom: 12px; color: var(--theme-primary); font-family: Georgia, 'Songti SC', serif; font-size: 11px; letter-spacing: .26em; text-align: center; }
.links-page h1 { font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif; font-size: clamp(38px, 6vw, 58px); font-weight: 500; letter-spacing: .05em; }
.subtitle { margin-bottom: 52px; letter-spacing: .06em; }
.link-card { margin-bottom: 20px; border: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; transition: border-color .25s ease, transform .25s ease; }
.link-card:hover { border-color: var(--theme-primary); box-shadow: none; transform: translateY(-4px); }
.link-card:focus-visible { outline: 2px solid var(--theme-primary); outline-offset: 3px; }
.link-content h3 { font-family: Georgia, 'Songti SC', serif; font-size: 19px; font-weight: 500; letter-spacing: .04em; }
.link-content :deep(.el-avatar) { border: 1px solid var(--theme-border); }
.link-icon { color: var(--theme-primary); }
.apply-section { margin-top: 64px; padding-top: 36px; border-top: 1px solid var(--theme-border); }
.apply-section :deep(.el-divider) { display: none; }
.apply-section h3 { font-family: Georgia, 'Songti SC', serif; font-size: 24px; font-weight: 500; }
@media (max-width: 900px) { .links-page .container { padding: 0 24px; } }
@media (max-width: 560px) { .links-page { padding: 38px 0 56px; } .links-page .container { padding: 0 18px; } .subtitle { margin-bottom: 36px; } }
</style>

