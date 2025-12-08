<template>
  <div class="home">
    <!-- Hero Carousel Section -->
    <section class="hero-carousel">
      <div class="container">
        <el-carousel 
          v-if="carouselItems.length > 0"
          :height="carouselHeight"
          :interval="5000"
          :arrow="carouselItems.length > 1 ? 'always' : 'never'"
          indicator-position="inside"
        >
          <el-carousel-item v-for="(item, index) in carouselItems" :key="index">
            <div 
              class="hero-slide" 
              :class="{ 'hero-slide-clickable': item.link }"
              :style="{
                background: item.background || item.backgroundGradient || 'var(--theme-hero-gradient)',
                backgroundImage: item.backgroundImage ? `url(${item.backgroundImage})` : 'none',
                backgroundSize: 'cover',
                backgroundPosition: 'center'
              }"
              @click="handleHeroSlideClick(item)"
            >
              <div class="hero-content">
                <h1 class="hero-title" v-if="item.title">{{ item.title }}</h1>
                <p class="hero-subtitle" v-if="item.subtitle">{{ item.subtitle }}</p>
              </div>
            </div>
          </el-carousel-item>
        </el-carousel>
        <!-- 默认内容（当没有配置轮播图时） -->
        <div v-else class="hero-default">
          <h1 class="hero-title">欢迎来到我的个人网站</h1>
          <p class="hero-subtitle">分享技术、记录生活、展示作品</p>
        </div>
      </div>
    </section>

    <!-- Main Content -->
    <section class="main-content">
      <div class="container">
        <el-row :gutter="30">
          <!-- Left: Articles List -->
          <el-col :xs="24" :lg="17">
            <div class="content-section">
              <div class="section-header">
                <h2>热门文章</h2>
                <el-link type="primary" @click="$router.push('/blog')">
                  查看全部 <el-icon><ArrowRight /></el-icon>
                </el-link>
              </div>

              <div class="article-list">
                <el-card 
                  v-for="article in articles" 
                  :key="article.id" 
                  class="article-item"
                  shadow="hover"
                  @click="$router.push(`/blog/${article.id}`)"
                >
                  <div class="article-content">
                    <div class="article-main">
                      <div class="article-header-info">
                        <el-tag v-if="article.is_top" type="danger" size="small" effect="dark">置顶</el-tag>
                        <el-tag v-if="article.category" size="small">{{ article.category.name }}</el-tag>
                      </div>
                      <h3 class="article-title">{{ article.title }}</h3>
                      <p class="article-summary">{{ article.summary }}</p>
                      <div class="article-meta">
                        <span class="meta-item">
                          <el-avatar :size="20" :src="article.author?.avatar" />
                          {{ article.author?.nickname || article.author?.username }}
                        </span>
                        <span class="meta-item">
                          <el-icon><Clock /></el-icon>
                          {{ formatDate(article.created_at) }}
                        </span>
                        <span class="meta-item">
                          <el-icon><View /></el-icon>
                          {{ article.view_count }}
                        </span>
                        <span class="meta-item" v-if="article.like_count">
                          <el-icon><Star /></el-icon>
                          {{ article.like_count }}
                        </span>
                      </div>
                    </div>
                    <div class="article-cover" v-if="article.cover">
                      <el-image :src="article.cover" fit="cover" />
                    </div>
                  </div>
                </el-card>
              </div>

              <div class="view-more">
                <el-button type="primary" plain @click="$router.push('/blog')">
                  查看更多文章
                </el-button>
              </div>
            </div>

            <!-- Featured Works -->
            <div class="content-section works-section">
              <div class="section-header">
                <h2>精选作品</h2>
                <el-link type="primary" @click="$router.push('/works')">
                  查看全部 <el-icon><ArrowRight /></el-icon>
                </el-link>
              </div>

              <el-row :gutter="20">
                <el-col :xs="24" :sm="12" :md="12" v-for="work in works" :key="work.id">
                  <el-card 
                    class="work-card work-card-clickable"
                    shadow="hover" 
                    @click="$router.push(`/works/${work.id}`)"
                  >
                    <div class="work-type-badge">
                      <el-tag :type="work.type === 'photography' ? 'warning' : 'primary'" size="small">
                        {{ work.type === 'photography' ? '📷' : '💻' }}
                      </el-tag>
                    </div>
                    <el-image :src="work.cover" class="work-cover" fit="cover" />
                    <div class="work-info">
                      <h4>{{ work.title }}</h4>
                      <p class="work-desc">{{ work.description }}</p>
                      <!-- 开源项目显示技术栈 -->
                      <div class="work-tech-stack" v-if="work.type === 'project' && work.tech_stack">
                        <el-tag
                          v-for="(tech, index) in getTechStack(work.tech_stack)"
                          :key="index"
                          size="small"
                          class="tech-tag"
                        >
                          {{ tech }}
                        </el-tag>
                      </div>
                    </div>
                  </el-card>
                </el-col>
              </el-row>
            </div>
          </el-col>

          <!-- Right: Sidebar -->
          <el-col :xs="24" :lg="7">
            <div class="sidebar">
              <!-- Stats Card -->
              <el-card class="sidebar-card stats-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><DataAnalysis /></el-icon>
                    <span>网站统计</span>
                  </div>
                </template>
                <div class="stats-grid">
                  <div class="stat-item">
                    <div class="stat-value">{{ stats.articleCount }}</div>
                    <div class="stat-label">文章</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ stats.workCount }}</div>
                    <div class="stat-label">作品</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ stats.categoryCount }}</div>
                    <div class="stat-label">分类</div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">{{ tags.length }}</div>
                    <div class="stat-label">标签</div>
                  </div>
                </div>
              </el-card>

              <!-- Recommended Articles -->
              <el-card class="sidebar-card" shadow="hover" v-if="recommendedArticles.length > 0">
                <template #header>
                  <div class="card-header">
                    <el-icon><Star /></el-icon>
                    <span>推荐文章</span>
                  </div>
                </template>
                <div class="recommended-list">
                  <div 
                    v-for="article in recommendedArticles" 
                    :key="article.id" 
                    class="recommended-item"
                    @click="$router.push(`/blog/${article.id}`)"
                  >
                    <h4>{{ article.title }}</h4>
                    <div class="recommended-meta">
                      <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
                      <span><el-icon><Star /></el-icon> {{ article.like_count }}</span>
                    </div>
                  </div>
                </div>
              </el-card>

              <!-- Recommended Works -->
              <el-card class="sidebar-card" shadow="hover" v-if="recommendedWorks.length > 0">
                <template #header>
                  <div class="card-header">
                    <el-icon><Picture /></el-icon>
                    <span>推荐作品</span>
                  </div>
                </template>
                <div class="recommended-works">
                  <div 
                    v-for="work in recommendedWorks" 
                    :key="work.id" 
                    class="recommended-work-item"
                    @click="$router.push(`/works/${work.id}`)"
                  >
                    <el-image :src="work.cover" class="work-thumb" fit="cover" />
                    <div class="work-title">{{ work.title }}</div>
                  </div>
                </div>
              </el-card>

              <!-- Tags Card -->
              <el-card class="sidebar-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><PriceTag /></el-icon>
                    <span>热门标签</span>
                  </div>
                </template>
                <div class="tags-cloud">
                  <el-tag
                    v-for="tag in tags"
                    :key="tag.id"
                    class="tag-item"
                    @click="$router.push(`/blog?tag_id=${tag.id}`)"
                  >
                    {{ tag.name }} ({{ tag.article_count }})
                  </el-tag>
                </div>
              </el-card>

              <!-- About Card -->
              <el-card class="sidebar-card about-card" shadow="hover">
                <template #header>
                  <div class="card-header">
                    <el-icon><User /></el-icon>
                    <span>关于本站</span>
                  </div>
                </template>
                <div class="about-content">
                  <p>分享技术文章、记录学习心得、展示个人作品。</p>
                  <el-button type="primary" plain size="small" @click="$router.push('/about')">
                    了解更多
                  </el-button>
                </div>
              </el-card>
            </div>
          </el-col>
        </el-row>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { 
  Reading, 
  Picture, 
  ArrowRight, 
  Clock, 
  View, 
  Star,
  DataAnalysis,
  PriceTag,
  User
} from '@element-plus/icons-vue'
import api from '@/utils/api'
import dayjs from 'dayjs'

