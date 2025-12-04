<template>
  <div class="blog-detail">
    <div class="container">
      <el-card v-if="article" class="article-card">
        <div class="article-header">
          <h1>{{ article.title }}</h1>
          <div class="article-meta">
            <el-tag v-if="article.category">{{ article.category.name }}</el-tag>
            <span><el-icon><User /></el-icon> {{ article.author?.nickname || article.author?.username }}</span>
            <span><el-icon><View /></el-icon> {{ article.view_count }}</span>
            <span><el-icon><Clock /></el-icon> {{ formatDate(article.created_at) }}</span>
          </div>
          <div class="article-tags">
            <el-tag v-for="tag in article.tags" :key="tag.id" size="small" type="info">
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
            ({{ article.like_count }})
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
            <el-avatar :src="comment.user?.avatar" />
            <div class="comment-content">
              <div class="comment-header">
                <span class="comment-author">{{ comment.user?.nickname || comment.nickname }}</span>
                <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
              </div>
              <p>{{ comment.content }}</p>
              <div class="comment-actions">
                <el-button text size="small" @click="toggleCommentLike(comment)">
                  <el-icon><Star /></el-icon> 
                  {{ comment.like_count || 0 }}
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
            </div>
          </div>
        </div>

        <div class="pagination" v-if="commentsTotal > 0">
          <el-pagination
            v-model:current-page="commentsPage"
            :page-size="10"
            :total="commentsTotal"
            layout="prev, pager, next"
            @current-change="loadComments"
          />
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Star, Collection, User, View, Clock } from '@element-plus/icons-vue'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import Vditor from 'vditor'
import 'vditor/dist/index.css'

const route = useRoute()
const userStore = useUserStore()

const article = ref(null)
const comments = ref([])
const commentsPage = ref(1)
const commentsTotal = ref(0)
const commentContent = ref('')
const submitting = ref(false)
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

const loadComments = async () => {
  try {
    const response = await api.get('/comments', {
      params: {
        article_id: route.params.id,
        page: commentsPage.value,
        page_size: 10
      }
    })
    comments.value = response.data.list || []
    commentsTotal.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load comments:', error)
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
    loadComments()
  } catch (error) {
    ElMessage.error('评论发表失败')
  } finally {
    submitting.value = false
  }
}

const handleDeleteComment = async (comment) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', { type: 'warning' })
    await api.delete(`/comments/${comment.id}`)
    ElMessage.success('删除成功')
    loadComments()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const toggleCommentLike = async (comment) => {
  try {
    await api.post(`/comments/${comment.id}/like`)
    comment.like_count = (comment.like_count || 0) + 1
    ElMessage.success('点赞成功')
  } catch (error) {
    // 如果已点赞，尝试取消点赞
    try {
      await api.delete(`/comments/${comment.id}/like`)
      comment.like_count = Math.max(0, (comment.like_count || 0) - 1)
      ElMessage.success('取消点赞成功')
    } catch (err) {
      ElMessage.error('操作失败')
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
  likeLoading.value = true
  try {
    if (isLiked.value) {
      await api.delete(`/articles/${route.params.id}/like`)
      ElMessage.success('取消点赞成功')
      isLiked.value = false
      article.value.like_count--
    } else {
      await api.post(`/articles/${route.params.id}/like`)
      ElMessage.success('点赞成功')
      isLiked.value = true
      article.value.like_count++
    }
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
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
  favoriteLoading.value = true
  try {
    if (isFavorited.value) {
      await api.delete(`/articles/${route.params.id}/favorite`)
      ElMessage.success('取消收藏成功')
      isFavorited.value = false
      article.value.favorite_count--
    } else {
      await api.post(`/articles/${route.params.id}/favorite`)
      ElMessage.success('收藏成功')
      isFavorited.value = true
      article.value.favorite_count++
    }
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    favoriteLoading.value = false
  }
}

onMounted(() => {
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
}

.article-card {
  margin-bottom: 30px;
}

.article-header {
  border-bottom: 1px solid var(--border-lighter);
  padding-bottom: 20px;
  margin-bottom: 30px;
}

.article-header h1 {
  margin-bottom: 15px;
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
  margin-bottom: 10px;
}

.comment-author {
  font-weight: 600;
  color: var(--text-primary);
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

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>

