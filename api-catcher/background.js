// background.js - Service Worker 后台脚本

// 上传服务器地址
const UPLOAD_URL = 'http://localhost/api_catcher/upload';

// 初始化
chrome.runtime.onInstalled.addListener(() => {
  console.log('[API Catcher] 插件已安装');

  // 初始化默认设置
  chrome.storage.local.set({
    isRecording: false,
    filterKeywords: [],
    serverConnected: true
  });
});

// 监听标签页更新，注入hook脚本
chrome.tabs.onUpdated.addListener(async (tabId, changeInfo, tab) => {
  // 只在页面加载完成时注入
  if (changeInfo.status !== 'complete' || !tab.url) return;

  // 跳过chrome://和chrome-extension://等内部页面
  if (tab.url.startsWith('chrome://') ||
      tab.url.startsWith('chrome-extension://') ||
      tab.url.startsWith('edge://') ||
      tab.url.startsWith('about:') ||
      tab.url.startsWith('file://')) {
    return;
  }

  try {
    // 使用chrome.scripting.executeScript在MAIN world中执行
    await chrome.scripting.executeScript({
      target: { tabId: tabId, allFrames: true },
      world: 'MAIN',
      injectImmediately: true,
      func: injectHookScript
    });
    console.log('[API Catcher] 脚本已注入到标签页:', tabId);
  } catch (error) {
    console.error('[API Catcher] 注入脚本失败:', error);
  }
});

