# 测试API接口
Write-Host "测试API接口..." -ForegroundColor Green

# 测试产品列表
Write-Host "1. 测试产品列表 GET /api/product/list" -ForegroundColor Yellow
Invoke-RestMethod -Uri "http://localhost:8080/api/product/list" -Method GET

Write-Host "`n2. 测试用户注册 POST /api/user/register" -ForegroundColor Yellow
$body = @{
    username = "testuser"
    password = "123456"
    email = "test@example.com"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/user/register" -Method POST -Body $body -ContentType "application/json"

Write-Host "`n3. 测试管理后台首页 GET /admin/" -ForegroundColor Yellow
Invoke-RestMethod -Uri "http://localhost:8081/admin/" -Method GET