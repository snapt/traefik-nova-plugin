# Traefik + Nova WAF on Kubernetes: Blog Demo

This guide will help you to install Nova and the Nova Traefik Plugin on 
Kubernetes. It's primarily meant as an example, for you to implement 
in your own configuration. 

We have included the following configurations / config sets: 

1. ```blog/``` a sample blog backend to deploy if you have no services
2. ```blog-ingress.yaml``` an ingress that sends traffic to the blog above
3. ```middleware.yaml``` the middleware this plugin provides
4. ```nova-deployment.yaml``` a deployment of the Nova WAF/WAAP to scan traffic on

In basic terms, we have a blog and an ingress using Traefik to send our requests to that blog.

We add a Nova deployment (using an AutoJoin Key) and then use the middleware to send 
traffic to Nova to be scanned. If it's safe, Traefik continues as normal. If it's 
dangerous a block page (or something else) is returned. 


### Setup Nova

You need to deploy at least one Nova container in order to scan content with Nova. 
This uses Nova AutoJoin to automatically connect to Nova, provision a new container, 
and load your WAF profile, rules, etc. 

1. Set up a Nova Traefik ADC with WAF enabled at https://nova.snapt.net
2. Get the AutoJoin key from https://nova.snapt.net/adcs/auto-join/keys

**Important!** Edit ```nova-deployment.yaml``` to enter your NOVA_AUTO_CONF key that you 
get from Nova. Also set the number of replicas you need (we recommend at least 2). 

Run the Nova deployment as follows:

```shell
$ kubectl apply -f nova-deployment.yaml
```


### Setup Nova Middleware

To load our middleware onto Traefik you need to enable several components: 

1. Traefik Pilot
2. The Traefik Nova Plugin
3. The middleware below. 

To enable Pilot you need to register and set up Pilot to work with your Traefik install: https://pilot.traefik.io/

**Plugins do not work without Pilot**

During this you will need to configure Pilot, and you can also enable the Nova plugins. Edit your Traefik deployment 
and add the following: 

```shell
$ kubectl edit deployment.v1.apps/traefik
```

```yaml
--pilot.token=YOUR_PILOT_TOKEN
--experimental.plugins.traefik-nova-plugin.modulename=github.com/snapt/traefik-nova-plugin
--experimental.plugins.traefik-nova-plugin.version=v1.1.1
```

If you are going to be using a public cloud Kubernetes we recommend also adding: 
```yaml
--entryPoints.web.proxyProtocol.insecure
```

And then using the PROXY protocol on your external Load Balancers in order to see client IPs. 

Once you have completed the above you can load the Nova middleware by running:

```shell
$ kubectl apply -f middleware.yaml
```

### Use Nova on an Ingress Controller

You can now add the following annotation to any Ingress in order to use the Nova middleware: 
```yaml
traefik.ingress.kubernetes.io/router.middlewares: default-nova@kubernetescrd
```

As an example, here is our demo Blog Ingress: 
```yaml
kind: Ingress
apiVersion: networking.k8s.io/v1
metadata:
  name: blog-ingress
  namespace: blog
  annotations:
    traefik.ingress.kubernetes.io/router.entrypoints: web
    traefik.ingress.kubernetes.io/router.middlewares: default-nova@kubernetescrd

spec:
  rules:
    - host: traefik.nova-adc.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: blog-svc
                port:
                  number: 80
```
