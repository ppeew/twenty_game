
cd ~/twenty_game/srv/user_srv
go build
# nohup 后台执行进程
nohup ./user_srv > ../../us.log 2>&1 &

cd ~/twenty_game/srv/game_srv
go build
nohup ./game_srv > ../../gs.log 2>&1 &

cd ~/twenty_game/web/game_web
go build
nohup ./game_web > ../../gw.log 2>&1 &

cd ~/twenty_game/web/user_web
go build
nohup ./user_web > ../../uw.log 2>&1 &

cd ~/twenty_game/web/file_web
go build
nohup ./file_web > ../../fw.log 2>&1 &

cd ~/twenty_game/web/process_web
go build
nohup ./process_web > ../../pw.log 2>&1 &