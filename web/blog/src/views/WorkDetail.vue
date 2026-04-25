<template>
  <div class="work-detail">
    <div class="container">
      <el-card v-if="work" class="detail-card">
        <template #header>
          <div class="card-header">
            <!-- 审核状态提示（仅在非已发布状态时显示） -->
            <div v-if="work.status !== 1" class="work-status-banner">
              <el-alert
                :type="getStatusAlertType(work.status)"
                :closable="false"
                show-icon
              >
                <template #title>
                  <div>
                    <div>{{ getStatusText(work.status) }}</div>
                    <div v-if="work.audit_message" class="audit-message-text">
                      {{ work.audit_message }}
                    </div>
                  </div>
                </template>
              </el-alert>
            </div>
            <!-- 编辑按钮（仅作者可见） -->
            <div v-if="isWorkOwner" class="card-header-actions">
              <el-button 
                type="primary" 
                :icon="Edit" 
                @click="handleEdit"
                size="default"
              >
                编辑作品
              </el-button>
            </div>
          </div>
        </template>
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
                      :preview-src-list="photoUrlList"
                      :initial-index="index"
                      preview-teleported
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
                <div id="work-description-preview-photography" v-if="work.description"></div>
                <p v-else class="album-description-empty">暂无描述</p>
                
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
                <!-- 如果是作品作者，显示"我的主页"按钮 -->
                <el-button 
                  v-if="isWorkOwner"
                  type="primary" 
                  class="full-width-btn"
                  @click="goToMyProfile"
                >
                  我的主页
                </el-button>
                <!-- 如果不是作品作者，显示关注/已关注按钮 -->
                <el-button 
                  v-else-if="userStore.isLoggedIn"
                  :type="isFollowing ? 'default' : 'primary'" 
                  class="full-width-btn"
                  :loading="followLoading"
                  @click="handleFollow"
                >
                  <el-icon v-if="!followLoading"><Plus v-if="!isFollowing" /><Check v-else /></el-icon>
                  {{ isFollowing ? '已关注' : '关注' }}
                </el-button>
                <!-- 未登录用户不显示按钮 -->
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
                <!-- 只有已发布的作品才显示点赞和收藏按钮 -->
                <div class="action-buttons" v-if="work && work.is_published">
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
                  class="clickable-avatar"
                  @click="goToUserProfile(work.author.id)"
                />
                <div class="author-info">
                  <span class="author-name" @click="goToUserProfile(work.author.id)">
                    {{ work.author.nickname || work.author.username }}
                  </span>
                  <span class="publish-time">{{ formatDate(work.created_at) }}</span>
                </div>
              </div>
            </div>

            <div class="work-content">
              <div id="work-description-preview-project" v-if="work.description"></div>
              <p v-else class="work-description-empty">暂无描述</p>
              
              <div class="tech-stack" v-if="work.tech_stack">
                <h3>技术栈</h3>
                <p>{{ work.tech_stack }}</p>
              </div>

              <div class="work-links" v-if="work.github_url || work.demo_url || work.link">
                <el-button 
                  v-if="work.github_url" 
                  type="primary"
                  size="large"
                  class="work-link-btn github-btn"
                  @click="openLink(work.github_url)"
                >
                  <el-icon><Link /></el-icon>
                  GitHub
                </el-button>
                <el-button 
                  v-if="work.link" 
                  type="success"
                  size="large"
                  class="work-link-btn project-btn"
                  @click="openLink(work.link)"
                >
                  <el-icon><Link /></el-icon>
                  项目主页
                </el-button>
                <el-button 
                  v-if="work.demo_url" 
                  plain
                  size="large"
                  class="work-link-btn demo-btn"
                  @click="openLink(work.demo_url)"
                >
                  <el-icon><Link /></el-icon>
                  在线演示
                </el-button>
              </div>

              <div class="work-stats">
                <span><el-icon><View /></el-icon> {{ work.view_count }}</span>
                <span><el-icon><ChatDotRound /></el-icon> {{ work.comment_count }}</span>
              </div>

              <!-- 点赞和收藏按钮（只有已发布的作品才显示） -->
              <div class="work-actions" v-if="work && work.is_published">
                <el-button 
                  :type="isLiked ? 'primary' : 'default'"
                  @click="handleLike"
                  :loading="liking"
                  size="default"
                  class="action-btn"
                >
                  <el-icon><Star /></el-icon>
                  {{ work.like_count || 0 }} {{ isLiked ? '已点赞' : '点赞' }}
                </el-button>
                <el-button 
                  :type="isFavorited ? 'warning' : 'default'"
                  @click="handleFavorite"
                  :loading="favoriting"
                  size="default"
                  class="action-btn"
                >
                  <el-icon><Star /></el-icon>
                  {{ work.favorite_count || 0 }} {{ isFavorited ? '已收藏' : '收藏' }}
                </el-button>
              </div>
            </div>
          </div>
        </template>

        <!-- 编辑按钮（底部） -->
        <div v-if="isWorkOwner" class="edit-actions-bottom">
          <el-button 
            type="primary" 
            :icon="Edit" 
            @click="handleEdit"
            size="large"
          >
            <el-icon><Edit /></el-icon>
            编辑作品
          </el-button>
        </div>

        <!-- 评论区（两种类型共用） -->
        <!-- 只有已发布的作品（is_published=true）且评论功能开启时才显示评论区 -->
        <el-divider v-if="workCommentEnabled && work && work.is_published" />
        
        <div v-if="workCommentEnabled && work && work.is_published" class="comment-section">
          <h3>评论 ({{ work.comment_count }})</h3>
          
          <!-- 发表评论 -->
          <el-form v-if="userStore.isLoggedIn" @submit.prevent="submitComment" class="comment-form">
            <el-input
              v-model="commentContent"
              type="textarea"
              :rows="4"
              placeholder="写下你的评论..."
              maxlength="500"
              show-word-limit
            />
            <el-button type="primary" @click="submitComment" :loading="submittingComment">发表评论</el-button>
          </el-form>
          <el-alert v-else type="info" :closable="false">
            请<el-link type="primary" @click="$router.push('/login')">登录</el-link>后发表评论
          </el-alert>

          <!-- 评论列表 -->
          <div class="comment-list">
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <el-avatar 
                :src="comment.user?.avatar" 
                class="clickable-avatar"
                @click="handleUserClick(comment.user_id)"
              />
              <div class="comment-content">
                <div class="comment-header">
                  <div class="comment-author-section">
                    <span class="comment-author">
                      {{ comment.user?.nickname || comment.nickname }}
                    </span>
                    <el-tag v-if="isCommentAuthor(comment)" type="warning" size="small" effect="plain" class="author-tag">
                      作者
                    </el-tag>
                  </div>
                  <div class="comment-meta-section">
                    <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                  </div>
                </div>
                <p>{{ comment.content }}</p>
                <div class="comment-actions">
                  <el-button 
                    text 
                    size="small" 
                    :type="comment.isLiked ? 'primary' : 'default'"
                    @click="handleCommentLike(comment)"
                    :loading="comment.likeLoading"
                  >
                    <el-icon><Star /></el-icon> 
                    {{ comment.like_count || 0 }}
                  </el-button>
                  <el-button 
                    v-if="userStore.isLoggedIn"
                    text 
                    size="small" 
                    @click="showReplyInput(comment)"
                  >
                    <el-icon><ChatDotRound /></el-icon> 回复
                  </el-button>
                  <el-button 
                    v-if="canDeleteComment(comment)"
                    text 
                    type="danger" 
                    size="small" 
                    @click="handleDeleteComment(comment)"
                  >
                    删除
                  </el-button>
                </div>
                
                <!-- 回复输入框 -->
                <div v-if="comment.showReply" class="reply-input">
                  <el-input
                    v-model="comment.replyContent"
                    type="textarea"
                    :rows="3"
                    placeholder="写下你的回复..."
                    maxlength="500"
                    show-word-limit
                  />
                  <div class="reply-actions">
                    <el-button size="small" @click="cancelReply(comment)">取消</el-button>
                    <el-button type="primary" size="small" @click="submitReply(comment)" :loading="comment.replying">
                      发表回复
                    </el-button>
                  </div>
                </div>
                
                <!-- 回复列表 -->
                <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
                  <div v-for="reply in getDisplayedReplies(comment)" :key="reply.id" class="reply-item">
                    <el-avatar 
                      :src="reply.user?.avatar" 
                      size="small"
                      class="clickable-avatar"
                      @click="handleUserClick(reply.user_id)"
                    />
                    <div class="reply-content">
                      <div class="reply-header">
                        <span class="reply-author">
                          {{ reply.user?.nickname || reply.nickname }}
                          <el-tag v-if="isCommentAuthor(reply)" type="warning" size="small" effect="plain" class="author-tag">
                            作者
                          </el-tag>
                          <span v-if="reply.parent_id" class="reply-to">
                            <template v-if="reply.parent_id === comment.id">
                              回复 @{{ comment.user?.nickname || comment.nickname }}
                            </template>
                            <template v-else>
                              回复 @{{ getReplyTargetName(comment, reply) }}
                            </template>
                          </span>
                        </span>
                        <span class="reply-time">{{ formatDate(reply.created_at) }}</span>
                      </div>
                      <p>{{ reply.content }}</p>
                      <div class="reply-actions">
                        <el-button 
                          text 
                          size="small" 
                          :type="reply.isLiked ? 'primary' : 'default'"
                          @click="handleCommentLike(reply)"
                          :loading="reply.likeLoading"
                        >
                          <el-icon><Star /></el-icon> 
                          {{ reply.like_count || 0 }}
                        </el-button>
                        <el-button 
                          v-if="userStore.isLoggedIn"
                          text 
                          size="small" 
                          @click="showReplyToReplyInput(comment, reply)"
                        >
                          <el-icon><ChatDotRound /></el-icon> 回复
                        </el-button>
                        <el-button 
                          v-if="canDeleteComment(reply)"
                          text 
                          type="danger" 
                          size="small" 
                          @click="handleDeleteComment(reply)"
                        >
                          删除
                        </el-button>
                      </div>
                      
                      <!-- 回复的回复输入框 -->
                      <div v-if="reply.showReply" class="reply-input">
                        <el-input
                          v-model="reply.replyContent"
                          type="textarea"
                          :rows="3"
                          placeholder="写下你的回复..."
                          maxlength="500"
                          show-word-limit
                        />
                        <div class="reply-actions">
                          <el-button size="small" @click="cancelReplyToReply(reply)">取消</el-button>
                          <el-button type="primary" size="small" @click="submitReplyToReply(comment, reply)" :loading="reply.replying">
                            发表回复
                          </el-button>
                        </div>
                      </div>
                    </div>
                  </div>
                  
                  <!-- 展开更多子评论按钮 -->
                  <div v-if="hasMoreReplies(comment)" class="load-more-replies">
                    <el-button 
                      text 
                      size="small" 
                      @click="loadMoreReplies(comment)"
                      :loading="comment.loadingReplies"
                    >
                      <el-icon><ArrowDown /></el-icon>
                      {{ comment.loadingReplies ? '加载中...' : `展开更多回复 (${comment.reply_count - getDisplayedReplies(comment).length} 条)` }}
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="load-more-container" v-if="hasMoreComments">
            <el-button 
              @click="loadMoreComments"
              :loading="loadingMoreComments"
              class="load-more-btn"
            >
              <el-icon v-if="!loadingMoreComments"><ArrowDown /></el-icon>
              {{ loadingMoreComments ? '加载中...' : '加载更多评论' }}
            </el-button>
            <div class="comment-count-info">
              已显示 {{ comments.length }} / {{ commentsTotal }} 条评论
            </div>
          </div>
          <div v-else-if="comments.length > 0" class="no-more-comments">
            已显示全部 {{ comments.length }} 条评论
          </div>
        </div>
      </el-card>

      <div v-else-if="loading" class="loading-container" v-loading="loading" element-loading-text="加载中...">
        <div class="loading-placeholder"></div>
      </div>
      <el-empty v-else description="作品不存在" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  View, ChatDotRound, Star, Link, 
  Location, Calendar, Picture, ArrowDown, Edit, Plus, Check
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { loadCodeTheme, loadHighlightTheme, getMarkdownTheme } from '@/utils/codeTheme'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const work = ref(null)
const loading = ref(true) // 添加加载状态
const comments = ref([])
const commentsPage = ref(1)
const commentsTotal = ref(0)
const commentContent = ref('')
const submittingComment = ref(false)
const loadingMoreComments = ref(false)
const hasMoreComments = ref(false)

