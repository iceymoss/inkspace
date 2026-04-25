<template>
  <div class="home">
    <!-- Hero Carousel Section -->
    <section class="hero-carousel">
      <div class="container-blog">
        <div
          v-if="carouselItems.length > 0"
          class="hero-slides"
          @mouseenter="pauseCarousel"
          @mouseleave="resumeCarousel"
        >
          <transition name="fade" mode="out-in">
            <div
              :key="currentSlide"
              class="hero-slide"
              :class="{ 'hero-slide-clickable': carouselItems[currentSlide]?.link }"
              :style="{
                background: carouselItems[currentSlide]?.background || carouselItems[currentSlide]?.backgroundGradient || 'var(--theme-hero-gradient)',
                backgroundImage: carouselItems[currentSlide]?.backgroundImage ? `url(${carouselItems[currentSlide].backgroundImage})` : 'none',
                backgroundSize: 'cover',
                backgroundPosition: 'center'
              }"
              @click="handleCarouselClick(carouselItems[currentSlide])"
            >
              <div class="hero-content">
                <h1 class="hero-title" v-if="carouselItems[currentSlide]?.title">{{ carouselItems[currentSlide].title }}</h1>
                <p class="hero-subtitle" v-if="carouselItems[currentSlide]?.subtitle">{{ carouselItems[currentSlide].subtitle }}</p>
              </div>
            </div>
          </transition>
          <div v-if="carouselItems.length > 1" class="carousel-dots">
            <button
              v-for="(_, i) in carouselItems"
              :key="i"
              class="carousel-dot"
              :class="{ active: currentSlide === i }"
              @click="currentSlide = i"
            />
          </div>
        </div>
        <div v-else class="hero-default">
          <h1 class="hero-title">欢迎来到我的个人网站</h1>
          <p class="hero-subtitle">分享技术、记录生活、展示作品</p>
        </div>
      </div>
    </section>

    <!-- Main Content -->
    <section class="main-content">
      <div class="container-blog">
        <div class="grid grid-cols-1 lg:grid-cols-12 gap-x-[30px]">
          <!-- Left: Articles List -->
          <div class="lg:col-span-7">
            <div class="content-section">
              <div class="section-header">
                <h2>热门文章</h2>
                <button class="section-link" @click="$router.push('/blog')">
                  查看全部 <ArrowRight class="inline w-4 h-4" />
                </button>
              </div>

              <div class="article-list">
                <Card
                  v-for="article in articles"
                  :key="article.id"
                  class="article-item card-skeuomorphic cursor-pointer"
                  @click="$router.push(`/blog/${article.id}`)"
                >
                  <CardContent class="p-md">
                    <div class="article-content">
                      <div class="article-main">
                        <div class="article-header-info">
                          <Badge v-if="article.is_top" variant="destructive">置顶</Badge>
                          <Badge v-if="article.category" variant="secondary">{{ article.category.name }}</Badge>
                        </div>
                        <h3 class="article-title">{{ article.title }}</h3>
                        <p class="article-summary">{{ article.summary }}</p>
                        <div class="article-meta">
                          <span class="meta-item">
                            <Avatar :size="20">
                              <AvatarImage :src="article.author?.avatar" />
                              <AvatarFallback>{{ (article.author?.nickname || '?')[0] }}</AvatarFallback>
                            </Avatar>
                            {{ article.author?.nickname || article.author?.username }}
                          </span>
                          <span class="meta-item">
                            <Clock class="w-4 h-4" />
                            {{ formatDate(article.created_at) }}
                          </span>
                          <span class="meta-item">
                            <Eye class="w-4 h-4" />
                            {{ article.view_count }}
                          </span>
                          <span class="meta-item" v-if="article.like_count">
                            <Star class="w-4 h-4" />
                            {{ article.like_count }}
                          </span>
                          <span class="meta-item" v-if="article.comment_count">
                            <MessageCircle class="w-4 h-4" />
                            {{ article.comment_count }}
                          </span>
                          <span class="meta-item" v-if="article.favorite_count">
                            <Bookmark class="w-4 h-4" />
                            {{ article.favorite_count }}
                          </span>
                        </div>
                      </div>
                      <div class="article-cover" v-if="article.cover">
                        <img :src="article.cover" alt="" class="w-full h-full object-cover" />
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>

              <div class="view-more">
                <Button variant="outline" @click="$router.push('/blog')">
                  查看更多文章
                </Button>
              </div>
            </div>

            <!-- Featured Works -->
            <div class="content-section works-section">
              <div class="section-header">
                <h2>精选作品</h2>
                <button class="section-link" @click="$router.push('/works')">
                  查看全部 <ArrowRight class="inline w-4 h-4" />
                </button>
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-sm">
                <Card
                  v-for="work in works"
                  :key="work.id"
                  class="work-card card-skeuomorphic cursor-pointer"
                  @click="navigateToWorkDetail(work.id, router)"
                >
                  <CardContent class="p-0">
                    <div class="work-type-badge">
                      <Badge :variant="work.type === 'photography' ? 'outline' : 'default'" class="text-xs">
                        {{ work.type === 'photography' ? '📷' : '💻' }}
                      </Badge>
                    </div>
                    <img :src="work.cover" :alt="work.title" class="work-cover object-cover w-full" />
                    <div class="work-info">
                      <h4>{{ work.title }}</h4>
                      <p class="work-desc">{{ work.description }}</p>
                      <div class="work-tech-stack" v-if="work.type === 'project' && work.tech_stack">
                        <Badge
                          v-for="(tech, index) in getTechStack(work.tech_stack)"
                          :key="index"
                          variant="secondary"
                          class="text-xs"
                        >
                          {{ tech }}
                        </Badge>
                      </div>
                      <div class="work-meta">
                        <span class="work-meta-item"><Eye class="w-3.5 h-3.5" /> {{ work.view_count || 0 }}</span>
                        <span class="work-meta-item"><Star class="w-3.5 h-3.5" /> {{ work.like_count || 0 }}</span>
                        <span class="work-meta-item"><MessageCircle class="w-3.5 h-3.5" /> {{ work.comment_count || 0 }}</span>
                        <span class="work-meta-item"><Bookmark class="w-3.5 h-3.5" /> {{ work.favorite_count || 0 }}</span>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>

          <!-- Right: Sidebar -->
          <div class="lg:col-span-5">
            <div class="sidebar">
              <!-- Stats Card -->
              <Card class="sidebar-card card-skeuomorphic-static">
                <CardHeader>
                  <div class="card-header">
                    <BarChart3 class="w-4 h-4" />
                    <span>网站统计</span>
                  </div>
                </CardHeader>
                <CardContent>
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
                </CardContent>
              </Card>

              <!-- Recommended Articles -->
              <Card class="sidebar-card card-skeuomorphic-static" v-if="recommendedArticles.length > 0">
                <CardHeader>
                  <div class="card-header">
                    <Star class="w-4 h-4" />
                    <span>推荐文章</span>
                  </div>
                </CardHeader>
                <CardContent>
                  <div class="recommended-list">
                    <div
                      v-for="article in recommendedArticles"
                      :key="article.id"
                      class="recommended-item"
                      @click="$router.push(`/blog/${article.id}`)"
                    >
                      <h4>{{ article.title }}</h4>
                      <div class="recommended-meta">
                        <span><Eye class="w-3 h-3" /> {{ article.view_count }}</span>
                        <span><Star class="w-3 h-3" /> {{ article.like_count || 0 }}</span>
                        <span><MessageCircle class="w-3 h-3" /> {{ article.comment_count || 0 }}</span>
                        <span><Bookmark class="w-3 h-3" /> {{ article.favorite_count || 0 }}</span>
                      </div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <!-- Recommended Works -->
              <Card class="sidebar-card card-skeuomorphic-static" v-if="recommendedWorks.length > 0">
                <CardHeader>
                  <div class="card-header">
                    <Image class="w-4 h-4" />
                    <span>推荐作品</span>
                  </div>
                </CardHeader>
                <CardContent>
                  <div class="recommended-works">
                    <div
                      v-for="work in recommendedWorks"
                      :key="work.id"
                      class="recommended-work-item"
                      @click="navigateToWorkDetail(work.id, router)"
                    >
                      <img :src="work.cover" :alt="work.title" class="work-thumb object-cover w-full" />
                      <div class="work-title">{{ work.title }}</div>
                    </div>
                  </div>
                </CardContent>
              </Card>

              <!-- Tags Card -->
              <Card class="sidebar-card card-skeuomorphic-static">
                <CardHeader>
                  <div class="card-header">
                    <Tag class="w-4 h-4" />
                    <span>热门标签</span>
                  </div>
                </CardHeader>
                <CardContent>
                  <div class="tags-cloud">
                    <Badge
                      v-for="tag in tags"
                      :key="tag.id"
                      variant="outline"
                      class="cursor-pointer tag-item"
                      @click="$router.push(`/blog?tag_id=${tag.id}`)"
                    >
                      {{ tag.name }} ({{ tag.article_count }})
                    </Badge>
                  </div>
                </CardContent>
              </Card>

              <!-- About Card -->
              <Card class="sidebar-card card-skeuomorphic-static">
                <CardHeader>
                  <div class="card-header">
                    <UserIcon class="w-4 h-4" />
                    <span>关于本站</span>
                  </div>
                </CardHeader>
                <CardContent>
                  <div class="about-content">
                    <p>分享技术文章、记录学习心得、展示个人作品。</p>
                    <Button variant="outline" size="sm" @click="$router.push('/about')">
                      了解更多
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { navigateToWorkDetail } from '@/utils/workNavigation'
import {
  BookOpen,
  Image,
  ArrowRight,
  Clock,
  Eye,
  Star,
  BarChart3,
  Tag,
  User as UserIcon,
  MessageCircle,
  Bookmark
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
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
const currentSlide = ref(0)
let carouselInterval = null

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD')

const getTechStack = (techStack) => {
  if (!techStack) return []
  return techStack.split(',').map(tech => tech.trim()).filter(tech => tech.length > 0)
}

const startCarousel = () => {
  if (carouselItems.value.length > 1) {
    carouselInterval = setInterval(() => {
      currentSlide.value = (currentSlide.value + 1) % carouselItems.value.length
    }, 5000)
  }
}

const pauseCarousel = () => {
  if (carouselInterval) {
    clearInterval(carouselInterval)
    carouselInterval = null
  }
}

const resumeCarousel = () => {
  startCarousel()
}

const loadCarousel = async () => {
  try {
    const response = await api.get('/settings/public')
    const carouselData = response.data?.home_carousel
    if (carouselData) {
      try {
        carouselItems.value = JSON.parse(carouselData)
        startCarousel()
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
  if (item?.link) {
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
      api.get('/articles?page=1&page_size=1'),
      api.get('/works?page=1&page_size=1'),
      api.get('/categories')
    ])
    
    articles.value = hotArticlesRes.data || []
    works.value = hotWorksRes.data || []
    tags.value = (tagsRes.data || []).slice(0, 15)
    recommendedArticles.value = recommendedArticlesRes.data || []
    recommendedWorks.value = recommendedWorksRes.data || []
    
    stats.value.articleCount = statsRes.data.total || 0
    stats.value.workCount = worksStatsRes.data.total || 0
    stats.value.categoryCount = (categoriesRes.data || []).length
  } catch (error) {
    console.error('Failed to load data:', error)
  }
}

onMounted(() => {
  loadCarousel()
  loadData()
})

onUnmounted(() => {
  pauseCarousel()
})
</script>

<style scoped>
.home {
  background-color: var(--theme-bg-secondary);
  @apply font-sans;
}

.hero-carousel {
  @apply p-0 mb-0;
}

.hero-slides {
  @apply relative overflow-hidden;
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-md);
  height: 320px;
}

.hero-slide {
  @apply w-full h-full flex items-center justify-center relative text-white text-center;
}

.hero-slide::before {
  content: '';
  @apply absolute inset-0 z-[1];
  background: rgba(0, 0, 0, 0.3);
}

.hero-slide-clickable {
  @apply cursor-pointer;
  transition: transform var(--transition-base);
}

.hero-slide-clickable:hover {
  transform: scale(1.01);
}

.hero-content {
  @apply relative z-[2] max-w-[800px] px-md;
}

.hero-default {
  background: var(--theme-hero-gradient);
  @apply text-white py-xl text-center;
  border-radius: var(--radius-lg);
  height: 320px;
  @apply flex flex-col justify-center items-center;
}

.hero-title {
  font-size: var(--font-size-3xl);
  @apply mb-sm text-white font-bold font-serif;
}

.hero-subtitle {
  @apply mb-md;
  font-size: var(--font-size-base);
  color: rgba(255, 255, 255, 0.9);
  line-height: var(--line-height-base);
}

.carousel-dots {
  @apply absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2 z-10;
}

.carousel-dot {
  @apply w-2.5 h-2.5 rounded-full bg-white/50 cursor-pointer transition-all;
  border: none;
}

.carousel-dot.active {
  @apply bg-white w-6;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.main-content {
  @apply py-md;
  padding-top: var(--spacing-md);
  padding-bottom: var(--spacing-xl);
}

.content-section {
  @apply mb-xl;
}

.section-header {
  @apply flex justify-between items-center mb-lg pb-md;
  border-bottom: 2px solid var(--theme-border);
}

.section-header h2 {
  font-size: var(--font-size-2xl);
  color: var(--theme-text-primary);
  @apply m-0 flex items-center gap-sm font-bold;
}

.section-header h2::before {
  content: '';
  @apply w-1 rounded-sm;
  height: var(--spacing-md);
  background: var(--theme-primary);
}

.section-link {
  @apply flex items-center gap-1 text-sm font-medium cursor-pointer bg-transparent border-none;
  color: var(--theme-primary);
  transition: color var(--transition-fast);
}

.section-link:hover {
  color: var(--theme-primary-hover);
}

.article-list {
  @apply flex flex-col gap-sm;
}

.article-item {
  height: 160px;
}

.article-item:hover {
  transform: translateY(-2px);
}

.article-content {
  @apply flex gap-md h-full items-center;
}

.article-main {
  @apply flex-1 min-w-0 flex flex-col justify-between overflow-hidden;
}

.article-header-info {
  @apply flex gap-sm mb-sm;
}

.article-title {
  @apply font-semibold m-0 mb-xs;
  font-size: var(--font-size-lg);
  color: var(--theme-text-primary);
  @apply overflow-hidden text-ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-tight);
  word-break: break-word;
  max-height: 1.4em;
}

.article-summary {
  color: var(--theme-text-secondary);
  @apply mb-sm;
  font-size: var(--font-size-sm);
  @apply overflow-hidden text-ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
  flex: 1;
  word-break: break-word;
  max-height: 3em;
}

.article-meta {
  @apply flex flex-wrap gap-md;
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
}

.meta-item {
  @apply flex items-center gap-xs;
}

.article-cover {
  @apply w-[160px] h-[120px] shrink-0 overflow-hidden;
  border-radius: var(--radius-md);
}

.view-more {
  @apply text-center mt-[30px];
}

.works-section {
  @apply mt-[50px];
}

.work-card {
  @apply relative w-full min-w-0 h-full;
}

.work-type-badge {
  @apply absolute top-sm right-sm z-10;
}

.work-cover {
  height: 200px;
  border-radius: var(--radius-sm) var(--radius-sm) 0 0;
  @apply shrink-0;
}

.work-info {
  @apply p-sm flex-1 flex flex-col min-w-0 w-full;
}

.work-info h4 {
  @apply m-0 mb-xs font-semibold;
  font-size: var(--font-size-lg);
  color: var(--theme-text-primary);
  @apply overflow-hidden text-ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-tight);
}

