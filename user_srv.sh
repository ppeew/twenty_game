#用户服务

cd ~/twenty_game/srv/user_srv
chmod 766 ./user_srv
go build
nohup ./user_srv > gw.log 2>&1 &
