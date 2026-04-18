let currentTabId = null;

document.addEventListener('DOMContentLoaded', async () => {
  const tab = await getCurrentTab();
  currentTabId = tab ? tab.id : null;

  await loadRecordState();
  await loadFilters();
  await loadApiList();
  await checkServerStatus();

  setupEventListeners();
});

async function getCurrentTab() {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  return tab;
}

async function loadRecordState() {
  const result = await chrome.storage.local.get(['recordState']);
  const isRecording = result.recordState || false;
  updateRecordButton(isRecording);
}

async function loadFilters() {
  const result = await chrome.storage.local.get(['filterList']);
  const filterList = result.filterList || [''];
  renderFilterList(filterList);
}

async function loadApiList() {
  if (!currentTabId) return;

  const storageKey = `apiList_${currentTabId}`;
  const result = await chrome.storage.local.get([storageKey]);
  const apiList = result[storageKey] || [];
  renderApiList(apiList);
}

async function checkServerStatus() {
  try {
    const response = await fetch('http://localhost/api_catcher/upload', {
      method: 'OPTIONS',
      mode: 'cors'
    });
    hideErrorMessage();
  } catch (error) {
    showErrorMessage('上传接口未连通，请检查服务器状态');
  }
}

function showErrorMessage(message) {
  const errorEl = document.getElementById('errorMessage');
  errorEl.textContent = message;
  errorEl.classList.add('show');
}

function hideErrorMessage() {
  const errorEl = document.getElementById('errorMessage');
  errorEl.classList.remove('show');
}

function updateRecordButton(isRecording) {
  const btn = document.getElementById('recordBtn');
  btn.textContent = isRecording ? '录制中' : '开始录制';
  btn.classList.remove('on', 'off');
  btn.classList.add(isRecording ? 'on' : 'off');
}

function renderFilterList(filterList) {
  const container = document.getElementById('filterList');
  container.innerHTML = '';

  filterList.forEach((filter, index) => {
    const item = document.createElement('div');
    item.className = 'filter-item';
    item.innerHTML = `
      <input type="text" class="filter-input" value="${escapeHtml(filter)}" data-index="${index}" placeholder="输入筛选关键词">
      <button class="filter-remove-btn" data-index="${index}">×</button>
    `;
    container.appendChild(item);
  });

  setupFilterEventListeners();
}

function renderApiList(apiList) {
  const container = document.getElementById('apiList');
  const countEl = document.getElementById('apiCount');

  countEl.textContent = `共 ${apiList.length} 条`;

  if (apiList.length === 0) {
    container.innerHTML = '<div class="empty-list">暂无捕获的接口数据</div>';
    return;
  }

  const sortedList = [...apiList].sort((a, b) => b.capture_time - a.capture_time);

  container.innerHTML = sortedList.map(api => `
    <div class="api-item" data-id="${api.id}">
      <div class="api-item-header">
        <span class="api-method ${api.method.toLowerCase()}">${api.method}</span>
        <span class="api-status ${api.status >= 200 && api.status < 300 ? 'success' : 'error'}">${api.status}</span>
        <span class="api-time">${formatTime(api.capture_time)}</span>
      </div>
      <div class="api-url">${escapeHtml(getPathFromUrl(api.url))}</div>
    </div>
  `).join('');
}

function setupEventListeners() {
  document.getElementById('recordBtn').addEventListener('click', toggleRecord);
  document.getElementById('addFilterBtn').addEventListener('click', addFilter);
  document.getElementById('clearBtn').addEventListener('click', clearApiList);

  chrome.storage.onChanged.addListener((changes, namespace) => {
    if (namespace === 'local' && currentTabId) {
      const storageKey = `apiList_${currentTabId}`;
      if (changes[storageKey]) {
        renderApiList(changes[storageKey].newValue || []);
      }
      if (changes.recordState) {
        updateRecordButton(changes.recordState.newValue);
      }
    }
  });
}

function setupFilterEventListeners() {
  document.querySelectorAll('.filter-input').forEach(input => {
    input.addEventListener('input', debounce(handleFilterInput, 300));
  });

  document.querySelectorAll('.filter-remove-btn').forEach(btn => {
    btn.addEventListener('click', handleFilterRemove);
  });
}

async function toggleRecord() {
  const result = await chrome.storage.local.get(['recordState']);
  const currentState = result.recordState || false;
  const newState = !currentState;

  await chrome.storage.local.set({ recordState: newState });
  updateRecordButton(newState);

  chrome.runtime.sendMessage({
    type: 'RECORD_STATE_CHANGED',
    isRecording: newState
  });
}

async function handleFilterInput(e) {
  const index = parseInt(e.target.dataset.index);
  const result = await chrome.storage.local.get(['filterList']);
  const filterList = result.filterList || [''];
  filterList[index] = e.target.value;
  await chrome.storage.local.set({ filterList });
}

async function handleFilterRemove(e) {
  const index = parseInt(e.target.dataset.index);
  const result = await chrome.storage.local.get(['filterList']);
  let filterList = result.filterList || [''];

  if (filterList.length <= 1) {
    filterList = [''];
  } else {
    filterList.splice(index, 1);
  }

  await chrome.storage.local.set({ filterList });
  renderFilterList(filterList);
}

async function addFilter() {
  const result = await chrome.storage.local.get(['filterList']);
  const filterList = result.filterList || [''];
  filterList.push('');
  await chrome.storage.local.set({ filterList });
  renderFilterList(filterList);
}

async function clearApiList() {
  if (!currentTabId) return;

  const storageKey = `apiList_${currentTabId}`;
  await chrome.storage.local.set({ [storageKey]: [] });
  renderApiList([]);
}

function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text || '';
  return div.innerHTML;
}

function getPathFromUrl(url) {
  try {
    const urlObj = new URL(url);
    return urlObj.pathname + urlObj.search;
  } catch {
    return url;
  }
}

function formatTime(timestamp) {
  const date = new Date(timestamp);
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  const seconds = date.getSeconds().toString().padStart(2, '0');
  return `${hours}:${minutes}:${seconds}`;
}

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