// 摄影相关
const carouselRef = ref()
const currentPhotoIndex = ref(0)
const carouselHeight = ref(600)

// 点赞和收藏状态
const isLiked = ref(false)
const isFavorited = ref(false)
const liking = ref(false)
const favoriting = ref(false)

// 关注状态
const isFollowing = ref(false)
const followLoading = ref(false)
const workCommentEnabled = ref(true) // 默认允许评论

const photos = computed(() => {
  if (!work.value || work.value.type !== 'photography') return []
  return work.value.images || []
})

// 所有图片的URL列表，用于图片预览
const photoUrlList = computed(() => {
  return photos.value.map(photo => photo.url)
})

const currentPhoto = computed(() => {
  return photos.value[currentPhotoIndex.value] || { metadata: {} }
})

// 判断当前用户是否是作品作者
const isWorkOwner = computed(() => {
  if (!userStore.isLoggedIn || !work.value || !userStore.user) {
    return false
  }
  
  // 检查作品作者ID是否等于当前用户ID
  // 优先使用 author_id，如果没有则使用 author.id
  const authorId = work.value.author_id || work.value.author?.id
  const userId = userStore.user.id
  
  if (!authorId || !userId) {
    return false
  }
  
  // 处理可能的类型不匹配（字符串 vs 数字）
  return Number(authorId) === Number(userId)
})

