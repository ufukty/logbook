FROM golang:1.18.5-alpine3.15

WORKDIR /app
COPY . .

RUN go build -o /bin/logbook-backend .

ENTRYPOINT [ "/bin/logbook-backend" ]