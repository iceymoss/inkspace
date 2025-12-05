# 多类型作品系统 - 部署指南

## 🚀 快速部署

### 步骤 1：运行数据库迁移

```bash
# 添加新字段
mysql -h localhost -u root -proot mysite < scripts/add_work_types.sql
```

### 步骤 2：重启后端服务

```bash
# 停止旧服务（Ctrl+C）

# 重新启动
make dev              # 用户服务
make dev-admin        # 管理服务
make dev-scheduler    # 定时任务
```

### 步骤 3：前端自动热更新

前端会自动重新加载，无需手动操作。

---

## 📊 数据库变更

### Works 表新增字段

```sql
type         VARCHAR(50) NOT NULL DEFAULT 'project'
metadata     TEXT NULL
daily_quota  BOOLEAN DEFAULT FALSE
```

### 索引

```sql
idx_type          ON works(type)
idx_author_date   ON works(author_id, created_at)
```

---

## 🎨 支持的作品类型

### 1. 开源项目 (project)

**特点：**
- 无每日限制
- 支持项目链接、GitHub、演示地址
- 显示技术栈标签

**字段：**
- `link`: 项目主页
- `github_url`: GitHub 仓库
- `demo_url`: 在线演示
- `tech_stack`: 技术栈（逗号分隔）

**示例数据：**
```json
{
  "type": "project",
  "title": "个人博客系统",
  "description": "基于 Go + Vue 的多用户博客",
  "tech_stack": "Go, Gin, Vue, MySQL, Redis",
  "github_url": "https://github.com/...",
  "demo_url": "https://demo.example.com"
}
```

### 2. 摄影作品 (photography)

**特点：**
- ⚠️ 每天最多发布 3 张
- 图片保留原图质量（不压缩）
- 支持EXIF参数展示

**Metadata 字段：**
```json
{
  "camera": "Canon EOS R5",
  "lens": "RF 24-70mm f/2.8",
  "focal_length": "50mm",
  "aperture": "f/2.8",
  "shutter_speed": "1/200s",
  "iso": "400",
  "location": "杭州西湖",
  "shooting_date": "2025-12-05"
}
```

**示例数据：**
```json
{
  "type": "photography",
  "title": "日落西湖",
  "description": "美丽的日落时分，湖面波光粼粼",
  "daily_quota": true,
  "metadata": {
    "camera": "Canon EOS R5",
    "lens": "RF 24-70mm f/2.8",
    "focal_length": "50mm",
    "aperture": "f/2.8",
    "shutter_speed": "1/200s",
    "iso": "400",
    "location": "杭州西湖",
    "shooting_date": "2025-12-05"
  }
}
```

---

## 🔌 API 接口

### 查询今日配额

```http
GET /api/works/quota
Authorization: Bearer <token>

响应：
{
  "code": 0,
  "data": {
    "used": 2,      // 今日已发布
    "limit": 3,     // 每日限制
    "remaining": 1  // 剩余配额
  }
}
```

### 创建作品

```http
POST /api/works
Authorization: Bearer <token>
Content-Type: application/json

// 开源项目
{
  "title": "项目名称",
  "type": "project",
  "description": "项目描述",
  "cover": "/uploads/...",
  "tech_stack": "Go, Vue",
  "github_url": "https://github.com/...",
  "status": 1
}

// 摄影作品
{
  "title": "作品名称",
  "type": "photography",
  "description": "作品描述",
  "cover": "/uploads/photos/...",
  "metadata": {
    "camera": "Canon EOS R5",
    "lens": "RF 24-70mm f/2.8",
    "focal_length": "50mm",
    "aperture": "f/2.8"
  },
  "status": 1
}
```

### 上传摄影作品原图

```http
POST /api/upload/photo
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <原图文件>

限制：
- 最大 20MB
- 格式：JPG, PNG
- 保留原图质量
```

---

## 🎨 前端展示

### 作品列表页

```
┌──────────────────────────┐
│ [💻 项目]    │ 右上角标签 │
│ ┌──────────────────────┐ │
│ │   封面图              │ │
│ └──────────────────────┘ │
│ 项目标题                 │
│ 项目描述...              │
│ 👁 150  💬 5             │
└──────────────────────────┘

┌──────────────────────────┐
│ [📷 摄影]    │ 右上角标签 │
│ ┌──────────────────────┐ │
│ │   高清照片            │ │
│ └──────────────────────┘ │
│ 作品标题                 │
│ 作品描述...              │
│ 👁 200  💬 10            │
└──────────────────────────┘
```

### 作品详情页

**开源项目：**
```
作品标题
[作者信息] [浏览量] [时间]

作品描述...

┌─ 项目信息 ─────────┐
│ 技术栈: Go | Vue   │
│ GitHub: ...        │
│ 在线演示: ...      │
└───────────────────┘

[图片展示]
[访问项目] 按钮
```

**摄影作品：**
```
作品标题
[作者信息] [浏览量] [时间]

作品描述...

┌─ 📷 摄影参数 ───────┐
│ 相机: Canon EOS R5  │
│ 镜头: RF 24-70mm    │
│ 焦段: 50mm          │
│ 光圈: f/2.8         │
│ 快门: 1/200s        │
│ ISO: 400            │
│ 拍摄地点: 杭州西湖   │
│ 拍摄日期: 2025-12-05│
└────────────────────┘

[高清图片展示]
```

---

## 🔧 管理后台

### 作品列表

