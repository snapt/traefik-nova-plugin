# Traefik + Nova WAF on Docker

We have included a docker-compose.yml file primarily for testing purposes. Typically, 
Nova is deployed in a Kubernetes environment with Traefik.

See [docker-compose.yml](docker-compose.yml) for configuration options.

### Setup Nova

You need to deploy at least one Nova container in order to scan content with Nova. 
This uses Nova AutoJoin to automatically connect to Nova, provision a new container, 
and load your WAF profile, rules, etc. 


1. Set up a Nova Traefik ADC with WAF enabled at https://nova.snapt.net
2. Get the AutoJoin key from https://nova.snapt.net/adcs/auto-join/keys
3. Configure your docker-compose.yml or set your environment variables - you must provide an autojoin key and a Traefik Pilot token.

### Running Docker

You can now launch the Docker environment. It can take 30-60 seconds to come online 
after your first launch.

1. docker-compose up
2. Go to http://localhost:8000/website, you see your request details
3. Go to http://localhost:8000/website?test=/etc/passwd, the request is blocked if Nova's WAF is enabled

```/website``` is setup a test page, but you can tinker with the backend from there.
