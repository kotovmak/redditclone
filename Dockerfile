FROM golang:1.16 AS build

COPY . /redditclone/
WORKDIR /redditclone/

RUN go mod download
RUN GOOS=linux go build -v -o ./build/app ./cmd/apiserver/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=build /redditclone/build/app .
COPY --from=build /redditclone/configs configs/
COPY --from=build /redditclone/web web/

CMD ["./app"]
