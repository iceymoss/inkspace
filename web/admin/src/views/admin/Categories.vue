<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold">分类管理</h2>
      <Button @click="showDialog()"><Plus class="h-4 w-4 mr-1" /> 新建分类</Button>
    </div>

    <Table>
      <TableHeader>
        <TableRow>
          <TableHead class="w-[80px]">ID</TableHead>
          <TableHead class="w-[100px]">Logo</TableHead>
          <TableHead class="w-[150px]">名称</TableHead>
          <TableHead class="w-[150px]">别名</TableHead>
          <TableHead>描述</TableHead>
          <TableHead class="w-[100px]">文章数</TableHead>
          <TableHead class="w-[150px]">操作</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        <TableRow v-for="row in categories" :key="row.id">
          <TableCell>{{ row.id }}</TableCell>
          <TableCell>
            <img v-if="row.logo" :src="row.logo" class="w-[50px] h-[50px] rounded object-cover" />
          </TableCell>
          <TableCell>{{ row.name }}</TableCell>
          <TableCell>{{ row.slug }}</TableCell>
          <TableCell>{{ row.description }}</TableCell>
          <TableCell>{{ row.article_count }}</TableCell>
          <TableCell>
            <div class="flex gap-2">
              <Button size="sm" variant="outline" @click="showDialog(row)">
                <Pencil class="h-3 w-3 mr-1" /> 编辑
              </Button>
              <Button size="sm" variant="destructive" @click="handleDelete(row)">
                <Trash2 class="h-3 w-3 mr-1" /> 删除
              </Button>
            </div>
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>

    <EmptyState v-if="!categories.length" title="暂无分类" description="点击上方按钮创建第一个分类" />

    <div class="flex items-center justify-between">
      <span class="text-sm text-muted-foreground">共 {{ total }} 条</span>
      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious :disabled="currentPage <= 1" @click="currentPage > 1 && handlePageChange(currentPage - 1)" />
          </PaginationItem>
          <PaginationItem v-for="(page, idx) in visiblePages" :key="idx">
            <PaginationLink v-if="page !== '...'" :is-active="page === currentPage" @click="handlePageChange(page)">
              {{ page }}
            </PaginationLink>
            <span v-else class="inline-flex items-center justify-center h-9 w-9 text-sm">...</span>
          </PaginationItem>
          <PaginationItem>
            <PaginationNext :disabled="currentPage >= totalPages" @click="currentPage < totalPages && handlePageChange(currentPage + 1)" />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>

    <Dialog v-model:open="dialogVisible">
      <DialogContent class="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>{{ isEdit ? '编辑分类' : '新建分类' }}</DialogTitle>
        </DialogHeader>
        <form class="space-y-4">
          <div class="space-y-2">
            <label class="text-sm font-medium">分类名称 <span class="text-destructive">*</span></label>
            <Input v-model="form.name" placeholder="请输入分类名称" />
            <p v-if="formErrors.name" class="text-sm text-destructive">{{ formErrors.name }}</p>
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">别名</label>
            <Input v-model="form.slug" placeholder="URL友好的别名" />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">分类Logo <span class="text-destructive">*</span></label>
            <ImageCropUpload
              v-model="form.logo"
              preview-size="120px"
              placeholder="点击上传Logo"
              tip="上传正方形Logo，系统会自动裁剪"
            />
            <p v-if="formErrors.logo" class="text-sm text-destructive">{{ formErrors.logo }}</p>
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">描述</label>
            <Textarea v-model="form.description" placeholder="分类描述" />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium">排序</label>
            <Input
              :model-value="form.sort"
              @update:model-value="val => form.sort = val === '' ? 0 : Number(val)"
              type="number"
              :min="0"
              class="w-[200px]"
            />
          </div>
        </form>
        <DialogFooter>
          <Button variant="outline" @click="dialogVisible = false">取消</Button>
          <Button @click="handleSubmit" :disabled="loading">
            <Loader2 v-if="loading" class="h-4 w-4 mr-1 animate-spin" /> 保存
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

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
import { Plus, Pencil, Trash2, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationPrevious, PaginationNext } from '@/components/ui/pagination'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { EmptyState } from '@/components/ui/empty-state'
import adminApi from '@/utils/adminApi'
import { useAdminStore } from '@/stores/admin'
import ImageCropUpload from '@/components/ImageCropUpload.vue'

const adminStore = useAdminStore()
const categories = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dialogVisible = ref(false)
const formRef = ref()
const loading = ref(false)
const editingId = ref(null)

const isEdit = computed(() => !!editingId.value)

const form = reactive({
  name: '',
  slug: '',
  description: '',
  logo: '',
  sort: 0
})

const formErrors = reactive({
  name: '',
  logo: ''
})

const rules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }],
  logo: [{ required: true, message: '请上传分类Logo', trigger: 'change' }]
}

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

const validateForm = () => {
  let valid = true
  formErrors.name = ''
  formErrors.logo = ''

  for (const rule of rules.name) {
    if (rule.required && !form.name) {
      formErrors.name = rule.message
      valid = false
      break
    }
  }
  for (const rule of rules.logo) {
    if (rule.required && !form.logo) {
      formErrors.logo = rule.message
      valid = false
      break
    }
  }
  return valid
}

const loadCategories = async () => {
  try {
    const response = await adminApi.get('/admin/categories', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value
      }
    })
    categories.value = response.data?.list || []
    total.value = response.data?.total || 0
    if (response.data?.page_size) {
      pageSize.value = response.data.page_size
    }
  } catch (error) {
    toast.error('加载失败')
  }
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadCategories()
}

const showDialog = (category = null) => {
  formErrors.name = ''
  formErrors.logo = ''
  if (category) {
    editingId.value = category.id
    Object.assign(form, {
      name: category.name,
      slug: category.slug,
      description: category.description,
      logo: category.logo,
      sort: category.sort || 0
    })
  } else {
    editingId.value = null
    Object.assign(form, { 
      name: '', 
      slug: '', 
      description: '', 
      logo: '',
      sort: 0
    })
  }
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!validateForm()) return

  loading.value = true
  try {
    if (isEdit.value) {
      await adminApi.put(`/admin/categories/${editingId.value}`, form)
      toast.success('更新成功')
    } else {
      await adminApi.post('/admin/categories', form)
      toast.success('创建成功')
    }
    dialogVisible.value = false
    loadCategories()
  } catch (error) {
    toast.error('保存失败')
  } finally {
    loading.value = false
  }
}

const handleDelete = async (category) => {
  try {
    await confirmDialog('确定要删除这个分类吗？', '提示')
    await adminApi.delete(`/admin/categories/${category.id}`)
    toast.success('删除成功')
    loadCategories()
  } catch (error) {
    if (error !== 'cancel') {
      toast.error('删除失败')
    }
  }
}

onMounted(() => {
  loadCategories()
})
</script>
