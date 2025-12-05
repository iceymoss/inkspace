# 作品类型设计文档

## 📐 数据库设计

### Works 表结构

```sql
-- 新增字段
type         VARCHAR(50) NOT NULL DEFAULT 'project'  -- 作品类型
metadata     TEXT                                    -- 类型专属元数据(JSON)
daily_quota  BOOLEAN DEFAULT FALSE                   -- 是否受每日配额限制
```

### 完整表结构

```go
type Work struct {
    ID           uint           `json:"id"`
    CreatedAt    time.Time      `json:"created_at"`
    UpdatedAt    time.Time      `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-"`
    
    // 基础字段（所有类型通用）
    Title        string         `json:"title"`
    Description  string         `json:"description"`
    Cover        string         `json:"cover"`
    Images       string         `json:"images"` // JSON array
    
    // 类型相关
    Type         string         `json:"type"` // project, photography, video, etc.
    Metadata     string         `json:"metadata"` // JSON 存储类型专属数据
    DailyQuota   bool           `json:"daily_quota"` // 是否受每日配额限制
    
    // 项目类型字段
    Link         string         `json:"link"`
    GithubURL    string         `json:"github_url"`
    DemoURL      string         `json:"demo_url"`
    TechStack    string         `json:"tech_stack"`
    
    // 公共字段
    AuthorID     uint           `json:"author_id"`
    Author       *User          `json:"author,omitempty"`
    Sort         int            `json:"sort"`
    ViewCount    int            `json:"view_count"`
    CommentCount int            `json:"comment_count"`
    Status       int            `json:"status"`
    IsRecommend  bool           `json:"is_recommend"`
}
```

---

## 🎭 支持的作品类型

### 1. 开源项目 (project)

**字段使用：**
- `title`, `description`, `cover`, `images`
- `link`, `github_url`, `demo_url`, `tech_stack`
- `metadata`: `{}`

**示例：**
```json
{
  "type": "project",
  "title": "个人博客系统",
  "tech_stack": "Go, Vue, MySQL",
  "github_url": "https://github.com/...",
  "demo_url": "https://demo.example.com"
}
```

### 2. 摄影作品 (photography)

**字段使用：**
- `title`, `description`, `cover`, `images`（原图，不压缩）
- `daily_quota`: `true`（受每日3张限制）
- `metadata`: JSON 存储摄影参数

**Metadata 结构：**
```json
{
  "camera": "Canon EOS R5",          // 相机型号
  "lens": "RF 24-70mm f/2.8",        // 镜头
  "focal_length": "50mm",            // 焦段
  "aperture": "f/2.8",               // 光圈
  "shutter_speed": "1/200s",         // 快门速度
  "iso": "400",                      // ISO
  "location": "杭州西湖",             // 拍摄地点
  "shooting_date": "2025-12-05"      // 拍摄日期
}
```

**示例：**
```json
{
  "type": "photography",
  "title": "日落西湖",
  "description": "美丽的日落时分",
  "daily_quota": true,
  "metadata": {
    "camera": "Canon EOS R5",
    "lens": "RF 24-70mm f/2.8",
    "focal_length": "50mm",
    "aperture": "f/2.8",
    "shutter_speed": "1/200s",
    "iso": "400"
  }
}
```

### 3. 未来可扩展类型

- **视频作品 (video)**: duration, resolution, codec
- **音乐作品 (music)**: duration, genre, instrument
- **设计作品 (design)**: tools, dimensions, format
- **写作作品 (writing)**: word_count, genre

---

## 🔧 实现细节

### 后端模型

```go
// Metadata 结构（用于摄影作品）
type PhotographyMetadata struct {
    Camera       string `json:"camera"`
    Lens         string `json:"lens"`
    FocalLength  string `json:"focal_length"`
    Aperture     string `json:"aperture"`
    ShutterSpeed string `json:"shutter_speed"`
    ISO          string `json:"iso"`
    Location     string `json:"location"`
    ShootingDate string `json:"shooting_date"`
}

// Work 模型添加
Type       string `gorm:"size:50;not null;default:'project';index:idx_type" json:"type"`
Metadata   string `gorm:"type:text" json:"metadata"`
DailyQuota bool   `gorm:"default:false" json:"daily_quota"`
```

### 服务层逻辑

```go
// CheckDailyQuota 检查每日配额
func (s *WorkService) CheckDailyQuota(userID uint, workType string) (bool, error) {
    if workType != "photography" {
        return true, nil // 非摄影作品不限制
    }
    
    today := time.Now().Format("2006-01-02")
    var count int64
    
    err := database.DB.Model(&models.Work{}).
        Where("author_id = ? AND type = ? AND DATE(created_at) = ?", 
              userID, workType, today).
        Count(&count).Error
    
    if err != nil {
        return false, err
    }
    
    return count < 3, nil // 每天最多3个
}
```

### 前端表单

```vue
<!-- 摄影作品表单 -->
<el-form v-if="form.type === 'photography'">
  <el-form-item label="相机型号">
    <el-input v-model="photoMetadata.camera" />
  </el-form-item>
  <el-form-item label="镜头">
    <el-input v-model="photoMetadata.lens" />
  </el-form-item>
  <el-form-item label="焦段">
    <el-input v-model="photoMetadata.focal_length" placeholder="例如: 50mm" />
  </el-form-item>
  <el-form-item label="光圈">
    <el-input v-model="photoMetadata.aperture" placeholder="例如: f/2.8" />
  </el-form-item>
  <el-form-item label="快门速度">
    <el-input v-model="photoMetadata.shutter_speed" placeholder="例如: 1/200s" />
  </el-form-item>
  <el-form-item label="ISO">
    <el-input v-model="photoMetadata.iso" placeholder="例如: 400" />
  </el-form-item>
</el-form>
```

---

## 🎯 我的建议

**推荐使用单表 + JSON 设计**，原因：

1. ✅ **灵活性高** - 新增类型只需要定义 metadata 结构
2. ✅ **代码复用** - 所有类型共享评论、点赞等功能
3. ✅ **维护简单** - 一套代码管理所有作品
4. ✅ **扩展容易** - 未来可以轻松添加新类型

**实现步骤：**
1. 为 works 表添加 `type`, `metadata`, `daily_quota` 字段
2. 后端添加类型验证和配额检查
3. 前端根据类型渲染不同表单
4. 摄影作品上传时保留原图（跳过压缩）

---

## 🚀 要我现在实现吗？

我可以立即为你实现完整的多类型作品系统，包括：
- ✅ 数据库迁移脚本
- ✅ 后端模型和服务
- ✅ 每日配额检查
- ✅ 前端类型选择和表单
- ✅ 摄影参数展示

需要我开始实现吗？还是先解决评论功能的数据库问题？
