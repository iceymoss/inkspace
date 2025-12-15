<template>
  <div class="blog">
    <div class="container">
      <div class="blog-header">
        <h1>博客文章</h1>
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
            <h3>分类</h3>
            <div class="category-list">
              <el-tag
                v-for="category in categories"
                :key="category.id"
                :type="selectedCategory === category.id ? 'primary' : 'info'"
                class="category-tag"
                @click="filterByCategory(category.id)"
              >
                {{ category.name }}
              </el-tag>
            </div>
          </el-card>

          <el-card class="sidebar-card">
            <h3>标签</h3>
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
import { ref, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
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

const filterByCategory = (categoryId) => {
  selectedCategory.value = selectedCategory.value === categoryId ? null : categoryId
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
  padding: 40px 0;
}

.blog-header {
  margin-bottom: 30px;
  text-align: center;
}

.blog-header h1 {
  margin-bottom: 20px;
}

.search-input {
  max-width: 500px;
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

.category-list, .tag-list {
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