const router = useRouter()
const articles = ref([])
const works = ref([])
const tags = ref([])
const recommendedArticles = ref([])
const recommendedWorks = ref([])
const stats = ref({
  articleCount: 0,
  workCount: 0,
  categoryCount: 0
})
const carouselItems = ref([])
const carouselHeight = ref('320px')

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

// 解析技术栈字符串（逗号分隔）
const getTechStack = (techStack) => {
  if (!techStack) return []
  return techStack.split(',').map(tech => tech.trim()).filter(tech => tech.length > 0)
}

const loadCarousel = async () => {
  try {
    const response = await api.get('/settings/public')
    const carouselData = response.data?.home_carousel
    if (carouselData) {
      try {
        carouselItems.value = JSON.parse(carouselData)
      } catch (e) {
        console.error('Failed to parse carousel data:', e)
        carouselItems.value = []
      }
    }
  } catch (error) {
    console.error('Failed to load carousel:', error)
  }
}

const handleCarouselClick = (item) => {
  if (item.link) {
    if (item.link.startsWith('http')) {
      window.open(item.link, '_blank')
    } else {
      router.push(item.link)
    }
  }
}

const loadData = async () => {
  try {
    const [hotArticlesRes, hotWorksRes, tagsRes, recommendedArticlesRes, recommendedWorksRes, statsRes, worksStatsRes] = await Promise.all([
      api.get('/articles/hot?limit=6'),
      api.get('/works/hot?limit=4'),
      api.get('/tags'),
      api.get('/articles/recommended?limit=3'),
      api.get('/works/recommended?limit=2'),
      api.get('/articles?page=1&page_size=1'), // 只获取统计数据
      api.get('/works?page=1&page_size=1') // 只获取统计数据
    ])
    
    articles.value = hotArticlesRes.data || []
    works.value = hotWorksRes.data || []
    tags.value = (tagsRes.data || []).slice(0, 15)
    recommendedArticles.value = recommendedArticlesRes.data || []
    recommendedWorks.value = recommendedWorksRes.data || []
    
    // Calculate stats
    stats.value.articleCount = statsRes.data.total || 0
    stats.value.workCount = worksStatsRes.data.total || 0
    stats.value.categoryCount = new Set(articles.value.map(a => a.category_id).filter(Boolean)).size
  } catch (error) {
    console.error('Failed to load data:', error)
  }
}

