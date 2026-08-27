# AI Subtitle Translator (VTT & SRT)

Đây là một ứng dụng dòng lệnh (CLI) được viết bằng Go, có chức năng tự động tìm kiếm và dịch các tệp phụ đề từ tiếng Anh sang tiếng Việt sử dụng AI.

Ứng dụng **hỗ trợ đồng thời cả hai định dạng phụ đề phổ biến là `.vtt` và `.srt`** dựa vào tham số truyền vào khi chạy. Đặc biệt, công cụ hiện đã hỗ trợ **Menu chọn cấu hình API tương tác**, cho phép người dùng dễ dàng chuyển đổi giữa các hãng AI hàng đầu như Google Gemini, OpenAI (ChatGPT), Anthropic (Claude) hoặc OpenRouter.

## Tính năng nổi bật

- **Đa dạng nguồn AI:** Chọn trực tiếp từ Menu tương tác giữa Gemini, ChatGPT, Claude, và OpenRouter mà không cần sửa code.
- **Tự động chọn Model:** Tự động gán model AI tối ưu (nhanh nhất, rẻ nhất) tương ứng với từng nhà cung cấp nếu người dùng không tự thiết lập.
- **Hỗ trợ đa định dạng:** Dịch mượt mà cả tệp WebVTT (`.vtt`) và SubRip (`.srt`).
- **Tự động quét thư mục đệ quy:** Tìm kiếm tất cả các tệp phụ đề tiếng Anh (có đuôi `.en.vtt` hoặc `.en.srt`) trong thư mục hiện tại và tất cả các thư mục con.
- **Bảo toàn cấu trúc gốc:** Giữ nguyên các định dạng gốc như dấu thời gian, số thứ tự, header `WEBVTT` và các dòng trống của tệp phụ đề.
- **Dịch "cuốn chiếu" thông minh:** Tệp phụ đề được chia thành các đoạn nhỏ (mặc định 40 block mỗi đoạn) để gửi lên API, giúp bảo đảm ngữ cảnh, tối ưu giới hạn token và chất lượng dịch thuật.
- **Tiếp tục dịch khi bị gián đoạn (Resume):** Tự động nhận diện và bỏ qua các tệp đã được dịch thành công sang tiếng Việt. Nếu tệp đích tồn tại nhưng chưa có tiếng Việt, ứng dụng sẽ dịch đè lên.
- **Tự động thử lại (Auto-Retry):** Xử lý tốt các tình trạng lỗi mạng, API chập chờn hoặc bị giới hạn tốc độ (rate limit - lỗi 429) với cơ chế tự động chờ và thử lại.
- **Tùy chỉnh linh hoạt:** Cho phép thêm các chỉ dẫn dịch thuật (System Prompt) từ file bên ngoài.

## Cài đặt & Yêu cầu

