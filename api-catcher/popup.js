// popup.js - 插件弹窗页面逻辑

// 全局状态
let currentTabId = null;
let apiList = [];
let isRecording = false;
let filterKeywords = [];
let serverConnected = true;

// DOM 元素
let recordToggle;
let filterListEl;
let apiListEl;
let apiCountEl;
let btnClear;
let btnAddFilter;
let serverStatusEl;

// 初始化
document.addEventListener('DOMContentLoaded', async () => {
  initElements();
  await getCurrentTabId();
  await loadState();
  bindEvents();
  render();
  startApiListSync();
});

// 初始化DOM元素引用
function initElements() {
  recordToggle = document.getElementById('recordToggle');
  filterListEl = document.getElementById('filterList');
  apiListEl = document.getElementById('apiList');
  apiCountEl = document.getElementById('apiCount');
  btnClear = document.getElementById('btnClear');
  btnAddFilter = document.getElementById('btnAddFilter');
  serverStatusEl = document.getElementById('serverStatus');
}

// 获取当前标签页ID
async function getCurrentTabId() {
  const tabs = await chrome.tabs.query({ active: true, currentWindow: true });
  if (tabs.length > 0) {
    currentTabId = tabs[0].id;
  }
}

// 从storage加载状态
async function loadState() {
  // 加载录制状态
  const recordResult = await chrome.storage.local.get('isRecording');
  isRecording = recordResult.isRecording || false;
  
  // 加载筛选关键词
  const filterResult = await chrome.storage.local.get('filterKeywords');
  filterKeywords = filterResult.filterKeywords || [];
  
  // 加载当前标签页的API列表
  if (currentTabId) {
    const apiResult = await chrome.storage.local.get(`apiList_${currentTabId}`);
    apiList = apiResult[`apiList_${currentTabId}`] || [];
  }
  
  // 加载服务器连接状态
  const serverResult = await chrome.storage.local.get('serverConnected');
  serverConnected = serverResult.serverConnected !== false;
}

// 绑定事件
function bindEvents() {
  // 录制开关
  recordToggle.addEventListener('change', async () => {
    isRecording = recordToggle.checked;
    await chrome.storage.local.set({ isRecording });
    
    // 通知content script更新录制状态
    await notifyContentScript('recordingChanged', { isRecording });
    render();
  });
  
  // 添加筛选条件
  btnAddFilter.addEventListener('click', addFilterInput);
  
  // 清空列表
  btnClear.addEventListener('click', clearApiList);
}

// 添加筛选输入框
function addFilterInput(keyword = '') {
  const filterItem = document.createElement('div');
  filterItem.className = 'filter-item';
  
  const input = document.createElement('input');
  input.type = 'text';
  input.className = 'filter-input';
  input.placeholder = '输入URL关键词，如: api/user';
  input.value = keyword;
  
  // 输入变化时保存
  input.addEventListener('input', debounce(saveFilters, 300));
  
  const removeBtn = document.createElement('button');
  removeBtn.className = 'btn-remove';
  removeBtn.innerHTML = '×';
  removeBtn.title = '删除';
  removeBtn.addEventListener('click', () => {
    filterItem.remove();
    saveFilters();
  });
  
  filterItem.appendChild(input);
  filterItem.appendChild(removeBtn);
  filterListEl.appendChild(filterItem);
}

// 保存筛选条件
async function saveFilters() {
  const inputs = filterListEl.querySelectorAll('.filter-input');
  filterKeywords = Array.from(inputs)
    .map(input => input.value.trim())
    .filter(keyword => keyword.length > 0);
  
  await chrome.storage.local.set({ filterKeywords });
  
  // 通知content script更新筛选条件
  await notifyContentScript('filtersChanged', { filterKeywords });
}

// 清空API列表
async function clearApiList() {
  if (currentTabId) {
    await chrome.storage.local.set({ [`apiList_${currentTabId}`]: [] });
    apiList = [];
    renderApiList();
  }
}

