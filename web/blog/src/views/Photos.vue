<template>
  <main class="photos-page">
    <header class="photos-hero">
      <div class="photos-shell hero-grid">
        <div>
          <div
            class="issue-line"
            aria-label="摄影刊物"
          >
            <span class="issue-no">LENS / 01</span>
            <span
              class="issue-rule"
              aria-hidden="true"
            />
            <span>Photography journal</span>
          </div>
          <h1>光落下的地方，<br><em>都有故事。</em></h1>
          <p class="hero-intro">
            收录行走途中看见的光影、面孔与地景。每一组照片，都是一次凝视留下的短章。
          </p>
        </div>
      </div>
    </header>

    <section
      class="photos-journal"
      aria-labelledby="photos-heading"
    >
      <div class="photos-shell">
        <div class="section-heading">
          <div>
            <p class="section-kicker">
              Archive
            </p>
            <h2 id="photos-heading">
              摄影集 <span>PHOTOGRAPHY</span>
            </h2>
          </div>
          <p
            v-if="!loading && !error"
            class="archive-count"
            aria-live="polite"
          >
            共 {{ total }} 组影像
          </p>
        </div>

        <div
          v-if="loading"
          class="photo-grid"
          aria-busy="true"
          aria-label="正在加载摄影作品"
        >
          <div
            v-for="index in pageSize"
            :key="index"
            class="photo-card skeleton-card"
          >
            <div class="skeleton-image" />
            <div class="skeleton-line" />
            <div class="skeleton-line skeleton-line-short" />
          </div>
        </div>

        <div
          v-else-if="error"
          class="state-panel"
          role="alert"
        >
          <p class="state-label">
            LOAD FAILED
          </p>
          <h2>影集暂时无法打开</h2>
          <p>请检查网络连接后重试。</p>
          <button
            type="button"
            class="retry-button"
            @click="loadPhotos"
          >
            重新加载
          </button>
        </div>

        <div
          v-else-if="photos.length === 0"
          class="state-panel"
        >
          <p class="state-label">
            NO PHOTOGRAPHS
          </p>
          <h2>这里还没有公开的摄影作品</h2>
          <p>新的光影故事发布后，会陈列在这里。</p>
        </div>

        <div
          v-else
          class="photo-grid"
        >
          <article
            v-for="work in photos"
            :key="work.id"
            class="photo-card"
          >
            <router-link
              class="photo-link"
              :to="`/works/${work.id}`"
              :aria-label="`查看摄影作品《${work.title}》`"
            >
              <figure>
                <div class="photo-media">
                  <img
                    v-if="getImageUrl(work)"
                    :src="getImageUrl(work)"
                    :alt="getImageAlt(work)"
                    loading="lazy"
                    decoding="async"
                    @error="handleImageError"
                  >
                  <span
                    class="image-fallback"
                    aria-hidden="true"
                  >NO IMAGE</span>
                  <span
                    class="film-grain"
                    aria-hidden="true"
                  />
                  <figcaption>
                    <span>{{ work.metadata?.location || '未标注地点' }}</span>
                    <span v-if="getYear(work)">{{ getYear(work) }}</span>
                  </figcaption>
                </div>
                <div class="photo-copy">
                  <div class="title-row">
                    <h3>{{ work.title }}</h3>
                    <span aria-hidden="true">↗</span>
                  </div>
                  <div class="photo-meta">
                    <span>{{ getAuthor(work) }}</span>
                    <span v-if="work.metadata?.location">{{ work.metadata.location }}</span>
                    <time
                      v-if="getYear(work)"
                      :datetime="getDateTime(work)"
                    >{{ getYear(work) }}</time>
                  </div>
                  <div
                    class="photo-counts"
                    :aria-label="getCountsLabel(work)"
                  >
                    <span>{{ getPhotoCount(work) }} 张</span>
                    <span>{{ formatCount(work.view_count) }} 浏览</span>
                    <span>{{ formatCount(work.like_count) }} 喜欢</span>
                    <span>{{ formatCount(work.comment_count) }} 评论</span>
                  </div>
                </div>
              </figure>
            </router-link>
          </article>
        </div>

        <nav
          v-if="!loading && !error && total > pageSize"
          class="pagination"
          aria-label="摄影作品分页"
        >
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="pageSize"
            :total="total"
            layout="prev, pager, next"
            @current-change="handlePageChange"
          />
        </nav>
      </div>
    </section>
  </main>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import api from '@/utils/api'

