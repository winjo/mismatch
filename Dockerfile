FROM golang:1.14 as builder

ENV GOPROXY="https://goproxy.cn,direct"
# https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host
# build for static link
ENV CGO_ENABLED=0
WORKDIR /app
COPY . /app
RUN go build .

# production stage
FROM alpine as production
WORKDIR /app
COPY --from=builder /app/mismatch /app
EXPOSE 9000
ENTRYPOINT ["/app/mismatch"]
