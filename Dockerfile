FROM golang

ENV GO111MODULE=on
ENV GOPROXY https://goproxy.io

#ENV \
#    PORT=9000 \
#    HOST=0.0.0.0
#EXPOSE 9000

WORKDIR $GOPATH/src
RUN mkdir twenty_game && cd twenty_game
COPY . twenty_game
RUN pwd
RUN cd twenty_game && ls

RUN cd twenty_game/cmd/Game-user_srv  \
    && go build -o user_srv main.go

RUN cd twenty_game/cmd/Game-game_srv  \
    && go build -o game_srv main.go

RUN cd twenty_game/cmd/Game-game_web  \
    && go build -o game_web main.go

RUN cd twenty_game/cmd/Game-user_web  \
    && go build -o user_web main.go

RUN cd twenty_game/cmd/Game-file_web \
    && go build -o file_web main.go

RUN cd twenty_game/cmd/Game-process_web \
    && go build -o process_web main.go

RUN cd twenty_game/cmd/Game-store_web  \
    && go build -o store_web main.go

RUN echo "build finish"

CMD ["nohup twenty_game/cmd/Game-user_srv/user_srv &","&& nohup twenty_game/cmd/Game-game_srv/game_srv &","&& nohup twenty_game/cmd/Game-process_web/process_web &","&& nohup twenty_game/cmd/Game-user_web/user_web &","&& nohup twenty_game/cmd/Game-game_web/game_web &","&& nohup twenty_game/cmd/Game-store_web/store_web &","&& nohup twenty_game/cmd/Game-file_web/file_web &"]

