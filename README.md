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

## NOTICE

Only run and use this on local machine or trusted network. as DNS over UDP (which is default option by most machine) is unsecured protocol