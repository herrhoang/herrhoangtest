# 个人记账系统

这是一个使用 Go 和 React 构建的个人记账系统，采用前后端分离架构，提供直观的用户界面和可靠的数据存储功能。

## 系统架构

### 前端架构

```
frontend/
├── src/
│   ├── components/          # 可复用组件
│   │   ├── AccountList      # 账户列表组件
│   │   └── TransactionList  # 交易列表组件
│   ├── pages/              # 页面组件
│   │   ├── Dashboard       # 仪表盘页面
│   │   ├── AccountPage     # 账户管理页面
│   │   └── TransactionPage # 交易管理页面
│   ├── services/           # API服务
│   │   └── api.ts         # API调用封装
│   └── types/             # TypeScript类型定义
└── public/                # 静态资源
```

### 后端架构

```
backend/
├── config/               # 配置文件
├── handlers/             # HTTP处理器
│   ├── account_handler   # 账户相关处理
│   └── transaction_handler # 交易相关处理
├── models/               # 数据模型
│   ├── account.go       # 账户模型
│   └── transaction.go   # 交易模型
├── database/            # 数据库配置
└── main.go             # 应用入口
```

## 功能模块

### 1. 账户管理

#### 前端实现
- `AccountPage.tsx`: 账户管理主页面
  - 展示账户列表
  - 添加/编辑账户表单
  - 余额显示和更新

#### 后端API
```typescript
interface AccountAPI {
  getAll: () => Promise<Account[]>;
  create: (data: { name: string, balance: number }) => Promise<Account>;
  update: (id: number, data: { name?: string, balance?: number }) => Promise<Account>;
}
```

### 2. 交易管理

#### 前端实现
- `TransactionPage.tsx`: 交易管理主页面
  - 交易记录列表
  - 添加新交易表单
  - 交易分类和筛选

#### 后端API
```typescript
interface TransactionAPI {
  getAll: () => Promise<Transaction[]>;
  create: (data: {
    account_id: number,
    amount: number,
    type: 'income' | 'expense',
    category_id: number,
    description?: string
  }) => Promise<Transaction>;
}
```

### 3. 数据展示

#### 前端实现
- `Dashboard.tsx`: 数据统计仪表盘
  - 总资产概览
  - 收支统计
  - 账户余额分布

## 数据流

1. 账户创建流程
```
前端 AccountPage
    → 调用 accountApi.create()
    → 后端 /accounts POST
    → 数据库创建记录
    → 返回新账户数据
    → 前端更新状态并显示
```

2. 交易记录流程
```
前端 TransactionPage
    → 调用 transactionApi.create()
    → 后端 /transactions POST
    → 数据库创建交易记录
    → 更新相关账户余额
    → 返回交易数据
    → 前端更新状态并显示
```

## 技术栈

### 前端
- React 18
- TypeScript 4.x
- Ant Design 5.x
- Axios 用于 API 调用

### 后端
- Go 1.21
- Gin Web框架
- GORM ORM框架
- SQLite 数据库

## 安装说明

### 后端安装
1. 安装 Go (1.21 或更高版本)
2. 进入后端目录：`cd backend`
3. 安装依赖：`go mod tidy`
4. 运行服务器：`go run main.go`

### 前端安装
1. 安装 Node.js (v16 或更高版本)
2. 进入前端目录：`cd frontend`
3. 安装依赖：`npm install`
4. 运行开发服务器：`npm start`

## 开发计划

1. [x] 基础后端API实现
2. [x] 前端界面开发
3. [ ] 用户认证
4. [ ] 数据可视化
5. [ ] 导出报表功能
6. [ ] 预算管理
7. [ ] 定期交易