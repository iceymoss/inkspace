<template>
  <div class="user-search-page">
    <div class="container">
      <div class="search-card">
        <h1>搜索用户</h1>
        <div class="search-bar">
          <el-input
            v-model="keyword"
            placeholder="请输入用户昵称或账号"
            clearable
            size="large"
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button :icon="Search" :loading="loading" @click="handleSearch">
                搜索
              </el-button>
            </template>
          </el-input>
        </div>
        <p class="search-tip">按昵称或账号搜索，将优先显示最佳匹配的用户</p>
      </div>

      <div v-if="topResults.length || userResults.length" class="results-card">
        <template v-if="topResults.length">
          <h2>最佳匹配</h2>
          <div class="user-list">
            <el-card
              v-for="user in topResults"
              :key="user.id"
              class="user-card"
              shadow="hover"
              @click="goToUser(user.id)"
            >
              <div class="user-card-header">
                <el-avatar :size="40" :src="user.avatar" />
                <div class="user-info">
                  <div class="user-name-line">
                    <span class="nickname">{{ user.nickname || user.username }}</span>
                    <el-tag
                      size="small"
                      type="success"
                      class="exact-tag"
                    >
                      最佳匹配
                    </el-tag>
                  </div>
                  <div class="username">@{{ user.username }}</div>
                  <p v-if="user.bio" class="bio">{{ user.bio }}</p>
                  <div class="user-stats">
                    <span>文章 {{ user.article_count }}</span>
                    <span>作品 {{ user.work_count }}</span>
                    <span>粉丝 {{ user.follower_count }}</span>
                  </div>
                </div>
              </div>
            </el-card>
          </div>
        </template>

        <template v-if="userResults.length">
          <h2 class="section-title">相关用户</h2>
          <div class="user-list">
            <el-card
              v-for="user in userResults"
              :key="user.id"
              class="user-card"
              shadow="hover"
              @click="goToUser(user.id)"
            >
              <div class="user-card-header">
                <el-avatar :size="40" :src="user.avatar" />
                <div class="user-info">
                  <div class="user-name-line">
                    <span class="nickname">{{ user.nickname || user.username }}</span>
                  </div>
                  <div class="username">@{{ user.username }}</div>
                  <p v-if="user.bio" class="bio">{{ user.bio }}</p>
                  <div class="user-stats">
                    <span>文章 {{ user.article_count }}</span>
                    <span>作品 {{ user.work_count }}</span>
                    <span>粉丝 {{ user.follower_count }}</span>
                  </div>
                </div>
              </div>
            </el-card>
          </div>
        </template>
      </div>

      <el-empty
        v-else-if="!loading && hasSearched"
        description="未找到匹配的用户"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()

const keyword = ref('')
const loading = ref(false)
const topResults = ref([])
const userResults = ref([])
const hasSearched = ref(false)

const handleSearch = async () => {
  const kw = keyword.value.trim()
  if (!kw) {
    ElMessage.warning('请输入要搜索的昵称或账号')
    return
  }

  loading.value = true
  hasSearched.value = true

  try {
    const response = await api.get('/users/search', {
      params: { keyword: kw, limit: 20 }
    })

    const payload = response.data || {}
    topResults.value = Array.isArray(payload.top) ? payload.top : []
    userResults.value = Array.isArray(payload.users) ? payload.users : []

    if (!topResults.value.length && !userResults.value.length) {
      topResults.value = []
      userResults.value = []
      return
    }
  } catch (error) {
    console.error('User search failed:', error)
    ElMessage.error('搜索用户失败')
  } finally {
    loading.value = false
  }
}

const goToUser = (id) => {
  router.push(`/users/${id}`)
}

onMounted(() => {
  const initialKeyword = route.query.keyword
  if (typeof initialKeyword === 'string' && initialKeyword.trim()) {
    keyword.value = initialKeyword.trim()
    handleSearch()
  }
})
</script>

<style scoped>
.user-search-page {
  padding: var(--spacing-md) 0 var(--spacing-xl);
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.search-card {
  background: var(--theme-bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-lg) var(--spacing-lg) var(--spacing-md);
  box-shadow: var(--shadow-sm);
  margin-bottom: var(--spacing-md);
}

.search-card h1 {
  margin: 0 0 var(--spacing-sm);
  font-size: var(--font-size-2xl);
  color: var(--theme-text-primary);
}

.search-bar {
  margin-bottom: var(--spacing-sm);
}

.search-tip {
  margin: 0;
  font-size: var(--font-size-sm);
  color: var(--theme-text-tertiary);
}

.results-card {
  background: var(--theme-bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-md);
  box-shadow: var(--shadow-sm);
}

.results-card h2 {
  margin: 0 0 var(--spacing-md);
  font-size: var(--font-size-xl);
  color: var(--theme-text-primary);
}

.user-list {
  display: grid;
  gap: var(--spacing-sm);
}

.user-card {
  margin-bottom: var(--spacing-sm);
  cursor: pointer;
  transition: transform var(--transition-base), box-shadow var(--transition-base);
}

.user-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.user-card-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name-line {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
}

.nickname {
  font-weight: 600;
  font-size: var(--font-size-base);
  color: var(--theme-text-primary);
}

.exact-tag {
  font-size: var(--font-size-xs);
}

.username {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
}

.bio {
  font-size: var(--font-size-sm);
  color: var(--theme-text-secondary);
  margin: 0 0 var(--spacing-sm);
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  line-height: var(--line-height-base);
}

.user-stats {
  display: flex;
  gap: var(--spacing-sm);
  font-size: var(--font-size-xs);
  color: var(--theme-text-tertiary);
}

.user-stats span {
  white-space: nowrap;
}

.section-title {
  margin-top: var(--spacing-md);
}
</style>


