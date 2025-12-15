<template>
  <div class="work-detail">
    <div class="container">
      <el-card v-if="work">
        <h1>{{ work.title }}</h1>
        
        <div class="work-meta">
          <span><el-icon><View /></el-icon> {{ work.view_count }} 浏览</span>
          <span><el-icon><Clock /></el-icon> {{ formatDate(work.created_at) }}</span>
        </div>

        <div class="work-description">
          <p>{{ work.description }}</p>
        </div>

        <div class="work-images">
          <el-image
            v-for="(image, index) in work.images"
            :key="index"
            :src="image"
            :preview-src-list="work.images"
            fit="cover"
            class="work-image"
          />
        </div>

        <div v-if="work.link" class="work-link">
          <el-button type="primary" :icon="Link" @click="openLink">访问项目</el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Link } from '@element-plus/icons-vue'
import api from '@/utils/api'
import dayjs from 'dayjs'

const route = useRoute()
const work = ref(null)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

const loadWork = async () => {
  try {
    const response = await api.get(`/works/${route.params.id}`)
    work.value = response.data
  } catch (error) {
    ElMessage.error('作品加载失败')
  }
}

const openLink = () => {
  if (work.value?.link) {
    window.open(work.value.link, '_blank')
  }
}

onMounted(() => {
  loadWork()
})
</script>

<style scoped>
.work-detail {
  padding: 40px 0;
}

.work-detail h1 {
  margin-bottom: 20px;
}

.work-meta {
  display: flex;
  gap: 20px;
  color: var(--text-secondary);
  margin-bottom: 30px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--border-lighter);
}

.work-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.work-description {
  margin-bottom: 30px;
  font-size: 16px;
  line-height: 1.8;
}

.work-images {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.work-image {
  width: 100%;
  height: 300px;
  border-radius: 8px;
  cursor: pointer;
}

.work-link {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid var(--border-lighter);
}
</style>

