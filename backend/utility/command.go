package utility

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

// 根據 Context、可執行檔路徑與參數，建立外部指令物件
//
// 此函式只負責驗證可執行檔是否存在，以及建立 *exec.Cmd；
// 不會立刻執行指令。呼叫端可選擇使用 cmd.Run()、cmd.Output()、
// cmd.CombinedOutput()，或 cmd.Start() 後再呼叫 cmd.Wait() 執行指令
//
// 參數：
//   - ctx：控制指令生命週期的 Context
//   - path：可執行檔名稱或完整路徑，例如 "go"、"ffmpeg"、"/opt/homebrew/bin/ffmpeg"
//   - args：傳遞給外部指令的參數列表，不包含 path 本身；例如執行 `go version` 時，args 為 []string{"version"}。
//
// 回傳值：
//   - *exec.Cmd：已建立但尚未執行的指令物件
//   - error：當 path 為空字串，或系統找不到對應的可執行檔時回傳錯誤
func CreateCommandContext(ctx context.Context, path string, args []string) (*exec.Cmd, error) {

	if path == "" {
		return nil, fmt.Errorf("command path cannot be empty")
	}

	executablePath, err := exec.LookPath(path)
	if err != nil {
		return nil, fmt.Errorf("command not found: %q: %w", path, err)
	}

	cmd := exec.CommandContext(ctx, executablePath, args...)
	cmd = HideWindow(cmd)

	return cmd, nil
}

// 建立可取消的 Context
//
// 參數：
//   - timeout：逾時時間。傳入 0 或負數代表不設逾時限制
//
// 回傳值：
//   - ctx：可傳給支援 context 的函式，例如 exec.CommandContext
//   - cancel：手動取消 ctx 的函式；使用完畢後應呼叫它以釋放相關資源
func CreateCancelContext(timeout time.Duration) (ctx context.Context, cancel context.CancelFunc) {

	if timeout > 0 {
		return context.WithTimeout(context.Background(), timeout)
	}

	return context.WithCancel(context.Background())
}

// 建立永不取消、沒有 deadline 的根 Context
//
// 回傳值：
//   - context.Context：沒有 timeout、不能手動取消的背景 Context；適合短暫測試或不需要取消機制的操作。
func CreateContext() context.Context {
	return context.Background()
}

// 執行外部指令，並即時串流讀取 stdout 與 stderr
//
// 每讀取到一段以 "\n"、"\r" 或 "\r\n" 結尾的輸出時，會呼叫 onOutput。callback 收到的 text 會保留原始結尾字元
//
// 參數：
//   - cmd：要執行的外部指令
//   - onOutput：每次讀取到輸出片段時呼叫。source 為 "stdout" 或 "stderr"；text 包含原始結尾字元，例如 "\n"、"\r" 或 "\r\n"
//
// 回傳值：
//   - error：建立 pipe、啟動程序、讀取輸出或等待程序結束失敗時回傳錯誤
func RunCommandStream(cmd *exec.Cmd, onOutput func(source string, text string)) error {

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	var group sync.WaitGroup
	errChannel := make(chan error, 2)

	// 讀取單一輸出串流，直到 EOF 或讀取錯誤
	readStream := func(source string, reader io.Reader) {

		defer group.Done()

		bufferedReader := bufio.NewReader(reader)

		for {
			text, readErr := ReadLineEnding(bufferedReader)

			if text != "" && onOutput != nil {
				onOutput(source, text)
			}

			if readErr == io.EOF {
				return
			}

			if readErr != nil {
				errChannel <- fmt.Errorf("read %s: %w", source, readErr)
				return
			}
		}
	}

	group.Add(2)

	// 必須並行讀取，避免其中一個 pipe 的緩衝區被寫滿。
	go readStream("stdout", stdout)
	go readStream("stderr", stderr)

	// 必須先讀完 pipes，再呼叫 cmd.Wait()。
	group.Wait()
	close(errChannel)

	// 收集 pipe 讀取錯誤。
	var readErr error
	for err := range errChannel {
		if err != nil {
			readErr = err
			break
		}
	}

	// 即使讀取失敗，仍必須 Wait，避免遺留子程序與系統資源。
	waitErr := cmd.Wait()

	if readErr != nil {
		return readErr
	}

	if waitErr != nil {
		return fmt.Errorf("command failed: %w", waitErr)
	}

	return nil
}

// 從 reader 讀取一段輸出，直到遇到 "\n"、"\r" 或 "\r\n"
//
// 回傳的字串會保留結尾字元：
//   - "hello\n"
//   - "hello\r"
//   - "hello\r\n"
//
// 若讀到 EOF 前仍有尚未結尾的資料，會回傳該資料與 io.EOF
func ReadLineEnding(reader *bufio.Reader) (string, error) {

	var result []byte

	for {
		b, err := reader.ReadByte()
		if err != nil {
			if len(result) > 0 {
				return string(result), err
			}

			return "", err
		}

		result = append(result, b)

		switch b {
		case '\n':
			return string(result), nil
		case '\r':
			next, err := reader.Peek(1)
			if err != nil {
				if err == io.EOF {
					return string(result), nil
				}
				return string(result), err
			}

			if next[0] == '\n' {
				_, _ = reader.ReadByte()
				result = append(result, '\n')
			}

			return string(result), nil
		}
	}
}
