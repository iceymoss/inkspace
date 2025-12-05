# 🔧 通知功能修复方案

## 🐛 发现的问题

从日志中发现了3个关键错误：

### 1. notifications 表有多余的 title 字段
```
Error 1364 (HY000): Field 'title' doesn't have a default value
```

**原因：** 数据库中的 `notifications` 表有一个 `title` 字段，但代码中的 Notification 模型没有这个字段。

### 2. works 表缺少 like_count 字段
```
Error 1054 (42S22): Unknown column 'like_count' in 'field list'
```

### 3. works 表缺少 favorite_count 字段
```
Error 1054 (42S22): Unknown column 'favorite_count' in 'field list'
```

---

## ✅ 修复步骤

### 执行修复脚本

```bash
cd /home/jeff/icey/open-source/inkspace

# 执行修复脚本
mysql -h localhost -u root -proot mysite < scripts/fix_notifications_and_works.sql
```

这个脚本会：
1. ✅ 删除 notifications 表的 title 字段
2. ✅ 为 works 表添加 like_count 字段
3. ✅ 为 works 表添加 favorite_count 字段
4. ✅ 测试通知插入
5. ✅ 验证修复结果

---

## 🧪 修复后测试

### 1. 不需要重启服务
修复后可以直接测试，无需重启。

### 2. 测试步骤

#### 测试评论通知
1. 用户A（id=2）登录
2. 访问用户B（id=1）的作品：http://127.0.0.1:3001/works/2
3. 发表评论
4. **后端日志应该看到：**
   ```
   ✅ 成功创建作品评论通知: 用户2 -> 用户1, 作品2
   ```
5. 用户B（id=1）登录
6. 访问通知中心：http://127.0.0.1:3001/dashboard/notifications
7. **应该看到通知！**

#### 测试点赞通知
1. 用户A登录
2. 访问用户B的作品
3. 点击"点赞"按钮
4. **后端日志应该看到：**
   ```
   ✅ 成功创建作品点赞通知: 用户2 -> 用户1, 作品2
   ```
5. 用户B查看通知

#### 测试收藏通知
1. 用户A登录
2. 访问用户B的作品
3. 点击"收藏"按钮
4. **后端日志应该看到：**
   ```
   ✅ 成功创建作品收藏通知: 用户2 -> 用户1, 作品2
   ```
5. 用户B查看通知

---

## 📊 预期日志

### 修复前（错误日志）
```
❌ 创建作品评论通知失败: 用户2 -> 用户1, 作品2, 错误: Error 1364: Field 'title' doesn't have a default value
```

### 修复后（成功日志）
```
✅ 成功创建作品评论通知: 用户2 -> 用户1, 作品2
✅ 成功创建作品点赞通知: 用户2 -> 用户1, 作品2
✅ 成功创建作品收藏通知: 用户2 -> 用户1, 作品2
```

---

## 🎯 验证清单

执行修复脚本后：

- [ ] notifications 表没有 title 字段
- [ ] works 表有 like_count 字段
- [ ] works 表有 favorite_count 字段
- [ ] 可以成功插入通知
- [ ] 评论/点赞/收藏后能看到成功日志
- [ ] 用户B能在通知中心看到通知
- [ ] 导航栏显示未读数量

---

## 💡 为什么会有这个问题？

可能的原因：
1. 数据库表是从旧版本迁移过来的，有额外字段
2. 之前手动创建表时添加了 title 字段
3. 使用了不同的建表脚本

---

## 🚀 快速修复

```bash
# 一键修复
cd /home/jeff/icey/open-source/inkspace
mysql -h localhost -u root -proot mysite < scripts/fix_notifications_and_works.sql

# 然后重新测试
```

修复后，所有通知功能应该正常工作！

