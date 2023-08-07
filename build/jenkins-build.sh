echo "jenkins-build"

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-user_srv
go build -o user_srv main.go 
nohup ./user_srv > ../../us.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-game_srv
go build -o game_srv main.go 
nohup ./game_srv > ../../gs.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-game_web
go build -o game_web main.go
nohup ./game_web > ../../gw.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-user_web
go build -o user_web main.go
nohup ./user_web > ../../uw.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-file_web
go build -o file_web main.go
nohup ./file_web > ../../fw.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-process_web
go build -o process_web main.go
nohup ./process_web > ../../pw.log 2>&1 &

cd /var/jenkins_home/workspace/twenty_game/cmd/Game-store_web
go build -o store_web main.go
nohup ./store_web > ../../sw.log 2>&1 &
