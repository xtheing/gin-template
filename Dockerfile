FROM golang:1.17 AS build
WORKDIR /app
COPY . /app
ENV GOPROXY https://goproxy.cn,direct
RUN CGO_ENABLED=0 go build -o /app/main -tags netgo -ldflags '-w -extldflags "-static"' main.go

# ------------------ 生成镜像 ------------------
FROM scratch
WORKDIR /app
COPY --from=build /app/config /app/config
COPY --from=build /app/main /app/main
EXPOSE 8080
CMD ["/app/main"]
