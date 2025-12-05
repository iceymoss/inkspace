<template>
  <div class="work-detail">
    <div class="container">
      <el-card v-if="work" class="detail-card">
        <!-- 摄影作品布局 -->
        <template v-if="work.type === 'photography'">
          <div class="photography-layout">
            <!-- 左侧：图片展示区域（70%） -->
            <div class="photo-section">
              <!-- 大图轮播 -->
              <div class="main-photo-container">
                <el-carousel 
                  ref="carouselRef"
                  :height="carouselHeight + 'px'"
                  arrow="always"
                  indicator-position="none"
                  @change="handlePhotoChange"
                >
                  <el-carousel-item v-for="(photo, index) in photos" :key="index">
                    <el-image 
                      :src="photo.url" 
                      :alt="photo.description || work.title"
                      fit="contain"
                      class="main-photo"
                      @load="handleImageLoad"
                    />
                  </el-carousel-item>
                </el-carousel>
                
                <!-- 照片计数 -->
                <div class="photo-counter">
                  {{ currentPhotoIndex + 1 }} / {{ photos.length }}
                </div>
              </div>

              <!-- 缩略图导航 -->
              <div class="thumbnail-nav" v-if="photos.length > 1">
                <div 
                  v-for="(photo, index) in photos" 
                  :key="index"
                  class="thumbnail-item"
                  :class="{ active: index === currentPhotoIndex }"
                  @click="setActivePhoto(index)"
                >
                  <el-image :src="photo.url" fit="cover" />
                </div>
              </div>

              <!-- 当前照片描述 -->
              <div class="photo-description" v-if="currentPhoto.description">
                <p>{{ currentPhoto.description }}</p>
              </div>

              <!-- 相册信息 -->
              <div class="album-info">
                <h2>{{ work.title }}</h2>
                <p class="album-description">{{ work.description }}</p>
                
                <div class="album-meta" v-if="work.metadata">
                  <div class="meta-item" v-if="work.metadata.location">
                    <el-icon><Location /></el-icon>
                    <span>{{ work.metadata.location }}</span>
                  </div>
                  <div class="meta-item" v-if="work.metadata.shooting_date">
                    <el-icon><Calendar /></el-icon>
                    <span>{{ work.metadata.shooting_date }}</span>
                  </div>
                  <div class="meta-item">
                    <el-icon><Picture /></el-icon>
                    <span>共 {{ photos.length }} 张照片</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 右侧：信息栏（30%） -->
            <div class="info-sidebar">
              <!-- 作者信息 -->
              <div class="author-card" v-if="work.author">
                <div class="author-header" @click="goToUserProfile(work.author.id)">
                  <el-avatar :size="60" :src="work.author.avatar" />
                  <div class="author-info">
                    <h3>{{ work.author.nickname || work.author.username }}</h3>
                    <p>{{ work.author.bio || '这个人很懒，什么都没留下' }}</p>
                  </div>
                </div>
                <el-button type="primary" style="width: 100%">关注</el-button>
              </div>

              <!-- 当前照片参数 -->
              <div class="photo-params-card" v-if="currentPhoto.metadata">
                <h3>📷 拍摄参数</h3>
                <div class="params-list">
                  <div class="param-item" v-if="currentPhoto.metadata.camera">
                    <span class="param-label">相机：</span>
                    <span class="param-value">{{ currentPhoto.metadata.camera }}</span>
                  </div>
                  <div class="param-item" v-if="currentPhoto.metadata.lens">
                    <span class="param-label">镜头：</span>
                    <span class="param-value">{{ currentPhoto.metadata.lens }}</span>
                  </div>
                  <div class="param-item" v-if="currentPhoto.metadata.focal_length">
                    <span class="param-label">焦段：</span>
                    <span class="param-value">{{ currentPhoto.metadata.focal_length }}</span>
                  </div>
                  <div class="param-item" v-if="currentPhoto.metadata.aperture">
                    <span class="param-label">光圈：</span>
                    <span class="param-value">{{ currentPhoto.metadata.aperture }}</span>
                  </div>
                  <div class="param-item" v-if="currentPhoto.metadata.shutter_speed">
                    <span class="param-label">快门：</span>
                    <span class="param-value">{{ currentPhoto.metadata.shutter_speed }}</span>
                  </div>
                  <div class="param-item" v-if="currentPhoto.metadata.iso">
                    <span class="param-label">ISO：</span>
                    <span class="param-value">{{ currentPhoto.metadata.iso }}</span>
                  </div>
                </div>
              </div>

              <!-- 互动统计 -->
              <div class="stats-card">
                <h3>互动数据</h3>
                <div class="stats-list">
                  <div class="stat-item">
                    <el-icon><View /></el-icon>
                    <span>{{ work.view_count }} 浏览</span>
                  </div>
                  <div class="stat-item">
                    <el-icon><ChatDotRound /></el-icon>
                    <span>{{ work.comment_count }} 评论</span>
                  </div>
                </div>
                <div class="action-buttons">
                  <el-button 
                    :type="isLiked ? 'primary' : 'default'"
                    @click="handleLike"
                    :loading="liking"
                  >
                    <el-icon><Star /></el-icon>
                    {{ work.like_count || 0 }} {{ isLiked ? '已点赞' : '点赞' }}
                  </el-button>
                  <el-button 
                    :type="isFavorited ? 'warning' : 'default'"
                    @click="handleFavorite"
                    :loading="favoriting"
                  >
                    <el-icon><Star /></el-icon>
                    {{ work.favorite_count || 0 }} {{ isFavorited ? '已收藏' : '收藏' }}
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </template>

        <!-- 项目作品布局 -->
        <template v-else>
          <div class="project-layout">
            <div class="work-header">
              <div class="work-title-section">
                <h1>{{ work.title }}</h1>
                <el-tag :type="work.type === 'photography' ? 'warning' : 'primary'">
                  {{ work.type === 'photography' ? '📷 摄影作品' : '💻 开源项目' }}
                </el-tag>
              </div>
              
              <div class="work-author" v-if="work.author">
                <el-avatar 
                  :size="40" 
                  :src="work.author.avatar"
                  @click="goToUserProfile(work.author.id)"
                  style="cursor: pointer"
                />
                <div class="author-info">
                  <span class="author-name" @click="goToUserProfile(work.author.id)">
                    {{ work.author.nickname || work.author.username }}
                  </span>
                  <span class="publish-time">{{ formatDate(work.created_at) }}</span>
                </div>
              </div>
            </div>

            <el-image :src="work.cover" fit="cover" class="work-cover" />

            <div class="work-content">
              <p class="work-description">{{ work.description }}</p>
              
              <div class="work-links" v-if="work.github_url || work.demo_url || work.link">
                <el-link 
                  v-if="work.github_url" 
                  :href="work.github_url" 
                  target="_blank" 
                  type="primary"
                  :icon="Link"
                >
                  GitHub
                </el-link>
                <el-link 
                  v-if="work.demo_url" 
                  :href="work.demo_url" 
                  target="_blank" 
                  type="success"
                  :icon="Link"
                >
                  在线演示
                </el-link>
                <el-link 
                  v-if="work.link" 
                  :href="work.link" 
                  target="_blank"
                  :icon="Link"
                >
                  项目主页
                </el-link>
              </div>

              <div class="tech-stack" v-if="work.tech_stack">
                <h3>技术栈</h3>
                <p>{{ work.tech_stack }}</p>
              </div>

              <div class="work-stats">
                <span><el-icon><View /></el-icon> {{ work.view_count }}</span>
                <span><el-icon><ChatDotRound /></el-icon> {{ work.comment_count }}</span>
              </div>
            </div>
          </div>
        </template>

        <!-- 评论区（两种类型共用） -->
        <el-divider />
        
        <div class="comment-section">
          <h3>评论 ({{ work.comment_count }})</h3>
          
          <!-- 发表评论 -->
          <div class="comment-form" v-if="userStore.isLoggedIn">
            <el-input
              v-model="commentContent"
              type="textarea"
              :rows="3"
              placeholder="写下你的评论..."
              maxlength="500"
              show-word-limit
            />
            <div class="comment-actions">
              <el-button type="primary" @click="submitComment" :loading="submittingComment">
                发表评论
              </el-button>
            </div>
          </div>
          <div v-else class="login-tip">
            <el-link type="primary" @click="$router.push('/login')">登录</el-link> 后发表评论
          </div>

          <!-- 评论列表 -->
          <div class="comments-list">
            <div 
              v-for="comment in comments" 
              :key="comment.id" 
              class="comment-item"
            >
              <el-avatar 
                :src="comment.user?.avatar" 
                :size="40"
                @click="goToUserProfile(comment.user?.id)"
                style="cursor: pointer"
              />
              <div class="comment-content">
                <div class="comment-header">
                  <div class="comment-author-section">
                    <span 
                      class="comment-author"
                      @click="goToUserProfile(comment.user?.id)"
                    >
                      {{ comment.user?.nickname || comment.user?.username }}
                    </span>
                    <el-tag v-if="comment.user_id === work?.author_id" type="warning" size="small">
                      作者
                    </el-tag>
                  </div>
                  <div class="comment-meta-section">
                    <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                    <el-button 
                      text 
                      size="small"
                      type="danger"
                      v-if="userStore.user?.id === comment.user_id"
                      @click="deleteComment(comment.id)"
                    >
                      删除
                    </el-button>
                  </div>
                </div>
                <p class="comment-text">{{ comment.content }}</p>
              </div>
            </div>
          </div>

          <div v-if="hasMoreComments" class="load-more">
            <el-button text @click="loadMoreComments">加载更多评论</el-button>
          </div>
        </div>
      </el-card>

      <el-empty v-else description="作品不存在" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  View, ChatDotRound, Star, Link, 
  Location, Calendar, Picture 
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const work = ref(null)
const comments = ref([])
const commentContent = ref('')
const submittingComment = ref(false)
const commentPage = ref(1)
const commentPageSize = ref(10)
const totalComments = ref(0)

