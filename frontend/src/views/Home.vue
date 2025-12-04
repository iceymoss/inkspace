<template>
  <div class="home">
    <section class="hero">
      <div class="container">
        <h1 class="hero-title">欢迎来到我的个人网站</h1>
        <p class="hero-subtitle">分享技术、记录生活、展示作品</p>
        <div class="hero-actions">
          <el-button type="primary" size="large" @click="$router.push('/blog')">阅读博客</el-button>
          <el-button size="large" @click="$router.push('/works')">查看作品</el-button>
        </div>
      </div>
    </section>

    <section class="latest-articles">
      <div class="container">
        <h2 class="section-title">最新文章</h2>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="8" v-for="article in articles" :key="article.id">
            <el-card class="article-card" shadow="hover" @click="$router.push(`/blog/${article.id}`)">
              <img v-if="article.cover" :src="article.cover" class="article-cover" />
              <div class="article-info">
                <h3>{{ article.title }}</h3>
                <p class="article-summary">{{ article.summary }}</p>
                <div class="article-meta">
                  <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                  <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </section>

    <section class="featured-works">
      <div class="container">
        <h2 class="section-title">精选作品</h2>
        <el-row :gutter="20">
          <el-col :xs="24" :sm="12" :md="6" v-for="work in works" :key="work.id">
            <el-card class="work-card" shadow="hover" @click="$router.push(`/works/${work.id}`)">
              <img :src="work.cover" class="work-cover" />
              <h4>{{ work.title }}</h4>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'
import dayjs from 'dayjs'

const articles = ref([])
const works = ref([])

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

onMounted(async () => {
  try {
    const [articlesRes, worksRes] = await Promise.all([
      api.get('/articles?page=1&page_size=6'),
      api.get('/works?page=1&page_size=4&status=1')
    ])
    articles.value = articlesRes.data.list || []
    works.value = worksRes.data.list || []
  } catch (error) {
    console.error('Failed to load data:', error)
  }
})
</script>

<style scoped>
.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 100px 0;
  text-align: center;
}

.hero-title {
  font-size: 3rem;
  margin-bottom: 20px;
  color: white;
}

.hero-subtitle {
  font-size: 1.5rem;
  margin-bottom: 40px;
  color: rgba(255, 255, 255, 0.9);
}

.hero-actions {
  display: flex;
  gap: 20px;
  justify-content: center;
}

.latest-articles,
.featured-works {
  padding: 60px 0;
}

.section-title {
  text-align: center;
  margin-bottom: 40px;
  font-size: 2rem;
}

.article-card {
  cursor: pointer;
  margin-bottom: 20px;
  height: 100%;
}

.article-cover {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 15px;
}

.article-info h3 {
  margin-bottom: 10px;
  font-size: 1.25rem;
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
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.work-card {
  cursor: pointer;
  margin-bottom: 20px;
}

.work-cover {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 4px;
  margin-bottom: 10px;
}

.work-card h4 {
  text-align: center;
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1.2rem;
  }
}
</style>

