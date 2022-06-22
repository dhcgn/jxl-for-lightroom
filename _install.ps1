$file = "jxl-for-lightroom.exe"
go build  -ldflags="-s -w" -o $file  .\cmd\jxl-for-lightroom\
Write-Host ("Site: {0:N1} MB" -f (gci $file | %{$_.Length / 1mb}))

$targetLightroomFolder = Join-Path $env:APPDATA  'Adobe\Lightroom\Modules\jxl-for-lightroom\'
if (Test-Path $targetLightroomFolder) {
    Write-Host "Lightroom module folder already exists"
} else {
    Write-Host "Creating Lightroom module folder"
    New-Item -Path $targetLightroomFolder -ItemType Directory
}

Copy-Item $file $targetLightroomFolder