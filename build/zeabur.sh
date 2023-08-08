nohup ./twenty_game/cmd/Game-user_srv/user_srv > ./twenty_game/us.log &
nohup ./twenty_game/cmd/Game-game_srv/game_srv &
nohup ./twenty_game/cmd/Game-process_web/process_web &
nohup ./twenty_game/cmd/Game-user_web/user_web > ./twenty_game/uw.log &
nohup ./twenty_game/cmd/Game-game_web/game_web &
nohup ./twenty_game/cmd/Game-store_web/store_web &
nohup ./twenty_game/cmd/Game-file_web/file_web &

netstat -ntlup | grep srv
netstat -ntlup | grep web

cat ./twenty_game/us.log
cat ./twenty_game/uw.log
echo "finish"