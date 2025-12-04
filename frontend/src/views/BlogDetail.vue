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

        <div class="article-content" v-html="renderedContent"></div>

        <div class="article-actions">
          <el-button :icon="Star" @click="handleLike" :disabled="liked">
            点赞 ({{ article.like_count }})
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
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Star } from '@element-plus/icons-vue'
import api from '@/utils/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import MarkdownIt from 'markdown-it'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

const route = useRoute()
const userStore = useUserStore()

const article = ref(null)
const comments = ref([])
const commentsPage = ref(1)
const commentsTotal = ref(0)
const commentContent = ref('')
const submitting = ref(false)
const liked = ref(false)

const md = new MarkdownIt({
  highlight: (str, lang) => {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value
      } catch (__) {}
    }
    return ''
  }
})

const renderedContent = computed(() => {
  return article.value ? md.render(article.value.content) : ''
})

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const loadArticle = async () => {
  try {
    const response = await api.get(`/articles/${route.params.id}`)
    article.value = response.data
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

const handleLike = async () => {
  try {
    await api.post(`/articles/${route.params.id}/like`)
    article.value.like_count++
    liked.value = true
    ElMessage.success('点赞成功')
  } catch (error) {
    ElMessage.error('点赞失败')
  }
}

onMounted(() => {
  loadArticle()
  loadComments()
})
</script>

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

.article-content {
  margin-bottom: 30px;
  line-height: 2;
  font-size: 16px;
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3),
.article-content :deep(h4),
.article-content :deep(h5),
.article-content :deep(h6) {
  margin: 20px 0 10px;
}

.article-content :deep(p) {
  margin-bottom: 15px;
}

.article-content :deep(code) {
  background-color: #f5f5f5;
  padding: 2px 6px;
  border-radius: 3px;
}

.article-content :deep(pre) {
  background-color: #f6f8fa;
  padding: 15px;
  border-radius: 6px;
  overflow-x: auto;
  margin: 15px 0;
}

.article-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 15px 0;
}

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

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>

