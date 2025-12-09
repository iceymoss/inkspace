<template>
  <div class="works">
    <div class="container">
      <!-- 筛选栏 -->
      <div class="works-filters">
        <h1 class="works-title">作品展示</h1>
        <div class="works-filters-controls">
          <el-segmented v-model="filterType" :options="typeOptions" @change="handleFilterChange" />
          <el-select v-model="sortBy" placeholder="排序方式" @change="handleFilterChange" style="width: 150px; margin-left: 15px;">
            <el-option label="默认排序" value="" />
            <el-option label="🔥 热度排序" value="hot" />
            <el-option label="⏰ 最新发布" value="time" />
            <el-option label="👁️ 最多浏览" value="view" />
            <el-option label="❤️ 最多点赞" value="like" />
          </el-select>
        </div>
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
                    {{ work.type === 'photography' ? '📷' : '💻' }}
                  </el-tag>
                </div>
              </div>
            </div>
          </div>
          
          <div class="work-info">
            <h3 class="work-title">{{ work.title }}</h3>
            <div class="work-meta">
              <div class="work-author">
                <el-avatar :size="20" :src="work.author?.avatar" />
                <span>{{ work.author?.nickname || work.author?.username }}</span>
              </div>
              <div class="work-stats">
                <span><el-icon><View /></el-icon> {{ work.view_count }}</span>
                <span v-if="work.like_count">
                  <el-icon><Star /></el-icon> {{ work.like_count }}
                </span>
                <span v-if="work.comment_count">
                  <el-icon><ChatDotRound /></el-icon> {{ work.comment_count }}
                </span>
                <span v-if="work.favorite_count">
                  <el-icon><Collection /></el-icon> {{ work.favorite_count }}
                </span>
              </div>
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
import { View, Star, ChatDotRound, Collection } from '@element-plus/icons-vue'
import api from '@/utils/api'

const router = useRouter()
const works = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const filterType = ref('all')
const sortBy = ref('time') // 默认最新发布

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
    
    if (sortBy.value) {
      params.sort = sortBy.value
    }
    
    const response = await api.get('/works', { params })
    works.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load works:', error)
  }
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadWorks()
}

onMounted(() => {
  loadWorks()
})
</script>

<style scoped>
.works {
  padding: 10px 0 40px 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.works-filters {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  margin-bottom: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  width: 100%;
}

.works-title {
  margin: 0;
  font-size: 2rem;
  flex-shrink: 0;
}

.works-filters-controls {
  display: flex;
  align-items: center;
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
  border-radius: 6px;
  overflow: hidden;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.masonry-item:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.12);
}

.work-image-container {
  position: relative;
  width: 100%;
  overflow: hidden;
  aspect-ratio: 4 / 3;
}

.work-image {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
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
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.work-title {
  font-size: 1rem;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5;
  color: var(--text-primary);
  font-weight: 500;
  min-height: 3em;
}

.work-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  padding-top: 4px;
  border-top: 1px solid #f0f0f0;
}

.work-author {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.875rem;
  color: var(--text-secondary);
  flex: 1;
  min-width: 0;
}

.work-author span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-weight: 500;
}

.work-stats {
  display: flex;
  gap: 10px;
  font-size: 0.8rem;
  color: #909399;
  flex-shrink: 0;
}

.work-stats span {
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.work-stats .el-icon {
  font-size: 0.9rem;
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
  
  .works-filters {
    flex-direction: column;
    align-items: stretch;
    gap: 15px;
  }
  
  .works-title {
    text-align: center;
  }
  
  .works-filters-controls {
    flex-direction: column;
    align-items: stretch;
    gap: 12px;
    width: 100%;
  }
  
  .works-filters-controls .el-select {
    width: 100% !important;
    margin-left: 0 !important;
  }
  
  .works-filters .el-select {
    width: 100% !important;
    margin-left: 0 !important;
  }
  
  .work-info {
    padding: 12px;
    gap: 10px;
  }
  
  .work-title {
    font-size: 0.9rem;
    min-height: 2.7em;
  }
  
  .work-meta {
    flex-direction: row;
    align-items: center;
    gap: 8px;
    padding-top: 6px;
  }
  
  .work-stats {
    gap: 8px;
    font-size: 0.75rem;
  }
  
  .work-stats span {
    gap: 2px;
  }
}
</style>