.work-desc {
  color: var(--theme-text-secondary);
  @apply m-0 mb-sm;
  font-size: var(--font-size-sm);
  @apply overflow-hidden text-ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
}

.work-meta {
  @apply flex gap-sm mt-xs flex-wrap;
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
}

.work-meta-item {
  @apply flex items-center gap-xs;
}

.work-tech-stack {
  @apply flex flex-wrap gap-1.5 mt-0;
}

.sidebar {
  @apply sticky;
  top: 80px;
}

.sidebar-card {
  @apply mb-md;
}

.card-header {
  @apply flex items-center gap-sm font-semibold;
  color: var(--theme-text-primary);
}

.stats-grid {
  @apply grid grid-cols-2 gap-md;
}

.stat-item {
  @apply text-center p-sm;
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
  @apply font-bold mb-xs;
  color: var(--theme-primary);
}

.stat-label {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.tags-cloud {
  @apply flex flex-wrap gap-sm;
}

.tag-item {
  @apply cursor-pointer;
  transition: all var(--transition-slow);
}

.tag-item:hover {
  transform: scale(1.05);
}

.recommended-list {
  @apply flex flex-col gap-sm;
}

.recommended-item {
  @apply cursor-pointer p-sm;
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
  @apply m-0 mb-xs font-semibold;
  font-size: var(--font-size-sm);
  color: var(--theme-text-primary);
  @apply overflow-hidden text-ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-tight);
}

