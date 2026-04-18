# 抓接口 - Chrome浏览器扩展

一个用于抓取页面接口数据并上传至服务器的Chrome浏览器扩展。

## 功能特性

- **请求捕获**：自动Hook XHR和Fetch请求，捕获完整的请求/响应信息
- **筛选过滤**：支持多个筛选关键词，只捕获符合条件的接口
- **录制开关**：可随时开启/关闭录制，关闭时不捕获不上传
- **接口列表**：实时展示已捕获的接口，按时间倒序排列
- **数据隔离**：不同标签页的数据相互隔离
- **自动上传**：录制开启时，每捕获一个请求立即异步上传

## 安装方法

1. 打开Chrome浏览器，进入扩展管理页面 `chrome://extensions/`
2. 开启右上角的"开发者模式"
3. 点击"加载已解压的扩展程序"
4. 选择 `api-catcher` 文件夹

## 使用说明

### 基本操作

1. 点击浏览器工具栏中的扩展图标打开控制面板
2. 设置筛选条件（可选）
3. 点击"开始录制"按钮开启录制
4. 在页面中进行操作，触发接口请求
5. 捕获的接口会实时显示在接口列表中

### 筛选条件

- 筛选条件支持URL关键词匹配
- 可以添加多个筛选条件，满足任一条件即捕获
- 筛选条件为空时，捕获所有请求
- 筛选条件会自动保存

### 接口列表

- 显示已捕获的接口数量
- 按时间倒序排列（最新的在最上面）
- 显示请求方法、状态码、时间、URL路径
- 点击"清空"按钮可清除当前标签页的所有数据

## 捕获的数据字段

| 字段 | 说明 |
|------|------|
| id | 唯一标识（时间戳+随机数） |
| url | 完整URL |
| method | 请求方法（GET/POST/PUT/DELETE等） |
| headers | 请求头 |
| params | URL查询参数 |
| request_body | 请求体 |
| response_body | 响应体 |
| status | HTTP状态码 |
| capture_time | 捕获时间戳（毫秒） |
| duration | 请求耗时（毫秒） |

## 上传接口

捕获的数据会上传到：`http://localhost/api_catcher/upload`

### 请求格式

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

## 文件结构

```
api-catcher/
├── manifest.json    # 扩展配置文件
├── popup.html       # 弹出窗口HTML
├── popup.css        # 弹出窗口样式
├── popup.js         # 弹出窗口逻辑
├── background.js    # 后台服务脚本
├── content.js       # 内容脚本（与页面通信）
├── injected.js      # 注入脚本（Hook XHR/Fetch）
├── api_doc.md       # API接口文档
└── README.md        # 说明文档
```

## 架构说明

由于Manifest V3的内容安全策略(CSP)限制，content script无法直接修改页面的XHR和Fetch对象。因此采用以下架构：

1. **injected.js**：注入到页面上下文，Hook XHR和Fetch，通过`window.postMessage`发送捕获的数据
2. **content.js**：监听来自injected.js的消息，转发给background.js
3. **background.js**：处理数据存储和上传逻辑

## 注意事项

1. 上传接口需要服务端支持CORS
2. 如果服务接口不可用，会在页面显示"上传接口未连通"
3. 每个标签页最多保存500条接口数据
4. 关闭标签页后，该标签页的数据会被清除

## 开发说明

- 使用Manifest V3规范
- 使用Chrome Storage API进行数据持久化
- 使用Content Script + Injected Script实现请求Hook
- 使用Service Worker作为后台脚本