// 获取作品状态文本
const getStatusText = (status) => {
  switch (status) {
    case 0:
      return '草稿'
    case 1:
      return '已发布'
    case 2:
      return '审核中'
    case 3:
      return '审核不通过'
    default:
      return '未知状态'
  }
}

// 获取状态警告类型
const getStatusAlertType = (status) => {
  switch (status) {
    case 0:
      return 'info' // 草稿 - 信息提示
    case 2:
      return 'warning' // 待审核 - 警告提示
    case 3:
      return 'error' // 审核不通过 - 错误提示
    default:
      return 'info'
  }
}

// 判断评论/回复的作者是否是作品作者
const isCommentAuthor = (comment) => {
  if (!work.value || !comment) return false
  
  // 获取作品作者ID
  const workAuthorId = work.value.author_id || work.value.author?.id
  if (!workAuthorId) return false
  
  // 获取评论作者ID
  const commentAuthorId = comment.user_id || comment.user?.id
  if (!commentAuthorId) return false
  
  // 比较（处理类型不匹配）
  return Number(workAuthorId) === Number(commentAuthorId)
}

// 跳转到编辑页面
const handleEdit = () => {
  if (!work.value) return
  router.push(`/dashboard/works/${work.value.id}/edit`)
}

// 渲染 Markdown 描述
const renderDescription = async () => {
  if (!work.value || !work.value.description) {
    return
  }
  
  // 加载代码主题配置
  const codeThemeValue = await loadCodeTheme()
  // 加载 highlight.js 主题样式
  await loadHighlightTheme(codeThemeValue)
  // 加载 Markdown 主题配置
  const mdTheme = await getMarkdownTheme()
  
  // 确保样式表完全加载后再渲染
  await new Promise(resolve => setTimeout(resolve, 150))
  
  nextTick(() => {
    // 根据作品类型选择不同的预览元素
    const previewId = work.value.type === 'photography' 
      ? 'work-description-preview-photography' 
      : 'work-description-preview-project'
    
    const previewDiv = document.getElementById(previewId)
    if (!previewDiv) {
      console.error(`Description preview element not found: ${previewId}`)
      return
    }
    
    // 清空之前的内容（仅在刷新时可能需要）
    previewDiv.innerHTML = ''
    
    console.log('Rendering work description with Vditor.preview, code theme:', codeThemeValue, 'markdown theme:', mdTheme)
    
    // 使用 Vditor.preview 进行渲染
    Vditor.preview(previewDiv, work.value.description, {
      mode: mdTheme || 'light',
      markdown: {
        toc: true,
        mark: true,
        footnotes: true,
        autoSpace: true,
      },
      hljs: {
        lineNumber: true,
        style: codeThemeValue || 'github', // 使用配置的代码主题
        enable: true
      },
      speech: {
        enable: false
      },
      anchor: 1,
      after: () => {
        console.log('Work description rendered successfully!')
      }
    })
  })
}

const loadWork = async (skipView = false) => {
  // 只有在首次加载时才显示 loading（skipView=true 表示刷新，不显示 loading）
  if (!skipView) {
    loading.value = true
  }
  try {
    const url = skipView 
      ? `/works/${route.params.id}?skip_view=true`
      : `/works/${route.params.id}`
    const response = await api.get(url)
    work.value = response.data
    // 渲染描述
    renderDescription()
  } catch (error) {
    // 如果作品不存在，work.value 保持为 null，loading 设为 false 以显示"作品不存在"
    work.value = null
    ElMessage.error('加载作品失败')
  } finally {
    loading.value = false
  }
}

const loadComments = async (append = false) => {
  try {
    const response = await api.get('/comments', {
      params: {
        work_id: route.params.id,
        page: commentsPage.value,
        page_size: 10
      }
    })
    const commentsList = response.data.list || []
    console.debug('Loaded comments:', commentsList.map(c => ({ id: c.id, like_count: c.like_count })))
    
    // 为每个评论添加状态和初始化回复
    const existingComments = append ? comments.value : []
    const newComments = commentsList.map(comment => {
      const existingComment = existingComments.find(c => c.id === comment.id)
      return {
        ...comment,
        like_count: comment.like_count !== undefined ? comment.like_count : 0,
        isLiked: false,
        likeLoading: false,
        showReply: existingComment?.showReply || false,
        replyContent: existingComment?.replyContent || '',
        replying: existingComment?.replying || false,
        replyTo: existingComment?.replyTo || null,
        loadingReplies: existingComment?.loadingReplies || false,
        repliesPage: existingComment?.repliesPage || 1,
        replies: (comment.replies || []).map(reply => {
          const existingReply = existingComment?.replies?.find(r => r.id === reply.id)
          return {
            ...reply,
            like_count: reply.like_count !== undefined ? reply.like_count : 0,
            isLiked: false,
            likeLoading: false,
            showReply: existingReply?.showReply || false,
            replyContent: existingReply?.replyContent || '',
            replying: existingReply?.replying || false,
            replyTo: existingReply?.replyTo || null
          }
        })
      }
    })
    
    // 如果是追加模式，合并评论；否则替换
    if (append) {
      comments.value = [...comments.value, ...newComments]
    } else {
      comments.value = newComments
    }
    
    commentsTotal.value = response.data.total || 0
    
    // 判断是否还有更多评论
    hasMoreComments.value = comments.value.length < commentsTotal.value
    
    // 检查每个评论的点赞状态
    const checkPromises = []
    for (const comment of newComments) {
      checkPromises.push(checkCommentLiked(comment))
      if (comment.replies && comment.replies.length > 0) {
        for (const reply of comment.replies) {
          checkPromises.push(checkCommentLiked(reply))
        }
      }
    }
    await Promise.all(checkPromises).catch(err => {
      console.error('Failed to check comment like status:', err)
    })
    console.debug('Comments loaded and like status checked:', comments.value.map(c => ({ id: c.id, like_count: c.like_count, isLiked: c.isLiked })))
  } catch (error) {
    console.error('Failed to load comments:', error)
  }
}

