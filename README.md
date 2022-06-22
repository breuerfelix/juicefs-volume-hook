# juicefs-volume-hook

This mutating webhook transforms all `volumeMounts` of ALL pods which use a `PersistentVolumeClaim` to use `mountPropagation: HostToContainer` which is needed in order for juiceFS volumes to recover.

## prerequisites

* have `cert-manager` installed in your cluster

## installation

```bash
kubectl create namespace juicefs
kubectl apply -f config/
```

# TODO

* Helmchart
* annotate only specific mounts (via annotation?)

# Nice to Know

I wrote a [Blog Post](https://breuer.dev/blog/kubernetes-webhooks) on how to create a minimal Kubernetes Admission Webhook like this one. Just check it out :)

