
cd ./twenty_game/cmd/Game-user_srv
nohup ./user_srv > ../../us.log &

cd ../Game-game_srv
nohup ./game_srv &

cd ../Game-process_web
nohup ./process_web &

cd ../Game-user_web
nohup ./user_web > ../../uw.log &

cd ../Game-game_web
nohup ./game_web &

cd ../Game-store_web
nohup ./store_web &

cd ../Game-file_web
nohup ./file_web &

ps -ef | grep srv
ps -ef | grep web

cd ../../
cat ./twenty_game/us.log
cat ./twenty_game/uw.log
echo "finish"