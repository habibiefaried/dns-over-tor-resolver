# Description
DNS server that having upstream resolver via TOR network. Only sole purpose to get IP from DNS

# Run

## Build docker

```
docker build . -t dnsserv
```

## Run

```
docker rm -f dnsserv 2>/dev/null && docker run --name dnsserv -dit -p 53:5353/udp dnsserv
```

Or just use generated one

```
docker rm -f dnsserv 2>/dev/null && docker run --name dnsserv -dit -p 53:5353/udp habibiefaried/dns-over-tor-resolver
```