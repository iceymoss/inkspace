<template>
  <div class="home">
    <!-- Hero Carousel Section -->
    <section class="hero-carousel">
      <div class="container">
        <TerminalHero
          v-if="isTerminal"
          :settings="terminalHeroSettings"
          :title="terminalHeroTitle"
          :stats="stats"
          :active="terminal.isOpen"
          @navigate="handleHeroLink"
          @activate="activateTerminal"
        />
        <el-carousel 
          v-else-if="carouselItems.length > 0"
          :height="carouselHeight"
          :interval="5000"
          :arrow="carouselItems.length > 1 ? 'always' : 'never'"
          indicator-position="inside"
        >
          <el-carousel-item
            v-for="(item, index) in carouselItems"
            :key="index"
          >
            <div 
              class="hero-slide" 
              :class="{ 'hero-slide-clickable': item.link }"
              :style="{
                background: item.background || item.backgroundGradient || 'var(--theme-hero-gradient)',
                backgroundImage: item.backgroundImage ? `url(${item.backgroundImage})` : 'none',
                backgroundSize: 'cover',
                backgroundPosition: 'center'
              }"
              @click="handleCarouselClick(item)"
            >
              <div class="hero-content">
                <h1
                  v-if="item.title"
                  class="hero-title"
                >
                  {{ item.title }}
                </h1>
                <p
                  v-if="item.subtitle"
                  class="hero-subtitle"
                >
                  {{ item.subtitle }}
                </p>
              </div>
            </div>
          </el-carousel-item>
        </el-carousel>
        <!-- 默认内容（当没有配置轮播图时） -->
        <div
          v-else
          class="editorial-hero"
        >
          <div class="issue-line">
            <span class="issue-number">{{ heroSettings.issue }}</span>
            <span class="issue-rule" />
            <span>{{ heroSettings.eyebrow }}</span>
          </div>
          <h1>
            {{ heroTitle.before }}<span v-if="heroTitle.accent">{{ heroTitle.accent }}</span>{{ heroTitle.after }}
          </h1>
          <p>{{ heroSettings.description }}</p>
          <div class="hero-actions">
            <a
              class="hero-primary"
              :href="safeHeroLink(heroSettings.primary_link)"
              @click="handleHeroLink($event, heroSettings.primary_link)"
            >{{ heroSettings.primary_text }}</a>
            <a
              class="hero-secondary"
              :href="safeHeroLink(heroSettings.secondary_link)"
              @click="handleHeroLink($event, heroSettings.secondary_link)"
            >{{ heroSettings.secondary_text }} <span aria-hidden="true">→</span></a>
          </div>
        </div>
      </div>
    </section>

    <!-- Main Content -->
    <section class="main-content">
      <div class="container">
        <el-row :gutter="30">
          <!-- Left: Articles List -->
          <el-col
            :xs="24"
            :lg="19"
          >
            <div class="content-section">
              <div class="section-header">
                <h2>热门文章</h2>
                <el-link
                  type="primary"
                  @click="$router.push('/blog')"
                >
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
                  <div class="article-content" :class="{ 'has-cover': article.cover }">
                    <div class="terminal-log-meta">
                      <span :class="`level-${getArticleLogLevel(article).toLowerCase()}`">{{ getArticleLogLevel(article) }}</span>
                      <time :datetime="article.created_at">{{ formatDate(article.created_at) }}</time>
                    </div>
                    <div class="article-main">
                      <div class="article-header-info">
                        <el-tag
                          v-if="article.is_top"
                          type="danger"
                          size="small"
                          effect="dark"
                        >
                          置顶
                        </el-tag>
                        <el-tag
                          v-if="article.category"
                          size="small"
                        >
                          {{ article.category.name }}
                        </el-tag>
                      </div>
                      <h3 class="article-title">
                        {{ article.title }}
                      </h3>
                      <p class="article-summary">
                        {{ article.summary }}
                      </p>
                      <div class="article-meta">
                        <span class="meta-item article-date">
                          <el-avatar
                            :size="20"
                            :src="article.author?.avatar"
                          />
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
                        <span
                          v-if="article.like_count"
                          class="meta-item"
                        >
                          <el-icon><Star /></el-icon>
                          {{ article.like_count }}
                        </span>
                        <span
                          v-if="article.comment_count"
                          class="meta-item"
                        >
                          <el-icon><ChatDotRound /></el-icon>
                          {{ article.comment_count }}
                        </span>
                        <span
                          v-if="article.favorite_count"
                          class="meta-item"
                        >
                          <el-icon><Collection /></el-icon>
                          {{ article.favorite_count }}
                        </span>
                      </div>
                    </div>
                    <div
                      v-if="article.cover"
                      class="article-cover"
                    >
                      <el-image
                        :src="article.cover"
                        fit="cover"
                      />
                    </div>
                  </div>
                </el-card>
              </div>

              <div class="view-more">
                <el-button
                  type="primary"
                  plain
                  @click="$router.push('/blog')"
                >
                  查看更多文章
                </el-button>
              </div>
            </div>

            <!-- Featured Works -->
            <div class="content-section works-section">
              <div class="section-header">
                <h2>精选作品</h2>
                <el-link
                  type="primary"
                  @click="$router.push('/works')"
                >
                  查看全部 <el-icon><ArrowRight /></el-icon>
                </el-link>
              </div>

              <el-row :gutter="20">
                <el-col
                  v-for="work in works"
                  :key="work.id"
                  :xs="24"
                  :sm="12"
                  :md="12"
                >
                  <el-card 
                    class="work-card work-card-clickable"
                    shadow="hover" 
                    @click="navigateToWorkDetail(work.id, router)"
                  >
                    <div class="work-type-badge">
                      <el-tag
                        :type="work.type === 'photography' ? 'warning' : 'primary'"
                        size="small"
                      >
                        {{ work.type === 'photography' ? '📷' : '💻' }}
                      </el-tag>
                    </div>
                    <el-image
                      :src="work.cover"
                      class="work-cover"
                      fit="cover"
                    />
                    <div class="work-info">
                      <h4>{{ work.title }}</h4>
                      <p class="work-desc">
                        {{ work.description }}
                      </p>
                      <!-- 开源项目显示技术栈 -->
                      <div
                        v-if="work.type === 'project' && work.tech_stack"
                        class="work-tech-stack"
                      >
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

            <div v-if="isTerminal" class="content-section terminal-captures">
              <div class="section-header">
                <h2>photos <small>— captures</small></h2>
                <router-link class="terminal-more" to="/photos">open gallery <el-icon><ArrowRight /></el-icon></router-link>
              </div>
              <div v-if="photos.length" class="terminal-photo-grid">
                <router-link v-for="photo in photos" :key="photo.id" :to="`/works/${photo.id}`" class="terminal-photo">
                  <img v-if="getWorkImage(photo)" :src="getWorkImage(photo)" :alt="photo.title" loading="lazy">
                  <span v-else class="terminal-photo-fallback">NO IMAGE</span>
                  <i aria-hidden="true" />
                  <div><strong>{{ photo.title }}</strong><span>{{ photo.metadata?.location || 'capture' }}</span></div>
                </router-link>
              </div>
              <div v-else-if="terminalSectionErrors.photos" class="terminal-empty terminal-error">
                [ captures unavailable ] <button type="button" @click="loadTerminalSections">retry</button>
              </div>
              <div v-else class="terminal-empty">[ no published captures ]</div>
            </div>

            <div v-if="isTerminal" class="content-section terminal-wiki">
              <div class="section-header">
                <h2>wiki <small>— public workspaces</small></h2>
                <router-link class="terminal-more" to="/wiki">open wiki/ <el-icon><ArrowRight /></el-icon></router-link>
              </div>
              <div v-if="wikiWorkspaces.length" class="terminal-wiki-panel">
                <router-link v-for="workspace in wikiWorkspaces" :key="workspace.id" :to="`/wiki/${workspace.id}`">
                  <span>▸ {{ workspace.name }}/</span>
                  <small>{{ workspace.doc_count }} docs · {{ workspace.author_name || 'anonymous' }}</small>
                </router-link>
              </div>
              <div v-else-if="terminalSectionErrors.wiki" class="terminal-empty terminal-error">
                [ wiki unavailable ] <button type="button" @click="loadTerminalSections">retry</button>
              </div>
              <div v-else class="terminal-empty">[ no public workspaces ]</div>
            </div>
          </el-col>

          <!-- Right: Sidebar -->
          <el-col
            :xs="24"
            :lg="5"
          >
            <div class="sidebar">
              <!-- Stats Card -->
              <el-card
                class="sidebar-card stats-card"
                shadow="hover"
              >
                <template #header>
                  <div class="card-header">
                    <el-icon><DataAnalysis /></el-icon>
                    <span>网站统计</span>
                  </div>
                </template>
                <div class="stats-grid">
                  <div class="stat-item">
                    <div class="stat-value">
                      {{ stats.articleCount }}
                    </div>
                    <div class="stat-label">
                      文章
                    </div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">
                      {{ stats.workCount }}
                    </div>
                    <div class="stat-label">
                      作品
                    </div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">
                      {{ stats.categoryCount }}
                    </div>
                    <div class="stat-label">
                      分类
                    </div>
                  </div>
                  <div class="stat-item">
                    <div class="stat-value">
                      {{ tags.length }}
                    </div>
                    <div class="stat-label">
                      标签
                    </div>
                  </div>
                </div>
              </el-card>

              <!-- Recommended Articles -->
              <el-card
                v-if="recommendedArticles.length > 0"
                class="sidebar-card"
                shadow="hover"
              >
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
              <el-card
                v-if="recommendedWorks.length > 0"
                class="sidebar-card"
                shadow="hover"
              >
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
                    <el-image
                      :src="work.cover"
                      class="work-thumb"
                      fit="cover"
                    />
                    <div class="work-title">
                      {{ work.title }}
                    </div>
                  </div>
                </div>
              </el-card>

              <!-- Tags Card -->
              <el-card
                class="sidebar-card"
                shadow="hover"
              >
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
              <el-card
                class="sidebar-card about-card"
                shadow="hover"
              >
                <template #header>
                  <div class="card-header">
                    <el-icon><User /></el-icon>
                    <span>关于本站</span>
                  </div>
                </template>
                <div class="about-content">
                  <p>分享技术文章、记录学习心得、展示个人作品。</p>
                  <el-button
                    type="primary"
                    plain
                    size="small"
                    @click="$router.push('/about')"
                  >
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
import { computed, reactive, ref, onMounted, watch } from 'vue'
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
import TerminalHero from '@/components/theme/TerminalHero.vue'
import { useAppearanceStore } from '@/stores/appearance'
import { useTerminalStore } from '@/stores/terminal'
import { getArticleLogLevel } from '@/utils/terminal/articleLogLevel'

