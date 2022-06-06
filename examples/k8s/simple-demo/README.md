# Traefik + Nova WAF on Kubernetes: Simple Demo

![Example](https://github.com/snapt/traefik-nova-plugin/blob/main/k8s/simple-demo/example-router.png?raw=true)

This simple demo creates a whoami container with a Nova WAF and Nova WAF middleware protecting the 
Traefik ingress to it. 

You can reach it by going to https://your-cluster/whoami/

Test a block by going to https://your-cluster/whoami/?test=/etc/passwd

If you receive a blank page after your first deployment, wait ~30-60 seconds for Nova to come online.


## Setup

You need to deploy at least one Nova container in order to scan content with Nova. 
This uses Nova AutoJoin to automatically connect to Nova, provision a new container, 
and load your WAF profile, rules, etc. 

1. Set up a Nova Traefik ADC with WAF enabled at https://nova.snapt.net
2. Get the AutoJoin key from https://nova.snapt.net/adcs/auto-join/keys

**Important!** Edit ```deployment.yaml``` to enter your NOVA_AUTO_CONF key that you 
get from Nova. 

Run the full deployment as follows:

```shell
$ kubectl apply -f deployment.yaml
```