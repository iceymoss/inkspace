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
                        <span class="meta-item" v-if="article.comment_count">
                          <el-icon><ChatDotRound /></el-icon>
                          {{ article.comment_count }}
                        </span>
                        <span class="meta-item" v-if="article.favorite_count">
                          <el-icon><Collection /></el-icon>
                          {{ article.favorite_count }}
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
                    @click="navigateToWorkDetail(work.id, router)"
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
                      <!-- 作品元数据 -->
                      <div class="work-meta">
                        <span class="work-meta-item">
                          <el-icon><View /></el-icon>
                          {{ work.view_count || 0 }}
                        </span>
                        <span class="work-meta-item">
                          <el-icon><Star /></el-icon>
                          {{ work.like_count || 0 }}
                        </span>
                        <span class="work-meta-item">
                          <el-icon><ChatDotRound /></el-icon>
                          {{ work.comment_count || 0 }}
                        </span>
                        <span class="work-meta-item">
                          <el-icon><Collection /></el-icon>
                          {{ work.favorite_count || 0 }}
                        </span>
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
                      <span><el-icon><Star /></el-icon> {{ article.like_count || 0 }}</span>
                      <span><el-icon><ChatDotRound /></el-icon> {{ article.comment_count || 0 }}</span>
                      <span><el-icon><Collection /></el-icon> {{ article.favorite_count || 0 }}</span>
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
                    @click="navigateToWorkDetail(work.id, router)"
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
import { navigateToWorkDetail } from '@/utils/workNavigation'
import { 
  Reading, 
  Picture, 
  ArrowRight, 
  Clock, 
  View, 
  Star,
  DataAnalysis,
  PriceTag,
  User,
  ChatDotRound,
  Collection
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
    const [hotArticlesRes, hotWorksRes, tagsRes, recommendedArticlesRes, recommendedWorksRes, statsRes, worksStatsRes, categoriesRes] = await Promise.all([
      api.get('/articles/hot?limit=6'),
      api.get('/works/hot?limit=4'),
      api.get('/tags'),
      api.get('/articles/recommended?limit=3'),
      api.get('/works/recommended?limit=2'),
      api.get('/articles?page=1&page_size=1'), // 只获取统计数据
      api.get('/works?page=1&page_size=1'), // 只获取统计数据
      api.get('/categories') // 获取所有分类
    ])
    
    articles.value = hotArticlesRes.data || []
    works.value = hotWorksRes.data || []
    tags.value = (tagsRes.data || []).slice(0, 15)
    recommendedArticles.value = recommendedArticlesRes.data || []
    recommendedWorks.value = recommendedWorksRes.data || []
    
    // Calculate stats
    stats.value.articleCount = statsRes.data.total || 0
    stats.value.workCount = worksStatsRes.data.total || 0
    stats.value.categoryCount = (categoriesRes.data || []).length // 使用所有分类的数量
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
  font-family: var(--font-sans);
}

.hero-carousel {
  padding: 0;
  margin-bottom: 0;
}

.hero-carousel .container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
}

.hero-carousel :deep(.el-carousel) {
  border-radius: var(--radius-md);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  margin-bottom: 0;
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
  color: #fff;
  text-align: center;
}

.hero-slide-clickable {
  cursor: pointer;
  transition: transform var(--transition-base);
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
  padding: 0 var(--spacing-md);
}

.hero-default {
  background: var(--theme-hero-gradient);
  color: #fff;
  padding: var(--spacing-xl) 0;
  text-align: center;
  border-radius: var(--radius-md);
  height: 320px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.hero-title {
  font-size: var(--font-size-3xl);
  margin-bottom: var(--spacing-sm);
  color: #fff;
  font-weight: 700;
  font-family: var(--font-serif);
}

.hero-subtitle {
  font-size: var(--font-size-base);
  margin-bottom: var(--spacing-md);
  color: rgba(255, 255, 255, 0.9);
  line-height: var(--line-height-base);
}

.main-content {
  padding: var(--spacing-md) 0 var(--spacing-xl) 0;
  margin-top: 0;
}

.content-section {
  margin-bottom: var(--spacing-xl);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-md);
  border-bottom: 2px solid var(--theme-border);
}

.section-header h2 {
  font-size: var(--font-size-2xl);
  color: var(--theme-text-primary);
  margin: 0;
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-weight: 700;
}

.section-header h2::before {
  content: '';
  width: 4px;
  height: var(--spacing-md);
  background: var(--theme-primary);
  border-radius: var(--radius-sm);
}

.article-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.article-list :deep(.el-card) {
  margin-bottom: 0;
}

.article-item {
  cursor: pointer;
  transition: all var(--transition-slow);
  box-shadow: var(--shadow-sm);
  height: 160px;
  margin-bottom: 0;
}

.article-item:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.article-item :deep(.el-card__body) {
  height: 100%;
  padding: var(--spacing-md);
  display: flex;
  flex-direction: column;
}

.article-content {
  display: flex;
  gap: var(--spacing-md);
  height: 100%;
  align-items: center;
}

