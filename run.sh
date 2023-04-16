cd ./srv/user_srv
go build
./user_srv &

cd -
cd ./web/user_web
go build
./user_web &

cd -
cd ./srv/game_srv
go build
./game_srv &

cd -
cd .web/game_web
go build
./game_web &