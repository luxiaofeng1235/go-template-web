# 上传接口测试脚本
# 使用方法: .\test_upload.ps1

$imagePath = "f:\素材图\桌面壁纸\ChMkJlf8lkmIabp2AAdCEFWIQCcAAWzfgEWv9EAB0Io449.jpg"
$uploadUrl = "http://localhost:8080/api/file/formimage"

# 检查文件是否存在
if (!(Test-Path $imagePath)) {
    Write-Host "错误: 文件不存在 - $imagePath" -ForegroundColor Red
    Write-Host "请修改脚本中的文件路径"
    exit 1
}

Write-Host "开始测试图片上传..." -ForegroundColor Green
Write-Host "文件路径: $imagePath"
Write-Host "上传接口: $uploadUrl"

try {
    # 使用 Invoke-RestMethod 上传文件
    $response = Invoke-RestMethod -Uri $uploadUrl -Method Post -InFile $imagePath -ContentType "multipart/form-data"
    
    Write-Host "上传成功!" -ForegroundColor Green
    Write-Host "响应内容:" -ForegroundColor Yellow
    $response | ConvertTo-Json -Depth 3
    
} catch {
    Write-Host "上传失败!" -ForegroundColor Red
    Write-Host "错误信息: $($_.Exception.Message)"
    
    # 显示详细错误信息
    if ($_.Exception.Response) {
        $statusCode = $_.Exception.Response.StatusCode
        Write-Host "HTTP状态码: $statusCode"
    }
}

Write-Host "`n测试完成，请检查服务器日志输出。"