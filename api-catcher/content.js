(function() {
  'use strict';

  let isRecording = false;
  let filterList = [''];

  function init() {
    loadState();
    setupMessageListener();
  }

  async function loadState() {
    try {
      const response = await chrome.runtime.sendMessage({ type: 'GET_RECORD_STATE' });
      isRecording = response?.isRecording || false;
    } catch (e) {
      isRecording = false;
    }

    try {
      const response = await chrome.runtime.sendMessage({ type: 'GET_FILTER_LIST' });
      filterList = response?.filterList || [''];
    } catch (e) {
      filterList = [''];
    }
  }

  function setupMessageListener() {
    chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
      if (message.type === 'RECORD_STATE_CHANGED') {
        isRecording = message.isRecording;
      }
    });

    window.addEventListener('message', (event) => {
      if (event.source !== window) return;
      if (event.data.type === 'API_CATCHER_REQUEST') {
        captureRequest(event.data.data);
      }
    });
  }

  function shouldCapture(url, filters) {
    const validFilters = filters.filter(f => f && f.trim() !== '');
    if (validFilters.length === 0) {
      return true;
    }
    return validFilters.some(filter => url.includes(filter.trim()));
  }

  function captureRequest(data) {
    if (!isRecording) return;

    if (!shouldCapture(data.url, filterList)) return;

    chrome.runtime.sendMessage({
      type: 'API_CAPTURED',
      data: data
    }).catch(() => {});
  }

  init();
})();
