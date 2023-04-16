cd ./twenty_game/srv/user_srv
go build
./user_srv &

cd -
cd ./twenty_game/web/user_web
go build
./user_web &

cd -
cd ./twenty_game/srv/game_srv
go build
./game_srv &

cd -
cd ./twenty_game/web/game_web
go build
./game_web &