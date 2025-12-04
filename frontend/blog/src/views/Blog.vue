<template>
  <div class="blog">
    <div class="container">
      <div class="blog-header">
        <div class="header-left">
          <h1>博客文章</h1>
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
                    <span><el-icon><User /></el-icon> {{ article.author?.nickname || article.author?.username }}</span>
                    <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                    <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
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
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { Search, User, View, Clock, ArrowDown } from '@element-plus/icons-vue'
import api from '@/utils/api'
import dayjs from 'dayjs'

const articles = ref([])
const categories = ref([])
const tags = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const searchKeyword = ref('')
const selectedCategory = ref(null)
const selectedTag = ref(null)

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

onMounted(() => {
  loadArticles()
  loadCategories()
  loadTags()
})
</script>

<style scoped>
.blog {
  padding: 0 0 20px;
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
  color: var(--text-primary);
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

.search-input {
  max-width: 400px;
}

.article-list {
  margin-bottom: 30px;
}

.article-item {
  margin-bottom: 20px;
  cursor: pointer;
}

.article-content {
  display: flex;
  gap: 20px;
}

.article-cover {
  width: 200px;
  height: 150px;
  object-fit: cover;
  border-radius: 4px;
  flex-shrink: 0;
}

.article-info {
  flex: 1;
}

.article-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.article-header h2 {
  font-size: 1.5rem;
  margin: 0;
}

.article-summary {
  color: var(--text-secondary);
  margin-bottom: 15px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.article-meta {
  display: flex;
  gap: 15px;
  color: var(--text-secondary);
  font-size: 14px;
  flex-wrap: wrap;
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
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
  background-color: #f5f7fa;
  transform: translateX(5px);
}

.category-item.active {
  background-color: #ecf5ff;
  border-color: #409eff;
  color: #409eff;
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

@media (max-width: 768px) {
  .article-content {
    flex-direction: column;
  }

  .article-cover {
    width: 100%;
  }
}
</style>

