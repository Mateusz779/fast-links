# Author:Mateusz Kędziora https://mkedziora.pl

FROM golang:1.23.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Etap minifikacji
FROM node:alpine AS minifier
WORKDIR /app
RUN npm install -g minify
COPY public/ ./public/
RUN find ./public -name "*.html" -exec minify -o {} {} \;
RUN find ./public/css -name "*.css" -exec minify -o {} {} \;
RUN find ./public/js -name "*.js" -exec minify -o {} {} \;

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=minifier /app/public/ ./public/
EXPOSE 8000
CMD ["./main"]