.recommended-meta {
  @apply flex gap-sm;
  font-size: var(--font-size-xs);
  color: var(--theme-text-secondary);
}

.recommended-meta span {
  @apply flex items-center gap-xs;
}

.recommended-works {
  @apply flex flex-col gap-sm;
}

.recommended-work-item {
  @apply cursor-pointer overflow-hidden;
  transition: all var(--transition-slow);
  border-radius: var(--radius-md);
}

.recommended-work-item:hover {
  transform: scale(1.02);
  box-shadow: var(--shadow-md);
}

.work-thumb {
  @apply w-full mb-xs;
  height: 120px;
}

.work-title {
  font-size: var(--font-size-sm);
  color: var(--theme-text-primary);
  @apply px-sm pb-sm;
  @apply overflow-hidden text-ellipsis whitespace-nowrap;
}

.about-content {
  @apply text-center;
}

.about-content p {
  color: var(--theme-text-secondary);
  @apply mb-sm;
  line-height: var(--line-height-relaxed);
}

@media (max-width: 1200px) {
  .sidebar {
    @apply static mt-xl;
  }
}

@media (max-width: 768px) {
  .hero-carousel {
    @apply pt-sm;
  }

  .hero-slides {
    height: 240px;
  }

  .hero-default {
    height: 240px;
    @apply py-[30px] px-md;
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
    @apply h-auto;
    min-height: 140px;
  }

  .article-content {
    @apply flex-col;
  }

  .article-cover {
    @apply w-full;
    height: 160px;
  }

  .section-header h2 {
    font-size: var(--font-size-xl);
  }

  .stats-grid {
    @apply gap-sm;
  }
}
</style>
