FROM golang:alpine as builder
RUN adduser -D -g '' appuser
WORKDIR /go/src/github.com/ekundo/godis/
COPY . .
RUN apk add --no-cache git
RUN go get -t -d -v ./...
RUN go test -v ./...
RUN go test -tags integration -v ./...
RUN CGO_ENABLED=0 go build -o /app github.com/ekundo/godis/server/main

FROM scratch
STOPSIGNAL SIGINT
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd
COPY --chown=appuser --from=builder /app /app
USER appuser
VOLUME [ "/work" ]
WORKDIR /work
EXPOSE 2121
ENTRYPOINT ["/app", "-host=0.0.0.0", "-port=2121", "-wal=true"]

# docker build -t juno . 
# docker run -v $(pwd):/work -p 2121:2121 -it --rm juno