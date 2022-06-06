### Install
Set up the blog and headless service. After this you will have 
two instances of our demo web container (the blog) and one 
service to use with Nova. 

```
./deploy.sh
```

### Uninstall 
Remove our whole blog namespace

```
kubectl delete all --all -n blog
kubectl delete namespace blog
```

### Check status
```
kubectl get pods,deployments,services -n blog
```