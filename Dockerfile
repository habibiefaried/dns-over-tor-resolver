FROM golang:1.18.2-alpine3.16 as builder

WORKDIR /app
COPY . /app/
RUN go build -o main && chmod +x main

FROM alpine:3.16.0

WORKDIR /app
COPY --from=builder /app/main .
RUN apk add --no-cache tor
ENTRYPOINT [ "/bin/sh" ]
CMD [ "-c", "/app/main" ]