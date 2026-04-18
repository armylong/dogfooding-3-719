// content.js - 内容脚本，负责监听页面消息并与background通信

// 监听来自injected hook脚本的消息（通过postMessage）
window.addEventListener('message', async (event) => {
  // 只处理来自当前窗口的消息
  if (event.source !== window) return;

  const { type, data } = event.data || {};

  if (type === 'API_CAPTURED') {
    // 转发给background.js
    try {
      await chrome.runtime.sendMessage({
        action: 'apiCaptured',
        data: data
      });
    } catch (error) {
      console.error('[API Catcher] 发送数据失败:', error);
    }
  }
});

// 监听来自popup/background的消息
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  const { action, isRecording, filterKeywords } = message;

  switch (action) {
    case 'recordingChanged':
      // 转发给页面的hook脚本
      window.postMessage({
        type: 'RECORDING_STATUS_CHANGED',
        isRecording: isRecording
      }, '*');
      break;

    case 'filtersChanged':
      // 转发给页面的hook脚本
      window.postMessage({
        type: 'FILTERS_CHANGED',
        filterKeywords: filterKeywords
      }, '*');
      break;
  }
});

console.log('[API Catcher] Content script 已加载');
