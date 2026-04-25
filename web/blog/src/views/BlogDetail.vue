<template>
  <div class="blog-detail">
    <div class="container">
      <Card v-if="article" class="article-card">
        <CardContent class="p-6">
          <div class="article-header">
            <div class="article-title-row">
              <h1>{{ article.title }}</h1>
              <Button
                v-if="isArticleOwner"
                variant="default"
                size="default"
                @click="handleEdit"
                class="edit-btn-header"
              >
                <Pencil class="mr-1 h-4 w-4" />
                编辑
              </Button>
            </div>
            <div class="article-meta">
              <Badge
                v-if="article.category"
                class="clickable-tag"
                @click="handleCategoryClick(article.category.id)"
              >
                {{ article.category.name }}
              </Badge>
              <span
                class="author-info clickable-author"
                @click="handleAuthorClick"
              >
                <Avatar class="h-6 w-6 author-avatar">
                  <AvatarImage :src="article.author?.avatar" />
                  <AvatarFallback>{{ (article.author?.nickname || article.author?.username || '?')[0] }}</AvatarFallback>
                </Avatar>
                {{ article.author?.nickname || article.author?.username }}
              </span>
              <span class="inline-flex items-center gap-1 text-sm"><Eye class="h-4 w-4" /> {{ article.view_count }}</span>
              <span class="inline-flex items-center gap-1 text-sm"><Clock class="h-4 w-4" /> {{ formatDate(article.created_at) }}</span>
            </div>
            <div class="article-tags">
              <Badge
                v-for="tag in article.tags"
                :key="tag.id"
                variant="secondary"
                class="clickable-tag"
                @click="handleTagClick(tag.id)"
              >
                {{ tag.name }}
              </Badge>
            </div>
          </div>

          <div class="article-content" id="article-preview" v-if="article"></div>
          <div v-else class="article-loading">正在加载内容...</div>

          <div class="article-actions">
            <Button
              :variant="isLiked ? 'default' : 'outline'"
              @click="handleLike"
              :disabled="likeLoading"
            >
              <Loader2 v-if="likeLoading" class="mr-1 h-4 w-4 animate-spin" />
              <Star v-else class="mr-1 h-4 w-4" />
              {{ isLiked ? '已点赞' : '点赞' }}
              ({{ article.like_count || 0 }})
            </Button>

            <Button
              v-if="userStore.isLoggedIn"
              :variant="isFavorited ? 'accent' : 'outline'"
              @click="handleFavorite"
              :disabled="favoriteLoading"
            >
              <Loader2 v-if="favoriteLoading" class="mr-1 h-4 w-4 animate-spin" />
              <Bookmark v-else class="mr-1 h-4 w-4" />
              {{ isFavorited ? '已收藏' : '收藏' }}
              ({{ article.favorite_count || 0 }})
            </Button>

            <Button
              v-if="isArticleOwner"
              variant="default"
              @click="handleEdit"
            >
              <Pencil class="mr-1 h-4 w-4" />
              编辑文章
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card v-if="articleCommentEnabled" class="comments-card">
        <CardContent class="p-6">
          <h3>评论区</h3>

          <form v-if="userStore.isLoggedIn" @submit.prevent="submitComment" class="comment-form">
            <Textarea
              v-model="commentContent"
              :rows="4"
              placeholder="写下你的评论..."
              maxlength="500"
            />
            <Button variant="default" @click="submitComment" :disabled="submitting" class="mt-2">
              <Loader2 v-if="submitting" class="mr-1 h-4 w-4 animate-spin" />
              发表评论
            </Button>
          </form>
          <Alert v-else>
            <AlertDescription>
              请<router-link to="/login" class="text-primary underline hover:no-underline">登录</router-link>后发表评论
            </AlertDescription>
          </Alert>

          <div class="comment-list">
            <div v-for="comment in comments" :key="comment.id" class="comment-item">
              <Avatar
                class="clickable-avatar h-10 w-10"
                @click="handleUserClick(comment.user_id)"
              >
                <AvatarImage :src="comment.user?.avatar" />
                <AvatarFallback>{{ (comment.user?.nickname || comment.nickname || '?')[0] }}</AvatarFallback>
              </Avatar>
              <div class="comment-content">
                <div class="comment-header">
                  <div class="comment-author-section">
                    <span class="comment-author">
                      {{ comment.user?.nickname || comment.nickname }}
                    </span>
                    <Badge v-if="isCommentAuthor(comment)" variant="accent" class="text-xs">
                      作者
                    </Badge>
                  </div>
                  <div class="comment-meta-section">
                    <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
                  </div>
                </div>
                <p>{{ comment.content }}</p>
                <div class="comment-actions">
                  <Button
                    variant="ghost"
                    size="sm"
                    :class="comment.isLiked ? 'text-primary' : ''"
                    @click="handleCommentLike(comment)"
                    :disabled="comment.likeLoading"
                  >
                    <Loader2 v-if="comment.likeLoading" class="mr-1 h-3 w-3 animate-spin" />
                    <Star v-else class="mr-1 h-3 w-3" />
                    {{ comment.like_count || 0 }}
                  </Button>
                  <Button
                    v-if="userStore.isLoggedIn"
                    variant="ghost"
                    size="sm"
                    @click="showReplyInput(comment)"
                  >
                    <MessageCircle class="mr-1 h-3 w-3" /> 回复
                  </Button>
                  <Button
                    v-if="canDeleteComment(comment)"
                    variant="ghost"
                    size="sm"
                    class="text-destructive hover:text-destructive"
                    @click="handleDeleteComment(comment)"
                  >
                    删除
                  </Button>
                </div>

                <!-- 回复输入框 -->
                <div v-if="comment.showReply" class="reply-input">
                  <Textarea
                    v-model="comment.replyContent"
                    :rows="3"
                    placeholder="写下你的回复..."
                    maxlength="500"
                  />
                  <div class="reply-actions">
                    <Button variant="outline" size="sm" @click="cancelReply(comment)">取消</Button>
                    <Button variant="default" size="sm" @click="submitReply(comment)" :disabled="comment.replying">
                      <Loader2 v-if="comment.replying" class="mr-1 h-3 w-3 animate-spin" />
                      发表回复
                    </Button>
                  </div>
                </div>

                <!-- 回复列表 -->
                <div v-if="comment.replies && comment.replies.length > 0" class="replies-list">
                  <div v-for="reply in getDisplayedReplies(comment)" :key="reply.id" class="reply-item">
                    <Avatar
                      class="clickable-avatar h-8 w-8"
                      @click="handleUserClick(reply.user_id)"
                    >
                      <AvatarImage :src="reply.user?.avatar" />
                      <AvatarFallback>{{ (reply.user?.nickname || reply.nickname || '?')[0] }}</AvatarFallback>
                    </Avatar>
                    <div class="reply-content">
                      <div class="reply-header">
                        <span class="reply-author">
                          {{ reply.user?.nickname || reply.nickname }}
                          <Badge v-if="isCommentAuthor(reply)" variant="accent" class="text-xs">
                            作者
                          </Badge>
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
                        <Button
                          variant="ghost"
                          size="sm"
                          :class="reply.isLiked ? 'text-primary' : ''"
                          @click="handleCommentLike(reply)"
                          :disabled="reply.likeLoading"
                        >
                          <Loader2 v-if="reply.likeLoading" class="mr-1 h-3 w-3 animate-spin" />
                          <Star v-else class="mr-1 h-3 w-3" />
                          {{ reply.like_count || 0 }}
                        </Button>
                        <Button
                          v-if="userStore.isLoggedIn"
                          variant="ghost"
                          size="sm"
                          @click="showReplyToReplyInput(comment, reply)"
                        >
                          <MessageCircle class="mr-1 h-3 w-3" /> 回复
                        </Button>
                        <Button
                          v-if="canDeleteComment(reply)"
                          variant="ghost"
                          size="sm"
                          class="text-destructive hover:text-destructive"
                          @click="handleDeleteComment(reply)"
                        >
                          删除
                        </Button>
                      </div>

                      <!-- 回复的回复输入框 -->
                      <div v-if="reply.showReply" class="reply-input">
                        <Textarea
                          v-model="reply.replyContent"
                          :rows="3"
                          placeholder="写下你的回复..."
                          maxlength="500"
                        />
                        <div class="reply-actions">
                          <Button variant="outline" size="sm" @click="cancelReplyToReply(reply)">取消</Button>
                          <Button variant="default" size="sm" @click="submitReplyToReply(comment, reply)" :disabled="reply.replying">
                            <Loader2 v-if="reply.replying" class="mr-1 h-3 w-3 animate-spin" />
                            发表回复
                          </Button>
                        </div>
                      </div>
                    </div>
                  </div>

                  <!-- 展开更多子评论按钮 -->
                  <div v-if="hasMoreReplies(comment)" class="load-more-replies">
                    <Button
                      variant="ghost"
                      size="sm"
                      @click="loadMoreReplies(comment)"
                      :disabled="comment.loadingReplies"
                    >
                      <Loader2 v-if="comment.loadingReplies" class="mr-1 h-3 w-3 animate-spin" />
                      <ChevronDown v-else class="mr-1 h-3 w-3" />
                      {{ comment.loadingReplies ? '加载中...' : `展开更多回复 (${comment.reply_count - getDisplayedReplies(comment).length} 条)` }}
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="load-more-container" v-if="hasMoreComments">
            <Button
              variant="outline"
              @click="loadMoreComments"
              :disabled="loadingMoreComments"
              class="load-more-btn"
            >
              <Loader2 v-if="loadingMoreComments" class="mr-1 h-4 w-4 animate-spin" />
              <ChevronDown v-else class="mr-1 h-4 w-4" />
              {{ loadingMoreComments ? '加载中...' : '加载更多评论' }}
            </Button>
            <div class="comment-count-info">
              已显示 {{ comments.length }} / {{ commentsTotal }} 条评论
            </div>
          </div>
          <div v-else-if="comments.length > 0" class="no-more-comments">
            已显示全部 {{ comments.length }} 条评论
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Confirmation Dialog -->
    <Dialog :open="showConfirmDialog" @update:open="onConfirmDialogUpdateOpen">
      <DialogContent>
        <DialogHeader>
          <DialogTitle>确认</DialogTitle>
          <DialogDescription>{{ confirmDialogMessage }}</DialogDescription>
        </DialogHeader>
        <DialogFooter>
          <Button variant="outline" @click="cancelConfirmDialog">取消</Button>
          <Button @click="confirmDialogCallback?.()">确认</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { Star, Bookmark, Eye, Clock, Pencil, MessageCircle, ChevronDown, Loader2 } from 'lucide-vue-next'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { loadCodeTheme, loadHighlightTheme, getMarkdownTheme } from '@/utils/codeTheme'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Textarea } from '@/components/ui/textarea'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'

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
const articleCommentEnabled = ref(true)

