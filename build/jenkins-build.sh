cd /var/jenkins_home/workspace/twenty_game/cmd/Game-user_srv
go build main.go -o user_srv
nohup ./user_srv > ../../us.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-game_srv
go build main.go -o game_srv
nohup ./game_srv > ../../gs.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-game_web
go build main.go -o game_web
nohup ./game_web > ../../gw.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-user_web
go build main.go -o user_web
nohup ./user_web > ../../uw.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-file_web
go build main.go -o file_web
nohup ./file_web > ../../fw.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-process_web
go build main.go -o process_web
nohup ./process_web > ../../pw.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-store_web
go build main.go -o store_web
nohup ./store_web > ../../sw.log 2>&1 &
