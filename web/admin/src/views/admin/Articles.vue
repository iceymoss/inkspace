<template>
  <div class="articles">
    <div class="page-header">
      <h2>文章管理</h2>
      <Button @click="$router.push('/articles/create')">
        <Plus class="h-4 w-4 mr-1" /> 新建文章
      </Button>
    </div>

    <Card>
      <CardContent class="pt-6">
        <div class="filter-bar">
          <form class="flex flex-wrap items-end gap-4">
            <div class="space-y-1">
              <label class="text-sm font-medium">关键字</label>
              <Input
                v-model="filters.keyword"
                placeholder="标题 / 内容"
                class="w-[220px]"
                @keyup.enter="handleSearch"
              />
            </div>

            <div class="space-y-1">
              <label class="text-sm font-medium">分类</label>
              <Select v-model="filters.categoryId" @update:model-value="handleFilterChange">
                <SelectTrigger class="w-[180px]">
                  <SelectValue placeholder="全部分类" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem v-for="cat in categories" :key="cat.id" :value="String(cat.id)">
                    {{ cat.name }}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-medium">状态</label>
              <Select :model-value="filters.status !== null && filters.status !== undefined && filters.status !== '' ? String(filters.status) : undefined" @update:model-value="filters.status = $event ? Number($event) : null; handleFilterChange()">
                <SelectTrigger class="w-[140px]">
                  <SelectValue placeholder="全部状态" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="1">已发布</SelectItem>
                  <SelectItem value="0">草稿</SelectItem>
                  <SelectItem value="2">私有</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div class="space-y-1">
              <label class="text-sm font-medium">排序</label>
              <Select v-model="filters.sortBy" @update:model-value="handleSortChange">
                <SelectTrigger class="w-[160px]">
                  <SelectValue placeholder="默认排序" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">默认（置顶 + 最新）</SelectItem>
                  <SelectItem value="time">最新发布</SelectItem>
                  <SelectItem value="hot">最热排序</SelectItem>
                  <SelectItem value="view_count">最多浏览</SelectItem>
                  <SelectItem value="like_count">最多点赞</SelectItem>
                  <SelectItem value="comment_count">最多评论</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div class="flex gap-2">
              <Button @click="handleSearch">查询</Button>
              <Button variant="outline" @click="resetFilters">重置</Button>
            </div>
          </form>
        </div>

        <div class="relative">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead class="w-[80px]">ID</TableHead>
                <TableHead class="min-w-[200px]">标题</TableHead>
                <TableHead class="w-[120px]">分类</TableHead>
                <TableHead class="w-[100px]">状态</TableHead>
                <TableHead class="w-[90px] text-center">推荐</TableHead>
                <TableHead class="w-[100px]">浏览</TableHead>
                <TableHead class="w-[180px]">创建时间</TableHead>
                <TableHead class="w-[350px]">操作</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              <TableRow v-for="row in articles" :key="row.id">
                <TableCell>{{ row.id }}</TableCell>
                <TableCell>{{ row.title }}</TableCell>
                <TableCell>
                  <Badge v-if="row.category" variant="secondary">{{ row.category.name }}</Badge>
                </TableCell>
                <TableCell>
                  <Badge :variant="row.status === 1 ? 'secondary' : 'outline'">
                    {{ row.status === 1 ? '已发布' : '草稿' }}
                  </Badge>
                </TableCell>
                <TableCell class="text-center">
                  <Badge :variant="row.is_recommend ? 'accent' : 'outline'">
                    {{ row.is_recommend ? '★ 推荐' : '-' }}
                  </Badge>
                </TableCell>
                <TableCell>{{ row.view_count }}</TableCell>
                <TableCell>{{ formatDate(row.created_at) }}</TableCell>
                <TableCell>
                  <div class="flex gap-2">
                    <Button size="sm" variant="outline" @click="$router.push(`/articles/${row.id}`)">查看</Button>
                    <Button size="sm" @click="$router.push(`/articles/${row.id}/edit`)">编辑</Button>
                    <Button
                      size="sm"
                      :variant="row.is_recommend ? 'outline' : 'secondary'"
                      @click="handleToggleRecommend(row)"
                    >
                      {{ row.is_recommend ? '取消推荐' : '推荐' }}
                    </Button>
                    <Button size="sm" variant="destructive" @click="handleDelete(row)">删除</Button>
                  </div>
                </TableCell>
              </TableRow>
            </TableBody>
          </Table>
        </div>

        <div class="pagination">
          <div class="flex items-center justify-between">
            <span class="text-sm text-muted-foreground">共 {{ total }} 条</span>
            <Pagination>
              <PaginationContent>
                <PaginationItem>
                  <PaginationPrevious :disabled="currentPage <= 1" @click="currentPage > 1 && (currentPage--; loadArticles())" />
                </PaginationItem>
                <PaginationItem v-for="(page, idx) in visiblePages" :key="idx">
                  <PaginationLink v-if="page !== '...'" :is-active="page === currentPage" @click="currentPage = page; loadArticles()">
                    {{ page }}
                  </PaginationLink>
                  <span v-else class="inline-flex items-center justify-center h-9 w-9 text-sm">...</span>
                </PaginationItem>
                <PaginationItem>
                  <PaginationNext :disabled="currentPage >= totalPages" @click="currentPage < totalPages && (currentPage++; loadArticles())" />
                </PaginationItem>
              </PaginationContent>
            </Pagination>
          </div>
        </div>
      </CardContent>
    </Card>

    <Dialog v-model:open="confirmDialogVisible">
      <DialogContent class="sm:max-w-[400px]">
        <DialogHeader>
          <DialogTitle>{{ confirmDialogConfig.title }}</DialogTitle>
        </DialogHeader>
        <p class="text-sm text-muted-foreground">{{ confirmDialogConfig.message }}</p>
        <DialogFooter>
          <Button variant="outline" @click="handleConfirmCancel">取消</Button>
          <Button variant="destructive" @click="handleConfirmOk">确定</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { toast } from 'vue-sonner'
