<template>
  <div class="works">
    <div class="container">
      <!-- 筛选栏 -->
      <div class="works-filters">
        <h1 class="works-title">作品分享</h1>
        <div class="works-filters-controls">
          <el-segmented v-model="filterType" :options="typeOptions" @change="handleFilterChange" />
          <el-select v-model="sortBy" placeholder="排序方式" @change="handleFilterChange" class="works-sort-select">
            <el-option label="默认排序" value="" />
            <el-option label="🔥 热度排序" value="hot" />
            <el-option label="⏰ 最新发布" value="time" />
            <el-option label="👁️ 最多浏览" value="view" />
            <el-option label="❤️ 最多点赞" value="like" />
          </el-select>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索作品（标题/描述）"
            clearable
            class="works-search-input"
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <el-icon class="search-icon" @click.stop="handleSearch"><Search /></el-icon>
            </template>
          </el-input>
        </div>
      </div>
      
      <!-- 瀑布流布局 -->
      <div class="masonry-grid">
        <div 
          v-for="work in works" 
          :key="work.id" 
          class="masonry-item"
          @click="handleWorkClick(work.id)"
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
import { View, Star, ChatDotRound, Collection, Search } from '@element-plus/icons-vue'
import api from '@/utils/api'
import { navigateToWorkDetail } from '@/utils/workNavigation'

const router = useRouter()
const works = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const filterType = ref('all')
const sortBy = ref('time') // 默认最新发布
const searchKeyword = ref('')

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

    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
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

const handleSearch = () => {
  currentPage.value = 1
  loadWorks()
}

// 处理作品点击，预加载数据后跳转
const handleWorkClick = (workId) => {
  navigateToWorkDetail(workId, router)
}

onMounted(() => {
  loadWorks()
})
</script>

<style scoped>
.works {
  padding: var(--spacing-sm) 0 var(--spacing-xl) 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.works-filters {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  margin-bottom: var(--spacing-md);
  background: var(--theme-bg-card);
  border-radius: var(--radius-md);
  box-shadow: var(--shadow-sm);
  width: 100%;
}

.works-title {
  margin: 0;
  font-size: var(--font-size-3xl);
  flex-shrink: 0;
  color: var(--theme-text-primary);
  font-weight: 700;
}

.works-filters-controls {
  display: flex;
  align-items: center;
}

.works-sort-select {
  width: 150px;
  margin-left: var(--spacing-md);
}

.search-icon {
  cursor: pointer;
}

.masonry-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-xl);
}

.masonry-item {
  cursor: pointer;
  background: var(--theme-bg-card);
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-sm);
  transition: all var(--transition-slow);
}

.masonry-item:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-md);
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
  transition: opacity var(--transition-slow);
}

.masonry-item:hover .work-overlay {
  opacity: 1;
}

.overlay-content {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  padding: var(--spacing-md);
}

.work-type-badge {
  display: flex;
  gap: var(--spacing-xs);
}

.work-info {
  padding: var(--spacing-md);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.work-title {
  font-size: var(--font-size-base);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
  color: var(--theme-text-primary);
  font-weight: 500;
  min-height: 3em;
}

.work-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: var(--spacing-sm);
  padding-top: var(--spacing-xs);
  border-top: 1px solid var(--theme-border-light);
}

.work-author {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
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
  gap: var(--spacing-sm);
  font-size: var(--font-size-xs);
  color: var(--theme-text-tertiary);
  flex-shrink: 0;
}

.work-stats span {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  white-space: nowrap;
}

.work-stats .el-icon {
  font-size: var(--font-size-sm);
}

.pagination {
  margin-top: var(--spacing-xl);
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .masonry-grid {
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: var(--spacing-md);
  }

  .works-filters {
    flex-direction: column;
    align-items: stretch;
    gap: var(--spacing-md);
  }

  .works-title {
    text-align: center;
  }

  .works-filters-controls {
    flex-direction: column;
    align-items: stretch;
    gap: var(--spacing-sm);
    width: 100%;
  }

  .works-filters-controls .el-select,
  .works-filters .el-select {
    width: 100%;
    margin-left: 0;
  }

  .works-sort-select {
    width: 100%;
    margin-left: 0;
  }

  .work-info {
    padding: var(--spacing-sm);
    gap: var(--spacing-sm);
  }

  .work-title {
    font-size: var(--font-size-sm);
    min-height: 2.7em;
  }

  .work-meta {
    flex-direction: row;
    align-items: center;
    gap: var(--spacing-sm);
    padding-top: 6px;
  }

  .work-stats {
    gap: var(--spacing-sm);
    font-size: var(--font-size-xs);
  }

  .work-stats span {
    gap: 2px;
  }
}
</style>

