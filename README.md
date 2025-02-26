# 个人记账系统

这是一个使用 Go 和 React 构建的个人记账系统。

## 功能特点

1. 账户管理
   - 创建和管理个人账户
   - 自定义账户名称和金额
   - 查看所有账户总余额

2. 收支记录
   - 记录收入和支出
   - 选择交易账户
   - 自动更新账户余额
   - 支持交易分类
   - 添加交易备注

## 后端技术栈

- Go
- Gin (Web框架)
- GORM (ORM框架)
- SQLite (数据库)

## 前端技术栈

- React
- TypeScript
- Ant Design (UI组件库)

## 安装说明

### 后端安装

1. 安装 Go (1.21 或更高版本)
2. 进入后端目录：`cd backend`
3. 安装依赖：`go mod tidy`
4. 运行服务器：`go run main.go`

### 前端安装

1. 安装 Node.js 和 npm
2. 进入前端目录：`cd frontend`
3. 安装依赖：`npm install`
4. 运行开发服务器：`npm start`

## API 接口

### 账户相关

- `POST /accounts` - 创建新账户
- `GET /accounts` - 获取所有账户
- `PUT /accounts/:id` - 更新账户信息

### 交易相关

- `POST /transactions` - 创建新交易
- `GET /transactions` - 获取交易记录

## 开发计划

1. [x] 基础后端API实现
2. [ ] 前端界面开发
3. [ ] 用户认证
4. [ ] 数据可视化
5. [ ] 导出报表功能