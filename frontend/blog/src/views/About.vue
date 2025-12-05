<template>
  <div class="about">
    <div class="container">
      <el-card>
        <h1>关于我</h1>
        
        <div class="about-content">
          <div class="profile">
            <el-avatar :size="120" src="" />
            <h2>Your Name</h2>
            <p class="bio">全栈开发工程师 | 技术爱好者</p>
          </div>

          <div class="intro">
            <h3>简介</h3>
            <p>
              欢迎来到我的个人网站！我是一名热爱技术的开发者，专注于Web开发和软件工程。
              这个网站用于分享我的技术文章、展示我的作品以及记录我的学习历程。
            </p>

            <h3>技能</h3>
            <div class="skills">
              <el-tag v-for="skill in skills" :key="skill" type="primary">{{ skill }}</el-tag>
            </div>

            <h3>联系方式</h3>
            <div class="contact">
              <p><el-icon><Message /></el-icon> Email: your.email@example.com</p>
              <p><el-icon><Link /></el-icon> GitHub: github.com/yourusername</p>
            </div>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const skills = ref([
  'Go', 'Gin', 'GORM', 'MySQL', 'Redis',
  'Vue.js', 'JavaScript', 'TypeScript', 'Docker'
])
</script>

<style scoped>
.about {
  padding: 40px 0;
  background-color: #f5f7fa;
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
</style>