// 加载更多评论
const loadMoreComments = async () => {
  if (loadingMoreComments.value || !hasMoreComments.value) return
  
  loadingMoreComments.value = true
  try {
    commentsPage.value++
    await loadComments(true) // 追加模式
  } finally {
    loadingMoreComments.value = false
  }
}

const canDeleteComment = (comment) => {
  if (!userStore.isLoggedIn || !userStore.user) return false
  
  // 管理员可以删除所有评论
  if (userStore.user.role === 'admin') return true
  
  // 作品作者可以删除自己作品下的所有评论
  if (isWorkOwner.value) return true
  
  // 获取评论的用户ID（支持多种格式）
  const commentUserId = comment.user_id || comment.user?.id || comment.userId
  
  // 如果评论没有用户ID（游客评论），只有管理员或作品作者可以删除
  if (!commentUserId || commentUserId === 0) {
    return false // 游客评论只能由管理员或作品作者删除（上面已检查）
  }
  
  // 检查是否是评论作者
  return Number(userStore.user.id) === Number(commentUserId)
}

const submitComment = async () => {
  if (!commentContent.value.trim()) {
    ElMessage.warning('请输入评论内容')
    return
  }

  submittingComment.value = true
  try {
    const response = await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: commentContent.value
    })
    // 如果返回的消息不是默认的 "success"，显示返回的消息
    // 注意：api 拦截器已经返回了 data 对象，所以直接使用 response.message
    if (response.message && response.message !== 'success') {
      ElMessage.success(response.message)
    }
    commentContent.value = ''
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
    // 重新加载作品以更新评论数（跳过浏览量增加）
    await loadWork(true)
  } catch (error) {
    ElMessage.error('评论发表失败')
  } finally {
    submittingComment.value = false
  }
}

// 显示回复输入框（回复评论）
const showReplyInput = (comment, replyTo = null) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  // 关闭其他评论的回复框
  comments.value.forEach(c => {
    if (c.id !== comment.id) {
      c.showReply = false
    }
  })
  
  comment.showReply = true
  comment.replyTo = replyTo || comment
  comment.replyContent = replyTo ? `@${replyTo.user?.nickname || replyTo.nickname} ` : ''
}

// 显示回复的回复输入框
const showReplyToReplyInput = (parentComment, replyTo) => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  // 关闭其他评论和回复的回复框
  comments.value.forEach(c => {
    c.showReply = false
    if (c.replies) {
      c.replies.forEach(r => {
        r.showReply = false
      })
    }
  })
  
  // 在回复上显示输入框
  replyTo.showReply = true
  replyTo.replyTo = replyTo
  replyTo.replyContent = `@${replyTo.user?.nickname || replyTo.nickname} `
}

// 取消回复
const cancelReply = (comment) => {
  comment.showReply = false
  comment.replyContent = ''
  comment.replyTo = null
}

// 提交回复
const submitReply = async (comment) => {
  if (!comment.replyContent || !comment.replyContent.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  // 去掉 @用户名称 前缀，只保留实际评论内容
  let content = comment.replyContent.trim()
  if (comment.replyTo) {
    const mentionPattern = new RegExp(`^@${comment.replyTo.user?.nickname || comment.replyTo.nickname}\\s*`, 'i')
    content = content.replace(mentionPattern, '')
  }

  if (!content) {
    ElMessage.warning('请输入回复内容')
    return
  }

  comment.replying = true
  try {
    const response = await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: content,
      parent_id: comment.replyTo ? comment.replyTo.id : comment.id
    })
    // 如果返回的消息不是默认的 "success"，显示返回的消息
    // 注意：api 拦截器已经返回了 data 对象，所以直接使用 response.message
    if (response.message && response.message !== 'success') {
      ElMessage.success(response.message)
    }
    comment.showReply = false
    comment.replyContent = ''
    comment.replyTo = null
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
    // 重新加载作品以更新评论数（跳过浏览量增加）
    await loadWork(true)
  } catch (error) {
    ElMessage.error('回复失败')
  } finally {
    comment.replying = false
  }
}

// 提交回复的回复
const submitReplyToReply = async (parentComment, reply) => {
  if (!reply.replyContent || !reply.replyContent.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }

  // 去掉 @用户名称 前缀，只保留实际评论内容
  let content = reply.replyContent.trim()
  if (reply.replyTo) {
    const mentionPattern = new RegExp(`^@${reply.replyTo.user?.nickname || reply.replyTo.nickname}\\s*`, 'i')
    content = content.replace(mentionPattern, '')
  }

  if (!content) {
    ElMessage.warning('请输入回复内容')
    return
  }

  reply.replying = true
  try {
    const response = await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: content,
      parent_id: reply.replyTo ? reply.replyTo.id : reply.id
    })
    // 如果返回的消息不是默认的 "success"，显示返回的消息
    // 注意：api 拦截器已经返回了 data 对象，所以直接使用 response.message
    if (response.message && response.message !== 'success') {
      ElMessage.success(response.message)
    }
    reply.showReply = false
    reply.replyContent = ''
    reply.replyTo = null
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
    // 重新加载作品以更新评论数（跳过浏览量增加）
    await loadWork(true)
  } catch (error) {
    ElMessage.error('回复失败')
  } finally {
    reply.replying = false
  }
}

// 取消回复的回复
const cancelReplyToReply = (reply) => {
  reply.showReply = false
  reply.replyContent = ''
  reply.replyTo = null
}

// 获取当前显示的子评论（所有子评论，包括回复的回复）
const getDisplayedReplies = (comment) => {
  if (!comment.replies) return []
  // 返回所有子评论（按时间排序，后端已经排序）
  return comment.replies
}

