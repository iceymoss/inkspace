<template>
  <div class="works">
    <div class="container">
      <h1>我的作品</h1>
      
      <el-row :gutter="20">
        <el-col :xs="24" :sm="12" :md="8" v-for="work in works" :key="work.id">
          <el-card class="work-card" shadow="hover" @click="$router.push(`/works/${work.id}`)">
            <img :src="work.cover" class="work-cover" />
            <div class="work-info">
              <h3>{{ work.title }}</h3>
              <p class="work-description">{{ work.description }}</p>
              <div class="work-meta">
                <span><el-icon><View /></el-icon> {{ work.view_count }}</span>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadWorks"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const works = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)

const loadWorks = async () => {
  try {
    const response = await api.get('/works', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        status: 1
      }
    })
    works.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load works:', error)
  }
}

onMounted(() => {
  loadWorks()
})
</script>

<style scoped>
.works {
  padding: 40px 0;
}

.works h1 {
  text-align: center;
  margin-bottom: 40px;
}

.work-card {
  cursor: pointer;
  margin-bottom: 20px;
  transition: transform 0.3s;
}

.work-card:hover {
  transform: translateY(-5px);
}

.work-cover {
  width: 100%;
  height: 250px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 15px;
}

.work-info h3 {
  margin-bottom: 10px;
  font-size: 1.25rem;
}

.work-description {
  color: var(--text-secondary);
  margin-bottom: 15px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.work-meta {
  display: flex;
  gap: 15px;
  color: var(--text-secondary);
  font-size: 14px;
}

.work-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.pagination {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}
</style>