新增"类型"列：
- 💻 项目（蓝色标签）
- 📷 摄影（黄色标签）

### 创建/编辑表单

**类型选择：**
- 单选框：开源项目 | 摄影作品

**动态表单：**
- 选择"开源项目" → 显示：链接、GitHub、演示地址、技术栈
- 选择"摄影作品" → 显示：相机参数表单 + 配额提示

**摄影作品提示：**
```
ℹ️ 摄影作品说明
• 每天最多发布3张摄影作品
• 图片将保留原图质量，不会压缩
• 建议上传高质量JPG或PNG格式
```

---

## 🛡️ 配额控制

### 检查逻辑

```go
// 后端自动检查
if workType == "photography" {
    // 查询今日已发布数量
    todayCount := GetTodayPhotographyCount(userID)
    
    if todayCount >= 3 {
        return Error("今日摄影作品发布数量已达上限（3张/天）")
    }
}
```

### 前端提示

创建时如果超出配额：
```
❌ 今日摄影作品发布数量已达上限（3张/天）
请明天再来发布更多精彩作品！
```

---

## 📁 文件存储

### 目录结构

```
uploads/
├── images/          # 普通图片（可能压缩）
│   └── 2025/01/02/
├── avatars/         # 用户头像
├── photos/          # 摄影作品原图（不压缩）⭐
│   └── 2025/01/02/
```

### 上传接口对比

| 接口 | 用途 | 大小限制 | 压缩 |
|-----|------|---------|-----|
| `/upload/image` | 普通图片 | 5MB | 可能 |
| `/upload/avatar` | 用户头像 | 2MB | 是 |
| `/upload/photo` | 摄影原图 | 20MB | ❌ 否 |

---

## 🎯 测试清单

### 1. 创建开源项目

- [ ] 管理后台 → 作品管理 → 新建
- [ ] 选择"开源项目"
- [ ] 填写项目信息（技术栈、GitHub等）
- [ ] 保存
- [ ] 前台查看是否显示正确

### 2. 创建摄影作品

- [ ] 管理后台 → 作品管理 → 新建
- [ ] 选择"摄影作品"
- [ ] 上传高质量照片（使用 /upload/photo）
- [ ] 填写摄影参数
- [ ] 保存
- [ ] 前台查看是否显示摄影参数

### 3. 配额测试

- [ ] 连续创建3张摄影作品
- [ ] 尝试创建第4张
- [ ] 应该提示"今日配额已用完"
- [ ] 等到第二天（或修改数据库日期测试）
- [ ] 应该可以继续发布

### 4. 前端展示

- [ ] 主页 - 作品卡片右上角显示类型标签
- [ ] 作品列表 - 显示类型标签和评论数
- [ ] 作品详情 - 根据类型显示不同信息
  - 项目：技术栈、链接等
  - 摄影：摄影参数表

### 5. 热门排名

- [ ] 摄影作品的评论会计入热度
- [ ] 热门作品列表同时包含项目和摄影

---

## 🔍 故障排查

### 问题：无法创建摄影作品

**检查：**
1. 数据库是否有 `type`, `metadata` 字段
2. 配额是否已用完
3. 图片是否超过20MB

### 问题：摄影参数不显示

**检查：**
1. metadata 是否正确保存
2. 前端是否正确解析 JSON
3. 浏览器控制台是否有错误

### 问题：配额限制不生效

**检查：**
1. 后端服务是否重启
2. 数据库 `created_at` 字段是否正确
3. 时区设置是否正确

---

## 📈 数据示例

### 数据库中的数据

```sql
-- 开源项目
INSERT INTO works (title, type, description, tech_stack, github_url) VALUES
('个人博客系统', 'project', '多用户博客', 'Go,Vue,MySQL', 'https://github.com/...');

-- 摄影作品
INSERT INTO works (title, type, description, metadata, daily_quota) VALUES
('日落西湖', 'photography', '美丽的日落', 
'{"camera":"Canon EOS R5","lens":"RF 24-70mm","aperture":"f/2.8"}', 
true);
```

---

## 🎉 功能总结

### ✅ 已实现功能

1. **多类型支持**
   - 开源项目
   - 摄影作品
   - 易于扩展更多类型

2. **摄影作品特性**
   - 每日配额限制（3张/天）
   - 原图质量保留（不压缩）
   - 完整的EXIF参数展示
   - 专用上传接口

3. **前端展示**
   - 类型标签（项目/摄影）
   - 动态表单（根据类型显示不同字段）
   - 参数卡片（摄影作品详情页）

4. **管理后台**
   - 类型选择器
   - 动态表单
   - 配额提示
   - 类型列显示

5. **API 支持**
   - 配额查询
   - 类型筛选
   - 元数据存储

---

## 🔄 未来扩展

轻松添加新类型，例如视频作品：

```go
// 1. 定义元数据结构
type VideoMetadata struct {
    Duration   string `json:"duration"`
    Resolution string `json:"resolution"`
    Codec      string `json:"codec"`
}

// 2. 前端添加类型选项
<el-radio label="video">视频作品</el-radio>

// 3. 添加对应表单
<template v-if="form.type === 'video'">
  <el-form-item label="时长">
    <el-input v-model="videoMetadata.duration" />
  </el-form-item>
</template>

// 4. 配置配额限制（如需要）
if req.Type == "video" {
    // 检查配额逻辑
}
```

---

**部署完成后，即可使用多类型作品系统！** 🎉

