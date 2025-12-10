<template>
  <div class="settings">
    <h2>系统配置</h2>

    <el-tabs v-model="activeGroup">
      <el-tab-pane label="首页轮播" name="carousel">
        <el-form :model="carouselSettings" label-width="150px">
          <el-form-item>
            <el-button type="primary" @click="addCarouselItem">添加轮播项</el-button>
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              配置首页轮播图（广告位/公告栏），支持多个轮播项，每个轮播项可以设置标题、副标题、背景、跳转链接等。点击整个轮播图区域会跳转到设置的链接
            </div>
          </el-form-item>
          
          <el-form-item v-for="(item, index) in carouselSettings.items" :key="index" :label="`轮播项 ${index + 1}`">
            <el-card style="margin-bottom: 20px;">
              <template #header>
                <div style="display: flex; justify-content: space-between; align-items: center;">
                  <span>轮播项 {{ index + 1 }}</span>
                  <el-button type="danger" size="small" @click="removeCarouselItem(index)">删除</el-button>
                </div>
              </template>
              
              <el-form-item label="标题">
                <el-input v-model="item.title" placeholder="例如：欢迎来到我的个人网站" />
              </el-form-item>
              
              <el-form-item label="副标题">
                <el-input v-model="item.subtitle" placeholder="例如：分享技术、记录生活、展示作品" />
              </el-form-item>
              
              <el-form-item label="背景图片">
                <div style="display: flex; gap: 10px; align-items: flex-start;">
                  <ImageCropUpload
                    v-model="item.backgroundImage"
                    :aspect-ratio="3.75"
                    :output-width="1200"
                    :output-height="320"
                    preview-size="300px"
                    placeholder="点击上传轮播图"
                    tip="上传图片将按照轮播图比例（1200x320）自动裁剪"
                    :max-size="5"
                  />
                  <div style="flex: 1;">
                    <el-input 
                      v-model="item.backgroundImage" 
                      placeholder="或直接输入图片URL" 
                      style="margin-bottom: 8px;"
                    />
                    <div style="font-size: 12px; color: #909399;">
                      支持上传图片或直接输入图片URL。上传的图片将按照轮播图比例（1200x320）自动裁剪
                    </div>
                  </div>
                </div>
              </el-form-item>
              
              <el-form-item label="背景渐变">
                <el-input v-model="item.backgroundGradient" placeholder="例如：linear-gradient(135deg, #667eea 0%, #764ba2 100%)" />
                <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                  如果设置了背景图片，背景渐变将作为遮罩层使用
                </div>
              </el-form-item>
              
              <el-form-item label="跳转链接">
                <el-input v-model="item.link" placeholder="点击轮播图跳转的链接地址（可选，留空则不跳转）" />
                <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                  支持相对路径（如：/blog）或绝对URL（如：https://example.com），点击整个轮播图区域会跳转到此链接
                </div>
              </el-form-item>
            </el-card>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="saveCarouselSettings" :loading="saving">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="关于页面" name="about">
        <el-form :model="aboutSettings" label-width="150px">
          <el-form-item label="页面标题">
            <el-input v-model="aboutSettings.title" placeholder="例如：关于我" />
          </el-form-item>
          
          <el-form-item label="头像">
            <div style="display: flex; gap: 10px; align-items: flex-start;">
              <ImageCropUpload
                v-model="aboutSettings.avatar"
                :aspect-ratio="1"
                :output-width="200"
                :output-height="200"
                preview-size="100px"
                placeholder="点击上传头像"
                tip="上传正方形头像，系统会自动裁剪为200x200"
                :max-size="2"
              />
              <div style="flex: 1;">
                <el-input 
                  v-model="aboutSettings.avatar" 
                  placeholder="或直接输入头像URL" 
                  style="margin-bottom: 8px;"
                />
              </div>
            </div>
          </el-form-item>
          
          <el-form-item label="姓名">
            <el-input v-model="aboutSettings.name" placeholder="例如：张三" />
          </el-form-item>
          
          <el-form-item label="个人简介">
            <el-input v-model="aboutSettings.bio" placeholder="例如：全栈开发工程师 | 技术爱好者" />
          </el-form-item>
          
          <el-form-item label="详细介绍">
            <el-input 
              v-model="aboutSettings.introduction" 
              type="textarea" 
              :rows="6"
              placeholder="输入详细介绍内容，支持换行"
            />
          </el-form-item>
          
          <el-form-item label="技能标签">
            <div style="margin-bottom: 10px;">
              <el-tag
                v-for="(skill, index) in aboutSettings.skills"
                :key="index"
                closable
                @close="removeSkill(index)"
                style="margin-right: 8px; margin-bottom: 8px;"
              >
                {{ skill }}
              </el-tag>
              <el-input
                v-if="skillInputVisible"
                ref="skillInputRef"
                v-model="skillInputValue"
                size="small"
                style="width: 120px;"
                @keyup.enter="addSkill"
                @blur="addSkill"
              />
              <el-button v-else size="small" @click="showSkillInput">+ 添加技能</el-button>
            </div>
            <div style="font-size: 12px; color: #909399;">
              点击"添加技能"按钮或按回车键添加技能标签
            </div>
          </el-form-item>
          
          <el-divider content-position="left">联系方式</el-divider>
          
          <el-form-item label="邮箱">
            <el-input v-model="aboutSettings.email" placeholder="例如：your.email@example.com" />
          </el-form-item>
          
          <el-form-item label="GitHub">
            <el-input v-model="aboutSettings.github" placeholder="例如：github.com/username 或 @username" />
          </el-form-item>
          
          <el-form-item label="微信">
            <el-input v-model="aboutSettings.wechat" placeholder="例如：your_wechat_id" />
          </el-form-item>
          
          <el-form-item label="QQ">
            <el-input v-model="aboutSettings.qq" placeholder="例如：123456789" />
          </el-form-item>
          
          <el-form-item label="微博">
            <el-input v-model="aboutSettings.weibo" placeholder="例如：weibo.com/username 或 @username" />
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" @click="saveAboutSettings" :loading="saving">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="网站信息" name="site">
        <el-form :model="siteSettings" label-width="120px">
          <el-form-item label="网站名称">
            <el-input v-model="siteSettings.site_name" />
          </el-form-item>
          <el-form-item label="网站描述">
            <el-input v-model="siteSettings.site_description" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="关键词">
            <el-input v-model="siteSettings.site_keywords" placeholder="关键词用逗号分隔" />
          </el-form-item>
          <el-form-item label="备案号">
            <el-input v-model="siteSettings.site_icp" />
          </el-form-item>
          <el-form-item label="版权信息">
            <el-input v-model="siteSettings.site_copyright" />
          </el-form-item>
          <el-form-item label="网站Logo">
            <div style="display: flex; gap: 10px; align-items: flex-start;">
              <ImageCropUpload
                v-model="siteSettings.site_logo"
                :aspect-ratio="0"
                :output-width="400"
                :output-height="400"
                preview-size="120px"
                placeholder="点击上传Logo"
                tip="上传Logo图片，支持任意比例裁切，系统会自动调整为400x400"
                :max-size="5"
              />
              <div style="flex: 1;">
                <el-input 
                  v-model="siteSettings.site_logo" 
                  placeholder="或直接输入Logo URL" 
                  style="margin-bottom: 8px;"
                />
                <div style="font-size: 12px; color: #909399;">
                  支持上传图片或直接输入图片URL。上传的图片支持任意比例裁切，系统会自动调整为400x400
                </div>
              </div>
            </div>
          </el-form-item>
          <el-divider content-position="left">友链申请</el-divider>
          <el-form-item label="申请标题">
            <el-input v-model="siteSettings.link_apply_title" placeholder="例如：申请友链" />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              友链申请区域的标题文字
            </div>
          </el-form-item>
          <el-form-item label="申请描述">
            <el-input v-model="siteSettings.link_apply_description" placeholder="例如：如果你也想交换友链，请联系我：" />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              友链申请区域的描述文字
            </div>
          </el-form-item>
          <el-form-item label="申请邮箱">
            <el-input v-model="siteSettings.link_apply_email" placeholder="例如：admin@example.com" />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              友链申请的联系邮箱，将在友链页面显示
            </div>
          </el-form-item>
          <el-divider content-position="left">管理后台</el-divider>
          <el-form-item label="管理后台URL">
            <el-input v-model="siteSettings.admin_backend_url" placeholder="例如：http://localhost:3002 或 https://admin.example.com" />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              管理后台的访问地址，将在用户中心页面的"管理后台"按钮中使用。支持相对路径（如：/admin）或绝对URL（如：http://localhost:3002）
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveSiteSettings" :loading="saving">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="功能设置" name="feature">
        <el-form :model="featureSettings" label-width="120px">
          <el-form-item label="开放注册">
            <el-switch v-model="featureSettings.register_enabled" />
          </el-form-item>
          <el-divider content-position="left">评论设置</el-divider>
          <el-form-item label="开放文章评论">
            <el-switch v-model="featureSettings.article_comment_enabled" />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              控制是否允许用户对文章进行评论
            </div>
          </el-form-item>
          <el-form-item label="开放作品评论">
            <el-switch v-model="featureSettings.work_comment_enabled" />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              控制是否允许用户对作品进行评论
            </div>
          </el-form-item>
          <el-form-item label="评论审核">
            <el-switch v-model="featureSettings.comment_audit" />
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              开启后，所有评论需要管理员审核后才能显示
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveFeatureSettings" :loading="saving">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="网站主题" name="theme">
        <el-form :model="themeSettings" label-width="150px">
          <el-form-item label="整体主题">
            <el-select v-model="themeSettings.site_theme" placeholder="选择网站整体主题" style="width: 300px">
              <el-option label="白天" value="day" />
              <el-option label="黑夜" value="night" />
              <el-option label="节假日" value="holiday" />
              <el-option label="哀悼日" value="mourning" />
            </el-select>
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              设置博客系统的整体主题风格，影响整个网站的配色方案
            </div>
          </el-form-item>
          
          <!-- 节假日主题自定义设置 -->
          <template v-if="themeSettings.site_theme === 'holiday'">
            <el-divider content-position="left">节假日主题设置</el-divider>
            
            <el-form-item label="节日类型">
              <el-select v-model="themeSettings.holiday_type" placeholder="选择节日类型" style="width: 300px">
                <el-option label="春节" value="spring_festival" />
                <el-option label="国庆节" value="national_day" />
                <el-option label="中秋节" value="mid_autumn" />
                <el-option label="元旦" value="new_year" />
                <el-option label="自定义" value="custom" />
              </el-select>
              <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                选择节日类型，系统会自动应用对应的配色方案
              </div>
            </el-form-item>

            <el-form-item label="背景主色">
              <el-color-picker v-model="themeSettings.holiday_bg_primary" />
              <el-input v-model="themeSettings.holiday_bg_primary" style="width: 200px; margin-left: 10px;" placeholder="#fff5f5" />
              <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                设置节假日主题的背景主色
              </div>
            </el-form-item>

            <el-form-item label="背景次色">
              <el-color-picker v-model="themeSettings.holiday_bg_secondary" />
              <el-input v-model="themeSettings.holiday_bg_secondary" style="width: 200px; margin-left: 10px;" placeholder="#ffe8e8" />
              <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                设置节假日主题的背景次色
              </div>
            </el-form-item>

            <el-form-item label="文字主色">
              <el-color-picker v-model="themeSettings.holiday_text_primary" />
              <el-input v-model="themeSettings.holiday_text_primary" style="width: 200px; margin-left: 10px;" placeholder="#8b1a1a" />
              <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                设置节假日主题的文字主色
              </div>
            </el-form-item>

            <el-form-item label="主色调">
              <el-color-picker v-model="themeSettings.holiday_primary" />
              <el-input v-model="themeSettings.holiday_primary" style="width: 200px; margin-left: 10px;" placeholder="#ff3333" />
              <div style="margin-top: 8px; color: #909399; font-size: 12px;">
                设置节假日主题的主色调（按钮、链接等）
              </div>
            </el-form-item>
          </template>

          <el-form-item>
            <el-button type="primary" @click="saveThemeSettings" :loading="saving">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="Markdown设置" name="markdown">
        <el-form :model="markdownSettings" label-width="150px">
          <el-form-item label="Markdown 主题风格">
            <el-select v-model="markdownSettings.markdown_theme" placeholder="选择 Markdown 主题风格" style="width: 300px">
              <el-option label="浅色主题" value="light" />
              <el-option label="深色主题" value="dark" />
            </el-select>
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              设置 Markdown 内容的整体主题风格（浅色/深色）
            </div>
          </el-form-item>
          <el-form-item label="代码高亮主题">
            <el-select v-model="markdownSettings.code_theme" placeholder="选择代码高亮主题" style="width: 300px">
              <el-option label="GitHub" value="github" />
              <el-option label="GitHub Dark" value="github-dark" />
              <el-option label="Atom One Dark" value="atom-one-dark" />
              <el-option label="Atom One Light" value="atom-one-light" />
              <el-option label="Monokai" value="monokai" />
              <el-option label="VS2015" value="vs2015" />
              <el-option label="VS" value="vs" />
              <el-option label="Xcode" value="xcode" />
              <el-option label="Dracula" value="dracula" />
              <el-option label="Nord" value="nord" />
              <el-option label="Solarized Dark" value="solarized-dark" />
              <el-option label="Solarized Light" value="solarized-light" />
              <el-option label="Tomorrow Night" value="tomorrow-night" />
              <el-option label="Tomorrow Night Blue" value="tomorrow-night-blue" />
              <el-option label="Tomorrow Night Bright" value="tomorrow-night-bright" />
              <el-option label="Tomorrow Night Eighties" value="tomorrow-night-eighties" />
              <el-option label="Default" value="default" />
            </el-select>
            <div style="margin-top: 8px; color: #909399; font-size: 12px;">
              设置代码块的高亮主题风格
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="saveMarkdownSettings" :loading="saving">保存</el-button>
          </el-form-item>
        </el-form>
      </el-tab-pane>

      <el-tab-pane label="所有配置" name="all">
        <el-button type="primary" @click="showEditDialog()" style="margin-bottom: 20px;">
          <el-icon><Plus /></el-icon> 新建配置
        </el-button>

        <el-table :data="allSettings" style="width: 100%;">
          <el-table-column prop="key" label="配置键" width="200" />
          <el-table-column prop="value" label="配置值" show-overflow-tooltip />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="group" label="分组" width="100" />
          <el-table-column label="公开" width="80">
            <template #default="{ row }">
              <el-tag :type="row.is_public ? 'success' : 'info'" size="small">
                {{ row.is_public ? '是' : '否' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" show-overflow-tooltip />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button size="small" @click="showEditDialog(row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleDelete(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="editDialogVisible" title="编辑配置" width="600px">
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="100px">
        <el-form-item label="配置键" prop="key">
          <el-input v-model="editForm.key" :disabled="isEditMode" />
        </el-form-item>
        <el-form-item label="配置值" prop="value">
          <el-input v-model="editForm.value" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="editForm.type">
            <el-option label="字符串" value="string" />
            <el-option label="整数" value="int" />
            <el-option label="布尔值" value="bool" />
            <el-option label="JSON" value="json" />
          </el-select>
        </el-form-item>
        <el-form-item label="分组">
          <el-input v-model="editForm.group" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editForm.description" />
        </el-form-item>
        <el-form-item label="公开">
          <el-switch v-model="editForm.is_public" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitEdit" :loading="editLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch, nextTick } from 'vue'
import ImageCropUpload from '@/components/ImageCropUpload.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import adminApi from '@/utils/adminApi'

const activeGroup = ref('site')
const allSettings = ref([])
const saving = ref(false)
const editDialogVisible = ref(false)
const editFormRef = ref()
const editLoading = ref(false)
const isEditMode = ref(false)
const skillInputVisible = ref(false)
const skillInputValue = ref('')
const skillInputRef = ref(null)

const siteSettings = reactive({
  site_name: '',
  site_description: '',
  site_keywords: '',
  site_icp: '',
  site_copyright: '',
  site_logo: '',
  link_apply_title: '',
  link_apply_description: '',
  link_apply_email: '',
  admin_backend_url: ''
})

const featureSettings = reactive({
  register_enabled: true,
  article_comment_enabled: true,
  work_comment_enabled: true,
  comment_audit: false
})

const themeSettings = reactive({
  site_theme: 'day', // 默认使用白天主题
  holiday_type: 'spring_festival', // 节假日类型
  holiday_bg_primary: '#fff5f5', // 节假日背景主色
  holiday_bg_secondary: '#ffe8e8', // 节假日背景次色
  holiday_text_primary: '#8b1a1a', // 节假日文字主色
  holiday_primary: '#ff3333' // 节假日主色调
})

const carouselSettings = reactive({
  items: [
    {
      title: '欢迎来到我的个人网站',
      subtitle: '分享技术、记录生活、展示作品',
      backgroundImage: '',
      backgroundGradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
      link: ''
    }
  ]
})

const aboutSettings = reactive({
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

const markdownSettings = reactive({
  markdown_theme: 'light', // 默认使用浅色主题
  code_theme: 'github' // 默认使用 github 主题
})

const editForm = reactive({
  key: '',
  value: '',
  type: 'string',
  group: 'general',
  description: '',
  is_public: false
})

const editRules = {
  key: [{ required: true, message: '请输入配置键', trigger: 'blur' }]
}

const loadAllSettings = async () => {
  try {
    const response = await adminApi.get('/admin/settings')
    allSettings.value = response.data || []
    
    // 分组加载到不同的设置对象
    allSettings.value.forEach(setting => {
      if (setting.group === 'site') {
        siteSettings[setting.key] = setting.value
      } else if (setting.group === 'feature') {
        if (setting.key === 'register_enabled' || setting.key === 'article_comment_enabled' || 
            setting.key === 'work_comment_enabled' || setting.key === 'comment_audit') {
          featureSettings[setting.key] = setting.value === '1' || setting.value === 'true'
        }
      } else if (setting.group === 'theme') {
        if (setting.key === 'site_theme') {
          themeSettings.site_theme = setting.value || 'day'
        } else if (setting.key === 'holiday_type') {
          themeSettings.holiday_type = setting.value || 'spring_festival'
        } else if (setting.key === 'holiday_bg_primary') {
          themeSettings.holiday_bg_primary = setting.value || '#fff5f5'
        } else if (setting.key === 'holiday_bg_secondary') {
          themeSettings.holiday_bg_secondary = setting.value || '#ffe8e8'
        } else if (setting.key === 'holiday_text_primary') {
          themeSettings.holiday_text_primary = setting.value || '#8b1a1a'
        } else if (setting.key === 'holiday_primary') {
          themeSettings.holiday_primary = setting.value || '#ff3333'
        } else {
          themeSettings[setting.key] = setting.value
        }
      } else if (setting.group === 'markdown') {
        markdownSettings[setting.key] = setting.value
      } else if (setting.group === 'carousel' || setting.key === 'home_carousel') {
        if (setting.key === 'home_carousel' && setting.value) {
          try {
            carouselSettings.items = JSON.parse(setting.value)
          } catch (e) {
            console.error('Failed to parse carousel data:', e)
          }
        }
      } else if (setting.group === 'about' || setting.key === 'about_page') {
        if (setting.key === 'about_page' && setting.value) {
          try {
            const data = JSON.parse(setting.value)
            Object.assign(aboutSettings, {
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
            })
          } catch (e) {
            console.error('Failed to parse about page data:', e)
          }
        }
      }
    })
  } catch (error) {
    ElMessage.error('加载失败')
  }
}

const saveSiteSettings = async () => {
  saving.value = true
  try {
    await adminApi.put('/admin/settings/batch', siteSettings)
    ElMessage.success('保存成功')
    loadAllSettings()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const saveFeatureSettings = async () => {
  saving.value = true
  try {
    // 转换布尔值为字符串
    const settings = {}
    Object.keys(featureSettings).forEach(key => {
      settings[key] = featureSettings[key] ? '1' : '0'
    })
    
    await adminApi.put('/admin/settings/batch', settings)
    ElMessage.success('保存成功')
    loadAllSettings()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const saveThemeSettings = async () => {
  saving.value = true
  try {
    const settings = {}
    Object.keys(themeSettings).forEach(key => {
      settings[key] = themeSettings[key]
    })
    
    await adminApi.put('/admin/settings/batch', settings)
    ElMessage.success('保存成功')
    loadAllSettings()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const saveMarkdownSettings = async () => {
  saving.value = true
  try {
    const settings = {}
    Object.keys(markdownSettings).forEach(key => {
      settings[key] = markdownSettings[key]
    })
    
    await adminApi.put('/admin/settings/batch', settings)
    ElMessage.success('保存成功')
    loadAllSettings()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const showEditDialog = (setting = null) => {
  if (setting) {
    isEditMode.value = true
    Object.assign(editForm, setting)
  } else {
    isEditMode.value = false
    Object.assign(editForm, {
      key: '',
      value: '',
      type: 'string',
      group: 'general',
      description: '',
      is_public: false
    })
  }
  editDialogVisible.value = true
}

const submitEdit = async () => {
  await editFormRef.value.validate(async (valid) => {
    if (!valid) return

    editLoading.value = true
    try {
      await adminApi.put('/admin/settings', editForm)
      ElMessage.success('保存成功')
      editDialogVisible.value = false
      loadAllSettings()
    } catch (error) {
      ElMessage.error('保存失败')
    } finally {
      editLoading.value = false
    }
  })
}

const handleDelete = async (setting) => {
  try {
    await ElMessageBox.confirm('确定要删除这个配置吗？', '提示', { type: 'warning' })
    await adminApi.delete(`/admin/settings/${setting.key}`)
    ElMessage.success('删除成功')
    loadAllSettings()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 节日类型预设颜色方案
const holidayPresets = {
  spring_festival: {
    holiday_bg_primary: '#fff5f5',
    holiday_bg_secondary: '#ffe8e8',
    holiday_text_primary: '#8b1a1a',
    holiday_primary: '#ff3333'
  },
  national_day: {
    holiday_bg_primary: '#fff0f0',
    holiday_bg_secondary: '#ffe0e0',
    holiday_text_primary: '#8b0000',
    holiday_primary: '#ff0000'
  },
  mid_autumn: {
    holiday_bg_primary: '#fff8e1',
    holiday_bg_secondary: '#ffecb3',
    holiday_text_primary: '#e65100',
    holiday_primary: '#ff9800'
  },
  new_year: {
    holiday_bg_primary: '#fff5f5',
    holiday_bg_secondary: '#ffe8e8',
    holiday_text_primary: '#c41e3a',
    holiday_primary: '#ff4444'
  }
}

// 监听节日类型变化，自动填充颜色
watch(() => themeSettings.holiday_type, (newType) => {
  if (newType && newType !== 'custom' && holidayPresets[newType]) {
    const preset = holidayPresets[newType]
    themeSettings.holiday_bg_primary = preset.holiday_bg_primary
    themeSettings.holiday_bg_secondary = preset.holiday_bg_secondary
    themeSettings.holiday_text_primary = preset.holiday_text_primary
    themeSettings.holiday_primary = preset.holiday_primary
  }
})

// 轮播图管理函数
const addCarouselItem = () => {
  carouselSettings.items.push({
    title: '',
    subtitle: '',
    backgroundImage: '',
    backgroundGradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    link: ''
  })
}

const removeCarouselItem = (index) => {
  carouselSettings.items.splice(index, 1)
}

const saveCarouselSettings = async () => {
  saving.value = true
  try {
    const settings = {
      home_carousel: JSON.stringify(carouselSettings.items)
    }
    await adminApi.put('/admin/settings/batch', settings)
    ElMessage.success('保存成功')
    loadAllSettings()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

// 关于页面管理函数
const showSkillInput = () => {
  skillInputVisible.value = true
  skillInputValue.value = ''
  nextTick(() => {
    skillInputRef.value?.focus()
  })
}

const addSkill = () => {
  const value = skillInputValue.value.trim()
  if (value && !aboutSettings.skills.includes(value)) {
    aboutSettings.skills.push(value)
  }
  skillInputVisible.value = false
  skillInputValue.value = ''
}

const removeSkill = (index) => {
  aboutSettings.skills.splice(index, 1)
}

const saveAboutSettings = async () => {
  saving.value = true
  try {
    const settings = {
      about_page: JSON.stringify(aboutSettings)
    }
    await adminApi.put('/admin/settings/batch', settings)
    ElMessage.success('保存成功')
    loadAllSettings()
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadAllSettings()
})
</script>

