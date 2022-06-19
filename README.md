# Description
DNS server that having upstream resolver via TOR network. Only sole purpose to get IP from DNS

# Run

## Build docker

```
docker build . -t dnsserv
```

## Run

```
docker rm -f dnsserv 2>/dev/null && docker run --name dnsserv -dit -p 127.0.0.1:53:5353/udp dnsserv
```

Or just use generated one

```
docker rm -f dnsserv 2>/dev/null && docker run --name dnsserv -dit -p 127.0.0.1:53:5353/udp habibiefaried/dns-over-tor-resolver
```

You can create your own config, put as a file, mount as volume to /app/config.yml. Default config is, what we have provided in docker image

## Feature

1. query via TOR network
2. blocking ads

```
dig @localhost ads.tiktok.com

;; ANSWER SECTION:
ads.tiktok.com.         3600    IN      A       0.0.0.0

;; Query time: 0 msec
;; SERVER: 127.0.0.1#53(localhost) (UDP)
;; WHEN: Sun Jun 19 01:09:57 UTC 2022
;; MSG SIZE  rcvd: 62
```

## NOTICE

Only run and use this on local machine or trusted network. as DNS over UDP (which is default option by most machine) is unsecured protocol