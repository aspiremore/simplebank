#build stage
FROM golang:1.21.3-alpine3.18 as builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main ./
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]