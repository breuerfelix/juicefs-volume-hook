# juicefs-volume-hook

This mutating webhook transforms all `volumeMounts` of ALL pods which use a `PersistentVolumeClaim` to use `mountPropagation: HostToContainer` which is needed in order for juiceFS volumes to recover.

## prerequisites

* have `helm` installed on your computer
* have `cert-manager` installed in your cluster

## installation

```bash
kubectl create namespace juicefs
helm repo add juicefs-volume-hook https://breuerfelix.github.io/juicefs-volume-hook
helm repo update
helm install juicefs-volume-hook juicefs-volume-hook/juicefs-volume-hook -n juicefs
```

# TODO

* annotate only specific mounts (via annotation?)

# Nice to Know

I wrote a [Blog Post](https://breuer.dev/blog/kubernetes-webhooks) on how to create a minimal Kubernetes Admission Webhook like this one. Just check it out :)

