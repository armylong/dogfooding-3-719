(function() {
  'use strict';

  const HOOK_SCRIPT_ID = 'api-catcher-hook-script';

  if (window[HOOK_SCRIPT_ID]) {
    return;
  }
  window[HOOK_SCRIPT_ID] = true;

  function generateId() {
    return `${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
  }

  function parseUrl(url) {
    try {
      const urlObj = new URL(url, window.location.origin);
      const params = {};
      urlObj.searchParams.forEach((value, key) => {
        params[key] = value;
      });
      return {
        fullUrl: urlObj.href,
        params
      };
    } catch (e) {
      return { fullUrl: url, params: {} };
    }
  }

  function parseHeaders(headers) {
    const result = {};
    if (headers instanceof Headers) {
      headers.forEach((value, key) => {
        result[key] = value;
      });
    } else if (typeof headers === 'object') {
      Object.assign(result, headers);
    }
    return result;
  }

  function parseBody(body) {
    if (!body) return null;
    try {
      if (typeof body === 'string') {
        try {
          return JSON.parse(body);
        } catch {
          return body;
        }
      }
      return body;
    } catch (e) {
      return null;
    }
  }

  function sendToContentScript(data) {
    window.postMessage({
      type: 'API_CATCHER_REQUEST',
      data: data
    }, '*');
  }

  function hookXHR() {
    const originalXHR = window.XMLHttpRequest;
    const XHR = function() {
      const xhr = new originalXHR();
      const requestData = {
        id: generateId(),
        url: '',
        method: '',
        headers: {},
        params: {},
        request_body: null,
        response_body: null,
        status: 0,
        capture_time: 0,
        duration: 0
      };
      let startTime = 0;
      let requestHeaders = {};

      const originalOpen = xhr.open;
      xhr.open = function(method, url, async, user, password) {
        requestData.method = method.toUpperCase();
        const urlInfo = parseUrl(url);
        requestData.url = urlInfo.fullUrl;
        requestData.params = urlInfo.params;
        return originalOpen.apply(this, arguments);
      };

      const originalSetRequestHeader = xhr.setRequestHeader;
      xhr.setRequestHeader = function(name, value) {
        requestHeaders[name] = value;
        return originalSetRequestHeader.apply(this, arguments);
      };

      const originalSend = xhr.send;
      xhr.send = function(body) {
        startTime = Date.now();
        requestData.capture_time = startTime;
        requestData.headers = { ...requestHeaders };
        requestData.request_body = parseBody(body);

        xhr.addEventListener('load', function() {
          requestData.duration = Date.now() - startTime;
          requestData.status = xhr.status;

          try {
            const contentType = xhr.getResponseHeader('Content-Type') || '';
            if (contentType.includes('application/json')) {
              requestData.response_body = JSON.parse(xhr.responseText);
            } else {
              requestData.response_body = xhr.responseText;
            }
          } catch (e) {
            requestData.response_body = xhr.responseText;
          }

          sendToContentScript(requestData);
        });

        xhr.addEventListener('error', function() {
          requestData.duration = Date.now() - startTime;
          requestData.status = 0;
          requestData.response_body = null;
          sendToContentScript(requestData);
        });

        return originalSend.apply(this, arguments);
      };

      return xhr;
    };

    XHR.prototype = originalXHR.prototype;
    XHR.DONE = originalXHR.DONE;
    XHR.HEADERS_RECEIVED = originalXHR.HEADERS_RECEIVED;
    XHR.LOADING = originalXHR.LOADING;
    XHR.OPENED = originalXHR.OPENED;
    XHR.UNSENT = originalXHR.UNSENT;

    window.XMLHttpRequest = XHR;
  }

  function hookFetch() {
    const originalFetch = window.fetch;

    window.fetch = async function(input, init = {}) {
      const startTime = Date.now();
      const requestData = {
        id: generateId(),
        url: '',
        method: 'GET',
        headers: {},
        params: {},
        request_body: null,
        response_body: null,
        status: 0,
        capture_time: startTime,
        duration: 0
      };

      let url;
      if (typeof input === 'string') {
        url = input;
      } else if (input instanceof Request) {
        url = input.url;
        requestData.method = input.method.toUpperCase();
        requestData.headers = parseHeaders(input.headers);
        if (input.body) {
          try {
            const clonedRequest = input.clone();
            const bodyText = await clonedRequest.text();
            requestData.request_body = parseBody(bodyText);
          } catch (e) {
            requestData.request_body = null;
          }
        }
      } else {
        url = input.url || '';
      }

      const urlInfo = parseUrl(url);
      requestData.url = urlInfo.fullUrl;
      requestData.params = urlInfo.params;

      if (init.method) {
        requestData.method = init.method.toUpperCase();
      }
      if (init.headers) {
        requestData.headers = parseHeaders(init.headers);
      }
      if (init.body && !requestData.request_body) {
        requestData.request_body = parseBody(init.body);
      }

      try {
        const response = await originalFetch.apply(this, arguments);
        requestData.duration = Date.now() - startTime;
        requestData.status = response.status;

        try {
          const clonedResponse = response.clone();
          const contentType = clonedResponse.headers.get('Content-Type') || '';
          if (contentType.includes('application/json')) {
            requestData.response_body = await clonedResponse.json();
          } else {
            requestData.response_body = await clonedResponse.text();
          }
        } catch (e) {
          requestData.response_body = null;
        }

        sendToContentScript(requestData);
        return response;
      } catch (error) {
        requestData.duration = Date.now() - startTime;
        requestData.status = 0;
        requestData.response_body = null;
        sendToContentScript(requestData);
        throw error;
      }
    };
  }

  hookXHR();
  hookFetch();
})();
