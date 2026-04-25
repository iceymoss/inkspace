<template>
  <div class="about">
    <div class="container">
      <Card class="card-skeuomorphic-static">
        <CardHeader class="text-center">
          <CardTitle>{{ aboutData.title || '关于我' }}</CardTitle>
        </CardHeader>
        <CardContent>
          <div class="about-content">
            <div class="profile" v-if="aboutData.avatar || aboutData.name || aboutData.bio">
              <Avatar class="h-[120px] w-[120px]" v-if="aboutData.avatar">
                <AvatarImage :src="aboutData.avatar" />
                <AvatarFallback>{{ aboutData.name?.charAt(0) || '?' }}</AvatarFallback>
              </Avatar>
              <h2 v-if="aboutData.name">{{ aboutData.name }}</h2>
              <p class="bio" v-if="aboutData.bio">{{ aboutData.bio }}</p>
            </div>

            <div class="intro">
              <h3 v-if="aboutData.introduction">简介</h3>
              <p v-if="aboutData.introduction" v-html="formatContent(aboutData.introduction)"></p>

              <h3 v-if="aboutData.skills && aboutData.skills.length > 0">技能</h3>
              <div class="skills" v-if="aboutData.skills && aboutData.skills.length > 0">
                <Badge v-for="skill in aboutData.skills" :key="skill">{{ skill }}</Badge>
              </div>

              <h3 v-if="hasContact">联系方式</h3>
              <div class="contact" v-if="hasContact">
                <p v-if="aboutData.email">
                  <Mail class="h-4 w-4 shrink-0" />
                  <a :href="`mailto:${aboutData.email}`">{{ aboutData.email }}</a>
                </p>
                <p v-if="aboutData.github">
                  <Link class="h-4 w-4 shrink-0" />
                  <a :href="formatGithubUrl(aboutData.github)" target="_blank" rel="noopener noreferrer">{{ aboutData.github }}</a>
                </p>
                <p v-if="aboutData.wechat">
                  <MessageSquare class="h-4 w-4 shrink-0" />
                  微信: {{ aboutData.wechat }}
                </p>
                <p v-if="aboutData.qq">
                  <MessageSquare class="h-4 w-4 shrink-0" />
                  QQ: {{ aboutData.qq }}
                </p>
                <p v-if="aboutData.weibo">
                  <Link class="h-4 w-4 shrink-0" />
                  <a :href="formatWeiboUrl(aboutData.weibo)" target="_blank" rel="noopener noreferrer">{{ aboutData.weibo }}</a>
                </p>
              </div>

              <div v-if="!aboutData.avatar && !aboutData.name && !aboutData.bio && !aboutData.introduction && (!aboutData.skills || aboutData.skills.length === 0) && !hasContact" class="empty-tip">
                <EmptyState title="暂无内容" description="暂无内容，请在管理后台配置关于页面信息" />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { Mail, Link, MessageSquare } from 'lucide-vue-next'
import { Card, CardHeader, CardContent, CardTitle } from '@/components/ui/card'
import { Avatar, AvatarImage, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { EmptyState } from '@/components/ui/empty-state'
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
  padding: var(--spacing-xl) 0;
  background-color: var(--theme-bg-secondary);
  min-height: 100vh;
}

.about-content {
  max-width: 800px;
  margin: 0 auto;
}

.profile {
  text-align: center;
  margin-bottom: var(--spacing-xl);
  padding-bottom: var(--spacing-xl);
  border-bottom: 1px solid var(--theme-border-light);
}

.profile h2 {
  margin: var(--spacing-lg) 0 var(--spacing-sm);
  font-size: var(--font-size-xl);
  color: var(--theme-text-primary);
}

.bio {
  color: var(--theme-text-secondary);
  font-size: var(--font-size-base);
  line-height: var(--line-height-base);
}

.intro h3 {
  margin: var(--spacing-lg) 0 var(--spacing-md);
  color: var(--theme-primary);
  font-size: var(--font-size-lg);
}

.intro p {
  line-height: var(--line-height-relaxed);
  margin-bottom: var(--spacing-lg);
  color: var(--theme-text-primary);
}

.skills {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-lg);
}

.contact p {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-sm);
  color: var(--theme-text-primary);
}

.contact a {
  color: var(--theme-primary);
  text-decoration: none;
  transition: color var(--transition-fast);
  cursor: pointer;
}

.contact a:hover {
  color: var(--theme-primary-hover);
  text-decoration: underline;
}
</style>
