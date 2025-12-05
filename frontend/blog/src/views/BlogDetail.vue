<template>
  <div class="blog-detail">
    <div class="container">
      <el-card v-if="article" class="article-card">
        <div class="article-header">
          <div class="article-title-row">
            <h1>{{ article.title }}</h1>
            <el-button 
              v-if="isArticleOwner"
              type="primary"
              size="default"
              @click="handleEdit"
              class="edit-btn-header"
            >
              <el-icon><Edit /></el-icon> 
              编辑
            </el-button>
          </div>
          <div class="article-meta">
            <el-tag 
              v-if="article.category" 
              class="clickable-tag"
              @click="handleCategoryClick(article.category.id)"
            >
              {{ article.category.name }}
            </el-tag>
            <span 
              class="author-info clickable-author"
              @click="handleAuthorClick"
            >
              <el-avatar 
                :size="24" 
                :src="article.author?.avatar" 
                class="author-avatar"
              />
              {{ article.author?.nickname || article.author?.username }}
            </span>
            <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
            <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
          </div>
          <div class="article-tags">
            <el-tag 
              v-for="tag in article.tags" 
              :key="tag.id" 
              size="small" 
              type="info"
              class="clickable-tag"
              @click="handleTagClick(tag.id)"
            >
              {{ tag.name }}
            </el-tag>
          </div>
        </div>

        <div class="article-content" id="article-preview" v-if="article"></div>
        <div v-else class="article-loading">正在加载内容...</div>

        <div class="article-actions">
          <el-button 
            :type="isLiked ? 'primary' : 'default'"
            @click="handleLike"
            :loading="likeLoading"
          >
            <el-icon><Star /></el-icon> 
            {{ isLiked ? '已点赞' : '点赞' }}
            ({{ article.like_count || 0 }})
          </el-button>
          
          <el-button 
            v-if="userStore.isLoggedIn"
            :type="isFavorited ? 'warning' : 'default'"
            @click="handleFavorite"
            :loading="favoriteLoading"
          >
            <el-icon><Collection /></el-icon> 
            {{ isFavorited ? '已收藏' : '收藏' }} 
            ({{ article.favorite_count || 0 }})
          </el-button>

          <el-button 
            v-if="isArticleOwner"
            type="primary"
            @click="handleEdit"
          >
            <el-icon><Edit /></el-icon> 
            编辑文章
          </el-button>
        </div>
      </el-card>

      <el-card class="comments-card">
        <h3>评论区</h3>
        
        <el-form v-if="userStore.isLoggedIn" @submit.prevent="submitComment" class="comment-form">
          <el-input
            v-model="commentContent"
            type="textarea"
            :rows="4"
            placeholder="写下你的评论..."
            maxlength="500"
            show-word-limit
          />
          <el-button type="primary" @click="submitComment" :loading="submitting">发表评论</el-button>
        </el-form>
        <el-alert v-else type="info" :closable="false">
          请<el-link type="primary" @click="$router.push('/login')">登录</el-link>后发表评论
        </el-alert>

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
                <div v-for="reply in getDirectReplies(comment)" :key="reply.id" class="reply-item">
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
                        <span v-if="reply.parent_id && reply.parent_id !== comment.id" class="reply-to">
                          回复 @{{ getReplyTargetName(comment, reply) }}
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
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Star, Collection, User, View, Clock, Edit, ChatDotRound, ArrowDown } from '@element-plus/icons-vue'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const article = ref(null)
const comments = ref([])
const commentsPage = ref(1)
const commentsTotal = ref(0)
const commentContent = ref('')
const submitting = ref(false)
const loadingMoreComments = ref(false)
const hasMoreComments = ref(false)
const isLiked = ref(false)
const likeLoading = ref(false)
const isFavorited = ref(false)
const favoriteLoading = ref(false)
const renderedContent = ref('')

// 渲染 Markdown 内容（支持 Mermaid）
const renderMarkdown = () => {
  if (!article.value || !article.value.content) {
    return
  }
  
  nextTick(() => {
    const previewDiv = document.getElementById('article-preview')
    if (!previewDiv) {
      console.error('Preview element not found')
      return
    }
    
    console.log('Rendering with Vditor.preview...')
    
    // 使用 Vditor.preview 进行渲染
    Vditor.preview(previewDiv, article.value.content, {
      mode: 'light',
      markdown: {
        toc: true,
        mark: true,
        footnotes: true,
        autoSpace: true,
      },
      hljs: {
        lineNumber: false,
        style: 'github'
      },
      speech: {
        enable: false
      },
      anchor: 1,
      after: () => {
        console.log('Vditor preview rendered successfully!')
      }
    })
  })
}

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

