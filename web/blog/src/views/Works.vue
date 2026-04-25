<template>
  <div class="works">
    <div class="container">
      <div class="works-filters">
        <h1 class="works-title">作品分享</h1>
        <div class="works-filters-controls">
          <div class="flex rounded-lg bg-muted p-1">
            <button
              v-for="opt in typeOptions"
              :key="opt.value"
              @click="filterType = opt.value; handleFilterChange()"
              :class="filterType === opt.value ? 'bg-background shadow-sm text-foreground' : 'text-muted-foreground hover:text-foreground'"
              class="px-3 py-1.5 text-sm font-medium rounded-md transition-all cursor-pointer"
            >
              {{ opt.label }}
            </button>
          </div>

          <Select v-model="sortBy" @update:model-value="handleFilterChange" class="works-sort-select">
            <SelectTrigger class="w-[150px]">
              <SelectValue placeholder="排序方式" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="">默认排序</SelectItem>
              <SelectItem value="hot">🔥 热度排序</SelectItem>
              <SelectItem value="time">⏰ 最新发布</SelectItem>
              <SelectItem value="view">👁️ 最多浏览</SelectItem>
              <SelectItem value="like">❤️ 最多点赞</SelectItem>
            </SelectContent>
          </Select>

          <div class="relative">
            <Input
              v-model="searchKeyword"
              placeholder="搜索作品（标题/描述）"
              @keyup.enter="handleSearch"
              class="pr-9"
            />
            <Button variant="ghost" size="sm" class="absolute right-0.5 top-0.5 h-9 w-9 p-0" @click="handleSearch">
              <Search class="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-lg">
        <div
          v-for="work in works"
          :key="work.id"
          class="card-skeuomorphic cursor-pointer overflow-hidden"
          @click="handleWorkClick(work.id)"
        >
          <div class="work-image-container">
            <img
              :src="work.cover"
              :alt="work.title"
              class="work-image"
              loading="lazy"
            />
            <div class="work-overlay">
              <div class="overlay-content">
                <div class="work-type-badge">
                  <Badge :variant="work.type === 'photography' ? 'secondary' : 'default'" class="text-xs">
                    {{ work.type === 'photography' ? '📷' : '💻' }}
                  </Badge>
                </div>
              </div>
            </div>
          </div>

          <div class="work-info">
            <h3 class="work-title">{{ work.title }}</h3>
            <div class="work-meta">
              <div class="work-author">
                <Avatar class="w-5 h-5">
                  <AvatarImage :src="work.author?.avatar" />
                </Avatar>
                <span>{{ work.author?.nickname || work.author?.username }}</span>
              </div>
              <div class="work-stats">
                <span><Eye class="w-4 h-4" /> {{ work.view_count }}</span>
                <span v-if="work.like_count">
                  <Star class="w-4 h-4" /> {{ work.like_count }}
                </span>
                <span v-if="work.comment_count">
                  <MessageCircle class="w-4 h-4" /> {{ work.comment_count }}
                </span>
                <span v-if="work.favorite_count">
                  <Bookmark class="w-4 h-4" /> {{ work.favorite_count }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-center gap-2 mt-xl" v-if="totalPages > 1">
        <Button variant="outline" size="sm" :disabled="currentPage <= 1" @click="currentPage--; loadWorks()">上一页</Button>
        <span class="text-sm text-muted-foreground">{{ currentPage }} / {{ totalPages }}</span>
        <Button variant="outline" size="sm" :disabled="currentPage >= totalPages" @click="currentPage++; loadWorks()">下一页</Button>
      </div>

      <EmptyState v-if="works.length === 0" title="暂无作品" description="还没有发布任何作品" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Eye, Star, MessageCircle, Bookmark } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Input } from '@/components/ui/input'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import { EmptyState } from '@/components/ui/empty-state'
import api from '@/utils/api'
import { navigateToWorkDetail } from '@/utils/workNavigation'

const router = useRouter()
const works = ref([])
const currentPage = ref(1)
const pageSize = ref(12)
const total = ref(0)
const filterType = ref('all')
const sortBy = ref('time')
const searchKeyword = ref('')

const typeOptions = [
  { label: '全部', value: 'all' },
  { label: '💻 项目', value: 'project' },
  { label: '📷 摄影', value: 'photography' }
]

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

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
  gap: var(--spacing-md);
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

.card-skeuomorphic:hover .work-overlay {
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

@media (max-width: 768px) {
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

  .works-filters-controls > * {
    width: 100%;
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