onMounted(() => {
  loadCarousel()
  loadData()
})
</script>

<style scoped>
.home {
  background-color: var(--theme-bg-secondary);
}

/* Hero Carousel Section */
.hero-carousel {
  padding: 0;
  margin-bottom: 0;
}

.hero-carousel .container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px 0 20px;
}

.hero-carousel :deep(.el-carousel) {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 4px 20px var(--theme-shadow);
  margin-bottom: 0;
}

.hero-carousel :deep(.el-carousel) {
  margin-bottom: 0 !important;
}

.hero-carousel :deep(.el-carousel__container) {
  height: 320px;
}

.hero-slide {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  color: white;
  text-align: center;
}

.hero-slide-clickable {
  cursor: pointer;
  transition: transform 0.2s;
}

.hero-slide-clickable:hover {
  transform: scale(1.01);
}

.hero-slide::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 1;
}

.hero-content {
  position: relative;
  z-index: 2;
  max-width: 800px;
  padding: 0 20px;
}

.hero-default {
  background: var(--theme-hero-gradient);
  color: white;
  padding: 40px 0;
  text-align: center;
  border-radius: 8px;
  height: 320px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.hero-title {
  font-size: 2rem;
  margin-bottom: 10px;
  color: white;
}

.hero-subtitle {
  font-size: 1rem;
  margin-bottom: 20px;
  color: rgba(255, 255, 255, 0.9);
}


/* Main Content */
.main-content {
  padding: 20px 0 40px 0;
  margin-top: 0;
}

.content-section {
  margin-bottom: 40px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25px;
  padding-bottom: 15px;
  border-bottom: 2px solid #e4e7ed;
}

.section-header h2 {
  font-size: 1.5rem;
  color: var(--theme-text-primary);
  margin: 0;
  display: flex;
  align-items: center;
  gap: 10px;
}

.section-header h2::before {
  content: '';
  width: 4px;
  height: 24px;
  background: var(--theme-hero-gradient);
  border-radius: 2px;
}

/* Article List */
.article-list {
  display: flex;
  flex-direction: column;
  gap: 8px !important;
}

.article-list :deep(.el-card) {
  margin-bottom: 0 !important;
}

.article-item {
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.08);
  height: 160px;
  margin-bottom: 0 !important;
}

.article-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 20px 0 rgba(0, 0, 0, 0.12);
}

.article-item :deep(.el-card__body) {
  height: 100%;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.article-content {
  display: flex;
  gap: 20px;
  height: 100%;
  align-items: center; /* 垂直居中对齐 */
}

.article-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: hidden; /* 防止内容溢出 */
}