const photos = ref([])
const currentPage = ref(1)
const pageSize = 12
const total = ref(0)
const loading = ref(true)
const error = ref(false)
let requestId = 0

const loadPhotos = async () => {
  const activeRequest = ++requestId
  loading.value = true
  error.value = false

  try {
    const response = await api.get('/works', {
      params: {
        type: 'photography',
        status: 1,
        page: currentPage.value,
        page_size: pageSize
      }
    })

    if (activeRequest !== requestId) return
    photos.value = response.data?.list || []
    total.value = response.data?.total || 0
  } catch (loadError) {
    if (activeRequest !== requestId) return
    console.error('Failed to load photography works:', loadError)
    photos.value = []
    total.value = 0
    error.value = true
  } finally {
    if (activeRequest === requestId) loading.value = false
  }
}

const getImageUrl = (work) => {
  if (work.cover) return work.cover
  const firstImage = work.images?.[0]
  return typeof firstImage === 'string' ? firstImage : firstImage?.url || ''
}

const getImageAlt = (work) => {
  const location = work.metadata?.location ? `，拍摄于${work.metadata.location}` : ''
  return `${work.title}${location}`
}

const getAuthor = (work) => work.author?.nickname || work.author?.username || '匿名摄影者'

const getYear = (work) => {
  const match = String(work.metadata?.shooting_date || '').match(/\d{4}/)
  return match?.[0] || ''
}

const getDateTime = (work) => work.metadata?.shooting_date || getYear(work)

const getPhotoCount = (work) => {
  const count = Number(work.metadata?.photo_count)
  return Number.isFinite(count) && count > 0 ? count : work.images?.length || 1
}

const formatCount = (count) => Number(count || 0).toLocaleString('zh-CN')

const getCountsLabel = (work) => [
  `${getPhotoCount(work)} 张照片`,
  `${formatCount(work.view_count)} 次浏览`,
  `${formatCount(work.like_count)} 个喜欢`,
  `${formatCount(work.comment_count)} 条评论`
].join('，')

const handleImageError = (event) => {
  event.currentTarget.hidden = true
  event.currentTarget.parentElement?.classList.add('has-image-error')
}

const handlePageChange = () => {
  loadPhotos()
  document.querySelector('.photos-journal')?.scrollIntoView({ behavior: 'smooth' })
}

onMounted(loadPhotos)
</script>

<style scoped>
.photos-page {
  min-height: 100vh;
  background: var(--bg, var(--theme-bg-primary));
  color: var(--ink, var(--theme-text-primary));
}

.photos-shell {
  width: min(1124px, calc(100% - 64px));
  margin: 0 auto;
}

.photos-hero {
  padding: 92px 0 84px;
}

.hero-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 64px;
  align-items: start;
}

.issue-line {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 32px;
  color: var(--sub, var(--theme-text-secondary));
  font-size: 11px;
  letter-spacing: .24em;
  text-transform: uppercase;
}

.issue-no {
  color: var(--accent, var(--theme-primary));
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  letter-spacing: .25em;
}

.issue-rule {
  width: 64px;
  height: 1px;
  background: var(--hairline, var(--theme-border));
}

.photos-hero h1 {
  max-width: 850px;
  margin: 0;
  color: var(--ink, var(--theme-text-primary));
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: clamp(42px, 6.3vw, 76px);
  font-weight: 500;
  line-height: 1.24;
  letter-spacing: .025em;
}

.photos-hero h1 em {
  color: var(--accent, var(--theme-primary));
  font-style: normal;
}

.hero-intro {
  max-width: 36em;
  margin: 32px 0 0;
  color: var(--sub, var(--theme-text-secondary));
  font-size: 15px;
  line-height: 1.9;
}

.photos-journal {
  padding: 72px 0 88px;
  border-top: 1px solid var(--hairline, var(--theme-border));
}

.section-heading {
  display: flex;
  align-items: end;
  justify-content: space-between;
  gap: 24px;
  margin-bottom: 42px;
}