// 判断是否有更多子评论
const hasMoreReplies = (comment) => {
  if (!comment.reply_count) return false
  const displayedCount = getDisplayedReplies(comment).length
  return comment.reply_count > displayedCount
}

// 加载更多子评论
const loadMoreReplies = async (comment) => {
  if (!comment || comment.loadingReplies) return
  
  comment.loadingReplies = true
  try {
    const nextPage = (comment.repliesPage || 1) + 1
    const response = await api.get(`/comments/replies/${comment.id}`, {
      params: {
        page: nextPage,
        page_size: 10
      }
    })
    
    const newReplies = (response.data.list || []).map(reply => ({
      ...reply,
      like_count: reply.like_count !== undefined ? reply.like_count : 0,
      isLiked: false,
      likeLoading: false,
      showReply: false,
      replyContent: '',
      replying: false,
      replyTo: null
    }))
    
    // 合并新回复到现有回复列表
    comment.replies = [...(comment.replies || []), ...newReplies]
    comment.repliesPage = nextPage
    
    // 检查新加载的回复的点赞状态
    for (const reply of newReplies) {
      await checkCommentLiked(reply)
    }
  } catch (error) {
    console.error('Failed to load more replies:', error)
    ElMessage.error('加载更多回复失败')
  } finally {
    comment.loadingReplies = false
  }
}

// 获取回复目标名称
const getReplyTargetName = (parentComment, reply) => {
  if (!reply.parent_id || reply.parent_id === parentComment.id) {
    return parentComment.user?.nickname || parentComment.nickname
  }
  // 查找被回复的评论
  const targetReply = parentComment.replies?.find(r => r.id === reply.parent_id)
  return targetReply ? (targetReply.user?.nickname || targetReply.nickname) : ''
}

// 检查评论点赞状态
const checkCommentLiked = async (comment) => {
  if (!comment || !comment.id) return
  
  try {
    const response = await api.get(`/comments/${comment.id}/is-liked`)
    // API拦截器返回的是 {code: 0, message: "success", data: {...}}
    // 所以需要访问 response.data.is_liked
    const isLiked = response.data?.is_liked ?? false
    comment.isLiked = isLiked
    console.debug(`Comment ${comment.id} like status:`, isLiked)
  } catch (error) {
    // 静默失败，不显示错误（可能是游客访问或网络问题）
    comment.isLiked = false
    console.debug(`Failed to check comment ${comment.id} like status:`, error)
  }
}

// 处理评论点赞
const handleCommentLike = async (comment) => {
  if (!comment || !comment.id || comment.likeLoading) return
  
  comment.likeLoading = true
  const oldIsLiked = comment.isLiked
  const commentId = comment.id // 保存ID，因为重新加载后comment对象可能变化
  
  try {
    if (comment.isLiked) {
      // 取消点赞
      console.debug(`Unliking comment ${comment.id}, current count: ${comment.like_count}`)
      const response = await api.delete(`/comments/${comment.id}/like`)
      // 先更新本地状态（乐观更新）
      comment.isLiked = false
      // 重新加载评论以获取服务器端的最新点赞数，确保数据一致性
      commentsPage.value = 1
      await loadComments(false)
      // 重新检查点赞状态（因为重新加载后comment对象可能变化）
      const updatedComment = comments.value.find(c => c.id === commentId) || 
                            comments.value.flatMap(c => c.replies || []).find(r => r.id === commentId)
      if (updatedComment) {
        await checkCommentLiked(updatedComment)
      }
      // 显示成功消息
      const message = response?.message || '取消点赞成功'
      ElMessage.success(message)
      console.debug(`Comment ${comment.id} unliked, new count:`, updatedComment?.like_count)
    } else {
      // 点赞
      console.debug(`Liking comment ${comment.id}, current count: ${comment.like_count}`)
      const response = await api.post(`/comments/${comment.id}/like`)
      // 先更新本地状态（乐观更新）
      comment.isLiked = true
      // 重新加载评论以获取服务器端的最新点赞数，确保数据一致性
      commentsPage.value = 1
      await loadComments(false)
      // 重新检查点赞状态（因为重新加载后comment对象可能变化）
      const updatedComment = comments.value.find(c => c.id === commentId) || 
                            comments.value.flatMap(c => c.replies || []).find(r => r.id === commentId)
      if (updatedComment) {
        await checkCommentLiked(updatedComment)
      }
      // 显示成功消息
      const message = response?.message || '点赞成功'
      ElMessage.success(message)
      console.debug(`Comment ${comment.id} liked, new count:`, updatedComment?.like_count)
    }
  } catch (error) {
    // 恢复原状态
    comment.isLiked = oldIsLiked
    console.error(`Failed to ${oldIsLiked ? 'unlike' : 'like'} comment ${comment.id}:`, error)
    // API拦截器已经显示了错误消息
    // 重新检查状态和重新加载评论以获取最新数据
    await checkCommentLiked(comment)
    // 重新加载评论列表以获取最新的点赞数
    commentsPage.value = 1
    await loadComments(false)
  } finally {
    comment.likeLoading = false
  }
}

const handleDeleteComment = async (comment) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', { type: 'warning' })
    await api.delete(`/comments/${comment.id}`)
    ElMessage.success('删除成功')
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
    // 重新加载作品以更新评论数（跳过浏览量增加）
    await loadWork(true)
  } catch (error) {
    if (error !== 'cancel') {
      const errorMessage = error.response?.data?.message || error.message || '删除失败'
      ElMessage.error(errorMessage)
    }
  }
}

// 处理用户点击 - 跳转到用户主页
const handleUserClick = (userId) => {
  if (!userId) return
  router.push(`/users/${userId}`)
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
    router.push(`/users/${userId}`)
  }
}

// 跳转到自己的主页
const goToMyProfile = () => {
  if (userStore.user?.id) {
    router.push(`/users/${userStore.user.id}`)
  }
}

// 检查关注状态
const checkFollowStatus = async () => {
  if (!userStore.isLoggedIn || !work.value || !work.value.author) {
    isFollowing.value = false
    return
  }

  // 如果是自己的作品，不需要检查关注状态
  if (isWorkOwner.value) {
    isFollowing.value = false
    return
  }

  try {
    const authorId = work.value.author_id || work.value.author?.id
    if (!authorId) return

    const response = await api.get(`/users/${authorId}/follow-stats`)
    isFollowing.value = response.data?.is_following || false
  } catch (error) {
    console.error('Failed to check follow status:', error)
    isFollowing.value = false
  }
}

