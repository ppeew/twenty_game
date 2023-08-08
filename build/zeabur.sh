
cd cmd/Game-user_srv
go mod tidy
go build -o user_srv main.go
nohup ./user_srv > ../../us.log 2>&1 &

cd cmd/Game-game_srv
go mod tidy
go build -o game_srv main.go
nohup ./game_srv > ../../gs.log 2>&1 &

cd cmd/Game-game_web
go mod tidy
go build -o game_web main.go
nohup ./game_web > ../../gw.log 2>&1 &

cd cmd/Game-user_web
go mod tidy
go build -o user_web main.go
nohup ./user_web > ../../uw.log 2>&1 &

cd cmd/Game-file_web
go mod tidy
go build -o file_web main.go
nohup ./file_web > ../../fw.log 2>&1 &

cd cmd/Game-process_web
go mod tidy
go build -o process_web main.go
nohup ./process_web > ../../pw.log 2>&1 &

cd cmd/Game-store_web
go mod tidy
go build -o store_web main.go
nohup ./store_web > ../../sw.log 2>&1 &

echo "build finish"
