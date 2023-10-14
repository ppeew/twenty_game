echo "jenkins-build start"

cd ~/twenty_game/cmd/Game-user_srv
go build -o user_srv main.go
nohup ./user_srv &

cd ~/twenty_game/cmd/Game-game_srv
go build -o game_srv main.go
nohup ./game_srv &

cd ~/twenty_game/cmd/Game-game_web
go build -o game_web main.go
nohup ./game_web &

cd ~/twenty_game/cmd/Game-user_web
go build -o user_web main.go
nohup ./user_web &

cd ~/twenty_game/cmd/Game-file_web
go build -o file_web main.go
nohup ./file_web &

cd ~/twenty_game/cmd/Game-process_web
go build -o process_web main.go
nohup ./process_web &

cd ~/twenty_game/cmd/Game-store_web
go build -o store_web main.go
nohup ./store_web &

cd ~/twenty_game/cmd/Game-hall_web
go build -o hall_web main.go
nohup ./hall_web &

cd ~/twenty_game/cmd/Game-admin_web
go build -o admin_web main.go
nohup ./admin_web &

echo "jenkis-build finish"