// 关注/取消关注
const handleFollow = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  if (!work.value || !work.value.author) {
    return
  }

  const authorId = work.value.author_id || work.value.author?.id
  if (!authorId) {
    return
  }

  // 不能关注自己
  if (Number(authorId) === Number(userStore.user?.id)) {
    ElMessage.warning('不能关注自己')
    return
  }

  followLoading.value = true
  try {
    if (isFollowing.value) {
      // 取消关注
      await api.delete(`/users/${authorId}/follow`)
      ElMessage.success('已取消关注')
      isFollowing.value = false
    } else {
      // 关注
      await api.post(`/users/${authorId}/follow`)
      ElMessage.success('关注成功')
      isFollowing.value = true
    }
  } catch (error) {
    const errorMsg = error.response?.data?.message || error.message || '操作失败'
    ElMessage.error(errorMsg)
    // 如果错误是"已经关注过该用户"或"未关注该用户"，刷新状态
    if (errorMsg.includes('已经关注过')) {
      isFollowing.value = true
    } else if (errorMsg.includes('未关注')) {
      isFollowing.value = false
    }
  } finally {
    followLoading.value = false
  }
}

// 打开外部链接
const openLink = (url) => {
  if (url) {
    window.open(url, '_blank')
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
  if (!userStore.isLoggedIn) {
    isLiked.value = false
    return
  }
  try {
    const response = await api.get(`/works/${route.params.id}/liked`)
    isLiked.value = response.data.liked || response.data.is_liked || false
  } catch (error) {
    console.error('Failed to check liked status:', error)
    isLiked.value = false
  }
}

const checkFavoritedStatus = async () => {
  if (!userStore.isLoggedIn) {
    isFavorited.value = false
    return
  }
  try {
    const response = await api.get(`/works/${route.params.id}/favorited`)
    isFavorited.value = response.data.favorited || response.data.is_favorited || false
  } catch (error) {
    console.error('Failed to check favorited status:', error)
    isFavorited.value = false
  }
}

const handleLike = async () => {
  if (!userStore.isLoggedIn) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }

  // 防止重复点击
  if (liking.value) {
    return
  }

  liking.value = true
  try {
    // 后端是 toggle 操作
    await api.post(`/works/${route.params.id}/like`)
    
    // 重新加载作品数据以获取服务器端的最新数量（跳过浏览量增加）
    await loadWork(true)
    // 重新检查点赞状态（从服务器获取最新状态）
    await checkLikedStatus()
    
    // 根据最新状态显示消息
    ElMessage.success(isLiked.value ? '点赞成功' : '取消点赞')
  } catch (error) {
    // 错误消息已在 api.js 拦截器中统一显示，这里不再重复显示
    // 重新检查状态和重新加载作品（跳过浏览量增加）
    await checkLikedStatus()
    await loadWork(true)
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

  // 防止重复点击
  if (favoriting.value) {
    return
  }

  favoriting.value = true
  try {
    // 根据当前状态选择操作
    if (isFavorited.value) {
      await api.delete(`/works/${route.params.id}/favorite`)
      ElMessage.success('取消收藏')
    } else {
      await api.post(`/works/${route.params.id}/favorite`)
      ElMessage.success('收藏成功')
    }
    
    // 重新加载作品数据以获取服务器端的最新数量（跳过浏览量增加）
    await loadWork(true)
    // 重新检查收藏状态（从服务器获取最新状态）
    await checkFavoritedStatus()
  } catch (error) {
    // 错误消息已在 api.js 拦截器中统一显示，这里不再重复显示
    // 重新检查状态和重新加载作品（跳过浏览量增加）
    await checkFavoritedStatus()
    await loadWork(true)
  } finally {
    favoriting.value = false
  }
}

// 监听 work 变化，重新渲染描述
watch(() => work.value?.description, () => {
  renderDescription()
})

// 监听 work 变化，检查关注状态
watch(() => work.value?.author?.id, () => {
  if (work.value) {
    checkFollowStatus()
  }
}, { immediate: false })

// 加载评论配置
const loadCommentSettings = async () => {
  try {
    const response = await api.get('/settings/public')
    const settings = response.data || {}
    // 检查作品评论是否开放（默认true，如果配置不存在或为'1'/'true'则允许）
    workCommentEnabled.value = settings.work_comment_enabled !== '0' && settings.work_comment_enabled !== 'false'
  } catch (error) {
    console.error('Failed to load comment settings:', error)
    // 默认允许评论（向后兼容）
    workCommentEnabled.value = true
  }
}

onMounted(async () => {
  if (userStore.isLoggedIn && !userStore.user) {
    await userStore.fetchProfile()
  }
  
  // 检查 sessionStorage 中是否有预加载的作品数据
  const workId = route.params.id
  const preloadedWorkStr = sessionStorage.getItem(`preloaded_work_${workId}`)
  
  if (preloadedWorkStr) {
    try {
      // 如果有预加载的数据，直接使用，不显示加载动画
      const preloadedWork = JSON.parse(preloadedWorkStr)
      work.value = preloadedWork
      loading.value = false
      // 清除 sessionStorage 中的数据（只使用一次）
      sessionStorage.removeItem(`preloaded_work_${workId}`)
      // 渲染描述
      renderDescription()
      // 加载评论配置
      await loadCommentSettings()
      // 只有在评论功能开启时才加载评论
      if (workCommentEnabled.value) {
        await loadComments()
      }
      checkLikedStatus()
      checkFavoritedStatus()
      checkFollowStatus()
    } catch (error) {
      // 如果解析失败，正常加载
      console.error('Failed to parse preloaded work data:', error)
      await loadCommentSettings()
      await loadWork()
      if (workCommentEnabled.value) {
        await loadComments()
      }
      checkLikedStatus()
      checkFavoritedStatus()
      checkFollowStatus()
    }
  } else {
    // 没有预加载数据，正常加载
    // 加载评论配置
    await loadCommentSettings()
    
    await loadWork()
    // 只有在评论功能开启时才加载评论
    if (workCommentEnabled.value) {
      await loadComments()
    }
    checkLikedStatus()
    checkFavoritedStatus()
    checkFollowStatus()
  }
  
  // 响应式调整轮播高度
  updateCarouselHeight()
  window.addEventListener('resize', updateCarouselHeight)
})

onUnmounted(() => {
  window.removeEventListener('resize', updateCarouselHeight)
  })
</script>

<style>
/* Vditor渲染样式需要全局作用域 */
@import 'vditor/dist/index.css';
</style>

<style scoped>
.work-detail {
  padding: var(--spacing-xl) 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.detail-card {
  max-width: 1400px;
  margin: 0 auto;
  box-shadow: var(--shadow-md);
  background-color: var(--theme-bg-card);
  border: 1px solid var(--theme-border-light);
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.work-status-banner {
  margin-bottom: var(--spacing-sm);
}

.work-status-banner :deep(.el-alert) {
  border-radius: var(--radius-md);
}

.audit-message-text {
  margin-top: var(--spacing-sm);
  font-size: var(--font-size-sm);
  font-weight: normal;
}

.loading-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: var(--spacing-md);
  min-height: 600px;
  position: relative;
}

.loading-placeholder {
  min-height: 600px;
}

.card-header-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.edit-actions-bottom {
  display: flex;
  justify-content: center;
  padding: var(--spacing-md) 0;
  margin-top: var(--spacing-md);
  border-top: 1px solid var(--theme-border-light);
}

.photography-layout {
  display: flex;
  gap: var(--spacing-lg);
}

.photo-section {
  flex: 7;
}

.main-photo-container {
  position: relative;
  background: var(--theme-bg-primary);
  border-radius: var(--radius-lg);
  overflow: hidden;
  margin-bottom: var(--spacing-md);
  box-shadow: var(--shadow-lg);
}

.main-photo {
  width: 100%;
  height: 100%;
}

.photo-counter {
  position: absolute;
  bottom: var(--spacing-lg);
  right: var(--spacing-lg);
  background: rgba(0, 0, 0, 0.75);
  color: var(--color-text-inverse);
  padding: 10px 20px;
  border-radius: var(--radius-full);
  font-size: var(--font-size-base);
  font-weight: 500;
  backdrop-filter: blur(10px);
  box-shadow: var(--shadow-md);
}

.thumbnail-nav {
  display: flex;
  gap: var(--spacing-sm);
  overflow-x: auto;
  padding: var(--spacing-sm) 0;
}

.thumbnail-item {
  flex-shrink: 0;
  width: 80px;
  height: 80px;
  border-radius: var(--radius-md);
  overflow: hidden;
  cursor: pointer;
  border: 3px solid transparent;
  transition: all var(--transition-slow);
  box-shadow: var(--shadow-sm);
}

.thumbnail-item:hover {
  border-color: var(--theme-primary);
  transform: translateY(-2px);
  box-shadow: 0 4px 12px color-mix(in srgb, var(--theme-primary) 30%, transparent);
}

.thumbnail-item.active {
  border-color: var(--theme-primary);
  box-shadow: 0 0 0 3px color-mix(in srgb, var(--theme-primary) 20%, transparent);
  transform: scale(1.05);
}

.thumbnail-item .el-image {
  width: 100%;
  height: 100%;
}

.photo-description {
  padding: var(--spacing-md);
  background: var(--theme-bg-secondary);
  border-radius: var(--radius-lg);
  margin-bottom: var(--spacing-lg);
  border-left: 4px solid var(--theme-primary);
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
  color: var(--theme-text-secondary);
}

.album-info {
  padding: var(--spacing-md) 0;
}

.album-info h2 {
  margin-bottom: var(--spacing-md);
  color: var(--theme-text-primary);
}

.album-description-empty {
  color: var(--theme-text-secondary);
  margin-bottom: var(--spacing-md);
  font-style: italic;
}

.album-meta {
  display: flex;
  gap: var(--spacing-md);
  flex-wrap: wrap;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  color: var(--theme-text-secondary);
}

.info-sidebar {
  flex: 3;
  display: flex;
  flex-direction: column;
  gap: var(--spacing-md);
}

.author-card,
.photo-params-card,
.stats-card {
  padding: var(--spacing-lg);
  background: var(--theme-bg-card);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--theme-border-light);
}

.author-header {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
  cursor: pointer;
  transition: opacity var(--transition-base);
}

.author-header:hover {
  opacity: 0.8;
}

.author-info h3 {
  margin: 0 0 5px 0;
  font-size: var(--font-size-lg);
  color: var(--theme-text-primary);
}

.author-info p {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.full-width-btn {
  width: 100%;
}

.photo-params-card h3,
.stats-card h3 {
  margin: 0 0 var(--spacing-md) 0;
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--theme-text-primary);
  padding-bottom: var(--spacing-sm);
  border-bottom: 2px solid var(--theme-border-light);
}

.params-list,
.stats-list {
  display: flex;
  flex-direction: column;
  gap: var(--spacing-sm);
}

.param-item {
  display: flex;
  justify-content: space-between;
  font-size: var(--font-size-sm);
  padding: var(--spacing-sm) 0;
  border-bottom: 1px solid var(--theme-border-light);
}

.param-item:last-child {
  border-bottom: none;
}

.param-label {
  color: var(--theme-text-tertiary);
  font-weight: 500;
}

.param-value {
  font-weight: 600;
  color: var(--theme-text-primary);
}

.stat-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.action-buttons {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-md);
}

