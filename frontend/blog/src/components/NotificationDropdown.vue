<template>
  <el-dropdown @command="handleCommand" @visible-change="onVisibleChange">
    <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
      <el-button circle>
        <el-icon><Bell /></el-icon>
      </el-button>
    </el-badge>
    
    <template #dropdown>
      <el-dropdown-menu class="notification-dropdown">
        <div class="dropdown-header">
          <span>通知</span>
          <el-button text size="small" @click="markAllAsRead" v-if="unreadCount > 0">
            全部已读
          </el-button>
        </div>
        
        <el-scrollbar max-height="400px" v-if="notifications.length > 0">
          <el-dropdown-item
            v-for="notification in notifications"
            :key="notification.id"
            :class="{ 'is-read': notification.is_read }"
            @click="handleNotificationClick(notification)"
          >
            <div class="notification-item">
              <el-avatar :size="32" :src="notification.from_user?.avatar" />
              <div class="notification-content">
                <div class="notification-title">{{ notification.title }}</div>
                <div class="notification-text">{{ notification.content }}</div>
                <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
              </div>
            </div>
          </el-dropdown-item>
        </el-scrollbar>
        
        <el-empty v-else description="暂无通知" :image-size="60" />
        
        <el-divider style="margin: 10px 0" />
        <el-dropdown-item @click="viewAll">
          <div style="text-align: center; color: var(--el-color-primary);">
            查看全部通知
          </div>
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '@/utils/api'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const router = useRouter()
const notifications = ref([])
const unreadCount = ref(0)

const formatTime = (time) => dayjs(time).fromNow()

const loadNotifications = async () => {
  try {
    const response = await api.get('/notifications', {
      params: { page: 1, page_size: 10 }
    })
    notifications.value = response.data.list || []
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

const markAllAsRead = async () => {
  try {
    await api.put('/notifications/read-all')
    unreadCount.value = 0
    notifications.value.forEach(n => n.is_read = true)
    ElMessage.success('已全部标记为已读')
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleNotificationClick = async (notification) => {
  if (!notification.is_read) {
    try {
      await api.put(`/notifications/${notification.id}/read`)
      notification.is_read = true
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    } catch (error) {
      console.error('Failed to mark as read:', error)
    }
  }
  
  if (notification.link) {
    router.push(notification.link)
  }
}

const onVisibleChange = (visible) => {
  if (visible) {
    loadNotifications()
  }
}

const viewAll = () => {
  router.push('/notifications')
}

const handleCommand = (command) => {
  if (command === 'view-all') {
    viewAll()
  }
}

onMounted(() => {
  loadUnreadCount()
  
  // 定时刷新未读数
  setInterval(() => {
    loadUnreadCount()
  }, 30000) // 每30秒刷新一次
})
</script>

<style scoped>
.notification-badge :deep(.el-badge__content) {
  font-size: 10px;
}

.notification-dropdown {
  min-width: 350px;
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 15px;
  font-weight: bold;
  border-bottom: 1px solid var(--el-border-color);
}

.notification-item {
  display: flex;
  gap: 10px;
  padding: 5px 0;
  width: 100%;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-title {
  font-weight: 500;
  margin-bottom: 4px;
}

.notification-text {
  font-size: 13px;
  color: var(--el-text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-time {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  margin-top: 4px;
}

.el-dropdown-menu__item.is-read {
  background-color: var(--el-fill-color-blank);
  opacity: 0.7;
}
</style>

