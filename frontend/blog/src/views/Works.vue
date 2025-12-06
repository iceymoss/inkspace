<template>
  <div class="works">
    <div class="container">
      <div class="works-header">
        <h1>作品展示</h1>
        <el-segmented v-model="filterType" :options="typeOptions" @change="handleTypeFilter" />
      </div>
      
      <!-- 瀑布流布局 -->
      <div class="masonry-grid">
        <div 
          v-for="work in works" 
          :key="work.id" 
          class="masonry-item"
          @click="router.push(`/works/${work.id}`)"
        >
          <div class="work-image-container">
            <el-image 
              :src="work.cover" 
              :alt="work.title"
              fit="cover"
              class="work-image"
              lazy
            />
            <div class="work-overlay">
              <div class="overlay-content">
                <div class="work-type-badge">
                  <el-tag :type="work.type === 'photography' ? 'warning' : 'primary'" size="small">
                    {{ work.type === 'photography' ? '📷 摄影' : '💻 项目' }}
                  </el-tag>
                  <el-tag v-if="work.type === 'photography' && work.metadata?.photo_count" type="info" size="small">
                    {{ work.metadata.photo_count }} 张
                  </el-tag>
                </div>
              </div>
            </div>
          </div>
          
          <div class="work-info">
            <h3 class="work-title">{{ work.title }}</h3>
            <div class="work-author">
              <el-avatar :size="24" :src="work.author?.avatar" />
              <span>{{ work.author?.nickname || work.author?.username }}</span>
            </div>
            <div class="work-stats">
              <span><el-icon><View /></el-icon> {{ work.view_count }}</span>
              <span v-if="work.comment_count > 0">
                <el-icon><ChatDotRound /></el-icon> {{ work.comment_count }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadWorks"
        />
      </div>

      <el-empty v-if="works.length === 0" description="暂无作品" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { View, ChatDotRound } from '@element-plus/icons-vue'
import api from '@/utils/api'

const router = useRouter()
const works = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const filterType = ref('all')

const typeOptions = [
  { label: '全部', value: 'all' },
  { label: '💻 项目', value: 'project' },
  { label: '📷 摄影', value: 'photography' }
]

const loadWorks = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      status: 1
    }
    
    if (filterType.value !== 'all') {
      params.type = filterType.value
    }
    
    const response = await api.get('/works', { params })
    works.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load works:', error)
  }
}

const handleTypeFilter = () => {
  currentPage.value = 1
  loadWorks()
}

onMounted(() => {
  loadWorks()
})
</script>

<style scoped>
.works {
  padding: 40px 0;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.works-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.works-header h1 {
  margin: 0;
  font-size: 2rem;
}

/* 瀑布流布局 */
.masonry-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}

.masonry-item {
  cursor: pointer;
  background: white;
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.masonry-item:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.15);
}

.work-image-container {
  position: relative;
  width: 100%;
  overflow: hidden;
}

.work-image {
  width: 100%;
  height: auto;
  min-height: 200px;
  display: block;
}

.work-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(to bottom, rgba(0,0,0,0.3) 0%, transparent 40%, transparent 60%, rgba(0,0,0,0.3) 100%);
  opacity: 0;
  transition: opacity 0.3s;
}

.masonry-item:hover .work-overlay {
  opacity: 1;
}

.overlay-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  padding: 15px;
}

.work-type-badge {
  display: flex;
  gap: 5px;
}

.work-info {
  padding: 15px;
}

.work-title {
  font-size: 1rem;
  margin: 0 0 10px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.4;
}

.work-author {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.work-stats {
  display: flex;
  gap: 15px;
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.work-stats span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.pagination {
  margin-top: 40px;
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .masonry-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 15px;
  }
  
  .works-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 15px;
  }
}
</style>

