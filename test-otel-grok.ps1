# 1. Generate sample Apache, Cisco ISE, and Cisco IOS log files for testing

@'
127.0.0.1 - frank [10/May/2024:11:22:33 -0700] "GET /apache_pb.gif HTTP/1.0" 200 2326 "http://example.com/start.html" "Mozilla/4.08 [en] (Win98; I ;Nav)"
'@ | Out-File -Encoding ascii apache.log

@'
<134>123: ise01.corp.lan: 2024-05-10 14:13:21 UTC: AuthZ: Passed authentication [User: alice] [Device: 00:11:22:33:44:55] [Policy: Default]
'@ | Out-File -Encoding ascii cisco_ise.log

@'
<189>123: ios-sw01.corp.lan:
'@ | Out-File -Encoding ascii cisco_ios.log

# 2. Start the OTEL Collector
$otelExe = ".\otelcol-custom.exe"
$otelConfig = "collector-config.yaml"

if (Test-Path $otelExe) {
    Write-Host "Starting OTEL Collector with config $otelConfig..." -ForegroundColor Green
    & $otelExe --config $otelConfig
} else {
    Write-Host "otelcol-custom.exe not found! Please check your collector binary." -ForegroundColor Red
    exit 1
}

# 3. (Manual Step): Review logging output in this console
Write-Host "Check the console output above for parsed fields."
Write-Host "Look for extracted fields such as clientip (Apache), User/Device/Policy (ISE), etc."
Write-Host "If you see parsed attributes, your Grok patterns and pipeline are working."

# Note for Unix (Linux/Mac): Convert this script to bash if needed.