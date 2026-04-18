const UPLOAD_URL = 'http://localhost/api_catcher/upload';
let serverAvailable = true;

chrome.runtime.onInstalled.addListener(() => {
  chrome.storage.local.set({
    recordState: false,
    filterList: ['']
  });
});

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.type === 'API_CAPTURED') {
    handleApiCaptured(message.data, sender.tab?.id);
  } else if (message.type === 'RECORD_STATE_CHANGED') {
    notifyContentScripts(message.isRecording);
  } else if (message.type === 'GET_RECORD_STATE') {
    chrome.storage.local.get(['recordState'], (result) => {
      sendResponse({ isRecording: result.recordState || false });
    });
    return true;
  } else if (message.type === 'GET_FILTER_LIST') {
    chrome.storage.local.get(['filterList'], (result) => {
      sendResponse({ filterList: result.filterList || [''] });
    });
    return true;
  }
});

async function handleApiCaptured(apiData, tabId) {
  if (!tabId) return;

  const result = await chrome.storage.local.get(['recordState', 'filterList']);
  const isRecording = result.recordState || false;
  const filterList = result.filterList || [''];

  if (!isRecording) return;

  if (!shouldCapture(apiData.url, filterList)) return;

  const storageKey = `apiList_${tabId}`;
  const storageResult = await chrome.storage.local.get([storageKey]);
  const apiList = storageResult[storageKey] || [];

  apiList.push(apiData);

  if (apiList.length > 500) {
    apiList.splice(0, apiList.length - 500);
  }

  await chrome.storage.local.set({ [storageKey]: apiList });

  uploadApiData(apiData, filterList);
}

function shouldCapture(url, filterList) {
  const validFilters = filterList.filter(f => f && f.trim() !== '');

  if (validFilters.length === 0) {
    return true;
  }

  return validFilters.some(filter => url.includes(filter.trim()));
}

async function uploadApiData(apiData, filterList) {
  try {
    const response = await fetch(UPLOAD_URL, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        filter_list: filterList.filter(f => f && f.trim() !== ''),
        api_data: apiData
      })
    });

    if (!response.ok) {
      console.error('Upload failed:', response.status);
      serverAvailable = false;
    } else {
      serverAvailable = true;
    }
  } catch (error) {
    console.error('Upload error:', error);
    serverAvailable = false;
  }
}

async function notifyContentScripts(isRecording) {
  const tabs = await chrome.tabs.query({});
  for (const tab of tabs) {
    try {
      await chrome.tabs.sendMessage(tab.id, {
        type: 'RECORD_STATE_CHANGED',
        isRecording
      });
    } catch (e) {
      // 忽略无法发送的标签页
    }
  }
}

chrome.tabs.onRemoved.addListener(async (tabId) => {
  const storageKey = `apiList_${tabId}`;
  await chrome.storage.local.remove([storageKey]);
});
