<template>
  <div class="relative mr-4">
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Button variant="ghost" size="icon" class="rounded-full">
          <Bell class="h-5 w-5" />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" class="w-[360px] p-0">
        <div class="dropdown-header">
          <span>通知</span>
          <button class="text-primary hover:underline text-sm cursor-pointer" @click="$router.push('/dashboard/notifications')">
            查看全部
          </button>
        </div>

        <div class="notifications-preview">
          <div
            v-for="notification in recentNotifications"
            :key="notification.id"
            class="notification-item"
            :class="{ unread: !notification.is_read }"
            @click="handleNotificationClick(notification)"
          >
            <Avatar class="h-8 w-8">
              <AvatarImage :src="notification.from_user?.avatar" />
              <AvatarFallback />
            </Avatar>
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

          <EmptyState
            v-if="recentNotifications.length === 0"
            description="暂无通知"
            class="py-6"
          />
        </div>
      </DropdownMenuContent>
    </DropdownMenu>
    <Badge v-if="unreadCount > 0" class="absolute -top-1 -right-1 min-w-[20px] h-5 flex items-center justify-center px-1">
      {{ unreadCount }}
    </Badge>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bell } from 'lucide-vue-next'
import { DropdownMenu, DropdownMenuTrigger, DropdownMenuContent } from '@/components/ui/dropdown-menu'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { EmptyState } from '@/components/ui/empty-state'
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
  if (!notification.is_read) {
    try {
      await api.put(`/notifications/${notification.id}/read`)
      loadUnreadCount()
      loadRecentNotifications()
    } catch (error) {
      console.error('Failed to mark as read:', error)
    }
  }

  if (notification.type === 'follow' && notification.from_user_id) {
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
.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--theme-border-light);
  font-weight: 600;
}

.notifications-preview {
  max-height: 400px;
  overflow-y: auto;
}

.notification-item {
  display: flex;
  gap: 12px;
  padding: var(--spacing-md) var(--spacing-lg);
  border-bottom: 1px solid var(--theme-border-light);
  cursor: pointer;
  transition: background-color var(--transition-slow);
}

.notification-item:hover {
  background-color: var(--theme-bg-hover);
}

.notification-item.unread {
  background-color: var(--theme-bg-secondary);
}

.notification-item:last-child {
  border-bottom: none;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-text {
  font-size: var(--font-size-sm);
  line-height: var(--line-height-base);
  margin-bottom: var(--spacing-xs);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.from-user {
  font-weight: 600;
  color: var(--theme-primary);
}

.notification-time {
  font-size: var(--font-size-xs);
  color: var(--theme-text-tertiary);
}
</style>
