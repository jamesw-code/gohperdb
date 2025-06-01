FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o gopherdb .
EXPOSE 5321
CMD ["./gopherdb"]