// 摄影相关
const carouselRef = ref()
const currentPhotoIndex = ref(0)
const carouselHeight = ref(600)

// 点赞和收藏状态
const isLiked = ref(false)
const isFavorited = ref(false)
const liking = ref(false)
const favoriting = ref(false)

const photos = computed(() => {
  if (!work.value || work.value.type !== 'photography') return []
  return work.value.images || []
})

const currentPhoto = computed(() => {
  return photos.value[currentPhotoIndex.value] || { metadata: {} }
})

const hasMoreComments = computed(() => {
  return comments.value.length < totalComments.value
})

const loadWork = async () => {
  try {
    const response = await api.get(`/works/${route.params.id}`)
    work.value = response.data
    totalComments.value = response.data.comment_count || 0
  } catch (error) {
    ElMessage.error('加载作品失败')
  }
}

const loadComments = async () => {
  try {
    const response = await api.get('/comments', {
      params: {
        work_id: route.params.id,
        page: commentPage.value,
        page_size: commentPageSize.value
      }
    })
    
    if (commentPage.value === 1) {
      comments.value = response.data.list || []
    } else {
      comments.value.push(...(response.data.list || []))
    }
  } catch (error) {
    console.error('Failed to load comments:', error)
  }
}

const loadMoreComments = () => {
  commentPage.value++
  loadComments()
}