const router = useRouter()
const appearance = useAppearanceStore()
const terminal = useTerminalStore()
const isTerminal = computed(() => appearance.activePreference.ui_theme === 'terminal')
const articles = ref([])
const works = ref([])
const tags = ref([])
const recommendedArticles = ref([])
const recommendedWorks = ref([])
const photos = ref([])
const wikiWorkspaces = ref([])
const terminalSectionsLoaded = ref(false)
const terminalSectionErrors = reactive({ photos: false, wiki: false })
const stats = ref({
  articleCount: 0,
  workCount: 0,
  categoryCount: 0
})
const carouselItems = ref([])
const carouselHeight = ref('320px')
const heroSettings = reactive({
  issue: `VOL. ${String(new Date().getFullYear()).slice(-2)}`,
  eyebrow: '持续记录 · 自由生长',
  title: '把日常的观察，写成可以停留的文字。',
  accent: '停留',
  description: '这里收录关于技术与设计的长文、正在生长的知识库，以及在生活与远方之间留下的作品。',
  primary_text: '开始阅读',
  primary_link: '/blog',
  secondary_text: '或先看看照片',
  secondary_link: '/photos'
})
const terminalHeroSettings = reactive({
  status: 'system online',
  eyebrow: '持续构建 · 持续记录',
  title: '在代码与生活之间，持续记录。',
  accent: '持续记录',
  description: '这里保存技术实践、正在构建的作品，以及偶尔偏离主线的生活片段。',
  primary_text: 'tail -f blog',
  primary_link: '/blog',
  secondary_text: 'ls projects/',
  secondary_link: '/works'
})
const splitHeroTitle = settings => {
  const title = settings.title
  const accent = settings.accent
  const index = accent ? title.indexOf(accent) : -1
  if (index < 0) return { before: title, accent: '', after: '' }
  return {
    before: title.slice(0, index),
    accent,
    after: title.slice(index + accent.length)
  }
}
const heroTitle = computed(() => splitHeroTitle(heroSettings))
const terminalHeroTitle = computed(() => splitHeroTitle(terminalHeroSettings))

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
    const heroData = response.data?.home_hero
    const terminalHeroData = response.data?.home_hero_terminal
    if (heroData) {
      try {
        const parsed = JSON.parse(heroData)
        if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
          Object.keys(heroSettings).forEach((key) => {
            if (typeof parsed[key] === 'string' && parsed[key].trim()) heroSettings[key] = parsed[key].trim()
          })
        }
      } catch (e) {
        console.error('Failed to parse home hero data:', e)
      }
    }
    if (terminalHeroData) {
      try {
        const parsed = JSON.parse(terminalHeroData)
        if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
          Object.keys(terminalHeroSettings).forEach((key) => {
            if (typeof parsed[key] === 'string' && parsed[key].trim()) terminalHeroSettings[key] = parsed[key].trim()
          })
        }
      } catch (e) {
        console.error('Failed to parse terminal home hero data:', e)
      }
    }
    if (carouselData) {
      try {
        const parsed = JSON.parse(carouselData)
        carouselItems.value = Array.isArray(parsed) ? parsed : []
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
      window.open(item.link, '_blank', 'noopener,noreferrer')
    } else {
      router.push(item.link)
    }
  }
}
const getWorkImage = work => work.cover || (typeof work.images?.[0] === 'string' ? work.images[0] : work.images?.[0]?.url) || ''

