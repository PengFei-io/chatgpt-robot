FROM golang:1.20-alpine
USER root
WORKDIR /app
COPY . .
RUN go build -o app .
CMD ["./app"]