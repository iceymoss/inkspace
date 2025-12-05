# 🚀 摄影相册系统 - 部署指南

## 📋 部署前检查

- [x] 所有代码已提交
- [x] 数据库脚本已准备
- [x] 文档已完善
- [x] 功能已测试

---

## 🗄️ 数据库迁移

### 步骤1：修复图片URL

```bash
cd /home/jeff/icey/open-source/inkspace
mysql -h localhost -u root -proot mysite < scripts/fix_image_urls.sql
```

**作用：** 将数据库中硬编码的 `http://localhost:8081` 替换为相对路径 `/uploads/...`

### 步骤2：创建点赞表

```bash
mysql -h localhost -u root -proot mysite < scripts/create_likes_table.sql
```

**作用：** 创建 `likes` 表，为 `works` 表添加 `like_count` 和 `favorite_count` 字段

### 步骤3：扩展收藏表

```bash
mysql -h localhost -u root -proot mysite < scripts/extend_favorites_works.sql
```

**作用：** 为 `favorites` 表添加 `work_id` 字段，支持作品收藏

---

## 🔧 后端部署

### 步骤1：编译检查

```bash
cd /home/jeff/icey/open-source/inkspace
go build ./...
```

**预期结果：** 无编译错误

### 步骤2：重启服务

```bash
# 用户服务（8081端口）
make dev

# 管理服务（8083端口）
make dev-admin

# 定时任务调度器
make dev-scheduler
```

---

## 🎨 前端部署

前端使用 Vite 开发服务器，会自动热更新，无需重启。

**访问地址：**
- 博客前端：http://127.0.0.1:3001
- 管理前端：http://127.0.0.1:3002

---

## ✅ 功能测试

### 1. 基础功能测试

```bash
# 测试作品列表
curl http://127.0.0.1:3001/api/works

# 测试作品详情
curl http://127.0.0.1:3001/api/works/1

# 测试配额查询（需登录）
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://127.0.0.1:3001/api/works/quota
```

### 2. 浏览器测试

#### 作品列表
1. 访问：http://127.0.0.1:3001/works
2. 检查瀑布流布局
3. 检查类型筛选
4. 点击作品进入详情

#### 作品详情
1. 访问：http://127.0.0.1:3001/works/1
2. 测试图片轮播
3. 点击缩略图切换
4. 查看摄影参数
5. 测试点赞按钮
6. 测试收藏按钮
7. 发表评论

#### 用户中心
1. 登录后访问：http://127.0.0.1:3001/dashboard/works
2. 点击"创建作品"
3. 选择"摄影作品"
4. 添加照片和参数
5. 发布作品
6. 编辑作品
7. 删除作品

---

## 🐛 常见问题

### 问题1：图片404错误

**症状：**
```
GET http://localhost:8081/uploads/xxx.jpg 404
```

**解决方案：**
```bash
# 执行图片URL修复脚本
mysql -h localhost -u root -proot mysite < scripts/fix_image_urls.sql
```

### 问题2：编译错误

**症状：**
```
cannot use images (variable of type []string) as []PhotoItem
```

**解决方案：**
```bash
# 确保所有文件已更新
git pull
go build ./...
```

### 问题3：点赞/收藏不生效

**症状：** 点击按钮无反应或报错

**解决方案：**
```bash
# 1. 检查数据库表是否创建
mysql -h localhost -u root -proot mysite -e "SHOW TABLES LIKE 'likes';"
mysql -h localhost -u root -proot mysite -e "DESC favorites;"

# 2. 如果表不存在，执行迁移脚本
mysql -h localhost -u root -proot mysite < scripts/create_likes_table.sql
mysql -h localhost -u root -proot mysite < scripts/extend_favorites_works.sql
```

### 问题4：配额限制不生效

**症状：** 可以无限创建摄影作品

**解决方案：**
检查 `work_service.go` 中的 `CheckDailyQuota` 方法是否正确调用。

---

## 📊 性能优化建议

### 1. 图片优化
- 使用CDN加速图片访问
- 启用浏览器缓存
- 考虑使用WebP格式

### 2. 数据库优化
- 为常用查询字段添加索引
- 定期清理软删除数据
- 使用Redis缓存热门作品

### 3. 前端优化
- 启用代码分割
- 使用懒加载
- 压缩静态资源

---

## 🔐 安全建议

### 1. 文件上传
- 限制文件大小
- 验证文件类型
- 防止路径遍历攻击

### 2. API安全
- 使用JWT认证
- 实现速率限制
- 防止SQL注入

### 3. 数据保护
- 定期备份数据库
- 加密敏感信息
- 实现访问日志

---

## 📈 监控建议

### 1. 应用监控
- 监控API响应时间
- 跟踪错误率
- 监控内存使用

### 2. 数据库监控
- 监控慢查询
- 跟踪连接数
- 监控磁盘使用

### 3. 用户行为
- 跟踪活跃用户
- 分析热门作品
- 监控上传量

---

## 🎯 部署检查清单

### 数据库
- [ ] 执行 `fix_image_urls.sql`
- [ ] 执行 `create_likes_table.sql`
- [ ] 执行 `extend_favorites_works.sql`
- [ ] 验证表结构正确

### 后端
- [ ] 代码编译成功
- [ ] 用户服务启动（8081）
- [ ] 管理服务启动（8083）
- [ ] 调度器启动
- [ ] API测试通过

### 前端
- [ ] 博客前端访问正常（3001）
- [ ] 管理前端访问正常（3002）
- [ ] 图片显示正常
- [ ] 点赞功能正常
- [ ] 收藏功能正常

### 功能测试
- [ ] 作品列表显示正常
- [ ] 作品详情显示正常
- [ ] 图片轮播正常
- [ ] 评论功能正常
- [ ] 点赞功能正常
- [ ] 收藏功能正常
- [ ] 用户中心正常
- [ ] 作品创建正常
- [ ] 配额限制生效

---

## 🎉 部署完成

所有步骤完成后，系统应该可以正常运行。

**访问地址：**
- 博客前端：http://127.0.0.1:3001
- 管理后台：http://127.0.0.1:3002
- 用户API：http://127.0.0.1:8081
- 管理API：http://127.0.0.1:8083

**测试账号：**
- 管理员：admin / admin123
- 普通用户：iceymoss / 123456

---

## 📞 技术支持

如遇到问题，请查看：
1. `docs/FINAL-SUMMARY.md` - 完整功能总结
2. `docs/PHOTOGRAPHY-ALBUM-SYSTEM.md` - 系统文档
3. `docs/LIKE-FAVORITE-NOTIFICATION-SYSTEM.md` - 点赞收藏文档
4. `CHANGELOG.md` - 更新日志

祝部署顺利！🚀