const showConfirmDialog = ref(false)
const confirmDialogMessage = ref('')
const confirmDialogCallback = ref(null)
let _confirmDialogReject = null

const confirmDialog = (message) => {
  return new Promise((resolve, reject) => {
    confirmDialogMessage.value = message
    _confirmDialogReject = reject
    confirmDialogCallback.value = () => {
      _confirmDialogReject = null
      showConfirmDialog.value = false
      resolve()
    }
    showConfirmDialog.value = true
  })
}

const onConfirmDialogUpdateOpen = (open) => {
  showConfirmDialog.value = open
  if (!open && _confirmDialogReject) {
    const rejectFn = _confirmDialogReject
    _confirmDialogReject = null
    rejectFn('cancel')
  }
}

const cancelConfirmDialog = () => {
  if (_confirmDialogReject) {
    const rejectFn = _confirmDialogReject
    _confirmDialogReject = null
    showConfirmDialog.value = false
    rejectFn('cancel')
  }
}

const renderMarkdown = async () => {
  if (!article.value || !article.value.content) {
    return
  }
  
  const codeThemeValue = await loadCodeTheme()
  await loadHighlightTheme(codeThemeValue)
  const mdTheme = await getMarkdownTheme()
  
  nextTick(() => {
    const previewDiv = document.getElementById('article-preview')
    if (!previewDiv) {
      console.error('Preview element not found')
      return
    }
    
    console.log('Rendering with Vditor.preview, code theme:', codeThemeValue, 'markdown theme:', mdTheme)
    
    Vditor.preview(previewDiv, article.value.content, {
      mode: mdTheme || 'light',
      markdown: {
        toc: true,
        mark: true,
        footnotes: true,
        autoSpace: true,
      },
      hljs: {
        lineNumber: true,
        style: codeThemeValue || 'github',
        enable: true
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

const isArticleOwner = computed(() => {
  if (!userStore.isLoggedIn || !article.value || !userStore.user) {
    return false
  }
  
  const authorId = article.value.author_id || article.value.author?.id
  const userId = userStore.user.id
  
  if (!authorId || !userId) {
    return false
  }
  
  return Number(authorId) === Number(userId)
})

const isCommentAuthor = (comment) => {
  if (!article.value || !comment) return false
  
  const articleAuthorId = article.value.author_id || article.value.author?.id
  if (!articleAuthorId) return false
  
  const commentAuthorId = comment.user_id || comment.user?.id
  if (!commentAuthorId) return false
  
  return Number(articleAuthorId) === Number(commentAuthorId)
}

const loadArticle = async () => {
  try {
    const response = await api.get(`/articles/${route.params.id}`)
    article.value = response.data
    
    renderedContent.value = 'loading'
    
    renderMarkdown()
  } catch (error) {
    toast.error('文章加载失败')
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
    
    if (append) {
      comments.value = [...comments.value, ...newComments]
    } else {
      comments.value = newComments
    }
    
    commentsTotal.value = response.data.total || 0
    
    hasMoreComments.value = comments.value.length < commentsTotal.value
    
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

const loadMoreComments = async () => {
  if (loadingMoreComments.value || !hasMoreComments.value) return
  
  loadingMoreComments.value = true
  try {
    commentsPage.value++
    await loadComments(true)
  } finally {
    loadingMoreComments.value = false
  }
}

const canDeleteComment = (comment) => {
  if (!userStore.isLoggedIn || !userStore.user) return false
  
  if (userStore.user.role === 'admin') return true
  
  if (isArticleOwner.value) return true
  
  const commentUserId = comment.user_id || comment.user?.id || comment.userId
  
  if (!commentUserId || commentUserId === 0) {
    return false
  }
  
  return Number(userStore.user.id) === Number(commentUserId)
}

const submitComment = async () => {
  if (!commentContent.value.trim()) {
    toast.warning('请输入评论内容')
    return
  }

  submitting.value = true
  try {
    const response = await api.post('/comments', {
      article_id: parseInt(route.params.id),
      content: commentContent.value
    })
    if (response.message && response.message !== 'success') {
      toast.success(response.message)
    }
    commentContent.value = ''
    commentsPage.value = 1
    await loadComments(false)
  } catch (error) {
    toast.error('评论发表失败')
  } finally {
    submitting.value = false
  }
}

const showReplyInput = (comment, replyTo = null) => {
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }
  
  comments.value.forEach(c => {
    if (c.id !== comment.id) {
      c.showReply = false
    }
  })
  
  comment.showReply = true
  comment.replyTo = replyTo || comment
  comment.replyContent = replyTo ? `@${replyTo.user?.nickname || replyTo.nickname} ` : ''
}

const showReplyToReplyInput = (parentComment, replyTo) => {
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }
  
  comments.value.forEach(c => {
    c.showReply = false
    if (c.replies) {
      c.replies.forEach(r => {
        r.showReply = false
      })
    }
  })
  
  replyTo.showReply = true
  replyTo.replyTo = replyTo
  replyTo.replyContent = `@${replyTo.user?.nickname || replyTo.nickname} `
}

const cancelReply = (comment) => {
  comment.showReply = false
  comment.replyContent = ''
  comment.replyTo = null
}

const submitReply = async (comment) => {
  if (!comment.replyContent || !comment.replyContent.trim()) {
    toast.warning('请输入回复内容')
    return
  }

  let content = comment.replyContent.trim()
  if (comment.replyTo) {
    const mentionPattern = new RegExp(`^@${comment.replyTo.user?.nickname || comment.replyTo.nickname}\\s*`, 'i')
    content = content.replace(mentionPattern, '')
  }

  if (!content) {
    toast.warning('请输入回复内容')
    return
  }

  comment.replying = true
  try {
    const response = await api.post('/comments', {
      article_id: parseInt(route.params.id),
      content: content,
      parent_id: comment.replyTo ? comment.replyTo.id : comment.id
    })
    if (response.message && response.message !== 'success') {
      toast.success(response.message)
    }
    comment.showReply = false
    comment.replyContent = ''
    comment.replyTo = null
    commentsPage.value = 1
    await loadComments(false)
  } catch (error) {
    toast.error('回复失败')
  } finally {
    comment.replying = false
  }
}

const submitReplyToReply = async (parentComment, reply) => {
  if (!reply.replyContent || !reply.replyContent.trim()) {
    toast.warning('请输入回复内容')
    return
  }

  let content = reply.replyContent.trim()
  if (reply.replyTo) {
    const mentionPattern = new RegExp(`^@${reply.replyTo.user?.nickname || reply.replyTo.nickname}\\s*`, 'i')
    content = content.replace(mentionPattern, '')
  }

  if (!content) {
    toast.warning('请输入回复内容')
    return
  }

  reply.replying = true
  try {
    const response = await api.post('/comments', {
      article_id: parseInt(route.params.id),
      content: content,
      parent_id: reply.replyTo ? reply.replyTo.id : reply.id
    })
    if (response.message && response.message !== 'success') {
      toast.success(response.message)
    }
    reply.showReply = false
    reply.replyContent = ''
    reply.replyTo = null
    commentsPage.value = 1
    await loadComments(false)
  } catch (error) {
    toast.error('回复失败')
  } finally {
    reply.replying = false
  }
}

const cancelReplyToReply = (reply) => {
  reply.showReply = false
  reply.replyContent = ''
  reply.replyTo = null
}

const getDirectReplies = (comment) => {
  if (!comment.replies) return []
  return comment.replies.filter(reply => reply.parent_id === comment.id)
}

const getDisplayedReplies = (comment) => {
  if (!comment.replies) return []
  return comment.replies
}

const hasMoreReplies = (comment) => {
  if (!comment.reply_count) return false
  const displayedCount = getDisplayedReplies(comment).length
  return comment.reply_count > displayedCount
}

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
    
    comment.replies = [...(comment.replies || []), ...newReplies]
    comment.repliesPage = nextPage
    
    for (const reply of newReplies) {
      await checkCommentLiked(reply)
    }
  } catch (error) {
    console.error('Failed to load more replies:', error)
    toast.error('加载更多回复失败')
  } finally {
    comment.loadingReplies = false
  }
}

const getReplyTargetName = (parentComment, reply) => {
  if (!reply.parent_id || reply.parent_id === parentComment.id) {
    return parentComment.user?.nickname || parentComment.nickname
  }
  const targetReply = parentComment.replies?.find(r => r.id === reply.parent_id)
  return targetReply ? (targetReply.user?.nickname || targetReply.nickname) : ''
}

const checkCommentLiked = async (comment) => {
  if (!comment || !comment.id) return
  
  try {
    const response = await api.get(`/comments/${comment.id}/is-liked`)
    const isLiked = response.data?.is_liked ?? false
    comment.isLiked = isLiked
    console.debug(`Comment ${comment.id} like status:`, isLiked)
  } catch (error) {
    comment.isLiked = false
    console.debug(`Failed to check comment ${comment.id} like status:`, error)
  }
}

const handleCommentLike = async (comment) => {
  if (!comment || !comment.id || comment.likeLoading) return
  
  comment.likeLoading = true
  const oldIsLiked = comment.isLiked
  const commentId = comment.id
  
  try {
    if (comment.isLiked) {
      console.debug(`Unliking comment ${comment.id}, current count: ${comment.like_count}`)
      const response = await api.delete(`/comments/${comment.id}/like`)
      comment.isLiked = false
      commentsPage.value = 1
      await loadComments(false)
      const updatedComment = comments.value.find(c => c.id === commentId) || 
                            comments.value.flatMap(c => c.replies || []).find(r => r.id === commentId)
      if (updatedComment) {
        await checkCommentLiked(updatedComment)
      }
      const message = response?.message || '取消点赞成功'
      toast.success(message)
      console.debug(`Comment ${comment.id} unliked, new count:`, updatedComment?.like_count)
    } else {
      console.debug(`Liking comment ${comment.id}, current count: ${comment.like_count}`)
      const response = await api.post(`/comments/${comment.id}/like`)
      comment.isLiked = true
      commentsPage.value = 1
      await loadComments(false)
      const updatedComment = comments.value.find(c => c.id === commentId) || 
                            comments.value.flatMap(c => c.replies || []).find(r => r.id === commentId)
      if (updatedComment) {
        await checkCommentLiked(updatedComment)
      }
      const message = response?.message || '点赞成功'
      toast.success(message)
      console.debug(`Comment ${comment.id} liked, new count:`, updatedComment?.like_count)
    }
  } catch (error) {
    comment.isLiked = oldIsLiked
    console.error(`Failed to ${oldIsLiked ? 'unlike' : 'like'} comment ${comment.id}:`, error)
    await checkCommentLiked(comment)
    commentsPage.value = 1
    await loadComments(false)
  } finally {
    comment.likeLoading = false
  }
}

const handleDeleteComment = async (comment) => {
  try {
    await confirmDialog('确定要删除这条评论吗？')
    await api.delete(`/comments/${comment.id}`)
    toast.success('删除成功')
    commentsPage.value = 1
    await loadComments(false)
    await loadArticle()
  } catch (error) {
    if (error !== 'cancel') {
      const errorMessage = error.response?.data?.message || error.message || '删除失败'
      toast.error(errorMessage)
    }
  }
}


const checkLiked = async () => {
  if (!userStore.isLoggedIn) {
    isLiked.value = false
    return
  }
  try {
    const response = await api.get(`/articles/${route.params.id}/is-liked`)
    isLiked.value = response.data.is_liked || response.data.liked || false
  } catch (error) {
    console.error('Failed to check like status:', error)
    isLiked.value = false
  }
}

const handleLike = async () => {
  if (!article.value) return
  
  if (!userStore.isLoggedIn) {
    toast.warning('请先登录')
    router.push('/login')
    return
  }
  
  likeLoading.value = true
  try {
    await api.post(`/articles/${route.params.id}/like`)
    
    isLiked.value = !isLiked.value
    
    if (isLiked.value) {
      toast.success('点赞成功')
      article.value.like_count = (article.value.like_count || 0) + 1
    } else {
      toast.success('取消点赞')
      article.value.like_count = Math.max(0, (article.value.like_count || 0) - 1)
    }
  } catch (error) {
    toast.error(error.response?.data?.message || '操作失败')
    await checkLiked()
    await loadArticle()
  } finally {
    likeLoading.value = false
  }
}

const checkFavorited = async () => {
  if (!userStore.isLoggedIn) {
    isFavorited.value = false
    return
  }
  
  try {
    const response = await api.get(`/articles/${route.params.id}/is-favorited`)
    isFavorited.value = response.data.is_favorited || response.data.favorited || false
  } catch (error) {
    console.error('Failed to check favorite status:', error)
    isFavorited.value = false
  }
}

const handleFavorite = async () => {
  if (!article.value || !userStore.isLoggedIn) return
  
  favoriteLoading.value = true
  try {
    if (isFavorited.value) {
      await api.delete(`/articles/${route.params.id}/favorite`)
      toast.success('取消收藏成功')
      isFavorited.value = false
      article.value.favorite_count = Math.max(0, (article.value.favorite_count || 0) - 1)
    } else {
      await api.post(`/articles/${route.params.id}/favorite`)
      toast.success('收藏成功')
      isFavorited.value = true
      article.value.favorite_count = (article.value.favorite_count || 0) + 1
    }
  } catch (error) {
    toast.error('操作失败')
    checkFavorited()
    loadArticle()
  } finally {
    favoriteLoading.value = false
  }
}

const handleEdit = () => {
  if (!article.value) return
  router.push(`/dashboard/articles/${article.value.id}/edit`)
}

const handleCategoryClick = (categoryId) => {
  if (!categoryId) return
  router.push({
    path: '/blog',
    query: { category_id: categoryId }
  })
}

const handleTagClick = (tagId) => {
  if (!tagId) return
  router.push({
    path: '/blog',
    query: { tag_id: tagId }
  })
}

const handleUserClick = (userId) => {
  if (!userId) return
  router.push(`/users/${userId}`)
}

const handleAuthorClick = () => {
  if (!article.value || !article.value.author) return
  const authorId = article.value.author_id || article.value.author.id
  if (!authorId) return
  router.push(`/users/${authorId}`)
}

const loadCommentSettings = async () => {
  try {
    const response = await api.get('/settings/public')
    const settings = response.data || {}
    articleCommentEnabled.value = settings.article_comment_enabled !== '0' && settings.article_comment_enabled !== 'false'
  } catch (error) {
    console.error('Failed to load comment settings:', error)
    articleCommentEnabled.value = true
  }
}

onMounted(async () => {
  if (userStore.isLoggedIn && !userStore.user) {
    await userStore.fetchProfile()
  }
  
  await loadCommentSettings()
  
  loadArticle()
  if (articleCommentEnabled.value) {
    loadComments()
  }
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
  padding: var(--spacing-xl) 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.blog-detail .container {
  max-width: 1200px;
}

.article-card {
  margin-bottom: var(--spacing-lg);
}

.comments-card {
}

.article-header {
  border-bottom: 1px solid var(--theme-border-light);
  padding-bottom: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
}

.article-title-row {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-md);
}

.article-title-row h1 {
  flex: 1;
  margin: 0;
  font-size: var(--font-size-2xl);
  color: var(--theme-text-primary);
  line-height: var(--line-height-tight);
}

.edit-btn-header {
  flex-shrink: 0;
}

.article-meta {
  display: flex;
  gap: var(--spacing-md);
  color: var(--theme-text-secondary);
  margin-bottom: var(--spacing-sm);
  flex-wrap: wrap;
}

.article-meta span {
  display: flex;
  align-items: center;
  gap: var(--spacing-xs);
  font-size: var(--font-size-sm);
}

.author-info {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.author-avatar {
  flex-shrink: 0;
}

.clickable-author {
  cursor: pointer;
  transition: all var(--transition-slow);
  padding: var(--spacing-xs) var(--spacing-sm);
  border-radius: var(--radius-sm);
}

.clickable-author:hover {
  background-color: var(--theme-bg-hover);
  color: var(--theme-primary);
}

.clickable-tag {
  cursor: pointer;
  transition: all var(--transition-slow);
}

.clickable-tag:hover {
  transform: translateY(-2px);
  opacity: 0.8;
}

.clickable-avatar {
  cursor: pointer;
  transition: all var(--transition-slow);
}

.clickable-avatar:hover {
  transform: scale(1.1);
  box-shadow: var(--shadow-sm);
}

.article-tags {
  display: flex;
  gap: var(--spacing-sm);
  flex-wrap: wrap;
}

.article-content {
  margin-bottom: var(--spacing-lg);
  min-height: 200px;
}

#article-preview {
  padding: var(--spacing-lg);
  max-width: 920px;
  margin: 0 auto;
  background-color: var(--theme-content-bg);
  border-radius: var(--radius-md);
  border: 1px solid var(--theme-border-light);
  box-shadow: var(--shadow-sm);
}

#article-preview :deep(code:not(pre code)) {
  background-color: var(--theme-bg-hover);
  color: var(--theme-text-primary);
  padding: 0.2em 0.4em;
  border-radius: var(--radius-sm);
  font-size: var(--font-size-xs);
}

.article-loading {
  padding: var(--spacing-xl) 0;
  text-align: center;
  color: var(--theme-text-tertiary);
  font-size: var(--font-size-sm);
}

.article-actions {
  text-align: center;
  padding-top: var(--spacing-md);
  border-top: 1px solid var(--theme-border-light);
  display: flex;
  gap: var(--spacing-sm);
  justify-content: center;
}

.comments-card h3 {
  margin-bottom: var(--spacing-md);
  font-size: var(--font-size-lg);
  color: var(--theme-text-primary);
  font-weight: 600;
}

.comment-form {
  margin-bottom: var(--spacing-lg);
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
  border-left: 2px solid var(--color-border-lighter);
}

.reply-item {
  display: flex;
  gap: var(--spacing-sm);
  padding: var(--spacing-md) 0;
  border-bottom: 1px solid var(--color-border-lighter);
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
}

.load-more-replies {
  margin-top: var(--spacing-sm);
  padding: var(--spacing-sm) 0;
  text-align: center;
  border-top: 1px solid var(--color-border-lighter);
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
</style>
