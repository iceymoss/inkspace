<template>
  <div class="notifications">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>消息通知</span>
          <div class="header-actions">
            <el-button text @click="markAllAsRead" v-if="unreadCount > 0">
              全部标记为已读
            </el-button>
            <el-button text type="danger" @click="deleteAllRead">
              清空已读
            </el-button>
          </div>
        </div>
      </template>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="全部" name="all">
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" />
        </el-tab-pane>
        <el-tab-pane label="未读" name="unread">
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" />
        </el-tab-pane>
      </el-tabs>

      <div class="notifications-list">
        <div 
          v-for="notification in notifications" 
          :key="notification.id"
          class="notification-item"
          :class="{ unread: !notification.is_read }"
          @click="handleNotificationClick(notification)"
        >
          <el-avatar 
            :src="notification.from_user?.avatar" 
            :size="50"
          />
          
          <div class="notification-content">
            <div class="notification-header">
              <span class="from-user">
                {{ notification.from_user?.nickname || notification.from_user?.username }}
              </span>
              <span class="notification-action">{{ notification.content }}</span>
              <el-tag 
                :type="getTypeTag(notification.type).type" 
                size="small"
                style="margin-left: 8px"
              >
                {{ getTypeTag(notification.type).label }}
              </el-tag>
            </div>
            
            <div class="notification-time">
              {{ formatDate(notification.created_at) }}
            </div>
          </div>

          <div class="notification-actions">
            <el-button 
              v-if="!notification.is_read"
              size="small"
              text
              @click.stop="markAsRead(notification.id)"
            >
              标记已读
            </el-button>
            <el-button 
              size="small"
              text
              type="danger"
              @click.stop="deleteNotification(notification.id)"
            >
              删除
            </el-button>
          </div>
        </div>

        <el-empty v-if="notifications.length === 0" description="暂无通知" />
      </div>

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadNotifications"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import api from '@/utils/api'

const router = useRouter()

const notifications = ref([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const activeTab = ref('all')
const unreadCount = ref(0)

const loadNotifications = async () => {
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (activeTab.value === 'unread') {
      params.only_unread = true
    }
    
    const response = await api.get('/notifications', { params })
    notifications.value = response.data.list || []
    total.value = response.data.total || 0
  } catch (error) {
    console.error('Failed to load notifications:', error)
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

const handleTabChange = () => {
  currentPage.value = 1
  loadNotifications()
}

const markAsRead = async (notificationId) => {
  try {
    await api.put(`/notifications/${notificationId}/read`)
    loadNotifications()
    loadUnreadCount()
  } catch (error) {
    ElMessage.error('操作失败')
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

const deleteNotification = async (notificationId) => {
  try {
    await api.delete(`/notifications/${notificationId}`)
    ElMessage.success('删除成功')
    loadNotifications()
    loadUnreadCount()
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const deleteAllRead = async () => {
  try {
    await ElMessageBox.confirm('确定要删除所有已读通知吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await api.delete('/notifications/read-all')
    ElMessage.success('删除成功')
    loadNotifications()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleNotificationClick = async (notification) => {
  // 标记为已读
  if (!notification.is_read) {
    await markAsRead(notification.id)
  }

  // 跳转到相关内容
  if (notification.article_id) {
    router.push(`/blog/${notification.article_id}`)
  } else if (notification.work_id) {
    router.push(`/works/${notification.work_id}`)
  }
}

const getTypeTag = (type) => {
  const tags = {
    comment: { label: '评论', type: 'primary' },
    like: { label: '点赞', type: 'danger' },
    favorite: { label: '收藏', type: 'warning' },
    follow: { label: '关注', type: 'success' },
    reply: { label: '回复', type: 'info' }
  }
  return tags[type] || { label: '通知', type: '' }
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

onMounted(() => {
  loadNotifications()
  loadUnreadCount()
})
</script>

<style scoped>
.notifications {
  max-width: 900px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.notifications-list {
  min-height: 400px;
}

.notification-item {
  display: flex;
  gap: 15px;
  padding: 20px;
  border-bottom: 1px solid #ebeef5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-item.unread {
  background-color: #ecf5ff;
}

.notification-item.unread:hover {
  background-color: #d9ecff;
}

.notification-content {
  flex: 1;
}

.notification-header {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.from-user {
  font-weight: 600;
  color: #303133;
  margin-right: 5px;
}

.notification-action {
  color: #606266;
}

.notification-time {
  font-size: 0.9rem;
  color: #909399;
}

.notification-actions {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>

