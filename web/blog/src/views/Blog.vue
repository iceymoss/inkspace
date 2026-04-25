<template>
  <div class="blog">
    <div class="container-blog">
      <div class="blog-header">
        <div class="flex items-center gap-md flex-wrap">
          <h1>博客文章</h1>

          <div class="flex rounded-lg bg-muted p-1 ml-lg">
            <button
              v-for="tab in rankTabs"
              :key="tab.value"
              @click="rankType = tab.value; handleRankChange()"
              :class="rankType === tab.value ? 'bg-background shadow-sm text-foreground' : 'text-muted-foreground hover:text-foreground'"
              class="px-3 py-1.5 text-sm font-medium rounded-md transition-all cursor-pointer"
            >
              {{ tab.label }}
            </button>
          </div>

          <DropdownMenu>
            <DropdownMenuTrigger>
              <Button variant="outline" class="gap-2">
                <img v-if="selectedCategoryData?.logo" :src="selectedCategoryData.logo" class="w-6 h-6 rounded object-cover" />
                <span>{{ selectedCategoryData?.name || '全部分类' }}</span>
                <ChevronDown class="w-4 h-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem @click="handleCategorySelect(null)">
                <span>全部分类</span>
              </DropdownMenuItem>
              <DropdownMenuItem
                v-for="cat in categories"
                :key="cat.id"
                @click="handleCategorySelect(cat.id)"
              >
                <div class="flex items-center gap-2 min-w-[200px]">
                  <img v-if="cat.logo" :src="cat.logo" class="w-8 h-8 rounded object-cover" />
                  <span>{{ cat.name }}</span>
                  <Badge variant="secondary" class="ml-auto">{{ cat.article_count }}</Badge>
                </div>
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>

        <div class="flex items-center gap-md flex-shrink-0">
          <Select v-model="sortBy" @update:model-value="handleSortChange">
            <SelectTrigger class="w-[140px]">
              <SelectValue placeholder="排序方式" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="hot">最热排序</SelectItem>
              <SelectItem value="time">最新发布</SelectItem>
              <SelectItem value="view_count">最多浏览</SelectItem>
              <SelectItem value="like_count">最多点赞</SelectItem>
              <SelectItem value="comment_count">最多评论</SelectItem>
            </SelectContent>
          </Select>

          <div class="relative max-w-[400px]">
            <Input v-model="searchKeyword" placeholder="搜索文章..." @keyup.enter="handleSearch" class="pr-9" />
            <Button variant="ghost" size="sm" class="absolute right-0.5 top-0.5 h-9 w-9 p-0" @click="handleSearch">
              <Search class="w-4 h-4" />
            </Button>
          </div>
        </div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-4 gap-5">
        <div class="md:col-span-3">
          <div class="flex flex-col gap-sm mb-xl">
            <div
              v-for="article in articles"
              :key="article.id"
              class="card-skeuomorphic cursor-pointer h-[140px] overflow-hidden article-card"
              @click="$router.push(`/blog/${article.id}`)"
            >
              <div class="h-full p-md">
                <div class="flex gap-md h-full items-center article-content">
                  <img v-if="article.cover" :src="article.cover" class="article-cover" />
                  <div class="article-info">
                    <div class="article-header">
                      <h2>{{ article.title }}</h2>
                      <Badge v-if="article.is_top" variant="destructive" class="text-xs flex-shrink-0">置顶</Badge>
                    </div>
                    <p class="article-summary">{{ article.summary }}</p>
                    <div class="article-meta">
                      <Badge v-if="article.category" variant="secondary" class="text-xs">{{ article.category.name }}</Badge>
                      <span class="author-info">
                        <Avatar class="w-5 h-5">
                          <AvatarImage :src="article.author?.avatar" />
                        </Avatar>
                        <span>{{ article.author?.nickname || article.author?.username }}</span>
                      </span>
                      <span><Eye class="w-4 h-4" /> {{ article.view_count }}</span>
                      <span><Clock class="w-4 h-4" /> {{ formatDate(article.created_at) }}</span>
                      <span v-if="article.like_count"><Star class="w-4 h-4" /> {{ article.like_count }}</span>
                      <span v-if="article.comment_count"><MessageCircle class="w-4 h-4" /> {{ article.comment_count }}</span>
                      <span v-if="article.favorite_count"><Bookmark class="w-4 h-4" /> {{ article.favorite_count }}</span>
                      <div class="article-bookmarks" v-if="article.tags && article.tags.length > 0">
                        <span
                          v-for="tag in article.tags.slice(0, 3)"
                          :key="tag.id"
                          class="bookmark-tag"
                          :style="{ backgroundColor: getTagColor(tag.name) }"
                          @click.stop="handleTagClick(tag.id)"
                        >
                          {{ tag.name }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="flex items-center justify-center gap-2 mt-xl" v-if="totalPages > 1">
            <Button variant="outline" size="sm" :disabled="currentPage <= 1" @click="currentPage--; loadArticles()">上一页</Button>
            <span class="text-sm text-muted-foreground">{{ currentPage }} / {{ totalPages }}</span>
            <Button variant="outline" size="sm" :disabled="currentPage >= totalPages" @click="currentPage++; loadArticles()">下一页</Button>
          </div>
        </div>

        <div class="md:col-span-1">
          <Card class="shadow-md mb-lg">
            <CardContent class="p-md pt-md">
              <h3 class="mb-md text-lg">热门标签</h3>
              <div class="flex flex-wrap gap-sm">
                <Badge
                  v-for="tag in tags"
                  :key="tag.id"
                  :variant="selectedTag === tag.id ? 'default' : 'outline'"
                  class="cursor-pointer"
                  @click="filterByTag(tag.id)"
                >
                  {{ tag.name }}
                </Badge>
              </div>
            </CardContent>
          </Card>

          <div class="ad-slots" v-if="ads.length > 0">
            <div
              v-for="(ad, index) in ads"
              :key="`ad-${ad.id || index}`"
              class="ad-slot"
              @click="handleAdClick(ad)"
            >
              <div class="ad-content">
                <img :src="ad.image" :alt="ad.title" @load="recordAdView(ad.id)" />
                <div v-if="ad.title" class="ad-title">{{ ad.title }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Search, Eye, Clock, ChevronDown, Star, MessageCircle, Bookmark } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage } from '@/components/ui/avatar'
