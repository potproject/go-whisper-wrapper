# Bash 
## Windows
```
GOOS=windows GOARCH=amd64 go build -o mp3_to_srt_whisper_windows_x64.exe .\main.go
```

## Linux
```
GOOS=linux GOARCH=amd64 go build -o mp3_to_srt_whisper_linux_x64 .\main.go
```

## Mac(Intel)
```
GOOS=darwin GOARCH=amd64 go build -o mp3_to_srt_whisper_mac_x64 .\main.go
```

## Mac(M1)
```
GOOS=darwin GOARCH=arm64 go build -o mp3_to_srt_whisper_mac_m1 .\main.go
```

# Powershell

## Windows
```
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -o mp3_to_srt_whisper_windows_x64.exe .\main.go
```

## Linux
```
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o mp3_to_srt_whisper_linux_x64 .\main.go
```

## Mac(Intel)
```
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o mp3_to_srt_whisper_mac_x64 .\main.go
```

## Mac(M1)
```
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o mp3_to_srt_whisper_mac_m1 .\main.go
```