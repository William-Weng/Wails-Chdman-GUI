<script lang="ts">
  import { onMount } from "svelte";
  import { Events } from "@wailsio/runtime";
  import { ParseFilePath, ExtractCD, CreateCD } from "../bindings/wails-chdman-gui/backend/chdmanservice";
  import { dialog } from "./utility/dialog";

  let isLoading = false;
  let imageName = "Empty.png";
  let progress = "0.0 %";

  onMount(() => {

    const unsubscribeDrop = Events.On("image-file-dropped", (event) => {
      void imageFileDroppedAction(event.data as ImageDroppedData);
    });

    const unsubscribeExtract = Events.On("extract-progress", (event) => {
      void parseExtractProgress(event.data as ExtractProgress);
    });

    const unsubscribeCreate = Events.On("create-progress", (event) => {
      void parseCreateProgress(event.data as CreateProgress);
    });

    window.addEventListener("contextmenu", disableContextMenu);

    return () => {
      unsubscribeDrop();
      unsubscribeCreate();
      unsubscribeExtract();
      window.removeEventListener("contextmenu", disableContextMenu);
    };
  });

  /**
   * 根據拖放檔案的副檔名，判斷要還原 CHD 或從 CUE 建立 CHD
   */
  async function imageFileDroppedAction(data: ImageDroppedData): Promise<void> {

    if (isLoading) { return; }

    const info = await ParseFilePath(data.path);
    const extension = info.ext.toLowerCase();

    imageName = "Empty.png";

    try {
      switch (extension) {
        case ".chd": case ".iso": await extractCD(data.path); break;
        case ".cue": await createCD(data.path); break;
        default: await dialog("warning", "不支援的檔案", `目前不支援「${extension || "無副檔名"}」格式。`); break;
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      reset();
      await dialog("warning", "檔案讀取失敗", message);
    }
  }

  /**
   * 將 CHD 或 ISO 還原為 CUE/BIN 等檔案
   */
  async function extractCD(filePath: string) {

    isLoading = true;
    imageName = "DVD.png";

    try {
      const messsage = await ExtractCD(filePath);
      reset();
      await dialog("info", "還原完成", messsage);
    } catch (err) {
      const error = err instanceof Error ? err.message : String(err);
      reset();
      await dialog("warning", "還原失敗", error);
    } finally {
      reset();
    }
  }

  /**
   * 從 CUE/BIN 建立 CHD
   */
  async function createCD(filePath: string) {

    imageName = "PCSX2.png";
    isLoading = true;

    try {
      const messsage = await CreateCD(filePath);
      reset();
      await dialog("info", "壓縮完成", messsage);
    } catch (err) {
      const error = err instanceof Error ? err.message : String(err);
      reset();
      await dialog("warning", "壓縮失敗", error);
    } finally {
      reset();
    }
  }

  /**
   * 接收後端 `extract-progress` 事件並更新解壓縮進度
   */
  function parseExtractProgress(data: ExtractProgress): void {
    progress = data.progress
  }

  /**
   * 接收後端 `create-progress` 事件並更新壓縮進度及壓縮比例
   */
  function parseCreateProgress(data: CreateProgress): void {
    progress = data.progress
  }

  /**
   * 停用 WebView 內的預設右鍵選單
   */
  function disableContextMenu(event: MouseEvent): void {
    event.preventDefault();
  }

  /**
   * 將畫面狀態還原為閒置狀態
   */
  function reset(): void {
    progress = "0.0 %"
    isLoading = false;
  }
</script>

<main>
    <section data-file-drop-target>
        <img src={imageName} alt="輸入圖片預覽" />
        {#if isLoading}
            <img class="loading-gif" src="Loading.gif" alt="處理中"/>
            <div class="loading-progress">{progress}</div>
        {/if}
    </section>
</main>
