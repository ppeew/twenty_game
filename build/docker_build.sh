#!sh

cd ~/twenty_game/cmd/Game-admin_web && docker build -t admin_web:1.0 .

cd ~/twenty_game/cmd/Game-file_web && docker build -t file_web:1.0 .

cd ~/twenty_game/cmd/Game-game_web && docker build -t game_web:1.0 .

cd ~/twenty_game/cmd/Game-hall_web && docker build -t hall_web:1.0 .

cd ~/twenty_game/cmd/Game-process_web && docker build -t process_web:1.0 .

cd ~/twenty_game/cmd/Game-hall_web && docker build -t hall_web:1.0 .

cd ~/twenty_game/cmd/Game-store_web && docker build -t store_web:1.0 .

cd ~/twenty_game/cmd/Game-user_web && docker build -t user_web:1.0 .

cd ~/twenty_game/cmd/Game-user_srv && docker build -t user_srv:1.0 .
