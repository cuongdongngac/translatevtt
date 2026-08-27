#!/bin/bash
# Chuyển đổi đệ quy tất cả các file .srt sang .vtt

# Tìm tất cả các file .srt
find . -type f -name "*.srt" | while read -r file; do
    # Lấy đường dẫn file bỏ đi phần đuôi .srt
    base="${file%.*}"
    
    echo "Đang chuyển đổi: $file -> $base.vtt"
    
    # Dùng ffmpeg để chuyển đổi (ẩn bớt log bằng -loglevel error)
    ffmpeg -loglevel error -y -i "$file" "$base.vtt"
    
    # (Tùy chọn) Bỏ comment dòng dưới đây nếu bạn muốn xóa file .srt gốc sau khi chuyển đổi xong
    # rm "$file"
done

echo "Hoàn tất chuyển đổi!"