// 通知content script
async function notifyContentScript(action, data) {
  try {
    const tabs = await chrome.tabs.query({});
    for (const tab of tabs) {
      try {
        await chrome.tabs.sendMessage(tab.id, { action, ...data });
      } catch (e) {
        // 忽略无法发送消息的标签页
      }
    }
  } catch (e) {
    console.error('通知content script失败:', e);
  }
}

// 渲染界面
function render() {
  renderRecordToggle();
  renderFilterList();
  renderApiList();
  renderServerStatus();
}

// 渲染录制开关
function renderRecordToggle() {
  recordToggle.checked = isRecording;
  const statusDot = document.querySelector('.status-dot');
  const statusText = document.querySelector('.status-text');
  
  if (isRecording) {
    statusDot.classList.add('recording');
    statusText.textContent = '录制中';
  } else {
    statusDot.classList.remove('recording');
    statusText.textContent = '已停止';
  }
}

// 渲染筛选列表
function renderFilterList() {
  filterListEl.innerHTML = '';
  if (filterKeywords.length === 0) {
    addFilterInput('');
  } else {
    filterKeywords.forEach(keyword => addFilterInput(keyword));
  }
}

// 渲染API列表
function renderApiList() {
  apiCountEl.textContent = apiList.length;
  btnClear.disabled = apiList.length === 0;
  
  if (apiList.length === 0) {
    apiListEl.innerHTML = `
      <div class="empty-state">
        <div class="empty-state-icon">📭</div>
        <div>暂无捕获的接口</div>
        <div style="font-size: 12px; margin-top: 4px;">开启录制后将自动捕获</div>
      </div>
    `;
    return;
  }
  
  apiListEl.innerHTML = '';
  // 按时间倒序排列
  const sortedList = [...apiList].sort((a, b) => b.capture_time - a.capture_time);
  
  sortedList.forEach(api => {
    const item = document.createElement('div');
    item.className = 'api-item';
    
    const methodClass = (api.method || 'get').toLowerCase();
    const timeStr = formatTime(api.capture_time);
    const urlObj = new URL(api.url);
    const displayUrl = urlObj.pathname + urlObj.search;
    
    item.innerHTML = `
      <div class="api-item-header">
        <span class="api-method ${methodClass}">${api.method || 'GET'}</span>
        <span class="api-time">${timeStr}</span>
      </div>
      <div class="api-url">${escapeHtml(displayUrl)}</div>
    `;
    
    // 点击显示详情（可扩展）
    item.addEventListener('click', () => {
      console.log('API详情:', api);
    });
    
    apiListEl.appendChild(item);
  });
}

// 渲染服务器状态
function renderServerStatus() {
  if (serverConnected) {
    serverStatusEl.classList.add('hidden');
  } else {
    serverStatusEl.classList.remove('hidden');
    serverStatusEl.innerHTML = `
      <span>⚠️</span>
      <span>上传接口未连通</span>
    `;
  }
}

// 同步API列表（定时从storage刷新）
function startApiListSync() {
  // 立即同步一次
  syncApiList();
  
  // 每500ms同步一次
  setInterval(syncApiList, 500);
}

// 同步API列表
async function syncApiList() {
  if (!currentTabId) return;
  
  const result = await chrome.storage.local.get([
    `apiList_${currentTabId}`,
    'serverConnected'
  ]);
  
  const newApiList = result[`apiList_${currentTabId}`] || [];
  const newServerConnected = result.serverConnected !== false;
  
  // 只在数据变化时重新渲染
  if (JSON.stringify(newApiList) !== JSON.stringify(apiList)) {
    apiList = newApiList;
    renderApiList();
  }
  
  if (newServerConnected !== serverConnected) {
    serverConnected = newServerConnected;
    renderServerStatus();
  }
}

// 格式化时间
function formatTime(timestamp) {
  const date = new Date(timestamp);
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  const ms = String(date.getMilliseconds()).padStart(3, '0');
  return `${hours}:${minutes}:${seconds}.${ms}`;
}

// HTML转义
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

// 防抖函数
function debounce(func, wait) {
  let timeout;
  return function executedFunction(...args) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  };
}