.action-buttons .el-button {
  flex: 1;
  height: 44px;
  font-size: var(--font-size-base);
  border-radius: var(--radius-md);
  font-weight: 500;
  transition: all var(--transition-slow);
  cursor: pointer;
}

.action-buttons .el-button:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.project-layout {
  max-width: 900px;
  margin: 0 auto;
}

.work-header {
  margin-bottom: var(--spacing-lg);
}

.work-title-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.work-title-section h1 {
  margin: 0;
  font-size: var(--font-size-3xl);
  color: var(--theme-text-primary);
}

.work-author {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  font-weight: 500;
  cursor: pointer;
  color: var(--theme-primary);
  transition: color var(--transition-fast);
}

.author-name:hover {
  color: var(--theme-primary-hover);
  text-decoration: underline;
}

.publish-time {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.work-content {
  line-height: var(--line-height-relaxed);
  padding: var(--spacing-lg);
  background-color: var(--theme-content-bg);
  border-radius: var(--radius-md);
  border: 1px solid var(--theme-border-light);
  box-shadow: var(--shadow-sm);
  margin: var(--spacing-md) 0;
}

.work-content :deep(pre) {
  line-height: var(--line-height-tight);
}

.work-description-empty {
  font-size: var(--font-size-lg);
  margin-bottom: var(--spacing-lg);
  color: var(--theme-text-secondary);
  font-style: italic;
}

#work-description-preview-photography {
  padding: var(--spacing-lg);
  background-color: var(--theme-content-bg);
  border-radius: var(--radius-md);
  border: 1px solid var(--theme-border-light);
  box-shadow: var(--shadow-sm);
  margin: var(--spacing-md) 0;
}

