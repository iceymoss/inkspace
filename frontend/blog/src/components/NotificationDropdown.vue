<template>
  <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
    <el-dropdown @command="handleCommand" trigger="click">
      <el-button text circle>
        <el-icon :size="20"><Bell /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu class="notification-dropdown">
          <div class="dropdown-header">
            <span>通知</span>
            <el-link type="primary" @click="$router.push('/dashboard/notifications')">
              查看全部
            </el-link>
          </div>
          
          <div class="notifications-preview">
            <div 
              v-for="notification in recentNotifications" 
              :key="notification.id"
              class="notification-item"
              :class="{ unread: !notification.is_read }"
              @click="handleNotificationClick(notification)"
            >
              <el-avatar :size="32" :src="notification.from_user?.avatar" />
              <div class="notification-content">
                <div class="notification-text">
                  <span class="from-user">
                    {{ notification.from_user?.nickname || notification.from_user?.username }}
                  </span>
                  {{ notification.content }}
                </div>
                <div class="notification-time">
                  {{ formatDate(notification.created_at) }}
                </div>
              </div>
            </div>
            
            <el-empty 
              v-if="recentNotifications.length === 0" 
              description="暂无通知" 
              :image-size="60"
            />
          </div>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </el-badge>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from '@element-plus/icons-vue'
import api from '@/utils/api'

const router = useRouter()

const unreadCount = ref(0)
const recentNotifications = ref([])

const loadUnreadCount = async () => {
  try {
    const response = await api.get('/notifications/unread-count')
    unreadCount.value = response.data.count || 0
  } catch (error) {
    console.error('Failed to load unread count:', error)
  }
}

const loadRecentNotifications = async () => {
  try {
    const response = await api.get('/notifications', {
      params: {
        page: 1,
        page_size: 5
      }
    })
    recentNotifications.value = response.data.list || []
  } catch (error) {
    console.error('Failed to load notifications:', error)
  }
}

const handleNotificationClick = async (notification) => {
  // 标记为已读
  if (!notification.is_read) {
    try {
      await api.put(`/notifications/${notification.id}/read`)
      loadUnreadCount()
      loadRecentNotifications()
    } catch (error) {
      console.error('Failed to mark as read:', error)
    }
  }

  // 跳转到相关内容
  if (notification.type === 'follow' && notification.from_user_id) {
    // 关注通知：跳转到关注者的个人主页
    router.push(`/users/${notification.from_user_id}`)
  } else if (notification.article_id) {
    router.push(`/blog/${notification.article_id}`)
  } else if (notification.work_id) {
    router.push(`/works/${notification.work_id}`)
  }
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

// 每30秒刷新一次
setInterval(() => {
  loadUnreadCount()
  loadRecentNotifications()
}, 30000)

onMounted(() => {
  loadUnreadCount()
  loadRecentNotifications()
})
</script>

<style scoped>
.notification-badge {
  margin-right: 15px;
}

.notification-dropdown {
  width: 360px;
  max-height: 500px;
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  border-bottom: 1px solid #ebeef5;
  font-weight: 600;
}

.notifications-preview {
  max-height: 400px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: 15px 20px;
  border-bottom: 1px solid #f5f5f5;
  cursor: pointer;
  transition: background-color 0.3s;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-item.unread {
  background-color: #ecf5ff;
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-text {
  font-size: 0.9rem;
  line-height: 1.5;
  margin-bottom: 5px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.from-user {
  font-weight: 600;
  color: #409eff;
}

.notification-time {
  font-size: 0.8rem;
  color: #909399;
}
</style>
