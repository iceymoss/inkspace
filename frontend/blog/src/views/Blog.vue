<template>
  <div class="blog">
    <div class="container">
      <div class="blog-header">
        <div class="header-left">
          <h1>博客文章</h1>
          
          <!-- 榜单选择器 -->
          <el-radio-group v-model="rankType" @change="handleRankChange" size="default" class="rank-selector">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="hot">热门</el-radio-button>
            <el-radio-button label="week">周榜</el-radio-button>
            <el-radio-button label="month">月榜</el-radio-button>
            <el-radio-button label="year">年榜</el-radio-button>
          </el-radio-group>
          
          <!-- 分类选择器 -->
          <el-dropdown trigger="click" @command="handleCategorySelect" class="category-dropdown">
            <el-button>
              <el-image 
                v-if="selectedCategoryData?.logo" 
                :src="selectedCategoryData.logo" 
                class="category-logo-small"
                fit="cover"
              />
              <span>{{ selectedCategoryData?.name || '全部分类' }}</span>
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :command="null">
                  <span>全部分类</span>
                </el-dropdown-item>
                <el-dropdown-item 
                  v-for="cat in categories" 
                  :key="cat.id" 
                  :command="cat.id"
                >
                  <div class="dropdown-category-item">
                    <el-image 
                      v-if="cat.logo" 
                      :src="cat.logo" 
                      class="category-logo-dropdown"
                      fit="cover"
                    />
                    <span>{{ cat.name }}</span>
                    <el-badge :value="cat.article_count" class="category-badge" />
                  </div>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        
        <div class="header-right">
          <!-- 排序选择器 -->
          <el-select v-model="sortBy" @change="handleSortChange" placeholder="排序方式" class="sort-select" size="default">
            <el-option label="最热排序" value="hot" />
            <el-option label="最新发布" value="time" />
            <el-option label="最多浏览" value="view_count" />
            <el-option label="最多点赞" value="like_count" />
            <el-option label="最多评论" value="comment_count" />
          </el-select>
          
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文章..."
            class="search-input"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button :icon="Search" @click="handleSearch" />
            </template>
          </el-input>
        </div>
      </div>

      <el-row :gutter="20">
        <el-col :xs="24" :md="18">
          <div class="article-list">
            <el-card
              v-for="article in articles"
              :key="article.id"
              class="article-item"
              shadow="hover"
              @click="$router.push(`/blog/${article.id}`)"
            >
              <div class="article-content">
                <img v-if="article.cover" :src="article.cover" class="article-cover" />
                <div class="article-info">
                  <div class="article-header">
                    <h2>{{ article.title }}</h2>
                    <el-tag v-if="article.is_top" type="danger" size="small">置顶</el-tag>
                  </div>
                  <p class="article-summary">{{ article.summary }}</p>
                  <div class="article-meta">
                    <el-tag v-if="article.category" size="small">{{ article.category.name }}</el-tag>
                    <span class="author-info">
                      <el-avatar :size="20" :src="article.author?.avatar" />
                      <span>{{ article.author?.nickname || article.author?.username }}</span>
                    </span>
                    <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                    <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
                    <!-- 书签样式的标签 -->
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
            </el-card>
          </div>

          <div class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :total="total"
              layout="prev, pager, next"
              @current-change="loadArticles"
            />
          </div>
        </el-col>

        <el-col :xs="24" :md="6">
          <el-card class="sidebar-card">
            <h3>热门标签</h3>
            <div class="tag-list">
              <el-tag
                v-for="tag in tags"
                :key="tag.id"
                :type="selectedTag === tag.id ? 'primary' : ''"
                size="small"
                class="tag-item"
                @click="filterByTag(tag.id)"
              >
                {{ tag.name }}
              </el-tag>
            </div>
          </el-card>

          <!-- 广告位 -->
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
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { Search, User, View, Clock, ArrowDown } from '@element-plus/icons-vue'
import api from '@/utils/api'
import dayjs from 'dayjs'

const route = useRoute()

const articles = ref([])
const categories = ref([])
const tags = ref([])
const ads = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const selectedCategory = ref(null)
const selectedTag = ref(null)
const rankType = ref('') // 默认全部
const sortBy = ref('time') // 默认最新发布

// 当前选中的分类数据
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
    
    // 榜单类型
    if (rankType.value) {
      params.rank_type = rankType.value
    }
    
    // 排序方式（榜单模式下也支持排序）
    if (sortBy.value) {
      params.sort_by = sortBy.value
      params.sort_order = 'desc' // 默认降序
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
    // 加载博客列表右侧广告位，根据API返回的数量动态显示
    const response = await api.get('/ads', { params: { code: 'blog_right_ad' } })
    ads.value = response.data || []
  } catch (error) {
    console.error('Failed to load ads:', error)
    ads.value = []
  }
}

