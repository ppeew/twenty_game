taskkill /F /IM user_srv_win.exe `
/IM game_srv_win.exe `
/IM user_web_win.exe `
/IM hall_web_win.exe `
/IM game_web_win.exe `
/IM process_web_win.exe `
/IM store_web_win.exe `
/IM file_web_win.exe


$currentFolder = Get-Item -Path "." | Select-Object -ExpandProperty Name
if ($currentFolder -eq "twenty_game")
{
    Write-Host "The current folder is 'twenty_game'."
    cd ./cmd/Game-user_srv
}
else
{
    Write-Host "The current folder is not 'twenty_game'."
    cd ../cmd/Game-user_srv
}

go build -o user_srv_win.exe main.go
Start-Process -FilePath "./user_srv_win.exe" -NoNewWindow

cd ../Game-game_srv
go build -o game_srv_win.exe main.go
Start-Process -FilePath "./game_srv_win.exe" -NoNewWindow

cd ../Game-user_web
go build -o user_web_win.exe main.go
Start-Process -FilePath "./user_web_win.exe" -NoNewWindow

cd ../Game-game_web
go build -o game_web_win.exe main.go
Start-Process -FilePath "./game_web_win.exe" -NoNewWindow

cd ../Game-hall_web
go build -o hall_web_win.exe main.go
Start-Process -FilePath "./hall_web_win.exe" -NoNewWindow

cd ../Game-store_web
go build -o store_web_win.exe main.go
Start-Process -FilePath "./store_web_win.exe" -NoNewWindow

cd ../Game-file_web
go build -o file_web_win.exe main.go
Start-Process -FilePath "./file_web_win.exe" -NoNewWindow

cd ../../