const submitComment = async () => {
  if (!commentContent.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  submittingComment.value = true
  try {
    await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: commentContent.value
    })
    
    ElMessage.success('评论成功')
    commentContent.value = ''
    commentPage.value = 1
    await loadComments()
    work.value.comment_count++
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '评论失败')
  } finally {
    submittingComment.value = false
  }
}

const deleteComment = async (commentId) => {
  try {
    await api.delete(`/comments/${commentId}`)
    ElMessage.success('删除成功')
    comments.value = comments.value.filter(c => c.id !== commentId)
    work.value.comment_count--
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const handlePhotoChange = (index) => {
  currentPhotoIndex.value = index
}

const setActivePhoto = (index) => {
  carouselRef.value.setActiveItem(index)
}

const handleImageLoad = () => {
  // 图片加载完成后可以调整高度
}

const goToUserProfile = (userId) => {
  if (userId) {
    router.push(`/user/${userId}`)
  }
}

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  const now = new Date()
  const diff = now - d
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`
  
  return d.toLocaleDateString('zh-CN')
}

// 响应式调整轮播高度
const updateCarouselHeight = () => {
  const width = window.innerWidth
  if (width < 768) {
    carouselHeight.value = 300
  } else if (width < 1200) {
    carouselHeight.value = 450
  } else {
    carouselHeight.value = 600
  }
}

const checkLikedStatus = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const response = await api.get(`/works/${route.params.id}/liked`)
    isLiked.value = response.data.liked
  } catch (error) {
    console.error('Failed to check liked status:', error)
  }
}

const checkFavoritedStatus = async () => {
  if (!userStore.isLoggedIn) return
  try {
    const response = await api.get(`/works/${route.params.id}/favorited`)
    isFavorited.value = response.data.favorited
  } catch (error) {
    console.error('Failed to check favorited status:', error)
  }
}

const handleLike = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  liking.value = true
  try {
    await api.post(`/works/${route.params.id}/like`)
    isLiked.value = !isLiked.value
    work.value.like_count = (work.value.like_count || 0) + (isLiked.value ? 1 : -1)
    ElMessage.success(isLiked.value ? '点赞成功' : '取消点赞')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    liking.value = false
  }
}

const handleFavorite = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  favoriting.value = true
  try {
    if (isFavorited.value) {
      await api.delete(`/works/${route.params.id}/favorite`)
    } else {
      await api.post(`/works/${route.params.id}/favorite`)
    }
    isFavorited.value = !isFavorited.value
    work.value.favorite_count = (work.value.favorite_count || 0) + (isFavorited.value ? 1 : -1)
    ElMessage.success(isFavorited.value ? '收藏成功' : '取消收藏')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
  } finally {
    favoriting.value = false
  }
}

onMounted(() => {
  loadWork()
  loadComments()
  checkLikedStatus()
  checkFavoritedStatus()
  updateCarouselHeight()
  window.addEventListener('resize', updateCarouselHeight)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateCarouselHeight)
})
</script>

<style scoped>
.work-detail {
  padding: 40px 0;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.detail-card {
  max-width: 1400px;
  margin: 0 auto;
}

/* ========== 摄影作品布局 ========== */
.photography-layout {
  display: flex;
  gap: 30px;
}

.photo-section {
  flex: 7;
}

.main-photo-container {
  position: relative;
  background: #000;
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 20px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.main-photo {
  width: 100%;
  height: 100%;
}

.photo-counter {
  position: absolute;
  bottom: 24px;
  right: 24px;
  background: rgba(0, 0, 0, 0.75);
  color: white;
  padding: 10px 20px;
  border-radius: 24px;
  font-size: 15px;
  font-weight: 500;
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.thumbnail-nav {
  display: flex;
  gap: 10px;
  overflow-x: auto;
  padding: 10px 0;
}

.thumbnail-item {
  flex-shrink: 0;
  width: 80px;
  height: 80px;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  border: 3px solid transparent;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.thumbnail-item:hover {
  border-color: #409eff;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.3);
}

.thumbnail-item.active {
  border-color: #409eff;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.2);
  transform: scale(1.05);
}

.thumbnail-item .el-image {
  width: 100%;
  height: 100%;
}

.photo-description {
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
  margin-bottom: 24px;
  border-left: 4px solid #409eff;
  font-size: 15px;
  line-height: 1.8;
  color: #606266;
}

.album-info {
  padding: 20px 0;
}

.album-info h2 {
  margin-bottom: 15px;
}

.album-description {
  color: var(--text-secondary);
  margin-bottom: 20px;
  line-height: 1.6;
}

.album-meta {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
}

/* 右侧信息栏 */
.info-sidebar {
  flex: 3;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.author-card,
.photo-params-card,
.stats-card {
  padding: 24px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  border: 1px solid #ebeef5;
}

.author-header {
  display: flex;
  gap: 15px;
  margin-bottom: 15px;
  cursor: pointer;
}

.author-info h3 {
  margin: 0 0 5px 0;
  font-size: 1.1rem;
}

.author-info p {
  margin: 0;
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.photo-params-card h3,
.stats-card h3 {
  margin: 0 0 20px 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: #303133;
  padding-bottom: 12px;
  border-bottom: 2px solid #f0f0f0;
}

.params-list,
.stats-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.param-item {
  display: flex;
  justify-content: space-between;
  font-size: 0.95rem;
  padding: 10px 0;
  border-bottom: 1px solid #f5f5f5;
}

.param-item:last-child {
  border-bottom: none;
}

.param-label {
  color: #909399;
  font-weight: 500;
}

.param-value {
  font-weight: 600;
  color: #303133;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.action-buttons {
  display: flex;
  gap: 12px;
  margin-top: 20px;
}

.action-buttons .el-button {
  flex: 1;
  height: 44px;
  font-size: 15px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s;
}

.action-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* ========== 项目作品布局 ========== */
.project-layout {
  max-width: 900px;
  margin: 0 auto;
}

.work-header {
  margin-bottom: 30px;
}

.work-title-section {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 20px;
}

.work-title-section h1 {
  margin: 0;
  font-size: 2rem;
}

.work-author {
  display: flex;
  align-items: center;
  gap: 12px;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: 500;
  cursor: pointer;
  color: var(--el-color-primary);
}

.author-name:hover {
  text-decoration: underline;
}

.publish-time {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.work-cover {
  width: 100%;
  border-radius: 8px;
  margin-bottom: 30px;
}

.work-content {
  line-height: 1.8;
}

.work-description {
  font-size: 1.1rem;
  margin-bottom: 25px;
}

.work-links {
  display: flex;
  gap: 15px;
  margin-bottom: 25px;
}

.tech-stack {
  margin-bottom: 25px;
}

.tech-stack h3 {
  margin-bottom: 10px;
}

.work-stats {
  display: flex;
  gap: 20px;
  color: var(--text-secondary);
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.work-stats span {
  display: flex;
  align-items: center;
  gap: 5px;
}

/* ========== 评论区 ========== */
.comment-section {
  margin-top: 40px;
}

.comment-section h3 {
  margin-bottom: 20px;
}

.comment-form {
  margin-bottom: 30px;
}

.comment-actions {
  margin-top: 10px;
  display: flex;
  justify-content: flex-end;
}

.login-tip {
  text-align: center;
  padding: 20px;
  color: var(--text-secondary);
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.comment-item {
  display: flex;
  gap: 15px;
  padding: 20px 0;
  border-bottom: 1px solid #f0f0f0;
}

.comment-item:last-child {
  border-bottom: none;
}

.comment-content {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.comment-author-section {
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-author {
  font-weight: 600;
  cursor: pointer;
  color: var(--el-color-primary);
  font-size: 15px;
}

.comment-author:hover {
  text-decoration: underline;
}

.comment-meta-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.comment-time {
  font-size: 0.85rem;
  color: var(--text-secondary);
}

.comment-text {
  margin: 0;
  line-height: 1.8;
  font-size: 15px;
  color: #303133;
}

.load-more {
  text-align: center;
  margin-top: 20px;
}

/* ========== 响应式 ========== */
@media (max-width: 1200px) {
  .photography-layout {
    flex-direction: column;
  }
  
  .photo-section,
  .info-sidebar {
    flex: 1;
  }
}

@media (max-width: 768px) {
  .work-title-section {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .thumbnail-nav {
    gap: 5px;
  }
  
  .thumbnail-item {
    width: 60px;
    height: 60px;
  }
}
</style>
