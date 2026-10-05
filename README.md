# [Wails-Chdman-GUI](https://github.com/William-Weng?tab=repositories&q=wails)

![Go](https://img.shields.io/badge/Go-1.27.1-00ADD8?logo=go&logoColor=white)
![Node.js](https://img.shields.io/badge/Node.js-20.19.2-339933?logo=nodedotjs&logoColor=white)
![Wails](https://img.shields.io/badge/Wails-v3.0.0--beta.26-DF0000?logo=wails&logoColor=white)
![LICENSE](https://img.shields.io/github/license/William-Weng/Wails-Chdman-GUI?style=flat&label=LICENSE&color=yellow)
![Tag](https://img.shields.io/github/v/tag/William-Weng/Wails-Chdman-GUI?style=flat&label=Tag)
![Stars](https://img.shields.io/github/stars/William-Weng/Wails-Chdman-GUI?style=flat&label=Stars)

一個使用 **Wails + Go + Svelte** 開發的 CHD 轉檔桌面工具。

本專案透過 [`chdman`](https://chdman.com/) 指令執行 CHD（Compressed Hunks of Data）相關轉檔工作，提供圖形化介面，讓使用者不必手動輸入長串終端機指令即可處理遊戲映像檔。

> 本工具本身不包含 `chdman`。執行轉檔前，請先在系統中安裝 `chdman`，並確保它可從終端機的 `PATH` 找到。

## 功能

- 以圖形介面選擇或拖放來源檔案
- 呼叫系統中的 `chdman` 執行轉檔
- 即時顯示 `chdman` 的標準輸出與錯誤輸出
- 支援長時間轉檔工作
- 轉檔時顯示 loading 狀態
- 可保留並顯示 CLI 原始輸出，包括 `\n`、`\r` 與 `\r\n`
- 使用 Go `context.Context` 管理外部程序生命週期
- 支援後續加入取消轉檔與進度條功能

## [建置指令整理](https://v3.wails.io/zh-tw/guides/build/building/)

| 目的 | 指令 |
| --- | --- |
| 建立新專案 | `wails3 init -n <專案名稱> -t <前端框架>` |
| 產生 bindings | `wails3 generate bindings` |
| 更新 Windows 與 macOS 專用圖示檔案 | `wails3 generate icons -input build/appicon.png -windowsfilename build/windows/icon.ico -macfilename build/darwin/icons.icns` |
| 更新建置資源 | `wails3 update build-assets -config build/config.yml -dir build` |
| 拉取（下載）用於跨平台交叉編譯的 Docker 映像檔 | `wails3 task setup:docker` |
| 建置 Windows x64 | `wails3 build GOOS=windows GOARCH=amd64` |
| 打包 macOS arm64 | `wails3 package GOOS=darwin GOARCH=arm64` |
| 打包 Linux arm64 | `wails3 build GOOS=linux GOARCH=arm64` |

## 安裝 chdman

| 平臺 | 位置 |
| --- | --- |
| Windows | `C:\Tools\MAME\chdman.exe` |
| macOS | `/opt/homebrew/bin/chdman` |
| Linux | `/opt/homebrew/bin/chdman` |

請先在終端機確認：

```bash
chdman -help
```

或：

```bash
chdman --help
```

若可以看到說明文字，代表 `chdman` 已安裝完成並且可供本工具使用。

### macOS

可透過 Homebrew 安裝提供 `chdman` 的套件：

```bash
brew install mame
```

安裝完成後確認：

```bash
chdman -help
```

Apple Silicon Mac 若找不到 `chdman`，可確認 Homebrew 路徑是否已加入 `PATH`：

```bash
echo $PATH
which chdman
```

常見位置包括：

```text
/opt/homebrew/bin/chdman
```

Intel Mac 常見位置：

```text
/usr/local/bin/chdman
```

### Windows

可安裝 MAME，或下載包含 `chdman.exe` 的 MAME 發行版本。

安裝後請將 `chdman.exe` 所在資料夾加入 Windows 的 `PATH` 環境變數，例如：

```text
C:\Tools\MAME
```

重新開啟 PowerShell 或 Command Prompt 後確認：

```powershell
chdman.exe -help
```

### Linux

不同發行版的套件名稱可能不同，通常可安裝 MAME 相關套件：

```bash
sudo apt install mame-tools
```

或：

```bash
sudo dnf install mame-tools
```

安裝後確認：

```bash
chdman -help
```

> 各平台的 MAME 套件名稱與是否內含 `chdman` 可能不同。若指令仍找不到，請確認 MAME / CHD tools 的安裝位置是否已加入 `PATH`。

## macOS 注意事項

目前 macOS 版本尚未使用 Apple Developer [證書簽署](https://www.cnblogs.com/Flat-White/p/17153264.html)及公證。

首次執行前，請在終端機執行：

```bash
xattr -cr "/Applications/Chdman-GUI.app"
```

如果應用程式不在 `/Applications`，請將路徑替換成實際位置。

## 注意事項

- 本專案不提供 ROM、BIOS、CHD、ISO、BIN/CUE 或任何受版權保護的遊戲檔案。
- 請只處理你擁有、備份或有合法使用權的資料。
- `chdman` 的轉檔速度與相容性取決於來源檔案格式、磁碟速度、CPU，以及使用的 `chdman` 版本。
- 輸出檔案可能覆寫既有檔案；使用前請確認輸出路徑。
- 轉檔期間請保留足夠磁碟空間，尤其是在建立或還原大型 CHD 時。
