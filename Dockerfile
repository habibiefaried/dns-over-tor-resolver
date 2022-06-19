FROM golang:1.18.2-alpine3.16 as builder

WORKDIR /app
COPY . /app/
RUN go build -o main && chmod +x main

FROM alpine:3.16.0

WORKDIR /app
COPY --from=builder /app/main .
COPY entrypoint.sh /app/entrypoint.sh
COPY config.yml /app/config.yml 
RUN apk add --no-cache tor
RUN wget https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts -O /app/hosts
RUN chmod +x /app/entrypoint.sh
ENTRYPOINT [ "/bin/sh" ]
CMD [ "-c", "/app/entrypoint.sh" ]