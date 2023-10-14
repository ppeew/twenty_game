echo "jenkins-build start"

cd ~/twenty_game/cmd/Game-user_srv
go build -o user_srv main.go
nohup ./user_srv > /dev/null 2>&1 &

cd ~/twenty_game/cmd/Game-game_srv
go build -o game_srv main.go
nohup ./game_srv > /dev/null 2>&1 &

cd ~/twenty_game/cmd/Game-game_web
go build -o game_web main.go
nohup ./game_web > /dev/null 2>&1 &

cd ~/twenty_game/cmd/Game-user_web
go build -o user_web main.go
nohup ./user_web > /dev/null 2>&1 &

cd ~/twenty_game/cmd/Game-file_web
go build -o file_web main.go
nohup ./file_web > /dev/null 2>&1 &

cd ~/twenty_game/cmd/Game-process_web
go build -o process_web main.go
nohup ./process_web > /dev/null 2>&1 &

cd ~/twenty_game/cmd/Game-store_web
go build -o store_web main.go
nohup ./store_web > /dev/null 2>&1 &

cd ~/twenty_game/cmd/Game-hall_web
go build -o hall_web main.go
nohup ./hall_web > /dev/null 2>&1 &

cd ~/twenty_game/cmd/Game-admin_web
go build -o admin_web main.go
nohup ./admin_web > /dev/null 2>&1 &

sleep 2s
echo "jenkis-build finish"