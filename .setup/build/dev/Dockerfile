FROM golang:1.18

WORKDIR /app

COPY . .

RUN go install github.com/githubnemo/CompileDaemon@latest && \
    go mod download

ENTRYPOINT ["CompileDaemon", "-directory=./", "-build=go build -o main main.go", "-command=./main"] 