.section-kicker,
.state-label {
  margin: 0 0 9px;
  color: var(--accent, var(--theme-primary));
  font-size: 10px;
  letter-spacing: .3em;
  line-height: 1;
  text-transform: uppercase;
}

.section-heading h2 {
  margin: 0;
  color: var(--ink, var(--theme-text-primary));
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: 28px;
  font-weight: 500;
  letter-spacing: .05em;
  line-height: 1.3;
}

.section-heading h2 span {
  margin-left: 14px;
  color: var(--sub, var(--theme-text-secondary));
  font-family: 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  font-size: 10px;
  font-weight: 400;
  letter-spacing: .22em;
}

.archive-count {
  margin: 0 0 3px;
  color: var(--sub, var(--theme-text-secondary));
  font-size: 12px;
  letter-spacing: .1em;
}

.photo-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 54px 20px;
  align-items: start;
}

.photo-card {
  min-width: 0;
}

.photo-card:nth-child(2) .photo-media {
  aspect-ratio: 4 / 6.15;
}

.photo-card:nth-child(5) {
  margin-top: -38px;
}

.photo-link {
  display: block;
  color: inherit;
}

.photo-link:hover {
  color: inherit;
}

.photo-card figure {
  margin: 0;
}

.photo-media {
  position: relative;
  aspect-ratio: 4 / 5;
  overflow: hidden;
  background: var(--bg-soft, var(--theme-bg-secondary));
}

.photo-media img {
  position: relative;
  z-index: 1;
  width: 100%;
  height: 100%;
  display: block;
  object-fit: cover;
  transition: transform .6s cubic-bezier(.2, .6, .2, 1), filter .35s ease;
}

.photo-link:hover .photo-media img {
  filter: saturate(.92) contrast(1.03);
  transform: scale(1.045);
}

.image-fallback {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  color: var(--sub, var(--theme-text-secondary));
  font-family: Georgia, serif;
  font-size: 11px;
  letter-spacing: .26em;
}

.film-grain {
  position: absolute;
  z-index: 2;
  inset: 0;
  pointer-events: none;
  opacity: .18;
  background-image: radial-gradient(rgba(255, 255, 255, .6) .6px, transparent .7px);
  background-size: 5px 5px;
  mix-blend-mode: overlay;
}

.photo-media figcaption {
  position: absolute;
  z-index: 3;
  inset: auto 0 0;
  display: flex;
  justify-content: space-between;
  gap: 16px;
  padding: 54px 18px 16px;
  background: linear-gradient(transparent, rgba(13, 15, 14, .72));
  color: #fff;
  font-size: 11px;
  letter-spacing: .16em;
  opacity: 0;
  transform: translateY(6px);
  transition: opacity .35s ease, transform .35s ease;
}

.photo-link:hover figcaption,
.photo-link:focus-visible figcaption,
.has-image-error figcaption {
  opacity: 1;
  transform: none;
}

.photo-copy {
  padding-top: 17px;
}

.title-row {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 16px;
}

.title-row h3 {
  margin: 0;
  overflow: hidden;
  color: var(--ink, var(--theme-text-primary));
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: 19px;
  font-weight: 500;
  line-height: 1.5;
  letter-spacing: .035em;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: color .25s ease;
}

.title-row > span {
  flex: none;
  color: var(--sub, var(--theme-text-secondary));
  font-size: 16px;
  transition: color .25s ease, transform .25s ease;
}

.photo-link:hover h3,
.photo-link:focus-visible h3,
.photo-link:hover .title-row > span,
.photo-link:focus-visible .title-row > span {
  color: var(--accent, var(--theme-primary));
}

.photo-link:hover .title-row > span,
.photo-link:focus-visible .title-row > span {
  transform: translate(3px, -3px);
}

.photo-meta,
.photo-counts {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0;
  color: var(--sub, var(--theme-text-secondary));
  font-size: 11px;
  line-height: 1.6;
  letter-spacing: .06em;
}

.photo-meta {
  margin-top: 7px;
}

.photo-counts {
  margin-top: 13px;
  padding-top: 11px;
  border-top: 1px solid var(--hairline, var(--theme-border));
}

