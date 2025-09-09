@echo off
echo 测试图片上传接口
echo.

set IMAGE_PATH=f:\素材图\桌面壁纸\ChMkJlf8lkmIabp2AAdCEFWIQCcAAWzfgEWv9EAB0Io449.jpg
set UPLOAD_URL=http://localhost:8080/api/file/formimage

echo 文件路径: %IMAGE_PATH%
echo 上传接口: %UPLOAD_URL%
echo.

if not exist "%IMAGE_PATH%" (
    echo 错误: 文件不存在 - %IMAGE_PATH%
    echo 请修改批处理文件中的文件路径
    pause
    exit /b 1
)

echo 开始上传...
curl -X POST "%UPLOAD_URL%" -F "file=@%IMAGE_PATH%" -H "Accept: application/json"

echo.
echo 测试完成，请检查服务器日志输出。
pause