// 判断当前用户是否是文章作者
const isArticleOwner = computed(() => {
  if (!userStore.isLoggedIn || !article.value || !userStore.user) {
    return false
  }
  
  // 检查文章作者ID是否等于当前用户ID
  // 优先使用 author_id，如果没有则使用 author.id
  const authorId = article.value.author_id || article.value.author?.id
  const userId = userStore.user.id
  
  if (!authorId || !userId) {
    return false
  }
  
  // 处理可能的类型不匹配（字符串 vs 数字）
  return Number(authorId) === Number(userId)
})

// 判断评论/回复的作者是否是文章作者
const isCommentAuthor = (comment) => {
  if (!article.value || !comment) return false
  
  // 获取文章作者ID
  const articleAuthorId = article.value.author_id || article.value.author?.id
  if (!articleAuthorId) return false
  
  // 获取评论作者ID
  const commentAuthorId = comment.user_id || comment.user?.id
  if (!commentAuthorId) return false
  
  // 比较（处理类型不匹配）
  return Number(articleAuthorId) === Number(commentAuthorId)
}

const loadArticle = async () => {
  try {
    const response = await api.get(`/articles/${route.params.id}`)
    article.value = response.data
    
    // 设置已加载标记
    renderedContent.value = 'loading'
    
    // 渲染 Markdown 内容
    renderMarkdown()
  } catch (error) {
    ElMessage.error('文章加载失败')
  }
}

