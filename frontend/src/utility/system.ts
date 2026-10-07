
/**
 * 停用 WebView 內的預設右鍵選單 (contextmenu)
 */
export function disableContextMenu(event: MouseEvent): void {
  event.preventDefault();
}

/**
 * 攔截鍵盤事件，禁止使用者用 Ctrl/Cmd + +/-/0 進行頁面縮放 (keydown)
 *
 * 適用於 Wails 前端（或其他 webview），搭配 CSS touch-action: none 一起使用，可以大幅壓制桌面端常見的縮放快捷鍵
 */
export function disableZoomKey(event: KeyboardEvent) {

  const isZoomControl = event.ctrlKey || event.metaKey;
  const isZoomKey = event.key === '+' || event.key === '-' || event.key === '0' || event.key === '=';

  if (isZoomControl && isZoomKey) {
    event.preventDefault();
  }
}

/**
 * 攔截滾輪事件，禁止使用者用 Ctrl + 滾輪進行頁面縮放 (wheel)
 *
 * 適用於桌面端瀏覽器 / webview（包含 Wails），搭配 CSS touch-action: none，與鍵盤縮放攔截一起使用，可大幅壓制常見的縮放操作
 */
export function disableCtrlWheel(event: WheelEvent) {
  if (event.ctrlKey) { event.preventDefault(); }
}