.article-header-info {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.article-title {
  font-size: 1.15rem;
  margin: 0 0 8px 0;
  color: var(--theme-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: 1.4;
  word-break: break-word;
  max-height: 1.4em; /* 确保只显示1行 */
}

.article-summary {
  color: var(--theme-text-secondary);
  font-size: 0.85rem;
  margin-bottom: 10px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.5;
  flex: 1;
  word-break: break-word;
  max-height: 3em; /* 确保只显示2行 (1.5 * 2 = 3em) */
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  color: var(--theme-text-secondary);
  font-size: 0.85rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 5px;
}

.article-cover {
  width: 160px;
  height: 120px;
  flex-shrink: 0;
  border-radius: 8px;
  overflow: hidden;
}

.article-cover .el-image {
  width: 100%;
  height: 100%;
}

.view-more {
  text-align: center;
  margin-top: 30px;
}

/* Works Section */
.works-section {
  margin-top: 50px;
}

.works-section :deep(.el-col) {
  display: flex;
  margin-bottom: 12px;
}

.work-card {
  margin-bottom: 12px;
  transition: all 0.3s;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.08);
  position: relative;
  width: 100%;
  min-width: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}

.work-card-clickable {
  cursor: pointer;
}

.work-card-clickable:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 20px 0 rgba(0, 0, 0, 0.12);
}

.work-card :deep(.el-card__body) {
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.work-type-badge {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 10;
}

.work-cover {
  width: 100%;
  max-width: 100%;
  height: 200px;
  border-radius: 4px 4px 0 0;
  flex-shrink: 0;
  overflow: hidden;
}

.work-cover :deep(.el-image) {
  width: 100%;
  height: 100%;
}

.work-cover :deep(.el-image__inner) {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.work-info {
  padding: 15px;
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  width: 100%;
  box-sizing: border-box;
}

.work-info h4 {
  margin: 0 0 8px 0;
  font-size: 1.1rem;
  color: var(--theme-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: 1.4;
}

.work-desc {
  color: var(--theme-text-secondary);
  font-size: 0.85rem;
  margin: 0 0 12px 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: 1.5;
}

.work-tech-stack {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 0;
}

.tech-tag {
  margin: 0;
  font-size: 0.75rem;
}

/* Sidebar */
.sidebar {
  position: sticky;
  top: 80px;
}

.sidebar-card {
  margin-bottom: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.08);
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  color: var(--theme-text-primary);
}

/* Stats Card */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.stat-item {
  text-align: center;
  padding: 15px;
  background: var(--theme-bg-card);
  border: 1px solid var(--theme-border-light);
  border-radius: 8px;
}

.stat-value {
  font-size: 1.8rem;
  font-weight: bold;
  color: var(--theme-primary);
  margin-bottom: 5px;
}

.stat-label {
  font-size: 0.85rem;
  color: var(--theme-text-secondary);
}

/* Tags Cloud */
.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.tag-item {
  cursor: pointer;
  transition: all 0.3s;
}

.tag-item:hover {
  transform: scale(1.05);
}

/* Recommended Articles */
.recommended-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.recommended-item {
  cursor: pointer;
  padding: 12px;
  border-radius: 8px;
  transition: all 0.3s;
  background: var(--theme-bg-secondary);
  border: 1px solid var(--theme-border-light);
}

.recommended-item:hover {
  background: var(--theme-bg-hover);
  transform: translateX(3px);
}

.recommended-item h4 {
  margin: 0 0 8px 0;
  font-size: 0.9rem;
  color: var(--theme-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: 1.4;
}

.recommended-meta {
  display: flex;
  gap: 12px;
  font-size: 0.75rem;
  color: var(--theme-text-secondary);
}

.recommended-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* Recommended Works */
.recommended-works {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.recommended-work-item {
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 8px;
  overflow: hidden;
}

.recommended-work-item:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.work-thumb {
  width: 100%;
  height: 120px;
  margin-bottom: 8px;
}

.work-title {
  font-size: 0.85rem;
  color: var(--theme-text-primary);
  padding: 0 8px 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* About Card */
.about-content {
  text-align: center;
}

.about-content p {
  color: var(--theme-text-secondary);
  margin-bottom: 15px;
  line-height: 1.6;
}

/* Responsive */
@media (max-width: 1200px) {
  .sidebar {
    position: static;
    margin-top: 40px;
  }
}

@media (max-width: 768px) {
  .hero-carousel {
    padding: 15px 0 0 0;
  }

  .hero-carousel :deep(.el-carousel__container) {
    height: 240px !important;
  }

  .hero-default {
    height: 240px;
    padding: 30px 20px;
  }
  
  .hero-title {
    font-size: 1.5rem;
  }

  .hero-subtitle {
    font-size: 0.9rem;
  }

  .main-content {
    padding: 30px 0;
  }

  .article-item {
    height: auto;
    min-height: 140px;
  }

  .article-content {
    flex-direction: column;
  }

  .article-cover {
    width: 100%;
    height: 160px;
  }

  .section-header h2 {
    font-size: 1.3rem;
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 15px;
  }
}
</style>
