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
  padding: 20px 0 40px;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.search-card {
  background: var(--theme-bg-card);
  border-radius: 10px;
  padding: 24px 24px 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  margin-bottom: 20px;
}

.search-card h1 {
  margin: 0 0 12px;
  font-size: 22px;
}

.search-bar {
  margin-bottom: 8px;
}

.search-tip {
  margin: 0;
  font-size: 13px;
  color: var(--theme-text-secondary);
}

.results-card {
  background: var(--theme-bg-card);
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.results-card h2 {
  margin: 0 0 16px;
  font-size: 18px;
}

.user-card {
  margin-bottom: 12px;
  cursor: pointer;
}

.user-card:hover {
  transform: translateY(-2px);
  transition: all 0.2s ease;
}

.user-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 8px;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-name-line {
  display: flex;
  align-items: center;
  gap: 8px;
}

.nickname {
  font-weight: 600;
  font-size: 15px;
}

.exact-tag {
  font-size: 11px;
}

.username {
  font-size: 13px;
  color: var(--theme-text-secondary);
}

.bio {
  font-size: 13px;
  color: var(--theme-text-secondary);
  margin: 0 0 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.user-stats {
  display: flex;
  gap: 10px;
  font-size: 12px;
  color: var(--theme-text-secondary);
}

.user-stats span {
  white-space: nowrap;
}
</style>


