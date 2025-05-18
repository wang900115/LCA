FROM golang:1.24.3-alpine


WORKDIR /usr/src/app


COPY . .


RUN go build -o app ./cmd/LCA/main.go


EXPOSE 8080


CMD ["./app"]
