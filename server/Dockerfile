FROM golang as builder

WORKDIR /app/server

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server

FROM scratch

WORKDIR /app

COPY --from=builder /app/server/server /app/server

CMD [ "/app/server" ]
EXPOSE 8080