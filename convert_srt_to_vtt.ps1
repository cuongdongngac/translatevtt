# Chuyển đổi đệ quy tất cả các file .srt sang .vtt trong PowerShell
Get-ChildItem -Path . -Filter *.srt -Recurse | ForEach-Object {
    $inputFile = $_.FullName
    $outputFile = [System.IO.Path]::ChangeExtension($inputFile, ".vtt")
    
    Write-Host "Đang chuyển đổi: $inputFile -> $outputFile"
    
    # Thực thi ffmpeg để chuyển đổi (ẩn bớt log)
    ffmpeg -loglevel error -y -i "$inputFile" "$outputFile"
    
    # (Tùy chọn) Bỏ comment dòng dưới đây nếu bạn muốn xóa file .srt gốc
    # Remove-Item "$inputFile"
}

Write-Host "Hoàn tất chuyển đổi!"
