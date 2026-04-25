<template>
  <div class="article-view">
    <div class="page-header">
      <h2>文章详情</h2>
      <div class="header-actions">
        <Button @click="handleEdit">编辑</Button>
        <Button variant="outline" @click="$router.back()">返回</Button>
      </div>
    </div>

    <Card v-if="article">
      <CardContent class="pt-6">
        <div v-if="loading" class="flex items-center justify-center py-8">
          <Loader2 class="h-6 w-6 animate-spin text-muted-foreground" />
        </div>
        <template v-else>
          <div class="article-header">
            <h1 class="article-title">{{ article.title }}</h1>
            <div class="article-meta">
              <Badge v-if="article.category" variant="secondary">{{ article.category.name }}</Badge>
              <Badge v-if="article.status === 1" variant="secondary">已发布</Badge>
              <Badge v-else variant="outline">草稿</Badge>
              <Badge v-if="article.is_top" variant="destructive">置顶</Badge>
              <Badge v-if="article.is_recommend" variant="accent">推荐</Badge>
              <span class="meta-item">作者：{{ article.author?.nickname || article.author?.username }}</span>
              <span class="meta-item">阅读：{{ article.view_count }}</span>
              <span class="meta-item">点赞：{{ article.like_count }}</span>
              <span class="meta-item">评论：{{ article.comment_count }}</span>
              <span class="meta-item">创建时间：{{ formatDate(article.created_at) }}</span>
            </div>
            <div v-if="article.tags && article.tags.length" class="article-tags">
              <Badge v-for="tag in article.tags" :key="tag.id" variant="outline">
                {{ tag.name }}
              </Badge>
            </div>
          </div>

          <Separator />

          <div v-if="article.summary" class="article-summary">
            <h3>摘要</h3>
            <p>{{ article.summary }}</p>
          </div>

          <div v-if="article.cover" class="article-cover">
            <h3>封面图</h3>
            <img :src="article.cover" class="max-w-full max-h-[400px] object-contain" />
          </div>

          <div class="article-content">
            <h3>内容</h3>
            <div id="article-preview" class="content-preview"></div>
          </div>
        </template>
      </CardContent>
    </Card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
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

  const codeThemeValue = await loadCodeTheme()
  await loadHighlightTheme(codeThemeValue)
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

    renderMarkdown()
  } catch (error) {
    console.error('Load article error:', error)
    toast.error('加载文章失败')
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