const handleAdClick = (ad) => {
  if (ad && ad.link) {
    // 记录点击
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

// 根据标签名称生成颜色（用于书签样式）
const getTagColor = (tagName) => {
  const colors = [
    '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399',
    '#409EFF', '#67C23A', '#E6A23C', '#F56C6C', '#909399'
  ]
  // 根据标签名称的hash值选择颜色
  let hash = 0
  for (let i = 0; i < tagName.length; i++) {
    hash = tagName.charCodeAt(i) + ((hash << 5) - hash)
  }
  return colors[Math.abs(hash) % colors.length]
}

// 监听 URL 参数变化，自动更新筛选条件
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
  
  // 从 URL 参数初始化筛选条件
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
  padding: 0 0 20px;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.blog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 20px;
  padding: 12px 0 16px;
  border-bottom: 1px solid #ebeef5;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 15px;
}

.blog-header h1 {
  font-size: 2rem;
  color: var(--theme-text-primary);
  margin: 0;
}

.category-dropdown {
  flex-shrink: 0;
}

.category-button-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.category-logo-small {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  flex-shrink: 0;
}

.dropdown-category-item {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 200px;
}

.category-logo-dropdown {
  width: 32px;
  height: 32px;
  border-radius: 4px;
  flex-shrink: 0;
}

.category-badge {
  margin-left: auto;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.rank-selector {
  margin-left: 20px;
}

.sort-select {
  width: 140px;
}

.search-input {
  max-width: 400px;
}

.article-list {
  margin-bottom: 30px;
  display: flex;
  flex-direction: column;
  gap: 8px !important;
}

.article-list :deep(.el-card) {
  margin-bottom: 0 !important;
}

.article-item {
  margin-bottom: 0 !important;
  cursor: pointer;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
  height: 140px;
  overflow: hidden;
  position: relative;
}

.article-item :deep(.el-card__body) {
  height: 100%;
  padding: 14px;
  display: flex;
  flex-direction: column;
}

.sidebar-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.article-content {
  display: flex;
  gap: 14px;
  height: 100%;
  align-items: center; /* 垂直居中对齐 */
}

.article-cover {
  width: 140px;
  height: 112px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.article-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0; /* 允许flex收缩 */
  min-height: 0;
  overflow: hidden; /* 防止内容溢出 */
}

.article-header {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 6px;
}

.article-header h2 {
  font-size: 1.2rem;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: 1.3;
  margin-bottom: 6px;
  color: var(--theme-text-primary);
  word-break: break-word;
  max-height: 1.3em; /* 确保只显示1行 */
  flex-shrink: 0; /* 防止标题被压缩 */
}

.article-summary {
  color: var(--theme-text-secondary);
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5; /* 使用标准行高 */
  font-size: 13px;
  word-break: break-word;
  max-height: 3em; /* 确保只显示2行 (1.5 * 2 = 3em) */
  flex: 1;
  min-height: 0; /* 允许flex收缩 */
}

.article-meta {
  display: flex;
  gap: 12px;
  color: var(--theme-text-secondary);
  font-size: 13px;
  flex-wrap: wrap;
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 6px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

.sidebar-card {
  margin-bottom: 20px;
}

.sidebar-card h3 {
  margin-bottom: 15px;
  font-size: 1.1rem;
}

.category-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.category-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
  border: 1px solid #ebeef5;
}

.category-item:hover {
  background-color: var(--theme-bg-secondary);
  transform: translateX(5px);
}

.category-item.active {
  background-color: var(--theme-bg-hover);
  border-color: var(--theme-primary);
  color: var(--theme-primary);
}

.category-logo {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  flex-shrink: 0;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.category-tag, .tag-item {
  cursor: pointer;
}

/* 书签样式标签 */
.article-bookmarks {
  display: flex;
  flex-direction: row; /* 横向排列 */
  gap: 4px; /* 标签之间的间距 */
  flex-wrap: wrap;
  margin-left: auto; /* 自动推到右侧 */
}

.bookmark-tag {
  position: relative;
  padding: 3px 10px 3px 6px;
  color: white;
  font-size: 11px;
  font-weight: 500;
  line-height: 1.3;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  white-space: nowrap; /* 防止文字换行 */
  /* 书签折角效果 */
  clip-path: polygon(0 0, calc(100% - 6px) 0, 100% 50%, calc(100% - 6px) 100%, 0 100%);
}

.bookmark-tag:hover {
  transform: translateX(-2px);
  box-shadow: 0 3px 6px rgba(0, 0, 0, 0.3);
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

/* 广告位样式 */
.ad-slots {
  margin-top: 20px;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ad-slot {
  width: 100%;
  aspect-ratio: 3 / 2; /* 宽度:高度 = 3:2，即高度是宽度的 2/3 */
  border: 1px solid #ebeef5;
  border-radius: 4px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  background-color: var(--theme-bg-primary);
}

.ad-slot:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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
  padding: 8px 12px;
  font-size: 12px;
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
    margin-left: 0; /* 移动端取消自动右对齐 */
    width: 100%; /* 占满整行 */
    margin-top: 8px; /* 添加顶部间距 */
  }
}
</style>

