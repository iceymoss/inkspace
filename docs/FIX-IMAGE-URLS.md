# 修复图片URL问题

## 🐛 问题描述

数据库中存储的图片URL包含硬编码的 `http://localhost:8081`，导致前端访问图片时出现404错误。

**错误示例：**
```
GET http://localhost:8081/uploads/avatars/xxx.jpg 404 (Not Found)
```

---

## ✅ 解决方案

### 1. 修复前端代码（已完成）

以下文件已修复，不再硬编码 `localhost:8081`：

- ✅ `frontend/blog/src/views/ProfileEdit.vue`
- ✅ `frontend/blog/src/components/ImageCropUpload.vue`
- ✅ `frontend/blog/src/components/VditorEditor.vue`
- ✅ `frontend/blog/vite.config.js` - 添加 `/uploads` 代理

### 2. 配置 Vite 代理（已完成）

**`frontend/blog/vite.config.js`：**
```javascript
server: {
  port: 3001,
  proxy: {
    '/api': {
      target: 'http://localhost:8081',
      changeOrigin: true
    },
    '/uploads': {
      target: 'http://localhost:8081',  // 静态文件代理
      changeOrigin: true
    }
  }
}
```

### 3. 修复数据库中的旧URL

**执行SQL脚本：**

```bash
cd /home/jeff/icey/open-source/inkspace
mysql -h localhost -u root -proot mysite < scripts/fix_image_urls.sql
```

**或手动执行SQL：**

```sql
USE mysite;

-- 修复用户头像URL
UPDATE users 
SET avatar = REPLACE(avatar, 'http://localhost:8081', '')
WHERE avatar LIKE 'http://localhost:8081%';

-- 修复作品封面URL
UPDATE works 
SET cover = REPLACE(cover, 'http://localhost:8081', '')
WHERE cover LIKE 'http://localhost:8081%';

-- 修复文章封面URL
UPDATE articles 
SET cover = REPLACE(cover, 'http://localhost:8081', '')
WHERE cover LIKE 'http://localhost:8081%';

-- 修复分类Logo URL
UPDATE categories 
SET logo = REPLACE(logo, 'http://localhost:8081', '')
WHERE logo LIKE 'http://localhost:8081%';

-- 查看修复结果
SELECT 'Users' as table_name, COUNT(*) as fixed_count 
FROM users 
WHERE avatar LIKE '/uploads%'
UNION ALL
SELECT 'Works', COUNT(*) 
FROM works 
WHERE cover LIKE '/uploads%'
UNION ALL
SELECT 'Articles', COUNT(*) 
FROM articles 
WHERE cover LIKE '/uploads%'
UNION ALL
SELECT 'Categories', COUNT(*) 
FROM categories 
WHERE logo LIKE '/uploads%';
```

---

## 🔍 验证修复

### 1. 检查数据库

```sql
-- 检查是否还有硬编码的URL
SELECT 'Users' as table_name, username, avatar 
FROM users 
WHERE avatar LIKE 'http://localhost:8081%'
UNION ALL
SELECT 'Works', title, cover 
FROM works 
WHERE cover LIKE 'http://localhost:8081%'
UNION ALL
SELECT 'Articles', title, cover 
FROM articles 
WHERE cover LIKE 'http://localhost:8081%';
```

### 2. 检查浏览器控制台

访问页面后，打开浏览器控制台（F12），应该没有 404 错误。

### 3. 检查图片显示

- ✅ 用户头像正常显示
- ✅ 作品封面正常显示
- ✅ 文章封面正常显示
- ✅ 分类Logo正常显示

---

## 📝 正确的URL格式

### 后端返回格式

```json
{
  "url": "/uploads/images/2025/12/05/xxx.jpg"
}
```

### 前端使用

```javascript
// ✅ 正确：直接使用相对路径
form.avatar = response.data.url

// ❌ 错误：不要拼接域名
form.avatar = `http://localhost:8081${response.data.url}`
```

### 数据库存储格式

```
/uploads/images/2025/12/05/xxx.jpg
/uploads/avatars/xxx.jpg
/uploads/photos/2025/12/05/xxx.jpg
```

---

## 🚀 工作原理

1. **后端**（8081端口）：
   - 提供 `/uploads` 静态文件服务
   - 返回相对路径（如 `/uploads/xxx.jpg`）

2. **前端**（3001端口）：
   - Vite 代理 `/uploads` 请求到后端
   - 浏览器请求 `http://localhost:3001/uploads/xxx.jpg`
   - Vite 转发到 `http://localhost:8081/uploads/xxx.jpg`

3. **数据库**：
   - 只存储相对路径
   - 不依赖具体域名和端口

---

## 🎯 优势

1. **灵活性**：不依赖固定域名和端口
2. **部署友好**：生产环境可以使用不同的域名
3. **开发便利**：开发和生产环境使用相同的路径格式
4. **CDN支持**：未来可以轻松切换到CDN

---

## 📌 注意事项

1. **新上传的图片**：自动使用正确格式（相对路径）
2. **旧数据**：需要执行SQL脚本修复
3. **Vite代理**：开发环境必须配置 `/uploads` 代理
4. **生产环境**：Nginx需要配置静态文件服务

---

## 🔧 生产环境配置

**Nginx配置示例：**

```nginx
server {
    listen 80;
    server_name example.com;

    # 前端静态文件
    location / {
        root /var/www/blog;
        try_files $uri $uri/ /index.html;
    }

    # API代理
    location /api {
        proxy_pass http://backend:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 静态文件代理
    location /uploads {
        proxy_pass http://backend:8081;
        proxy_set_header Host $host;
    }
}
```

---

## ✅ 完成清单

- [x] 修复前端上传组件
- [x] 配置 Vite 代理
- [x] 创建SQL修复脚本
- [x] 更新文档

现在图片应该可以正常显示了！🎉

