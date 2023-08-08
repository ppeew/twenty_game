FROM golang

ENV GO111MODULE=on
ENV GOPROXY https://goproxy.io

#ENV \
#    PORT=9000 \
#    HOST=0.0.0.0
#EXPOSE 9000

WORKDIR $GOPATH/src
COPY . .
RUN pwd
RUN ls
#RUN go mod init
#RUN go mod tidy

#RUN go build -o main .
#CMD ["./main", "&"]

RUN echo "zeabur-build start"
RUN cd cmd/Game-user_srv
RUN go build -o user_srv main.go
CMD nohup ./user_srv > ../../us.log 2>&1 &

RUN cd ~/twenty_game/cmd/Game-game_srv
RUN go build -o game_srv main.go
RUN nohup ./game_srv > ../../gs.log 2>&1 &

RUN cd ~/twenty_game/cmd/Game-game_web
RUN go build -o game_web main.go
RUN nohup ./game_web > ../../gw.log 2>&1 &

RUN cd ~/twenty_game/cmd/Game-user_web
RUN go build -o user_web main.go
RUN nohup ./user_web > ../../uw.log 2>&1 &

RUN cd ~/twenty_game/cmd/Game-file_web
RUN go build -o file_web main.go
RUN nohup ./file_web > ../../fw.log 2>&1 &

RUN cd ~/twenty_game/cmd/Game-process_web
RUN go build -o process_web main.go
RUN nohup ./process_web > ../../pw.log 2>&1 &

RUN cd ~/twenty_game/cmd/Game-store_web
RUN go build -o store_web main.go
RUN nohup ./store_web > ../../sw.log 2>&1 &

RUN echo "zeabur-build finish"



