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

      <el-empty v-else description="作品不存在" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  View, ChatDotRound, Star, Link, 
  Location, Calendar, Picture, ArrowDown 
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const work = ref(null)
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

const photos = computed(() => {
  if (!work.value || work.value.type !== 'photography') return []
  return work.value.images || []
})

const currentPhoto = computed(() => {
  return photos.value[currentPhotoIndex.value] || { metadata: {} }
})

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

const loadWork = async () => {
  try {
    const response = await api.get(`/works/${route.params.id}`)
    work.value = response.data
  } catch (error) {
    ElMessage.error('加载作品失败')
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
  if (!userStore.isLoggedIn) return false
  return userStore.user.id === comment.user_id || userStore.user.role === 'admin'
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
    ElMessage.success('评论发表成功')
    commentContent.value = ''
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
    // 重新加载作品以更新评论数
    await loadWork()
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
    await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: content,
      parent_id: comment.replyTo ? comment.replyTo.id : comment.id
    })
    ElMessage.success('回复成功')
    comment.showReply = false
    comment.replyContent = ''
    comment.replyTo = null
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
    // 重新加载作品以更新评论数
    await loadWork()
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
    await api.post('/comments', {
      work_id: parseInt(route.params.id),
      content: content,
      parent_id: reply.replyTo ? reply.replyTo.id : reply.id
    })
    ElMessage.success('回复成功')
    reply.showReply = false
    reply.replyContent = ''
    reply.replyTo = null
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
    // 重新加载作品以更新评论数
    await loadWork()
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
    // 重新加载作品以更新评论数
    await loadWork()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
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

  liking.value = true
  try {
    // 后端是 toggle 操作
    await api.post(`/works/${route.params.id}/like`)
    
    // toggle 状态
    isLiked.value = !isLiked.value
    
    // 重新加载作品数据以获取服务器端的最新数量
    await loadWork()
    // 重新检查点赞状态
    await checkLikedStatus()
    
    ElMessage.success(isLiked.value ? '点赞成功' : '取消点赞')
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
    // 重新检查状态和重新加载作品
    await checkLikedStatus()
    await loadWork()
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
    // 根据当前状态选择操作
    if (isFavorited.value) {
      await api.delete(`/works/${route.params.id}/favorite`)
      isFavorited.value = false
      ElMessage.success('取消收藏')
    } else {
      await api.post(`/works/${route.params.id}/favorite`)
      isFavorited.value = true
      ElMessage.success('收藏成功')
    }
    
    // 重新加载作品数据以获取服务器端的最新数量
    await loadWork()
    // 重新检查收藏状态
    await checkFavoritedStatus()
  } catch (error) {
    ElMessage.error(error.response?.data?.message || '操作失败')
    // 重新检查状态和重新加载作品
    await checkFavoritedStatus()
    await loadWork()
  } finally {
    favoriting.value = false
  }
}

onMounted(async () => {
  if (userStore.isLoggedIn && !userStore.user) {
    await userStore.fetchProfile()
  }
  
  await loadWork()
  await loadComments()
  checkLikedStatus()
  checkFavoritedStatus()
  
  // 响应式调整轮播高度
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

.comment-form .el-button {
  margin-top: 10px;
}

.comment-list {
  margin-top: 30px;
}

.comment-item {
  display: flex;
  gap: 15px;
  padding: 20px 0;
  border-bottom: 1px solid var(--border-lighter);
}

.comment-content {
  flex: 1;
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
  color: var(--text-primary);
  font-size: 15px;
}

.comment-meta-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.comment-time {
  color: var(--text-secondary);
  font-size: 14px;
}

.comment-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.reply-input {
  margin-top: 15px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 6px;
}

.reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 10px;
}

.replies-list {
  margin-top: 20px;
  padding-left: 20px;
  border-left: 2px solid #e4e7ed;
}

.reply-item {
  display: flex;
  gap: 10px;
  padding: 15px 0;
  border-bottom: 1px solid #f0f0f0;
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
  margin-bottom: 8px;
}

.reply-author {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.reply-to {
  color: var(--text-secondary);
  font-weight: normal;
  font-size: 12px;
  margin-left: 5px;
}

.author-tag {
  /* 已在 comment-author-section 中通过 gap 控制间距 */
}

.reply-time {
  color: var(--text-secondary);
  font-size: 12px;
}

.load-more-container {
  margin-top: 30px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
}

.load-more-btn {
  min-width: 200px;
}

.comment-count-info {
  font-size: 14px;
  color: var(--text-secondary);
}

.no-more-comments {
  margin-top: 30px;
  text-align: center;
  color: var(--text-secondary);
  font-size: 14px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.load-more-replies {
  margin-top: 10px;
  padding: 10px 0;
  text-align: center;
  border-top: 1px solid #f0f0f0;
}

.clickable-avatar {
  cursor: pointer;
  transition: all 0.3s;
}

.clickable-avatar:hover {
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
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
