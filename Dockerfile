FROM golang

COPY app .
ENTRYPOINT ["./app"]
