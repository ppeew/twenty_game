cd %~dp0/../srv/user_srv
start /B cmd /C "go run . & exit"

cd %~dp0/../srv/game_srv
start /B cmd /C "go run . & exit"

cd %~dp0/../web/game_web
start /B cmd /C "go run . & exit"

cd %~dp0/../web/user_web
start /B cmd /C "go run . & exit"

cd %~dp0/../web/file_web
start /B cmd /C "go run . & exit"

cd %~dp0/../web/process_web
start /B cmd /C "go run . & exit"