.photo-meta span:not(:last-child)::after,
.photo-meta time:not(:last-child)::after,
.photo-counts span:not(:last-child)::after {
  content: '·';
  margin: 0 8px;
  color: var(--hairline, var(--theme-border));
}

.state-panel {
  padding: 92px 24px 98px;
  border-top: 1px solid var(--hairline, var(--theme-border));
  border-bottom: 1px solid var(--hairline, var(--theme-border));
  text-align: center;
}

.state-panel h2 {
  margin: 14px 0 12px;
  color: var(--ink, var(--theme-text-primary));
  font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif;
  font-size: 25px;
  font-weight: 500;
}

.state-panel > p:not(.state-label) {
  margin: 0;
  color: var(--sub, var(--theme-text-secondary));
  font-size: 14px;
}

.retry-button {
  margin-top: 28px;
  padding: 11px 24px;
  border: 1px solid var(--ink, var(--theme-text-primary));
  background: var(--ink, var(--theme-text-primary));
  color: var(--bg, var(--theme-bg-primary));
  cursor: pointer;
  font: inherit;
  font-size: 13px;
  letter-spacing: .12em;
  transition: background .25s ease, border-color .25s ease;
}

.retry-button:hover {
  border-color: var(--accent, var(--theme-primary));
  background: var(--accent, var(--theme-primary));
}

.skeleton-image,
.skeleton-line {
  background: linear-gradient(100deg, var(--bg-soft, var(--theme-bg-secondary)) 30%, var(--surface, var(--theme-bg-card)) 50%, var(--bg-soft, var(--theme-bg-secondary)) 70%);
  background-size: 220% 100%;
  animation: skeleton-shift 1.4s ease-in-out infinite;
}

.skeleton-image {
  aspect-ratio: 4 / 5;
}

.skeleton-card:nth-child(2) .skeleton-image {
  aspect-ratio: 4 / 6.15;
}

.skeleton-line {
  width: 70%;
  height: 15px;
  margin-top: 18px;
}

.skeleton-line-short {
  width: 42%;
  height: 10px;
  margin-top: 10px;
}

@keyframes skeleton-shift {
  to { background-position-x: -220%; }
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 76px;
  padding-top: 30px;
  border-top: 1px solid var(--hairline, var(--theme-border));
}

.pagination :deep(.el-pager li),
.pagination :deep(button) {
  background: transparent !important;
  font-family: Georgia, serif;
}

@media (max-width: 900px) {
  .photos-shell {
    width: min(100% - 48px, 760px);
  }

  .photos-hero {
    padding: 72px 0 66px;
  }

  .hero-grid {
    grid-template-columns: 1fr;
  }

  .photos-journal {
    padding: 60px 0 72px;
  }

  .photo-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 48px 20px;
  }

  .photo-card:nth-child(5) {
    margin-top: 0;
  }
}

@media (max-width: 560px) {
  .photos-shell {
    width: calc(100% - 30px);
  }

  .photos-hero {
    padding: 54px 0 50px;
  }

  .issue-line {
    gap: 10px;
    margin-bottom: 24px;
    font-size: 9px;
    letter-spacing: .18em;
  }

  .issue-rule {
    width: 28px;
  }

  .photos-hero h1 {
    font-size: clamp(36px, 11vw, 48px);
    line-height: 1.3;
  }

  .hero-intro {
    margin-top: 24px;
    font-size: 14px;
  }

  .photos-journal {
    padding: 48px 0 60px;
  }

  .section-heading {
    align-items: start;
    margin-bottom: 30px;
  }

  .section-heading h2 span {
    display: block;
    margin: 7px 0 0;
  }

  .photo-grid {
    grid-template-columns: 1fr;
    gap: 44px;
  }

  .photo-card:nth-child(2) .photo-media,
  .skeleton-card:nth-child(2) .skeleton-image {
    aspect-ratio: 4 / 5;
  }

  .photo-media figcaption {
    opacity: 1;
    transform: none;
  }

  .pagination {
    margin-top: 58px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .photo-media img,
  .photo-media figcaption,
  .title-row h3,
  .title-row > span,
  .retry-button {
    transition: none;
  }

  .photo-link:hover .photo-media img,
  .photo-link:hover .title-row > span {
    transform: none;
  }

  .skeleton-image,
  .skeleton-line {
    animation: none;
  }
}
</style>
