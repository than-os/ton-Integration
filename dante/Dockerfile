FROM golang:alpine AS build-env
WORKDIR /root/sentinel
COPY dante.go /root/sentinel/
RUN apk update && apk add ca-certificates git && go get github.com/fatih/color github.com/jasonlvhit/gocron && cd /root/sentinel && go build -o dante dante.go

FROM alpine
WORKDIR /root/sentinel
COPY --from=build-env /root/sentinel/dante /root/sentinel/dante
EXPOSE 3030
ENTRYPOINT ["./dante"]
