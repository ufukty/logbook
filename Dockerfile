FROM ufukty/go-argon2:alpine-libargon2-v20190702

WORKDIR /app
COPY . .

RUN go build -o /bin/logbook-backend .

ENTRYPOINT [ "/bin/logbook-backend" ]