// 注入到页面的hook脚本函数
function injectHookScript() {
  // 如果已经注入，跳过
  if (window.__api_catcher_injected__) return;
  window.__api_catcher_injected__ = true;

  // 生成唯一ID
  function generateId() {
    return Date.now() + '_' + Math.random().toString(36).substr(2, 9);
  }

  // 解析URL参数
  function parseParams(url) {
    try {
      const urlObj = new URL(url);
      const params = {};
      urlObj.searchParams.forEach((value, key) => {
        params[key] = value;
      });
      return params;
    } catch (e) {
      return {};
    }
  }

  // 状态
  let isRecording = false;
  let filterKeywords = [];

  // 检查URL是否匹配筛选条件
  function shouldCapture(url) {
    if (!isRecording) return false;
    if (!filterKeywords || filterKeywords.length === 0) return true;
    return filterKeywords.some(keyword => keyword && url.includes(keyword));
  }

  // 发送捕获的数据
  function sendApiData(data) {
    window.postMessage({ type: 'API_CAPTURED', data: data }, '*');
  }

  // 安全地克隆对象
  function safeClone(obj) {
    try {
      if (obj === null || obj === undefined) return null;
      if (typeof obj === 'string') return obj;
      if (obj instanceof FormData) {
        const result = {};
        obj.forEach((value, key) => { result[key] = value; });
        return result;
      }
      if (obj instanceof URLSearchParams) {
        const result = {};
        obj.forEach((value, key) => { result[key] = value; });
        return result;
      }
      if (obj instanceof Blob) return '[Blob]';
      if (obj instanceof ArrayBuffer) return '[ArrayBuffer]';
      if (obj instanceof ReadableStream) return '[ReadableStream]';
      return JSON.parse(JSON.stringify(obj));
    } catch (e) {
      return '[Unable to serialize]';
    }
  }

  // 尝试解析JSON
  function tryParseJson(text) {
    try {
      return JSON.parse(text);
    } catch (e) {
      return text;
    }
  }

  // ==================== Hook XMLHttpRequest ====================
  const OriginalXMLHttpRequest = window.XMLHttpRequest;

  function HookedXMLHttpRequest() {
    const xhr = new OriginalXMLHttpRequest();
    const self = this;

    const originalOpen = xhr.open;
    const originalSend = xhr.send;
    const originalSetRequestHeader = xhr.setRequestHeader;

    let requestInfo = {
      id: generateId(),
      method: 'GET',
      url: '',
      headers: {},
      params: {},
      request_body: null,
      response_body: null,
      status: 0,
      capture_time: 0,
      duration: 0
    };

    let startTime = 0;

    xhr.open = function(method, url, async, user, password) {
      requestInfo.method = method.toUpperCase();
      requestInfo.url = url;
      requestInfo.params = parseParams(url);
      return originalOpen.apply(xhr, arguments);
    };

    xhr.setRequestHeader = function(header, value) {
      requestInfo.headers[header] = value;
      return originalSetRequestHeader.apply(xhr, arguments);
    };

    xhr.send = function(body) {
      if (!shouldCapture(requestInfo.url)) {
        return originalSend.apply(xhr, arguments);
      }

      requestInfo.request_body = safeClone(body);
      startTime = Date.now();
      requestInfo.capture_time = startTime;

      const onLoad = function() {
        requestInfo.duration = Date.now() - startTime;
        requestInfo.status = xhr.status;
        try {
          requestInfo.response_body = tryParseJson(xhr.responseText);
        } catch (e) {
          requestInfo.response_body = '[Unable to read response]';
        }
        sendApiData(requestInfo);
      };

      xhr.addEventListener('load', onLoad);
      xhr.addEventListener('error', onLoad);
      xhr.addEventListener('abort', onLoad);
      xhr.addEventListener('timeout', onLoad);

      return originalSend.apply(xhr, arguments);
    };

    // 代理属性
    ['readyState', 'response', 'responseText', 'responseType', 'responseURL',
     'responseXML', 'status', 'statusText', 'timeout', 'upload', 'withCredentials',
     'onload', 'onloadstart', 'onloadend', 'onerror', 'onabort', 'ontimeout',
     'onprogress', 'onreadystatechange'].forEach(prop => {
      Object.defineProperty(self, prop, {
        get: function() { return xhr[prop]; },
        set: function(value) { xhr[prop] = value; }
      });
    });

    self.abort = function() { return xhr.abort(); };
    self.getAllResponseHeaders = function() { return xhr.getAllResponseHeaders(); };
    self.getResponseHeader = function(header) { return xhr.getResponseHeader(header); };
    self.overrideMimeType = function(mimeType) { return xhr.overrideMimeType(mimeType); };

    return self;
  }

  HookedXMLHttpRequest.prototype = OriginalXMLHttpRequest.prototype;
  HookedXMLHttpRequest.UNSENT = OriginalXMLHttpRequest.UNSENT;
  HookedXMLHttpRequest.OPENED = OriginalXMLHttpRequest.OPENED;
  HookedXMLHttpRequest.HEADERS_RECEIVED = OriginalXMLHttpRequest.HEADERS_RECEIVED;
  HookedXMLHttpRequest.LOADING = OriginalXMLHttpRequest.LOADING;
  HookedXMLHttpRequest.DONE = OriginalXMLHttpRequest.DONE;
  window.XMLHttpRequest = HookedXMLHttpRequest;

  // ==================== Hook Fetch ====================
  const originalFetch = window.fetch;

  window.fetch = async function(input, init) {
    const startTime = Date.now();

    let url = '';
    let method = 'GET';
    let headers = {};
    let requestBody = null;

    if (typeof input === 'string') {
      url = input;
    } else if (input instanceof Request) {
      url = input.url;
      method = input.method || 'GET';
      input.headers.forEach((value, key) => { headers[key] = value; });
    }

    if (init) {
      method = init.method || method;
      if (init.headers) {
        if (init.headers instanceof Headers) {
          init.headers.forEach((value, key) => { headers[key] = value; });
        } else if (typeof init.headers === 'object') {
          headers = { ...headers, ...init.headers };
        }
      }
      requestBody = init.body;
    }

    const shouldCap = shouldCapture(url);
    const requestId = generateId();

    const responsePromise = originalFetch.apply(window, arguments);
    if (!shouldCap) return responsePromise;

    const requestInfo = {
      id: requestId,
      method: method.toUpperCase(),
      url: url,
      headers: headers,
      params: parseParams(url),
      request_body: safeClone(requestBody),
      response_body: null,
      status: 0,
      capture_time: startTime,
      duration: 0
    };

    return responsePromise.then(async (response) => {
      requestInfo.duration = Date.now() - startTime;
      requestInfo.status = response.status;

      try {
        const clonedResponse = response.clone();
        const contentType = clonedResponse.headers.get('content-type') || '';
        if (contentType.includes('application/json')) {
          requestInfo.response_body = await clonedResponse.json();
        } else if (contentType.includes('text/')) {
          const text = await clonedResponse.text();
          requestInfo.response_body = tryParseJson(text);
        } else {
          requestInfo.response_body = '[Binary or non-text response]';
        }
      } catch (e) {
        requestInfo.response_body = '[Unable to read response]';
      }

      sendApiData(requestInfo);
      return response;
    }).catch((error) => {
      requestInfo.duration = Date.now() - startTime;
      requestInfo.status = 0;
      requestInfo.response_body = { error: error.message };
      sendApiData(requestInfo);
      throw error;
    });
  };

  window.fetch.prototype = originalFetch.prototype;

  // 监听来自content script的消息
  window.addEventListener('message', (event) => {
    if (event.source !== window) return;
    const { type, isRecording: newRecordingState, filterKeywords: newFilters } = event.data || {};
    if (type === 'RECORDING_STATUS_CHANGED') isRecording = newRecordingState;
    if (type === 'FILTERS_CHANGED') filterKeywords = newFilters || [];
  });

  console.log('[API Catcher] Hook脚本已加载');
}