const handleHeroLink = (event, link) => {
  if (/^https?:\/\//.test(link)) return
  event.preventDefault()
  if (link?.startsWith('/')) router.push(link)
}

const safeHeroLink = link => /^https?:\/\//.test(link) || link?.startsWith('/') ? link : '#'

const activateTerminal = sourceRect => terminal.open(sourceRect)

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

const loadTerminalSections = async () => {
  terminalSectionErrors.photos = false
  terminalSectionErrors.wiki = false
  const [photosResult, wikiResult] = await Promise.allSettled([
    api.get('/works', { params: { type: 'photography', status: 1, page: 1, page_size: 3 } }),
    api.get('/wiki/workspaces', { params: { page: 1, page_size: 4 } })
  ])
  terminalSectionErrors.photos = photosResult.status === 'rejected'
  terminalSectionErrors.wiki = wikiResult.status === 'rejected'
  photos.value = photosResult.status === 'fulfilled' ? photosResult.value.data?.list || [] : []
  wikiWorkspaces.value = wikiResult.status === 'fulfilled' ? wikiResult.value.data?.list || [] : []
  terminalSectionsLoaded.value = true
}

onMounted(() => {
  loadCarousel()
  loadData()
  if (isTerminal.value) loadTerminalSections()
})

watch(isTerminal, (terminal) => {
  if (terminal && !terminalSectionsLoaded.value) loadTerminalSections()
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
  width: 124px;
  height: 84px;
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

.work-meta {
  display: flex;
  gap: 12px;
  margin-top: 8px;
  color: var(--theme-text-secondary);
  font-size: 0.85rem;
  flex-wrap: wrap;
}

.work-meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
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

.terminal-log-meta {
  display: none;
}

.terminal-photo-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 18px; }
.terminal-photo { position: relative; display: block; aspect-ratio: 3 / 2; overflow: hidden; border: 1px solid var(--line); border-radius: 14px; background: var(--panel); color: #afc4de; }
.terminal-photo img { width: 100%; height: 100%; object-fit: cover; transition: transform .5s cubic-bezier(.2, .6, .2, 1); }
.terminal-photo:hover img { transform: scale(1.05); }
.terminal-photo > i { position: absolute; inset: 0; background: repeating-linear-gradient(0deg, transparent 0 3px, rgba(255, 255, 255, .025) 3px 4px); }
.terminal-photo > div { position: absolute; inset: auto 0 0; display: flex; justify-content: space-between; gap: 12px; padding: 30px 13px 10px; background: linear-gradient(transparent, rgba(5, 8, 15, .85)); font: 11px var(--terminal-mono); }
.terminal-photo > div strong { overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.terminal-photo > div span { flex: none; }
.terminal-photo-fallback { position: absolute; inset: 0; display: grid; place-items: center; color: var(--sub); font: 11px var(--terminal-mono); }
.terminal-wiki-panel { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 1px; overflow: hidden; border: 1px solid var(--line); border-radius: 14px; background: var(--line); box-shadow: var(--terminal-shadow); }
.terminal-wiki-panel a { display: grid; gap: 8px; padding: 22px; background: var(--panel); color: var(--accent); font-family: var(--terminal-mono); text-decoration: none; }
.terminal-wiki-panel a:hover { background: var(--panel-2); }
.terminal-wiki-panel small { color: var(--sub); }
.terminal-empty { padding: 34px; border: 1px dashed var(--line); border-radius: 14px; color: var(--sub); font-family: var(--terminal-mono); text-align: center; }
.terminal-more { display: inline-flex; align-items: center; gap: 4px; color: var(--accent); font-family: var(--terminal-mono); text-decoration: none; }
.terminal-error { color: var(--amber); }
.terminal-error button { margin-left: 8px; padding: 5px 9px; border: 1px solid var(--line); border-radius: 7px; background: var(--panel); color: var(--accent); font: inherit; cursor: pointer; }

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

@media (max-width: 900px) {
  .terminal-photo-grid { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .terminal-wiki-panel { grid-template-columns: 1fr; }
}

/* Magazine adaptation */
.home {
  background: var(--theme-bg-primary);
  color: var(--theme-text-primary);
}

.home .container,
.hero-carousel .container {
  max-width: 1060px;
  padding-right: 32px;
  padding-left: 32px;
}

.issue-line {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 42px 0 22px;
  color: var(--theme-text-secondary);
  font-size: 11px;
  letter-spacing: .24em;
}

.issue-number {
  color: var(--theme-primary);
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: 13px;
  letter-spacing: .12em;
}

.issue-rule {
  width: 64px;
  border-top: 1px solid var(--theme-border);
}

.issue-directory {
  display: flex;
  gap: 18px;
  margin-left: auto;
}

.issue-directory a {
  color: var(--theme-text-secondary);
  letter-spacing: .12em;
}

.issue-directory a:hover {
  color: var(--theme-primary);
}

.hero-carousel :deep(.el-carousel) {
  border: 1px solid var(--theme-border);
  border-radius: 0;
  box-shadow: none;
}

.hero-slide-clickable {
  transition: filter .25s ease;
}

.hero-slide-clickable:hover {
  filter: saturate(.88);
  transform: none;
}

.hero-title {
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: clamp(34px, 5vw, 62px);
  font-weight: 500;
  letter-spacing: .04em;
  line-height: 1.25;
}

.hero-subtitle {
  letter-spacing: .08em;
}

.hero-default {
  align-items: flex-start;
  height: 320px;
  padding: 56px 8%;
  border: 1px solid var(--theme-border);
  border-radius: 0;
  text-align: left;
}

.editorial-hero {
  min-height: 390px;
  padding: 48px 0 58px;
  border-bottom: 3px double var(--theme-text-primary);
}

.editorial-hero .issue-line {
  padding: 0 0 32px;
}

.editorial-hero h1 {
  max-width: 800px;
  margin: 0;
  color: var(--theme-text-primary);
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: clamp(42px, 6.5vw, 72px);
  font-weight: 500;
  letter-spacing: .02em;
  line-height: 1.25;
}

.editorial-hero h1 span {
  color: var(--theme-primary);
}

.editorial-hero > p {
  max-width: 38em;
  margin: 24px 0 0;
  color: var(--theme-text-secondary);
  font-size: 16px;
  line-height: 1.85;
}

.hero-actions {
  display: flex;
  align-items: center;
  gap: 28px;
  margin-top: 32px;
}

.hero-primary {
  padding: 11px 24px;
  background: var(--theme-text-primary);
  color: var(--theme-bg-primary);
  font-size: 13px;
  text-decoration: none;
  transition: background .25s ease;
}

.hero-primary:hover {
  background: var(--theme-primary);
}

.hero-secondary {
  padding-bottom: 4px;
  border-bottom: 1px solid var(--theme-text-primary);
  color: var(--theme-text-primary);
  font-size: 13px;
  text-decoration: none;
}

.hero-primary:focus-visible,
.hero-secondary:focus-visible {
  outline: 2px solid var(--theme-primary);
  outline-offset: 3px;
}

.main-content {
  padding: 68px 0 76px;
}

.content-section,
.sidebar-card {
  margin-bottom: 52px;
}

.section-header {
  margin-bottom: 18px;
  padding-bottom: 18px;
  border-bottom: 1px solid var(--theme-border);
}

.section-header h2 {
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: 26px;
  font-weight: 500;
  letter-spacing: .06em;
}

.section-header h2::before {
  display: none;
}

.section-header small {
  margin-left: 10px;
  color: var(--theme-text-secondary);
  font-family: inherit;
  font-size: 11px;
  font-weight: 400;
  letter-spacing: .2em;
}

.article-list {
  gap: 0 !important;
  padding: 10px 32px;
  border: 1px solid var(--theme-border);
  background: var(--theme-bg-card);
}

.article-item,
.sidebar-card {
  border: 0 !important;
  border-bottom: 1px solid var(--theme-border) !important;
  border-radius: 0;
  box-shadow: none;
  background: transparent !important;
}

.work-card {
  border: 1px solid var(--theme-border) !important;
  border-radius: 0;
  background: var(--theme-bg-card) !important;
  box-shadow: none;
}

.article-item {
  height: auto;
  min-height: 116px;
  transition: padding-left .25s ease;
}

.article-item:hover {
  padding-left: 10px;
  box-shadow: none;
  transform: none;
}

.article-item:hover .article-title,
.work-card-clickable:hover h4,
.recommended-item:hover h4 {
  color: var(--theme-primary);
}

.article-item :deep(.el-card__body) {
  padding: 16px 8px;
}

.article-title,
.work-info h4,
.recommended-item h4,
.work-title {
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-weight: 500;
  letter-spacing: .03em;
}

.work-cover {
  border-radius: 0;
  filter: saturate(.82) contrast(.96);
}

.article-cover {
  align-self: center;
  border-radius: 8px;
  filter: saturate(.82) contrast(.96);
}

.work-card-clickable:hover {
  border-color: var(--theme-primary);
  box-shadow: none;
  transform: translateY(-4px);
}

.work-info {
  padding: 22px 24px 24px;
}

.sidebar-card :deep(.el-card__header),
.sidebar-card :deep(.el-card__body) {
  padding-right: 2px;
  padding-left: 2px;
}

.sidebar-card :deep(.el-card__header) {
  padding-top: 22px;
  padding-bottom: 16px;
  border-top: 1px solid var(--theme-border);
}

.card-header {
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-weight: 500;
  letter-spacing: .06em;
}

.stat-item,
.recommended-item {
  background: transparent;
  border: 0;
  border-radius: 0;
}

.stats-grid {
  gap: 0;
  border-top: 1px solid var(--theme-border);
  border-left: 1px solid var(--theme-border);
}

.stat-item {
  padding: 18px 12px;
  border-right: 1px solid var(--theme-border);
  border-bottom: 1px solid var(--theme-border);
}

.recommended-list {
  gap: 0;
}

.recommended-item {
  padding: 16px 2px;
  border-bottom: 1px solid var(--theme-border);
}

.stat-value {
  color: var(--theme-primary);
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-weight: 500;
}

.recommended-item:hover {
  padding-left: 8px;
  background: var(--theme-bg-hover);
  transform: none;
}

.recommended-work-item {
  border-radius: 0;
}

.recommended-work-item:hover {
  box-shadow: none;
  transform: translateX(4px);
}

.home :deep(.el-button),
.home :deep(.el-tag) {
  border-radius: 1px;
}

.home :deep(.el-card:focus-visible),
.issue-directory a:focus-visible {
  outline: 2px solid var(--theme-primary);
  outline-offset: 3px;
}

@media (max-width: 900px) {
  .home .container,
  .hero-carousel .container {
    padding-right: 24px;
    padding-left: 24px;
  }

  .sidebar {
    position: static;
    margin-top: 42px;
  }
}

@media (max-width: 560px) {
  .article-item :deep(.el-card__body) { padding: 24px 2px; }
  .work-info { padding: 20px 20px 22px; }
  .article-list { padding: 4px 16px; }
}

@media (max-width: 560px) {
  .home .container,
  .hero-carousel .container {
    padding-right: 18px;
    padding-left: 18px;
  }

  .issue-line {
    flex-wrap: wrap;
    gap: 10px;
    padding-top: 28px;
  }

  .issue-rule {
    width: 32px;
  }

  .issue-directory {
    width: 100%;
    margin: 8px 0 0;
  }

  .hero-carousel {
    padding-top: 0;
  }

  .hero-title {
    font-size: 32px;
  }

  .editorial-hero {
    min-height: 0;
    padding: 34px 0 42px;
  }

  .editorial-hero .issue-line {
    padding-top: 0;
    padding-bottom: 24px;
  }

  .editorial-hero h1 {
    font-size: 40px;
  }

  .editorial-hero > p {
    font-size: 15px;
  }

  .hero-actions {
    align-items: flex-start;
    flex-direction: column;
    gap: 18px;
  }

  .main-content {
    padding: 48px 0;
  }

  .section-header {
    align-items: flex-end;
  }

  .section-header small {
    display: block;
    margin: 5px 0 0;
  }

  .article-content {
    align-items: stretch;
  }

  .article-item :deep(.el-card__body) {
    padding: 24px 0;
  }

  .terminal-photo-grid { grid-template-columns: 1fr; }
}
</style>
