<template>
  <div class="notifications-page">
    <div class="container">
      <div class="page-header">
        <h1>我的通知</h1>
        <div class="actions">
          <el-button @click="markAllAsRead" v-if="unreadCount > 0">
            <el-icon><Check /></el-icon> 全部已读
          </el-button>
        </div>
      </div>

      <el-tabs v-model="activeTab" @tab-change="onTabChange">
        <el-tab-pane label="全部" name="all" />
        <el-tab-pane label="未读" name="unread" />
        <el-tab-pane label="已读" name="read" />
      </el-tabs>

      <div class="notification-list">
        <el-card 
          v-for="notification in notifications" 
          :key="notification.id"
          class="notification-card"
          :class="{ 'is-unread': !notification.is_read }"
          @click="handleClick(notification)"
        >
          <div class="notification-content">
            <el-avatar :src="notification.from_user?.avatar" :size="50" />
            <div class="notification-body">
              <div class="notification-header">
                <span class="notification-type">{{ getTypeLabel(notification.type) }}</span>
                <el-tag v-if="!notification.is_read" type="primary" size="small">未读</el-tag>
              </div>
              <h4>{{ notification.title }}</h4>
              <p>{{ notification.content }}</p>
              <div class="notification-footer">
                <span class="notification-time">{{ formatTime(notification.created_at) }}</span>
                <el-button 
                  text 
                  type="danger" 
                  size="small" 
                  @click.stop="handleDelete(notification)"
                >
                  删除
                </el-button>
              </div>
            </div>
          </div>
        </el-card>
      </div>

      <el-empty v-if="notifications.length === 0" description="暂无通知" />

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadNotifications"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Check } from '@element-plus/icons-vue'
import api from '@/utils/api'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()

const notifications = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const unreadCount = ref(0)
const activeTab = ref('all')

const formatTime = (time) => dayjs(time).fromNow()

const getTypeLabel = (type) => {
  const labels = {
    comment: '评论',
    reply: '回复',
    like: '点赞',
    system: '系统',
    mention: '提及'
  }
  return labels[type] || type
}

const loadNotifications = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (activeTab.value === 'unread') {
      params.is_read = false
    } else if (activeTab.value === 'read') {
      params.is_read = true
    }

    const response = await api.get('/notifications', { params })
    notifications.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const loadUnreadCount = async () => {
  try {
    const response = await api.get('/notifications/unread-count')
    unreadCount.value = response.data.count || 0
  } catch (error) {
    console.error('Failed to load unread count:', error)
  }
}

const markAllAsRead = async () => {
  try {
    await api.put('/notifications/read-all')
    ElMessage.success('已全部标记为已读')
    loadNotifications()
    loadUnreadCount()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleClick = async (notification) => {
  if (!notification.is_read) {
    try {
      await api.put(`/notifications/${notification.id}/read`)
      notification.is_read = true
      loadUnreadCount()
    } catch (error) {
      console.error('Failed to mark as read:', error)
    }
  }

  if (notification.link) {
    router.push(notification.link)
  }
}

const handleDelete = async (notification) => {
  try {
    await ElMessageBox.confirm('确定要删除这条通知吗？', '提示', { type: 'warning' })
    await api.delete(`/notifications/${notification.id}`)
    ElMessage.success('删除成功')
    loadNotifications()
    loadUnreadCount()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const onTabChange = () => {
  currentPage.value = 1
  loadNotifications()
}

onMounted(() => {
  loadNotifications()
  loadUnreadCount()
})
</script>

<style scoped>
.notifications-page {
  padding: 40px 0;
  background-color: #f5f7fa;
  min-height: 100vh;
}

.notification-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.notification-list {
  margin-top: 20px;
}

.notification-card {
  margin-bottom: 15px;
  cursor: pointer;
  transition: transform 0.2s;
}

.notification-card:hover {
  transform: translateX(5px);
}

.notification-card.is-unread {
  border-left: 3px solid var(--el-color-primary);
  background-color: #f0f9ff;
}

.notification-content {
  display: flex;
  gap: 15px;
}

.notification-body {
  flex: 1;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.notification-type {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.notification-body h4 {
  margin: 0 0 8px 0;
  font-size: 16px;
}

.notification-body p {
  margin: 0 0 10px 0;
  color: var(--el-text-color-regular);
}

.notification-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notification-time {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
}

.pagination {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}
</style>

