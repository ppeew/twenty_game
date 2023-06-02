cd ~/twenty_game/srv/game_srv
go build
nohup ./game_srv &
sleep 5

cd ~/twenty_game/srv/user_srv
go build
nohup ./user_srv &
sleep 5

cd ~/twenty_game/web/game_web
go build
nohup ./game_web &
sleep 5

cd ~/twenty_game/web/user_web
go build
nohup ./user_web &
sleep 5