# ==== OpenTelemetry Grok Log Validation Script ====
$EnableDebug = $true

$TestDir      = ".\test"
$OtelExe      = ".\otelcol-custom.exe"
$OtelConfig   = "collector-config.yaml"
$LogExport    = "otel-test-output.log"
$LogError     = "otel-test-error.log"
$TestCasesCsv = "$TestDir\test_cases.csv"

# Ensure the test directory exists
if (-not (Test-Path $TestDir)) {
    New-Item -ItemType Directory -Path $TestDir | Out-Null
}

# Ensure test cases CSV exists
if (-not (Test-Path $TestCasesCsv)) {
    Write-Host "ERROR: Test case file not found: $TestCasesCsv" -ForegroundColor Red
    exit 2
}

# Clean up any prior log files in the test dir
Get-ChildItem "$TestDir\*.log" -ErrorAction SilentlyContinue | Remove-Item -Force
Remove-Item $LogExport, $LogError -ErrorAction SilentlyContinue

# Import test cases and only process non-empty valid records
$TestCases = Import-Csv $TestCasesCsv | Where-Object {
    $_.LogFile -and $_.LogLine -and ($_.LogFile -ne $null) -and ($_.LogLine -ne $null) -and ($_.LogFile.Trim() -ne "") -and ($_.LogLine.Trim() -ne "")
}

if ($EnableDebug) {
    Write-Host "=== Parsed Test Cases (filtered) ==="
    $TestCases | Format-Table
    Write-Host "===================================="
}

# Write sample logs for each valid test case
foreach ($case in $TestCases) {
    $LogFile = Join-Path $TestDir $case.LogFile
    if ($EnableDebug) {
        Write-Host ("Writing test log: LogFile='{0}', LogLine='{1}'" -f $LogFile, $case.LogLine)
    }
    Add-Content -Path $LogFile -Value $case.LogLine
    Write-Host ("Added test log to {0}: {1}" -f $LogFile, $case.LogLine)
}

Write-Host "`nLaunching OpenTelemetry Collector...`n" -ForegroundColor Cyan
$p = Start-Process -FilePath $OtelExe -ArgumentList "--config", $OtelConfig `
    -RedirectStandardOutput $LogExport -RedirectStandardError $LogError -NoNewWindow -PassThru
Start-Sleep -Seconds 8

if ($null -ne $p) {
    if (Get-Process -Id $p.Id -ErrorAction SilentlyContinue) {
        try { Stop-Process -Id $p.Id -Force }
        catch { Write-Host "Collector process could not be stopped or was already exited." -ForegroundColor Yellow }
    } else {
        Write-Host "Collector process was already exited." -ForegroundColor Yellow
    }
}

# Merge stdout and stderr
if (Test-Path $LogError) { Get-Content $LogError | Add-Content $LogExport; Remove-Item $LogError }

$AllPassed = $true
$Counter   = 0

foreach ($case in $TestCases) {
    $Counter++
    $ReqFields = $case.RequiredFields -split ','

    Write-Host ("`nTest #{0}: {1}" -f $Counter, $case.LogFile) -ForegroundColor Yellow
    Write-Host ("Log: {0}" -f $case.LogLine)
    $Fail = $false
    foreach ($field in $ReqFields) {
        if ($EnableDebug) { Write-Host "Looking for: $field" -ForegroundColor Gray }
        if (-not (Select-String -Path $LogExport -SimpleMatch -Pattern "$field")) {
            Write-Host ("  [FAIL] Field missing: {0}" -f $field) -ForegroundColor Red
            $Fail = $true
        } else {
            Write-Host ("  [ OK ] Field found: {0}" -f $field) -ForegroundColor Green
        }
    }
    if ($Fail) {
        Write-Host "[FAIL] Test #$Counter FAILED" -ForegroundColor Red
        $AllPassed = $false
    } else {
        Write-Host "[ OK ] Test #$Counter PASSED" -ForegroundColor Green
    }
}

if ($AllPassed) {
    Write-Host "`nALL LOG PARSING TESTS PASSED!" -ForegroundColor Green
    exit 0
} else {
    Write-Host "`nSome tests FAILED. Check $LogExport for details." -ForegroundColor Red
    exit 1
}