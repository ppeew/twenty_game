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
RUN ls

ENTRYPOINT ["twenty_game/build/zeabur.sh"]


