<template>
  <div class="user-search-page">
    <div class="container">
      <div class="search-card">
        <div class="page-kicker">
          PEOPLE · INDEX
        </div>
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
              <el-button
                :icon="Search"
                :loading="loading"
                @click="handleSearch"
              >
                搜索
              </el-button>
            </template>
          </el-input>
        </div>
        <p class="search-tip">
          按昵称或账号搜索，将优先显示最佳匹配的用户
        </p>
      </div>

      <div
        v-if="topResults.length || userResults.length"
        class="results-card"
      >
        <template v-if="topResults.length">
          <h2>最佳匹配</h2>
          <div class="user-list">
            <el-card
              v-for="user in topResults"
              :key="user.id"
              class="user-card"
              shadow="hover"
              tabindex="0"
              @click="goToUser(user.id)"
              @keyup.enter="goToUser(user.id)"
            >
              <div class="user-card-header">
                <el-avatar
                  :size="40"
                  :src="user.avatar"
                />
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
                  <div class="username">
                    @{{ user.username }}
                  </div>
                  <p
                    v-if="user.bio"
                    class="bio"
                  >
                    {{ user.bio }}
                  </p>
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
          <h2 class="section-title">
            相关用户
          </h2>
          <div class="user-list">
            <el-card
              v-for="user in userResults"
              :key="user.id"
              class="user-card"
              shadow="hover"
              tabindex="0"
              @click="goToUser(user.id)"
              @keyup.enter="goToUser(user.id)"
            >
              <div class="user-card-header">
                <el-avatar
                  :size="40"
                  :src="user.avatar"
                />
                <div class="user-info">
                  <div class="user-name-line">
                    <span class="nickname">{{ user.nickname || user.username }}</span>
                  </div>
                  <div class="username">
                    @{{ user.username }}
                  </div>
                  <p
                    v-if="user.bio"
                    class="bio"
                  >
                    {{ user.bio }}
                  </p>
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
import { ref, watch } from 'vue'
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
let searchRequestId = 0

const searchUsers = async () => {
  const activeRequest = ++searchRequestId
  const kw = keyword.value.trim()
  loading.value = true
  hasSearched.value = true

  try {
    const response = await api.get('/users/search', {
      params: { keyword: kw, limit: 20 }
    })
    if (activeRequest !== searchRequestId) return

    const payload = response.data || {}
    topResults.value = Array.isArray(payload.top) ? payload.top : []
    userResults.value = Array.isArray(payload.users) ? payload.users : []

    if (!topResults.value.length && !userResults.value.length) {
      topResults.value = []
      userResults.value = []
      return
    }
  } catch (error) {
    if (activeRequest !== searchRequestId) return
    console.error('User search failed:', error)
    ElMessage.error('搜索用户失败')
  } finally {
    if (activeRequest === searchRequestId) loading.value = false
  }
}

const handleSearch = () => {
  const kw = keyword.value.trim()
  if (!kw) {
    ElMessage.warning('请输入要搜索的昵称或账号')
    return
  }

  const query = { ...route.query, keyword: kw }
  const target = { path: route.path, query }
  if (router.resolve(target).fullPath === route.fullPath) {
    searchUsers()
    return
  }
  router.push(target)
}

const goToUser = (id) => {
  router.push(`/users/${id}`)
}

watch(() => route.query.keyword, value => {
  const nextKeyword = typeof value === 'string' ? value.trim() : ''
  keyword.value = nextKeyword
  if (nextKeyword) {
    searchUsers()
  } else {
    topResults.value = []
    userResults.value = []
    hasSearched.value = false
  }
}, { immediate: true })
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

/* Magazine adaptation */
.user-search-page { padding: 62px 0 80px; background: var(--theme-bg-primary); }
.user-search-page .container { max-width: 860px; padding: 0 32px; }
.search-card, .results-card { padding: 0 0 32px; background: transparent; border-bottom: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; }
.search-card { margin-bottom: 42px; }
.page-kicker { margin-bottom: 10px; color: var(--theme-primary); font-family: Georgia, 'Songti SC', serif; font-size: 11px; letter-spacing: .26em; }
.search-card h1 { margin-bottom: 24px; font-family: Georgia, 'Songti SC', 'Noto Serif SC', SimSun, serif; font-size: clamp(38px, 6vw, 54px); font-weight: 500; letter-spacing: .04em; }
.search-bar :deep(.el-input__wrapper), .search-bar :deep(.el-input-group__append), .user-search-page :deep(.el-tag) { border-radius: 1px; box-shadow: none; }
.results-card h2 { padding-bottom: 12px; border-bottom: 1px solid var(--theme-border); font-family: Georgia, 'Songti SC', serif; font-size: 24px; font-weight: 500; letter-spacing: .05em; }
.section-title { margin-top: 44px !important; }
.user-card { margin: 0; border: 0; border-bottom: 1px solid var(--theme-border); border-radius: 0; box-shadow: none; transition: padding-left .25s ease; }
.user-card:hover { padding-left: 8px; box-shadow: none; transform: none; }
.user-card:focus-visible { outline: 2px solid var(--theme-primary); outline-offset: 3px; }
.nickname { font-family: Georgia, 'Songti SC', serif; font-size: 18px; font-weight: 500; }
@media (max-width: 900px) { .user-search-page .container { padding: 0 24px; } }
@media (max-width: 560px) { .user-search-page { padding: 38px 0 56px; } .user-search-page .container { padding: 0 18px; } .search-card h1 { font-size: 36px; } .user-stats { flex-wrap: wrap; } }
</style>


