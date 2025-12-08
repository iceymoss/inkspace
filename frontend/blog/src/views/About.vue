<template>
  <div class="about">
    <div class="container">
      <el-card>
        <h1>{{ aboutData.title || '关于我' }}</h1>
        
        <div class="about-content">
          <div class="profile" v-if="aboutData.avatar || aboutData.name || aboutData.bio">
            <el-avatar :size="120" :src="aboutData.avatar" v-if="aboutData.avatar" />
            <h2 v-if="aboutData.name">{{ aboutData.name }}</h2>
            <p class="bio" v-if="aboutData.bio">{{ aboutData.bio }}</p>
          </div>

          <div class="intro">
            <h3 v-if="aboutData.introduction">简介</h3>
            <p v-if="aboutData.introduction" v-html="formatContent(aboutData.introduction)"></p>

            <h3 v-if="aboutData.skills && aboutData.skills.length > 0">技能</h3>
            <div class="skills" v-if="aboutData.skills && aboutData.skills.length > 0">
              <el-tag v-for="skill in aboutData.skills" :key="skill" type="primary">{{ skill }}</el-tag>
            </div>

            <h3 v-if="hasContact">联系方式</h3>
            <div class="contact" v-if="hasContact">
              <p v-if="aboutData.email">
                <el-icon><Message /></el-icon> 
                <a :href="`mailto:${aboutData.email}`">{{ aboutData.email }}</a>
              </p>
              <p v-if="aboutData.github">
                <el-icon><Link /></el-icon> 
                <a :href="formatGithubUrl(aboutData.github)" target="_blank" rel="noopener noreferrer">{{ aboutData.github }}</a>
              </p>
              <p v-if="aboutData.wechat">
                <el-icon><ChatLineRound /></el-icon> 
                微信: {{ aboutData.wechat }}
              </p>
              <p v-if="aboutData.qq">
                <el-icon><ChatLineRound /></el-icon> 
                QQ: {{ aboutData.qq }}
              </p>
              <p v-if="aboutData.weibo">
                <el-icon><Link /></el-icon> 
                <a :href="formatWeiboUrl(aboutData.weibo)" target="_blank" rel="noopener noreferrer">{{ aboutData.weibo }}</a>
              </p>
            </div>
            
            <div v-if="!aboutData.avatar && !aboutData.name && !aboutData.bio && !aboutData.introduction && (!aboutData.skills || aboutData.skills.length === 0) && !hasContact" class="empty-tip">
              <el-empty description="暂无内容，请在管理后台配置关于页面信息" />
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Message, Link, ChatLineRound } from '@element-plus/icons-vue'
import api from '@/utils/api'

const aboutData = ref({
  title: '关于我',
  avatar: '',
  name: '',
  bio: '',
  introduction: '',
  skills: [],
  email: '',
  github: '',
  wechat: '',
  qq: '',
  weibo: ''
})

const hasContact = computed(() => {
  return aboutData.value.email || aboutData.value.github || 
         aboutData.value.wechat || aboutData.value.qq || aboutData.value.weibo
})

const formatContent = (content) => {
  if (!content) return ''
  // 将换行符转换为<br>
  return content.replace(/\n/g, '<br>')
}

const formatGithubUrl = (github) => {
  if (!github) return ''
  if (github.startsWith('http')) {
    return github
  }
  return `https://github.com/${github.replace(/^@?/, '')}`
}

const formatWeiboUrl = (weibo) => {
  if (!weibo) return ''
  if (weibo.startsWith('http')) {
    return weibo
  }
  return `https://weibo.com/${weibo}`
}

const loadAboutData = async () => {
  try {
    const response = await api.get('/settings/public')
    // API返回格式: { code: 0, data: { about_page: "...", ... } }
    const settings = response.data || {}
    const aboutJson = settings.about_page
    
    if (aboutJson) {
      try {
        const data = JSON.parse(aboutJson)
        aboutData.value = {
          title: data.title || '关于我',
          avatar: data.avatar || '',
          name: data.name || '',
          bio: data.bio || '',
          introduction: data.introduction || '',
          skills: data.skills || [],
          email: data.email || '',
          github: data.github || '',
          wechat: data.wechat || '',
          qq: data.qq || '',
          weibo: data.weibo || ''
        }
      } catch (e) {
        console.error('Failed to parse about page data:', e)
      }
    }
  } catch (error) {
    console.error('Failed to load about page data:', error)
  }
}

onMounted(() => {
  loadAboutData()
})
</script>

<style scoped>
.about {
  padding: 40px 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.about .el-card {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.about h1 {
  text-align: center;
  margin-bottom: 40px;
}

.about-content {
  max-width: 800px;
  margin: 0 auto;
}

.profile {
  text-align: center;
  margin-bottom: 40px;
  padding-bottom: 40px;
  border-bottom: 1px solid var(--border-lighter);
}

.profile h2 {
  margin: 20px 0 10px;
}

.bio {
  color: var(--text-secondary);
  font-size: 16px;
}

.intro h3 {
  margin: 30px 0 15px;
  color: var(--primary-color);
}

.intro p {
  line-height: 1.8;
  margin-bottom: 20px;
}

.skills {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 20px;
}

.contact p {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 10px;
}

.contact a {
  color: var(--primary-color);
  text-decoration: none;
}

.contact a:hover {
  text-decoration: underline;
}
</style>

