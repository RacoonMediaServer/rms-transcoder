FROM golang as builder
WORKDIR /src/service
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.Version=`git tag --sort=-version:refname | head -n 1`" -o rms-transcoder -a -installsuffix cgo rms-transcoder.go

FROM jrottenberg/ffmpeg:6.0-alpine
RUN apk --no-cache add ca-certificates tzdata
RUN mkdir /app
WORKDIR /app
COPY --from=builder /src/service/rms-transcoder .
COPY --from=builder /src/service/configs/rms-transcoder.json /etc/rms/
ENTRYPOINT ["./rms-transcoder"]