# 前端启动指南

## 🚀 快速启动

### 在WSL终端中执行

```bash
# 进入前端目录
cd /home/jeff/icey/open-source/mysite/frontend

# 安装依赖（首次）
npm install

# 启动开发服务器
npm run dev
```

**访问**：http://localhost:3000

---

## 📦 如果没有安装Node.js

### 在WSL中安装Node.js

```bash
# 安装nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash

# 重启终端，然后安装Node.js
nvm install 18
nvm use 18

# 验证安装
node -v
npm -v
```

---

## 🔧 可用命令

```bash
npm install        # 安装依赖
npm run dev        # 开发模式
npm run build      # 构建生产版本
npm run preview    # 预览构建结果
```

---

## 🌐 访问地址

- 开发环境：http://localhost:3000
- 后端API：http://localhost:8080（需先启动后端）

---

## ✅ 启动成功标志

终端显示：
```
VITE v5.x.x  ready in xxx ms

➜  Local:   http://localhost:3000/
➜  Network: use --host to expose
```