const loadComments = async (append = false) => {
  try {
    const response = await api.get('/comments', {
      params: {
        article_id: route.params.id,
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

  submitting.value = true
  try {
    await api.post('/comments', {
      article_id: parseInt(route.params.id),
      content: commentContent.value
    })
    ElMessage.success('评论发表成功')
    commentContent.value = ''
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
  } catch (error) {
    ElMessage.error('评论发表失败')
  } finally {
    submitting.value = false
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

  comment.replying = true
  try {
    await api.post('/comments', {
      article_id: parseInt(route.params.id),
      content: comment.replyContent.trim(),
      parent_id: comment.replyTo ? comment.replyTo.id : comment.id
    })
    ElMessage.success('回复成功')
    comment.showReply = false
    comment.replyContent = ''
    comment.replyTo = null
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
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

  reply.replying = true
  try {
    await api.post('/comments', {
      article_id: parseInt(route.params.id),
      content: reply.replyContent.trim(),
      parent_id: reply.replyTo ? reply.replyTo.id : reply.id
    })
    ElMessage.success('回复成功')
    reply.showReply = false
    reply.replyContent = ''
    reply.replyTo = null
    // 重置为第一页并重新加载
    commentsPage.value = 1
    await loadComments(false)
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

// 获取直接回复（parent_id = comment.id 的回复）
const getDirectReplies = (comment) => {
  if (!comment.replies) return []
  return comment.replies.filter(reply => reply.parent_id === comment.id)
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
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}


const checkLiked = async () => {
  try {
    const response = await api.get(`/articles/${route.params.id}/is-liked`)
    isLiked.value = response.data.is_liked
  } catch (error) {
    console.error('Failed to check like status:', error)
  }
}

const handleLike = async () => {
  if (!article.value) return
  
  likeLoading.value = true
  try {
    if (isLiked.value) {
      await api.delete(`/articles/${route.params.id}/like`)
      ElMessage.success('取消点赞成功')
      isLiked.value = false
      // 确保数字正确更新
      article.value.like_count = Math.max(0, (article.value.like_count || 0) - 1)
    } else {
      await api.post(`/articles/${route.params.id}/like`)
      ElMessage.success('点赞成功')
      isLiked.value = true
      // 确保数字正确更新
      article.value.like_count = (article.value.like_count || 0) + 1
    }
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
    // 如果操作失败，重新检查状态和重新加载文章
    checkLiked()
    loadArticle()
  } finally {
    likeLoading.value = false
  }
}

const checkFavorited = async () => {
  if (!userStore.isLoggedIn) return
  
  try {
    const response = await api.get(`/articles/${route.params.id}/is-favorited`)
    isFavorited.value = response.data.is_favorited
  } catch (error) {
    console.error('Failed to check favorite status:', error)
  }
}

const handleFavorite = async () => {
  if (!article.value || !userStore.isLoggedIn) return
  
  favoriteLoading.value = true
  try {
    if (isFavorited.value) {
      await api.delete(`/articles/${route.params.id}/favorite`)
      ElMessage.success('取消收藏成功')
      isFavorited.value = false
      // 确保数字正确更新
      article.value.favorite_count = Math.max(0, (article.value.favorite_count || 0) - 1)
    } else {
      await api.post(`/articles/${route.params.id}/favorite`)
      ElMessage.success('收藏成功')
      isFavorited.value = true
      // 确保数字正确更新
      article.value.favorite_count = (article.value.favorite_count || 0) + 1
    }
  } catch (error) {
    ElMessage.error('操作失败')
    // 如果操作失败，重新检查状态和重新加载文章
    checkFavorited()
    loadArticle()
  } finally {
    favoriteLoading.value = false
  }
}

// 跳转到文章编辑页面
const handleEdit = () => {
  if (!article.value) return
  router.push(`/dashboard/articles/${article.value.id}/edit`)
}

// 处理分类点击 - 跳转到博客列表页并筛选该分类
const handleCategoryClick = (categoryId) => {
  if (!categoryId) return
  router.push({
    path: '/blog',
    query: { category_id: categoryId }
  })
}

// 处理标签点击 - 跳转到博客列表页并筛选该标签
const handleTagClick = (tagId) => {
  if (!tagId) return
  router.push({
    path: '/blog',
    query: { tag_id: tagId }
  })
}

// 处理用户点击 - 跳转到用户主页
const handleUserClick = (userId) => {
  if (!userId) return
  router.push(`/users/${userId}`)
}

// 处理文章作者点击 - 跳转到作者主页
const handleAuthorClick = () => {
  if (!article.value || !article.value.author) return
  const authorId = article.value.author_id || article.value.author.id
  if (!authorId) return
  router.push(`/users/${authorId}`)
}

onMounted(async () => {
  // 如果已登录但用户信息未加载，先获取用户信息
  if (userStore.isLoggedIn && !userStore.user) {
    await userStore.fetchProfile()
  }
  
  loadArticle()
  loadComments()
  checkLiked()
  checkFavorited()
})
</script>

<style>
/* Vditor渲染样式需要全局作用域 */
@import 'vditor/dist/index.css';
</style>

<style scoped>
.blog-detail {
  padding: 40px 0;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.article-card {
  margin-bottom: 30px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.comments-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.article-header {
  border-bottom: 1px solid var(--border-lighter);
  padding-bottom: 20px;
  margin-bottom: 30px;
}

.article-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20px;
  margin-bottom: 15px;
}

.article-title-row h1 {
  flex: 1;
  margin: 0;
}

.edit-btn-header {
  flex-shrink: 0;
}

.article-meta {
  display: flex;
  gap: 15px;
  color: var(--text-secondary);
  margin-bottom: 10px;
  flex-wrap: wrap;
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: 5px;
}

.author-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.author-avatar {
  flex-shrink: 0;
}

.clickable-author {
  cursor: pointer;
  transition: all 0.3s;
  padding: 4px 8px;
  border-radius: 4px;
}

.clickable-author:hover {
  background-color: #f5f7fa;
  color: #409eff;
}

.clickable-tag {
  cursor: pointer;
  transition: all 0.3s;
}

.clickable-tag:hover {
  transform: translateY(-2px);
  opacity: 0.8;
}

.clickable-avatar {
  cursor: pointer;
  transition: all 0.3s;
}

.clickable-avatar:hover {
  transform: scale(1.1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
}

.article-tags {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

/* 文章内容渲染样式 */
.article-content {
  margin-bottom: 30px;
  min-height: 200px;
}

#article-preview {
  padding: 20px 0;
}

/* 加载状态 */
.article-loading {
  padding: 40px 0;
  text-align: center;
  color: #999;
  font-size: 14px;
}

/* vditor-reset样式由全局CSS提供 */

.article-actions {
  text-align: center;
  padding-top: 20px;
  border-top: 1px solid var(--border-lighter);
}

.comments-card h3 {
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

.reply-actions {
  display: flex;
  gap: 10px;
  margin-top: 8px;
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
</style>

