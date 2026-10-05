import { Dialogs } from "@wailsio/runtime";

/**
   * 顯示原生系統對話框
   *
   *  - 依照傳入的 type 決定要顯示一般提示、錯誤或警告對話框；title 會顯示在對話框標題，message 則是對話框內容
   *
   * @param type 對話框類型：info、error 或 warning
   * @param title 對話框標題
   * @param message 對話框顯示的訊息內容
   */
export async function dialog(type: DialogType, title: string, message: string) {

    switch (type) {
        case "info": await Dialogs.Info({ Title: title, Message: message }); break;
        case "error": await Dialogs.Error({ Title: title, Message: message }); break;
        case "warning": await Dialogs.Warning({ Title: title, Message: message }); break;
        default: await Dialogs.Info({ Title: title, Message: message });
    }
}
