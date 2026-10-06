/**
 * 支援的原生系統對話框類型
 * - info：一般資訊或操作成功
 * - error：操作失敗或發生錯誤
 * - warning：提醒使用者注意或補齊必要輸入
 */
type DialogType = "info" | "error" | "warning";

/**
 * 圖片拖放區域的 ID
 * - leftSlot：左側外框圖
 * - rightSlot：右側內容圖
 */
type SlotId = "leftSlot" | "rightSlot";

/**
 * 從 Wails「image-file-dropped」事件取得的圖片拖放資料
 * - path：圖片在本機檔案系統中的完整路徑
 */
type ImageDroppedData = {
  path: string;
};

/**
 * CHD 解壓縮（Extracting）時的進度資料
 *
 * - progress：進度
 */
type ExtractProgress = {
  progress: string;
}

/**
 * CHD 建立／壓縮（Compressing）時的進度資料
 *
 * - progress：進度
 * - ratio：壓縮比例
 */
type CreateProgress = {
  progress: string;
  ratio: string;
}
