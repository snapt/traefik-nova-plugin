# Traefik Nova Plugin

![Banner](./img/banner.png)

Traefik plugin to proxy requests to Snapt Nova for evaluation against the WAF. 

## Usage (docker-compose.yml)

See [docker-compose.yml](docker-compose.yml)

1. docker-compose up
2. Go to http://localhost:8000/website, you see your request details
3. Go to http://localhost:8000/website?test=/etc/passwd, the request is blocked if Nova's WAF is enabled

## How it works

This adds a middleware plugin to Traefik which proxies requests to a Nova container before 
sending them to your backend. If Nova determines the request should be blocked 
then it returns a block, otherwise it allows Traefik to continue as it would have.

This requires an AutoJoin key from a Traefik-based ADC you have already added 
on Nova (https://nova.snapt.net) and, naturally, requires that you enable the 
WAF. You can use learning mode on Nova to see what would be blocked. 

