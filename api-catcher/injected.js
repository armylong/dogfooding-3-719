// injected.js - 注入到页面上下文的脚本，用于Hook XHR和Fetch

(function() {
  'use strict';

  // 状态
  let isRecording = false;
  let filterKeywords = [];

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

  // 检查URL是否匹配筛选条件
  function shouldCapture(url) {
    // 未开启录制时不捕获
    if (!isRecording) {
      return false;
    }

    // 筛选条件为空时，捕获所有请求
    if (!filterKeywords || filterKeywords.length === 0) {
      return true;
    }

    // 检查URL是否包含任一关键词
    return filterKeywords.some(keyword => {
      if (!keyword) return false;
      return url.includes(keyword);
    });
  }

  // 发送捕获的数据
  function sendApiData(data) {
    window.postMessage({
      type: 'API_CAPTURED',
      data: data
    }, '*');
  }

  // 安全地克隆对象
  function safeClone(obj) {
    try {
      if (obj === null || obj === undefined) {
        return null;
      }
      if (typeof obj === 'string') {
        return obj;
      }
      if (obj instanceof FormData) {
        const result = {};
        obj.forEach((value, key) => {
          result[key] = value;
        });
        return result;
      }
      if (obj instanceof URLSearchParams) {
        const result = {};
        obj.forEach((value, key) => {
          result[key] = value;
        });
        return result;
      }
      if (obj instanceof Blob) {
        return '[Blob]';
      }
      if (obj instanceof ArrayBuffer) {
        return '[ArrayBuffer]';
      }
      if (obj instanceof ReadableStream) {
        return '[ReadableStream]';
      }
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

    // 保存原始方法
    const originalOpen = xhr.open;
    const originalSend = xhr.send;
    const originalSetRequestHeader = xhr.setRequestHeader;

    // 请求信息
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

    // 重写open方法
    xhr.open = function(method, url, async, user, password) {
      requestInfo.method = method.toUpperCase();
      requestInfo.url = url;
      requestInfo.params = parseParams(url);
      return originalOpen.apply(xhr, arguments);
    };

    // 重写setRequestHeader方法
    xhr.setRequestHeader = function(header, value) {
      requestInfo.headers[header] = value;
      return originalSetRequestHeader.apply(xhr, arguments);
    };

    // 重写send方法
    xhr.send = function(body) {
      // 检查是否应该捕获
      if (!shouldCapture(requestInfo.url)) {
        return originalSend.apply(xhr, arguments);
      }

      requestInfo.request_body = safeClone(body);
      startTime = Date.now();
      requestInfo.capture_time = startTime;

      // 监听加载完成
      const onLoad = function() {
        requestInfo.duration = Date.now() - startTime;
        requestInfo.status = xhr.status;

        // 获取响应头
        try {
          const responseHeaders = xhr.getAllResponseHeaders();
          // 可以选择解析响应头
        } catch (e) {}

        // 获取响应体
        try {
          const responseText = xhr.responseText;
          requestInfo.response_body = tryParseJson(responseText);
        } catch (e) {
          requestInfo.response_body = '[Unable to read response]';
        }

        // 发送数据
        sendApiData(requestInfo);
      };

      xhr.addEventListener('load', onLoad);
      xhr.addEventListener('error', onLoad);
      xhr.addEventListener('abort', onLoad);
      xhr.addEventListener('timeout', onLoad);

      return originalSend.apply(xhr, arguments);
    };

    // 代理属性和方法
    const properties = [
      'readyState', 'response', 'responseText', 'responseType', 'responseURL',
      'responseXML', 'status', 'statusText', 'timeout', 'upload', 'withCredentials',
      'onload', 'onloadstart', 'onloadend', 'onerror', 'onabort', 'ontimeout',
      'onprogress', 'onreadystatechange'
    ];

    properties.forEach(prop => {
      Object.defineProperty(self, prop, {
        get: function() {
          return xhr[prop];
        },
        set: function(value) {
          xhr[prop] = value;
        }
      });
    });

    // 代理方法
    self.abort = function() { return xhr.abort(); };
    self.getAllResponseHeaders = function() { return xhr.getAllResponseHeaders(); };
    self.getResponseHeader = function(header) { return xhr.getResponseHeader(header); };
    self.overrideMimeType = function(mimeType) { return xhr.overrideMimeType(mimeType); };

    return self;
  }

  // 复制原型链
  HookedXMLHttpRequest.prototype = OriginalXMLHttpRequest.prototype;
  HookedXMLHttpRequest.UNSENT = OriginalXMLHttpRequest.UNSENT;
  HookedXMLHttpRequest.OPENED = OriginalXMLHttpRequest.OPENED;
  HookedXMLHttpRequest.HEADERS_RECEIVED = OriginalXMLHttpRequest.HEADERS_RECEIVED;
  HookedXMLHttpRequest.LOADING = OriginalXMLHttpRequest.LOADING;
  HookedXMLHttpRequest.DONE = OriginalXMLHttpRequest.DONE;

  // 替换全局XMLHttpRequest
  window.XMLHttpRequest = HookedXMLHttpRequest;

  // ==================== Hook Fetch ====================

  const originalFetch = window.fetch;

  window.fetch = async function(input, init) {
    const startTime = Date.now();

    // 解析请求信息
    let url = '';
    let method = 'GET';
    let headers = {};
    let requestBody = null;

    if (typeof input === 'string') {
      url = input;
    } else if (input instanceof Request) {
      url = input.url;
      method = input.method || 'GET';
      headers = {};
      input.headers.forEach((value, key) => {
        headers[key] = value;
      });
    }

    if (init) {
      method = init.method || method;
      if (init.headers) {
        if (init.headers instanceof Headers) {
          init.headers.forEach((value, key) => {
            headers[key] = value;
          });
        } else if (typeof init.headers === 'object') {
          headers = { ...headers, ...init.headers };
        }
      }
      requestBody = init.body;
    }

    // 检查是否应该捕获
    const shouldCap = shouldCapture(url);

    // 生成请求ID
    const requestId = generateId();

    // 调用原始fetch
    const responsePromise = originalFetch.apply(window, arguments);

    // 如果不捕获，直接返回
    if (!shouldCap) {
      return responsePromise;
    }

    // 准备请求信息
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

    // 处理响应
    return responsePromise.then(async (response) => {
      requestInfo.duration = Date.now() - startTime;
      requestInfo.status = response.status;

      // 克隆响应以读取body
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

      // 发送数据
      sendApiData(requestInfo);

      return response;
    }).catch((error) => {
      // 请求失败也记录
      requestInfo.duration = Date.now() - startTime;
      requestInfo.status = 0;
      requestInfo.response_body = { error: error.message };

      sendApiData(requestInfo);

      throw error;
    });
  };

  // 复制fetch的属性
  window.fetch.prototype = originalFetch.prototype;

  // ==================== 监听来自content script的消息 ====================

  window.addEventListener('message', (event) => {
    if (event.source !== window) return;

    const { type, isRecording: newRecordingState, filterKeywords: newFilters } = event.data || {};

    if (type === 'RECORDING_STATUS_CHANGED') {
      isRecording = newRecordingState;
      console.log('[API Catcher] 录制状态:', isRecording ? '开启' : '关闭');
    }

    if (type === 'FILTERS_CHANGED') {
      filterKeywords = newFilters || [];
      console.log('[API Catcher] 筛选条件:', filterKeywords);
    }
  });

  console.log('[API Catcher] 注入脚本已加载');

})();
