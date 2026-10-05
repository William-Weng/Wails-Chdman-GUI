<script lang="ts">
  import { onMount } from "svelte";
  import { Events } from "@wailsio/runtime";
  import { ParseFilePath, ExtractCD, CreateCD } from "../bindings/wails-chdman-gui/backend/chdmanservice";
  import { dialog } from "./utility/dialog";

  let isLoading = false;
  let imageName = "Empty.png";

  onMount(() => {

    const unsubscribeDrop = Events.On("image-file-dropped", (event) => {
      void imageFileDroppedAction(event.data as ImageDroppedData);
    });

    window.addEventListener("contextmenu", disableContextMenu);

    return () => {
      unsubscribeDrop();
      window.removeEventListener("contextmenu", disableContextMenu);
    };
  });

  async function imageFileDroppedAction(data: ImageDroppedData): Promise<void> {

    if (isLoading) { return; }

    const info = await ParseFilePath(data.path);
    const extension = info.ext.toLowerCase();

    imageName = "Empty.png";

    if (extension == ".chd") {
      extractCD(data.path); return;
    }

    if (extension == ".iso") {
      extractCD(data.path); return;
    }

    if (extension == ".cue") {
      createCD(data.path); return;
    }
  }

  async function extractCD(filePath: string) {

    isLoading = true;
    imageName = "DVD.png";

    try {
      const messsage = await ExtractCD(filePath);
      isLoading = false;
      await dialog("info", "還原完成", messsage);
    } catch (err) {
      const error = err instanceof Error ? err.message : String(err);
      isLoading = false;
      await dialog("warning", "還原失敗", error);
    }
  }

  async function createCD(filePath: string) {

    imageName = "PCSX2.png";
    isLoading = true;

    try {
      const messsage = await CreateCD(filePath);
      isLoading = false;
      await dialog("info", "壓縮完成", messsage);
    } catch (err) {
      const error = err instanceof Error ? err.message : String(err);
      isLoading = false;
      await dialog("warning", "壓縮失敗", error);
    }
  }

  function disableContextMenu(event: any) {
      event.preventDefault();
  }
</script>

<main>
    <section data-file-drop-target>
        <img src={imageName} alt="輸入圖片預覽" />
        {#if isLoading}
            <img class="loading-gif" src="Loading.gif" alt="處理中"/>
        {/if}
    </section>
</main>