#work-description-preview-project {
  padding: 0;
  margin: 0;
}

#work-description-preview-photography :deep(code:not(pre code)),
#work-description-preview-project :deep(code:not(pre code)) {
  background-color: var(--theme-bg-hover);
  color: var(--theme-text-primary);
  padding: 0.2em 0.4em;
  border-radius: var(--radius-sm);
  font-size: 85%;
}

#work-description-preview-photography :deep(pre),
#work-description-preview-project :deep(pre) {
  line-height: var(--line-height-tight);
}

#work-description-preview-photography :deep(pre code),
#work-description-preview-project :deep(pre code) {
  line-height: inherit;
}

.work-links {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
  flex-wrap: wrap;
}

.work-link-btn {
  min-width: 140px;
  height: 48px;
  font-size: var(--font-size-lg);
  font-weight: 600;
  border-radius: var(--radius-md);
  transition: all var(--transition-slow);
  box-shadow: var(--shadow-sm);
  cursor: pointer;
}

.work-link-btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.github-btn {
  background: linear-gradient(135deg, var(--theme-text-primary) 0%, var(--theme-text-secondary) 100%);
  border: none;
  color: var(--color-text-inverse);
}

.github-btn:hover {
  background: linear-gradient(135deg, var(--theme-text-secondary) 0%, var(--theme-text-primary) 100%);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-md);
}

.project-btn {
  background: linear-gradient(135deg, var(--theme-accent) 0%, var(--theme-primary) 100%);
  border: none;
  color: var(--color-text-inverse);
}

.project-btn:hover {
  background: linear-gradient(135deg, var(--theme-primary) 0%, var(--theme-accent) 100%);
  color: var(--color-text-inverse);
  box-shadow: var(--shadow-md);
}

.demo-btn {
  border: 2px solid var(--theme-primary);
  color: var(--theme-primary);
  background: var(--theme-bg-card);
}

.demo-btn:hover {
  background: var(--theme-primary);
  color: var(--color-text-inverse);
  border-color: var(--theme-primary);
  box-shadow: var(--shadow-md);
}

.work-link-btn .el-icon {
  margin-right: 6px;
  font-size: var(--font-size-lg);
}

.tech-stack {
  margin-bottom: var(--spacing-lg);
}

.tech-stack h3 {
  margin-bottom: var(--spacing-sm);
  color: var(--theme-text-primary);
}

.work-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-md);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--theme-border-light);
  justify-content: center;
}

.action-btn {
  min-width: 120px;
  height: 40px;
  font-size: var(--font-size-sm);
  font-weight: 500;
  border-radius: var(--radius-md);
  transition: all var(--transition-slow);
  cursor: pointer;
}

.action-btn:hover {
  transform: translateY(-1px);
  box-shadow: var(--shadow-sm);
}

.action-btn .el-icon {
  margin-right: 4px;
  font-size: var(--font-size-base);
}

.work-stats {
  display: flex;
  gap: var(--spacing-md);
  color: var(--theme-text-secondary);
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--theme-border-light);
}

.work-stats span {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
}

.comment-section {
  margin-top: var(--spacing-xl);
}

.comment-section h3 {
  margin-bottom: var(--spacing-md);
  color: var(--theme-text-primary);
}

.comment-form {
  margin-bottom: var(--spacing-lg);
}

.comment-form .el-button {
  margin-top: var(--spacing-sm);
  cursor: pointer;
}

.comment-list {
  margin-top: var(--spacing-lg);
}

.comment-item {
  display: flex;
  gap: var(--spacing-md);
  padding: var(--spacing-md) 0;
  border-bottom: 1px solid var(--theme-border-light);
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-sm);
}

.comment-author-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.comment-author {
  font-weight: 600;
  color: var(--theme-text-primary);
  font-size: var(--font-size-base);
}

.comment-meta-section {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.comment-time {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
}

.comment-actions {
  display: flex;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.comment-actions .el-button {
  cursor: pointer;
}

.reply-input {
  margin-top: var(--spacing-md);
  padding: var(--spacing-md);
  background: var(--theme-bg-secondary);
  border-radius: var(--radius-md);
}

.reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-sm);
  margin-top: var(--spacing-sm);
}

.replies-list {
  margin-top: var(--spacing-md);
  padding-left: var(--spacing-md);
  border-left: 2px solid var(--theme-border-light);
}

.reply-item {
  display: flex;
  gap: var(--spacing-sm);
  padding: var(--spacing-md) 0;
  border-bottom: 1px solid var(--theme-border-light);
}

.reply-item:last-child {
  border-bottom: none;
}

.reply-content {
  flex: 1;
}

.reply-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: var(--spacing-sm);
}

.reply-author {
  font-weight: 600;
  color: var(--theme-text-primary);
  font-size: var(--font-size-sm);
}

.reply-to {
  color: var(--theme-text-secondary);
  font-weight: normal;
  font-size: var(--font-size-xs);
  margin-left: var(--spacing-xs);
}

.reply-time {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-xs);
}

.load-more-container {
  margin-top: var(--spacing-lg);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--spacing-sm);
}

.load-more-btn {
  min-width: 200px;
  cursor: pointer;
}

.comment-count-info {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.no-more-comments {
  margin-top: var(--spacing-lg);
  text-align: center;
  color: var(--theme-text-secondary);
  font-size: var(--font-size-sm);
  padding: var(--spacing-md);
  background-color: var(--theme-bg-secondary);
  border-radius: var(--radius-sm);
}

.load-more-replies {
  margin-top: var(--spacing-sm);
  padding: var(--spacing-sm) 0;
  text-align: center;
  border-top: 1px solid var(--theme-border-light);
}

.clickable-avatar {
  cursor: pointer;
  transition: all var(--transition-slow);
}

.clickable-avatar:hover {
  transform: scale(1.1);
  box-shadow: var(--shadow-sm);
}

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
    gap: var(--spacing-xs);
  }

  .thumbnail-item {
    width: 60px;
    height: 60px;
  }
}
</style>
