FROM golang:1.14.4-alpine as builder

WORKDIR /app
COPY . .
# binary build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /go/bin/app .

EXPOSE 8000
CMD ["go", "run", "main.go"]

# docker run -p 8000:8000 --env-file .env baitregistry.azurecr.io/tsaving:20200825.1
# docker run -p 8000:8000 baitregistry.azurecr.io/tsaving:20200825.1

# FROM alpine:latest
# COPY --from=builder /go/bin/app /go/bin/app
# COPY --from=builder /app/.env /go/bin/.env
# EXPOSE 8000
# ENTRYPOINT ["/go/bin/app"]