import { Input } from '@/components/ui/input'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent, DropdownMenuItem } from '@/components/ui/dropdown-menu'
import api from '@/utils/api'
import dayjs from 'dayjs'

const route = useRoute()

const articles = ref([])
const categories = ref([])
const tags = ref([])
const ads = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const searchKeyword = ref('')
const selectedCategory = ref(null)
const selectedTag = ref(null)
const rankType = ref('')
const sortBy = ref('time')

const rankTabs = [
  { label: '全部', value: '' },
  { label: '热门', value: 'hot' },
  { label: '周榜', value: 'week' },
  { label: '月榜', value: 'month' },
  { label: '年榜', value: 'year' },
]

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const selectedCategoryData = computed(() => {
  if (!selectedCategory.value) return null
  return categories.value.find(c => c.id === selectedCategory.value)
})

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

const loadArticles = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (searchKeyword.value) params.keyword = searchKeyword.value
    if (selectedCategory.value) params.category_id = selectedCategory.value
    if (selectedTag.value) params.tag_id = selectedTag.value

    if (rankType.value) {
      params.rank_type = rankType.value
    }

    if (sortBy.value) {
      params.sort_by = sortBy.value
      params.sort_order = 'desc'
    }

    const response = await api.get('/articles', { params })
    articles.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load articles:', error)
  }
}

