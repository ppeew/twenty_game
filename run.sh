cd /usr/local/go_pro/twenty_game/srv/user_srv
go build
./user_srv

cd ../game_srv
go build
./game_srv

cd /usr/local/go_pro/twenty_game/web/user_web
go build
./user_web

cd ../game_web
go build
./game_web