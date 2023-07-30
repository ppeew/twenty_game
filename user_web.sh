#用户端

cd ~/twenty_game/web/user_web
chmod 766 ./user_web
go build
nohup ./user_web > gw.log 2>&1 &