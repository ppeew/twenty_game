
cd ~/twenty_game/cmd/user_srv
go build
nohup ./user_srv > ../../us.log 2>&1 &

cd ~/twenty_game/cmd/game_srv
go build
nohup ./game_srv > ../../gs.log 2>&1 &

cd ~/twenty_game/cmd/game_web
go build
nohup ./game_web > ../../gw.log 2>&1 &

cd ~/twenty_game/cmd/user_web
go build
nohup ./user_web > ../../uw.log 2>&1 &

cd ~/twenty_game/cmd/file_web
go build
nohup ./file_web > ../../fw.log 2>&1 &

cd ~/twenty_game/cmd/process_web
go build
nohup ./process_web > ../../pw.log 2>&1 &

cd ~/twenty_game/cmd/store_web
go build
nohup ./store_web > ../../sw.log 2>&1 &
