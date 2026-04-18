# 抓接口 - Chrome浏览器扩展

一个用于捕获网页XHR/Fetch请求并上传至服务器的Chrome浏览器扩展。

## 功能特性

- **请求捕获**: Hook页面的XMLHttpRequest和Fetch请求，捕获完整请求信息
- **筛选功能**: 支持多个关键词筛选，只捕获符合条件的请求
- **录制开关**: 可控制是否捕获请求，状态持久化
- **接口列表**: 实时展示已捕获的接口，按时间倒序排列
- **数据上传**: 自动将捕获的数据上传至指定服务器
- **标签隔离**: 不同标签页的捕获数据相互隔离
- **连接检测**: 自动检测服务器连接状态，异常时提示

## 安装方法

### 开发者模式安装

1. 打开Chrome浏览器，访问 `chrome://extensions/`
2. 开启右上角的"开发者模式"
3. 点击"加载已解压的扩展程序"
4. 选择 `api-catcher` 文件夹
5. 安装完成，浏览器工具栏会出现扩展图标

## 文件结构

```
api-catcher/
├── manifest.json      # 扩展配置文件
├── popup.html         # 弹窗页面
├── popup.css          # 弹窗样式
├── popup.js           # 弹窗逻辑
├── background.js      # 后台脚本(Service Worker)
├── content.js         # 内容脚本
├── injected.js        # 注入页面的脚本
├── api_doc.md         # API接口文档
└── icons/             # 图标文件夹
    ├── icon16.png
    ├── icon48.png
    └── icon128.png
```

## 使用说明

### 1. 开启录制

点击浏览器工具栏的扩展图标，打开弹窗界面，点击"录制开关"开启录制。

### 2. 设置筛选条件

在"筛选条件"区域添加关键词，例如：
- `api/user` - 只捕获包含api/user的URL
- `api/order` - 只捕获包含api/order的URL

筛选条件为空时，将捕获所有请求。

### 3. 查看捕获的接口

在"接口列表"区域查看已捕获的接口，显示：
- 请求方法 (GET/POST/PUT/DELETE等)
- 请求时间
- URL路径

### 4. 清空列表

点击"清空"按钮可清空当前标签页的捕获列表。

### 5. 服务器连接状态

当上传服务器不可用时，界面底部会显示"上传接口未连通"提示。

## 数据格式

### 捕获的数据字段

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 唯一标识（时间戳+随机数） |
| url | string | 完整请求URL |
| method | string | 请求方法（GET/POST/PUT/DELETE等） |
| headers | object | 请求头信息 |
| params | object | URL查询参数 |
| request_body | any | 请求体数据 |
| response_body | any | 响应体数据 |
| status | number | HTTP状态码 |
| capture_time | number | 捕获时间戳（毫秒） |
| duration | number | 请求耗时（毫秒） |

### 上传接口

- **URL**: `http://localhost/api_catcher/upload`
- **Method**: POST
- **Content-Type**: application/json

### 请求体格式

```json
{
  "filter_list": ["api/user", "api/order"],
  "api_data": {
    "id": "1712345678901_abc123",
    "url": "https://example.com/api/user/info",
    "method": "GET",
    "headers": {},
    "params": {},
    "request_body": null,
    "response_body": {},
    "status": 200,
    "capture_time": 1712345678901,
    "duration": 120
  }
}
```

## 技术实现

### 请求捕获原理

1. **injected.js** 注入到页面上下文，重写 `XMLHttpRequest` 和 `fetch`
2. 拦截请求和响应，收集完整数据
3. 通过 `window.postMessage` 将数据发送给 content.js
4. **content.js** 转发给 **background.js**
5. **background.js** 保存数据并上传至服务器

### 状态同步

- 录制状态和筛选条件存储在 `chrome.storage.local`
- 通过 `chrome.runtime.sendMessage` 实现组件间通信
- 标签页关闭时自动清理对应数据

### 异常处理

- 上传失败时自动标记服务器状态
- 每10秒自动检测服务器连接
- 页面加载时异步上传，不阻塞页面

## 注意事项

1. **图标文件**: 需要准备 `icons/icon16.png`、`icons/icon48.png`、`icons/icon128.png` 三个图标文件
2. **服务器地址**: 默认上传到 `http://localhost/api_catcher/upload`，如需修改请编辑 `background.js`
3. **跨域问题**: 确保服务器允许跨域请求（CORS）
4. **数据存储**: 每个标签页最多保留最近100条记录

## 浏览器兼容性

- Chrome 88+ (Manifest V3)
- Edge 88+ (基于Chromium)

## 更新日志

### v1.0.0
- 初始版本
- 实现XHR和Fetch请求捕获
- 实现筛选和录制功能
- 实现数据上传和列表展示