import { Plus } from 'lucide-vue-next'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '@/components/ui/select'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Badge } from '@/components/ui/badge'
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationPrevious, PaginationNext } from '@/components/ui/pagination'
import adminApi from '@/utils/adminApi'
import dayjs from 'dayjs'

const articles = ref([])
const categories = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const visiblePages = computed(() => {
  const pages = []
  const t = totalPages.value
  if (t === 0) return pages
  const c = currentPage.value
  if (t <= 7) {
    for (let i = 1; i <= t; i++) pages.push(i)
    return pages
  }
  pages.push(1)
  if (c > 3) pages.push('...')
  const start = Math.max(2, c - 1)
  const end = Math.min(t - 1, c + 1)
  for (let i = start; i <= end; i++) pages.push(i)
  if (c < t - 2) pages.push('...')
  pages.push(t)
  return pages
})

const filters = ref({
  keyword: '',
  categoryId: null,
  status: null,
  sortBy: '',
  sortOrder: 'desc'
})

const confirmDialogVisible = ref(false)
const confirmDialogConfig = reactive({
  title: '提示',
  message: '',
  onConfirm: null,
  onCancel: null,
})

const confirmDialog = (message, title = '提示') => {
  return new Promise((resolve, reject) => {
    confirmDialogConfig.title = title
    confirmDialogConfig.message = message
    confirmDialogConfig.onConfirm = resolve
    confirmDialogConfig.onCancel = () => reject('cancel')
    confirmDialogVisible.value = true
  })
}

const handleConfirmOk = () => {
  confirmDialogVisible.value = false
  confirmDialogConfig.onConfirm?.()
}

const handleConfirmCancel = () => {
  confirmDialogVisible.value = false
  confirmDialogConfig.onCancel?.()
}

const formatDate = (date) => dayjs(date).format('YYYY-MM-DD HH:mm')

const loadArticles = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    if (filters.value.keyword) {
      params.keyword = filters.value.keyword
    }
    if (filters.value.categoryId) {
      params.category_id = filters.value.categoryId
    }
    if (filters.value.status !== null && filters.value.status !== undefined && filters.value.status !== '') {
      params.status = filters.value.status
    }
    if (filters.value.sortBy) {
      params.sort_by = filters.value.sortBy
      params.sort_order = filters.value.sortOrder || 'desc'
    }

    const response = await adminApi.get('/admin/articles', {
      params
    })
    articles.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    toast.error('加载失败')
  }
}

const loadCategories = async () => {
  try {
    const res = await adminApi.get('/admin/categories', {
      params: { page: 1, page_size: 100 }
    })
    categories.value = res.data?.list || []
  } catch (error) {
    console.error('加载分类失败', error)
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadArticles()
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadArticles()
}

const handleSortChange = () => {
  currentPage.value = 1
  loadArticles()
}

const resetFilters = () => {
  filters.value = {
    keyword: '',
    categoryId: null,
    status: null,
    sortBy: '',
    sortOrder: 'desc'
  }
  currentPage.value = 1
  loadArticles()
}

const handleToggleRecommend = async (row) => {
  try {
    await adminApi.put(`/admin/articles/${row.id}/recommend`, {
      is_recommend: !row.is_recommend
    })
    toast.success(row.is_recommend ? '已取消推荐' : '设置推荐成功')
    loadArticles()
  } catch (error) {
    toast.error('操作失败')
  }
}

const handleDelete = async (row) => {
  try {
    await confirmDialog('确定要删除这篇文章吗？', '提示')
    await adminApi.delete(`/admin/articles/${row.id}`)
    toast.success('删除成功')
    loadArticles()
  } catch (error) {
    if (error !== 'cancel') {
      toast.error('删除失败')
    }
  }
}

onMounted(() => {
  loadCategories()
  loadArticles()
})
</script>

<style scoped>
.articles {
  min-height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--spacing-lg);
}

.page-header h2 {
  font-size: var(--font-size-2xl);
  color: var(--color-text-primary);
  line-height: var(--line-height-tight);
}

.filter-bar {
  margin-bottom: var(--spacing-md);
}

.pagination {
  margin-top: var(--spacing-lg);
  display: flex;
  justify-content: flex-end;
}
</style>
