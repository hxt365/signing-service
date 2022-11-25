FROM golang:1.16

RUN mkdir /code

WORKDIR /code

COPY . .

RUN go mod download

RUN go get github.com/githubnemo/CompileDaemon

ENTRYPOINT chmod 777 ./wait-for-it.sh && ./wait-for-it.sh ${DATABASE_HOST}:${DATABASE_PORT} && CompileDaemon --build="go build main.go" --command=./main
