<template>
  <div class="article-view">
    <div class="page-header">
      <h2>文章详情</h2>
      <div class="header-actions">
        <el-button @click="handleEdit" type="primary">编辑</el-button>
        <el-button @click="$router.back()">返回</el-button>
      </div>
    </div>

    <el-card v-if="article" v-loading="loading">
      <div class="article-header">
        <h1 class="article-title">{{ article.title }}</h1>
        <div class="article-meta">
          <el-tag v-if="article.category" size="small">{{ article.category.name }}</el-tag>
          <el-tag v-if="article.status === 1" type="success" size="small">已发布</el-tag>
          <el-tag v-else type="info" size="small">草稿</el-tag>
          <el-tag v-if="article.is_top" type="danger" size="small">置顶</el-tag>
          <el-tag v-if="article.is_recommend" type="warning" size="small">推荐</el-tag>
          <span class="meta-item">作者：{{ article.author?.nickname || article.author?.username }}</span>
          <span class="meta-item">阅读：{{ article.view_count }}</span>
          <span class="meta-item">点赞：{{ article.like_count }}</span>
          <span class="meta-item">评论：{{ article.comment_count }}</span>
          <span class="meta-item">创建时间：{{ formatDate(article.created_at) }}</span>
        </div>
        <div v-if="article.tags && article.tags.length" class="article-tags">
          <el-tag v-for="tag in article.tags" :key="tag.id" size="small" type="info">
            {{ tag.name }}
          </el-tag>
        </div>
      </div>

      <el-divider />

      <div v-if="article.summary" class="article-summary">
        <h3>摘要</h3>
        <p>{{ article.summary }}</p>
      </div>

      <div v-if="article.cover" class="article-cover">
        <h3>封面图</h3>
        <el-image :src="article.cover" fit="contain" style="max-width: 100%; max-height: 400px;" />
      </div>

      <div class="article-content">
        <h3>内容</h3>
        <div id="article-preview" class="content-preview"></div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import adminApi from '@/utils/adminApi'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import dayjs from 'dayjs'
import { loadCodeTheme, loadHighlightTheme, getMarkdownTheme } from '@/utils/codeTheme'

const route = useRoute()
const router = useRouter()

const article = ref(null)
const loading = ref(false)

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const handleEdit = () => {
  router.push(`/articles/${route.params.id}/edit`)
}

const renderMarkdown = async () => {
  if (!article.value || !article.value.content) return
  
  // 加载代码主题配置
  const codeThemeValue = await loadCodeTheme()
  // 加载 highlight.js 主题样式
  await loadHighlightTheme(codeThemeValue)
  // 加载 Markdown 主题配置
  const mdTheme = await getMarkdownTheme()
  
  nextTick(() => {
    const previewDiv = document.getElementById('article-preview')
    if (!previewDiv) {
      console.error('Preview element not found')
      return
    }
    
    Vditor.preview(previewDiv, article.value.content, {
      mode: mdTheme || 'light',
      markdown: {
        toc: true,
        mark: true,
        footnotes: true,
        autoSpace: true,
      },
      hljs: {
        lineNumber: false,
        style: codeThemeValue || 'github'
      },
      speech: {
        enable: false
      },
      anchor: 1,
      after: () => {
        console.log('Article content rendered')
      }
    })
  })
}

const loadArticle = async () => {
  loading.value = true
  try {
    const response = await adminApi.get(`/admin/articles/${route.params.id}`)
    article.value = response.data
    
    // 渲染 Markdown 内容
    renderMarkdown()
  } catch (error) {
    console.error('Load article error:', error)
    ElMessage.error('加载文章失败')
    router.back()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadArticle()
})
</script>

<style scoped>
.article-view {
  padding: var(--spacing-lg);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.header-actions {
  display: flex;
  gap: var(--spacing-sm);
}

.article-header {
  margin-bottom: var(--spacing-lg);
}

.article-title {
  font-size: var(--font-size-2xl);
  font-weight: 600;
  margin: 0 0 var(--spacing-md);
  color: var(--color-text-primary);
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  color: var(--color-text-tertiary);
  font-size: var(--font-size-sm);
  margin-bottom: 12px;
}

.meta-item {
  display: flex;
  align-items: center;
}

.article-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
}

.article-summary {
  margin-bottom: var(--spacing-lg);
}

.article-summary h3 {
  font-size: var(--font-size-lg);
  margin-bottom: 12px;
  color: var(--color-text-primary);
}

.article-summary p {
  font-size: var(--font-size-base);
  line-height: var(--line-height-relaxed);
  color: var(--color-text-secondary);
}

.article-cover {
  margin-bottom: var(--spacing-lg);
}

.article-cover h3 {
  font-size: var(--font-size-lg);
  margin-bottom: 12px;
  color: var(--color-text-primary);
}

.article-content {
  margin-top: var(--spacing-lg);
}

.article-content h3 {
  font-size: var(--font-size-lg);
  margin-bottom: var(--spacing-md);
  color: var(--color-text-primary);
}

.content-preview {
  padding: var(--spacing-lg);
  background: var(--color-bg-secondary);
  border-radius: var(--radius-sm);
  min-height: 200px;
}

.content-preview :deep(.vditor-reset) {
  background: transparent;
}
</style>