.article-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  overflow: hidden;
}

.article-header-info {
  display: flex;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
}

.article-title {
  font-size: var(--font-size-lg);
  margin: 0 0 var(--spacing-xs) 0;
  color: var(--theme-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-tight);
  word-break: break-word;
  max-height: 1.4em;
  font-weight: 600;
}

.article-summary {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  margin-bottom: var(--spacing-sm);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
  flex: 1;
  word-break: break-word;
  max-height: 3em;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-md);
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.article-cover {
  width: 160px;
  height: 120px;
  flex-shrink: 0;
  border-radius: var(--radius-md);
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

.works-section {
  margin-top: 50px;
}

.works-section :deep(.el-col) {
  display: flex;
  margin-bottom: var(--spacing-sm);
}

.work-card {
  margin-bottom: var(--spacing-sm);
  transition: all var(--transition-slow);
  box-shadow: var(--shadow-sm);
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
  box-shadow: var(--shadow-md);
}

.work-card :deep(.el-card__body) {
  padding: 0;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.work-type-badge {
  position: absolute;
  top: var(--spacing-sm);
  right: var(--spacing-sm);
  z-index: 10;
}

.work-cover {
  width: 100%;
  max-width: 100%;
  height: 200px;
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
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
  padding: var(--spacing-sm);
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  width: 100%;
  box-sizing: border-box;
}

.work-info h4 {
  margin: 0 0 var(--spacing-xs) 0;
  font-size: var(--font-size-lg);
  color: var(--theme-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-tight);
  font-weight: 600;
}

.work-desc {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  margin: 0 0 var(--spacing-sm) 0;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
}

.work-meta {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-xs);
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  flex-wrap: wrap;
}

.work-meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.work-tech-stack {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 0;
}

.tech-tag {
  margin: 0;
  font-size: var(--font-size-xs);
}

.sidebar {
  position: sticky;
  top: 80px;
}

.sidebar-card {
  margin-bottom: var(--spacing-md);
  box-shadow: var(--shadow-sm);
  transition: box-shadow var(--transition-base);
}

.sidebar-card:hover {
  box-shadow: var(--shadow-md);
}

.card-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  font-weight: 600;
  color: var(--theme-text-primary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--spacing-md);
}

.stat-item {
  text-align: center;
  padding: var(--spacing-sm);
  background: var(--theme-bg-card);
  border: 1px solid var(--theme-border-light);
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
}

.stat-item:hover {
  border-color: var(--theme-primary);
}

.stat-value {
  font-size: var(--font-size-3xl);
  font-weight: bold;
  color: var(--theme-primary);
  margin-bottom: var(--spacing-xs);
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.tags-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.tag-item {
  cursor: pointer;
  transition: all var(--transition-slow);
}

.tag-item:hover {
  transform: scale(1.05);
}

.recommended-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.recommended-item {
  cursor: pointer;
  padding: var(--spacing-sm);
  border-radius: var(--radius-md);
  transition: all var(--transition-slow);
  background: var(--theme-bg-secondary);
  border: 1px solid var(--theme-border-light);
}

.recommended-item:hover {
  background: var(--theme-bg-hover);
  border-color: var(--theme-primary);
  transform: translateX(3px);
}

.recommended-item h4 {
  margin: 0 0 var(--spacing-xs) 0;
  font-size: var(--font-size-sm);
  color: var(--theme-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-tight);
  font-weight: 600;
}

.recommended-meta {
  display: flex;
  gap: var(--spacing-sm);
  font-size: var(--font-size-xs);
  color: var(--theme-text-secondary);
}

.recommended-meta span {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.recommended-works {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.recommended-work-item {
  cursor: pointer;
  transition: all var(--transition-slow);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.recommended-work-item:hover {
  transform: scale(1.02);
  box-shadow: var(--shadow-md);
}

.work-thumb {
  width: 100%;
  height: 120px;
  margin-bottom: var(--spacing-xs);
}

.work-title {
  font-size: var(--font-size-sm);
  color: var(--theme-text-primary);
  padding: 0 var(--spacing-sm) var(--spacing-sm);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.about-content {
  text-align: center;
}

.about-content p {
  color: var(--theme-text-secondary);
  margin-bottom: var(--spacing-sm);
  line-height: var(--line-height-relaxed);
}

@media (max-width: 1200px) {
  .sidebar {
    position: static;
    margin-top: var(--spacing-xl);
  }
}

@media (max-width: 768px) {
  .hero-carousel {
    padding: var(--spacing-sm) 0 0 0;
  }

  .hero-carousel :deep(.el-carousel__container) {
    height: 240px !important;
  }

  .hero-default {
    height: 240px;
    padding: 30px var(--spacing-md);
  }

  .hero-title {
    font-size: var(--font-size-2xl);
  }

  .hero-subtitle {
    font-size: var(--font-size-sm);
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
    font-size: var(--font-size-xl);
  }

  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: var(--spacing-sm);
  }
}
</style>