1. Bạn cần cài đặt [Go](https://go.dev/dl/).
2. Tải mã nguồn về máy:
   ```bash
   git clone <repo_url>
   cd translate
   ```
3. Cài đặt các gói phụ thuộc (nếu có):
   ```bash
   go mod tidy
   ```
4. Build ứng dụng:
   ```bash
   # Build cho Windows
   go build -o wtranslate.exe main.go
   
   # Build cho Linux/macOS
   GOOS="linux" GOARCH="amd64" go build -o ltranslate main.go
   ```

## Hướng dẫn sử dụng

### 1. Cấu hình API Key và Model

Tạo một tệp có tên `api.txt` tại thư mục chứa file thực thi (`wtranslate.exe` hoặc `ltranslate`). Cấu trúc của tệp như sau:
- **Dòng 1:** API Key của bạn (Lấy từ Google AI Studio, OpenAI, Anthropic, hoặc OpenRouter - Tùy thuộc vào việc bạn định chọn hãng nào ở Menu) (Bắt buộc).
- **Dòng 2:** (Tùy chọn) Tên model AI. Nếu để trống, ứng dụng sẽ tự động chọn model mặc định tối ưu nhất tuỳ theo hãng AI mà bạn chọn ở bước sau.

*Lưu ý:* Ứng dụng cũng hỗ trợ lấy API Key qua biến môi trường `API_KEY` hoặc `OPENROUTER_API_KEY`.

### 2. Thêm chỉ dẫn dịch thuật chuyên ngành (Tùy chọn)

Nếu bạn muốn AI chú ý một số quy tắc dịch đặc biệt (ví dụ: giữ nguyên thuật ngữ y khoa, không dịch tên phần mềm...), hãy tạo/sửa tệp `instruction.txt` trong cùng thư mục. Nội dung trong tệp này sẽ được gửi kèm như một hướng dẫn bổ sung cho AI.

### 3. Chạy ứng dụng

Mở Terminal / Command Prompt tại thư mục chứa công cụ và chạy:

**Dịch các tệp `.vtt` (Mặc định):**
```bash
# Trên Windows
.\wtranslate.exe

# Trên Linux
./ltranslate
```

**Dịch các tệp `.srt`:**
```bash
# Trên Windows
.\wtranslate.exe -format srt

# Trên Linux
./ltranslate -format srt
```

**Menu tương tác sẽ xuất hiện:**
Ngay khi chạy, phần mềm sẽ hiển thị một menu để bạn chọn nguồn AI:
```text
======================================
    CHỌN NGUỒN CUNG CẤP AI (API)
======================================
1. Gemini
2. ChatGPT (OpenAI)
3. Claude (Anthropic)
4. OpenRouter
--------------------------------------
Vui lòng nhập số (1-4) và nhấn Enter: 
```
Gõ con số tương ứng với API Key mà bạn đang có và nhấn Enter. Ứng dụng sẽ tự động bắt đầu quét và dịch tệp.

---

## Tiện ích: Chuyển đổi SRT sang VTT hàng loạt bằng FFmpeg

Trong trường hợp bạn có rất nhiều tệp `.srt` và muốn chuyển đổi tất cả sang định dạng `.vtt` trước hoặc sau khi dịch, bạn có thể sử dụng `ffmpeg`.

Lệnh cơ bản cho 1 file:
```bash
ffmpeg -i input.srt output.vtt
```

### Script chuyển đổi đệ quy hàng loạt

Dưới đây là 2 đoạn script (dành cho Windows và Linux/macOS) giúp tự động tìm quét đệ quy (quét cả các thư mục con bên trong) để chuyển đổi toàn bộ file `.srt` sang `.vtt`.

#### Dành cho Windows (PowerShell)

Bạn có thể tạo file `convert_srt_to_vtt.ps1` với nội dung sau:

```powershell
# Quét đệ quy (-Recurse) tất cả các thư mục con để tìm file .srt
Get-ChildItem -Path . -Filter *.srt -Recurse | ForEach-Object {
    $inputFile = $_.FullName
    $outputFile = [System.IO.Path]::ChangeExtension($inputFile, ".vtt")
    
    Write-Host "Đang chuyển đổi: $inputFile -> $outputFile"
    
    # Thực thi ffmpeg
    ffmpeg -loglevel error -y -i "$inputFile" "$outputFile"
    
    # (Tùy chọn) Bỏ comment dòng dưới để xóa file .srt gốc
    # Remove-Item "$inputFile"
}
Write-Host "Hoàn tất!"
```
*Cách chạy:* Chuột phải vào file `.ps1` chọn **Run with PowerShell** hoặc mở cửa sổ PowerShell gõ `.\convert_srt_to_vtt.ps1`.

#### Dành cho Linux / macOS / Git Bash (Bash Script)

Lệnh `find .` mặc định đã tự động quét đệ quy vào tất cả các thư mục con. Bạn tạo file `convert_srt_to_vtt.sh`:

```bash
#!/bin/bash
# Lệnh find mặc định tự động tìm đệ quy vào mọi thư mục con
find . -type f -name "*.srt" | while read -r file; do
    base="${file%.*}"
    echo "Đang chuyển đổi: $file -> $base.vtt"
    ffmpeg -loglevel error -y -i "$file" "$base.vtt"
    
    # rm "$file"
done
echo "Hoàn tất!"
```
