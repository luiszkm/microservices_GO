FROM golang:1.20 AS builder

WORKDIR /app
COPY ./src ./
RUN CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -o server .

CMD [ "tail", "-f", "/dev/null" ]
# ENTRYPOINT [ "/server" ] 