// 监听来自content script的消息
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  const { action, data } = message;

  switch (action) {
    case 'apiCaptured':
      handleApiCaptured(data, sender.tab.id);
      break;

    case 'checkRecordingStatus':
      chrome.storage.local.get('isRecording').then(result => {
        sendResponse({ isRecording: result.isRecording || false });
      });
      return true;

    case 'checkFilterKeywords':
      chrome.storage.local.get('filterKeywords').then(result => {
        sendResponse({ filterKeywords: result.filterKeywords || [] });
      });
      return true;
  }
});

// 处理捕获的API数据
async function handleApiCaptured(apiData, tabId) {
  await saveApiToList(apiData, tabId);
  uploadApiData(apiData);
}

// 保存API到列表
async function saveApiToList(apiData, tabId) {
  const storageKey = `apiList_${tabId}`;
  try {
    const result = await chrome.storage.local.get(storageKey);
    let apiList = result[storageKey] || [];
    apiList.push(apiData);
    if (apiList.length > 100) apiList = apiList.slice(-100);
    await chrome.storage.local.set({ [storageKey]: apiList });
  } catch (error) {
    console.error('[API Catcher] 保存API列表失败:', error);
  }
}

// 上传API数据到服务器
async function uploadApiData(apiData) {
  try {
    const result = await chrome.storage.local.get('filterKeywords');
    const filterKeywords = result.filterKeywords || [];

    const payload = {
      filter_list: filterKeywords,
      api_data: apiData
    };

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 5000);

    const response = await fetch(UPLOAD_URL, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
      signal: controller.signal
    });

    clearTimeout(timeoutId);

    if (response.ok) {
      await chrome.storage.local.set({ serverConnected: true });
    } else {
      await chrome.storage.local.set({ serverConnected: false });
    }
  } catch (error) {
    await chrome.storage.local.set({ serverConnected: false });
  }
}

// 标签页关闭时清理数据
chrome.tabs.onRemoved.addListener(async (tabId) => {
  const storageKey = `apiList_${tabId}`;
  try {
    await chrome.storage.local.remove(storageKey);
  } catch (error) {
    console.error('[API Catcher] 清理标签页数据失败:', error);
  }
});

// 监听storage变化，同步给所有标签页
chrome.storage.onChanged.addListener((changes, namespace) => {
  if (namespace !== 'local') return;

  chrome.tabs.query({}, (tabs) => {
    tabs.forEach(tab => {
      if (changes.isRecording) {
        chrome.tabs.sendMessage(tab.id, {
          action: 'recordingChanged',
          isRecording: changes.isRecording.newValue
        }).catch(() => {});
      }
      if (changes.filterKeywords) {
        chrome.tabs.sendMessage(tab.id, {
          action: 'filtersChanged',
          filterKeywords: changes.filterKeywords.newValue
        }).catch(() => {});
      }
    });
  });
});

// 定期检查服务器连接状态
setInterval(async () => {
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 3000);
    await fetch(UPLOAD_URL, { method: 'OPTIONS', signal: controller.signal });
    clearTimeout(timeoutId);

    const result = await chrome.storage.local.get('serverConnected');
    if (result.serverConnected === false) {
      await chrome.storage.local.set({ serverConnected: true });
    }
  } catch (error) {
    const result = await chrome.storage.local.get('serverConnected');
    if (result.serverConnected !== false) {
      await chrome.storage.local.set({ serverConnected: false });
    }
  }
}, 10000);
