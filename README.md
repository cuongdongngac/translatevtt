# VTT/SRT Subtitle Translator

Đây là một ứng dụng dòng lệnh (CLI) được viết bằng Go, tự động tìm và dịch các tệp phụ đề (định dạng `.vtt` hoặc `.srt`) từ tiếng Anh sang tiếng Việt sử dụng AI thông qua OpenRouter API.

## Tính năng nổi bật

- **Tự động quét thư mục:** Tìm kiếm đệ quy toàn bộ thư mục hiện tại để tìm các tệp phụ đề tiếng Anh có đuôi `.en.vtt` hoặc `.en.srt`.
- **Dịch tự động:** Dịch phụ đề sang tiếng Việt và lưu dưới dạng `.vi.vtt` hoặc `.vi.srt`.
- **Bảo toàn định dạng:** Giữ nguyên các định dạng gốc (thời gian, dấu hiệu WEBVTT, v.v.) của tệp phụ đề gốc.
- **Tiếp tục dịch thông minh:** Tự động bỏ qua các tệp đã được dịch (nhận diện bằng cách kiểm tra các ký tự tiếng Việt trong tệp đích). Nếu tệp đích tồn tại nhưng không có tiếng Việt, ứng dụng sẽ ghi đè dịch lại.
- **Chia nhỏ để dịch:** Chia tệp phụ đề ra nhiều đoạn nhỏ (mặc định 40 blocks/lần) để đảm bảo chất lượng dịch thuật của AI và tránh giới hạn token.
- **Tùy chỉnh model AI:** Mặc định sử dụng `google/gemini-flash-1.5`, nhưng bạn có thể dễ dàng thay đổi sang các model khác của OpenRouter.
- **Xử lý lỗi tự động:** Tích hợp tính năng tự động thử lại (retry) khi gặp lỗi mạng, lỗi API hoặc bị rate limit (lỗi 429).
- **Hỗ trợ thêm chỉ dẫn ngữ cảnh:** Cho phép cung cấp thêm chỉ dẫn dịch thuật (context/instruction) cho AI qua tệp văn bản.

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
   go build -o vtt-translator.exe main.go
   
   # Build cho Linux/macOS
   go build -o translator main.go
   ```

## Hướng dẫn sử dụng

### 1. Cấu hình API Key và Model

Tạo một tệp có tên `api.txt` tại thư mục chứa file thực thi. Cấu trúc của tệp như sau:
- **Dòng 1:** OpenRouter API Key của bạn (Bắt buộc).
- **Dòng 2:** (Tùy chọn) Tên model AI trên OpenRouter bạn muốn sử dụng. Nếu để trống, mặc định sẽ là `google/gemini-flash-1.5`.

*Lưu ý:* Bạn cũng có thể cung cấp API Key qua biến môi trường `OPENROUTER_API_KEY`. Tuy nhiên, nếu bạn chạy bằng file exe trong thư mục có file `api.txt`, nó sẽ ưu tiên đọc file `api.txt`.

### 2. Thêm chỉ dẫn dịch thuật (Tùy chọn)

Nếu bạn muốn AI chú ý một số quy tắc dịch đặc biệt (ví dụ: không dịch tên phần mềm, giữ nguyên thuật ngữ y khoa...), hãy tạo/sửa tệp `instruction.txt` trong cùng thư mục. Nội dung trong tệp này sẽ được gửi kèm như một hướng dẫn bổ sung (system prompt) cho AI.

### 3. Chạy ứng dụng

Để dịch các tệp `.vtt` (Mặc định):
```bash
./vtt-translator.exe
```

Để dịch các tệp `.srt`:
```bash
./vtt-translator.exe -format srt
```

Ứng dụng sẽ bắt đầu quét từ thư mục hiện tại của bạn và dịch tất cả các tệp `.en.vtt` hoặc `.en.srt` nó tìm thấy.

## Logic đặt tên file

- File gốc cần dịch: `<tên file>.en.vtt` hoặc `<tên file>.en.srt`
- File kết quả: `<tên file>.vi.vtt` hoặc `<tên file>.vi.srt`