const loadCategories = async () => {
  try {
    const response = await api.get('/categories')
    categories.value = response.data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

const loadTags = async () => {
  try {
    const response = await api.get('/tags')
    tags.value = response.data || []
  } catch (error) {
    console.error('Failed to load tags:', error)
  }
}

const loadAds = async () => {
  try {
    const response = await api.get('/ads', { params: { code: 'blog_right_ad' } })
    ads.value = response.data || []
  } catch (error) {
    console.error('Failed to load ads:', error)
    ads.value = []
  }
}

const handleAdClick = (ad) => {
  if (ad && ad.link) {
    if (ad.id) {
      api.post(`/ads/${ad.id}/click`).catch(() => {})
    }
    window.open(ad.link, '_blank')
  }
}

const recordAdView = (adId) => {
  if (adId) {
    api.post(`/ads/${adId}/view`).catch(() => {})
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadArticles()
}

const handleCategorySelect = (categoryId) => {
  selectedCategory.value = categoryId
  selectedTag.value = null
  currentPage.value = 1
  loadArticles()
}

const filterByTag = (tagId) => {
  selectedTag.value = selectedTag.value === tagId ? null : tagId
  selectedCategory.value = null
  currentPage.value = 1
  loadArticles()
}

const handleRankChange = () => {
  currentPage.value = 1
  loadArticles()
}

const handleSortChange = () => {
  currentPage.value = 1
  loadArticles()
}

const handleTagClick = (tagId) => {
  selectedTag.value = tagId
  selectedCategory.value = null
  currentPage.value = 1
  loadArticles()
}

const getTagColor = (tagName) => {
  const colors = [
    'var(--color-info)', 'var(--color-success)', 'var(--color-warning)',
    'var(--color-danger)', 'var(--color-text-tertiary)',
    'var(--color-info)', 'var(--color-success)', 'var(--color-warning)',
    'var(--color-danger)', 'var(--color-text-tertiary)'
  ]
  let hash = 0
  for (let i = 0; i < tagName.length; i++) {
    hash = tagName.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
}

watch(() => route.query, (newQuery) => {
  if (newQuery.category_id) {
    selectedCategory.value = parseInt(newQuery.category_id)
    selectedTag.value = null
  } else if (newQuery.tag_id) {
    selectedTag.value = parseInt(newQuery.tag_id)
    selectedCategory.value = null
  }
  currentPage.value = 1
  loadArticles()
}, { immediate: false })

onMounted(() => {
  loadCategories()
  loadTags()
  loadAds()

  if (route.query.category_id) {
    selectedCategory.value = parseInt(route.query.category_id)
  } else if (route.query.tag_id) {
    selectedTag.value = parseInt(route.query.tag_id)
  }

  loadArticles()
})
</script>

<style scoped>
.blog {
  padding: 0 0 var(--spacing-lg);
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.blog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  gap: var(--spacing-lg);
  padding: var(--spacing-md) 0 var(--spacing-md);
  border-bottom: 1px solid var(--theme-border-light);
}

.article-card {
  transition: box-shadow var(--transition-base), transform var(--transition-base);
}

.article-card:hover {
  box-shadow: var(--hover-lift-shadow), var(--card-inset-shadow);
  transform: translateY(var(--hover-lift));
}

.article-content {
  display: flex;
  gap: var(--spacing-md);
  height: 100%;
  align-items: center;
}

.article-cover {
  width: 140px;
  height: 112px;
  object-fit: cover;
  border-radius: var(--radius-sm);
  flex-shrink: 0;
}

.article-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
}

.article-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  margin-bottom: var(--spacing-xs);
}

.article-header h2 {
  font-size: var(--font-size-lg);
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-tight);
  margin-bottom: var(--spacing-xs);
  color: var(--theme-text-primary);
  word-break: break-word;
  max-height: 1.3em;
  flex-shrink: 0;
}

.article-summary {
  color: var(--theme-text-secondary);
  margin-bottom: var(--spacing-sm);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
  font-size: var(--font-size-sm);
  word-break: break-word;
  max-height: 3em;
  flex: 1;
  min-height: 0;
}

.article-meta {
  display: flex;
  gap: var(--spacing-md);
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  flex-wrap: wrap;
  align-items: center;
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.author-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.article-bookmarks {
  display: flex;
  flex-direction: row;
  gap: var(--spacing-xs);
  flex-wrap: wrap;
  margin-left: auto;
}

.bookmark-tag {
  position: relative;
  padding: var(--spacing-xs) var(--spacing-sm) var(--spacing-xs) var(--spacing-xs);
  color: white;
  font-size: var(--font-size-xs);
  font-weight: 500;
  line-height: var(--line-height-tight);
  cursor: pointer;
  transition: all var(--transition-slow);
  box-shadow: var(--shadow-sm);
  white-space: nowrap;
  clip-path: polygon(0 0, calc(100% - 6px) 0, 100% 50%, calc(100% - 6px) 100%, 0 100%);
}

.bookmark-tag:hover {
  transform: translateX(-2px);
  box-shadow: var(--shadow-md);
  opacity: 0.9;
}

.bookmark-tag::after {
  content: '';
  position: absolute;
  right: 0;
  top: 0;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 0 0 6px 6px;
  border-color: transparent transparent rgba(0, 0, 0, 0.15) transparent;
}

.ad-slots {
  margin-top: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.ad-slot {
  width: 100%;
  aspect-ratio: 3 / 2;
  border: 1px solid var(--theme-border-light);
  border-radius: var(--radius-sm);
  overflow: hidden;
  cursor: pointer;
  transition: all var(--transition-slow);
  background-color: var(--theme-bg-primary);
}

.ad-slot:hover {
  box-shadow: var(--shadow-sm);
  transform: translateY(-2px);
}

.ad-content {
  width: 100%;
  height: 100%;
  position: relative;
}

.ad-content img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
}

.ad-title {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
  color: white;
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-sm);
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 768px) {
  .article-content {
    flex-direction: column;
  }

  .article-cover {
    width: 100%;
  }

  .article-bookmarks {
    margin-left: 0;
    width: 100%;
    margin-top: var(--spacing-sm);
  }
}
</style>
