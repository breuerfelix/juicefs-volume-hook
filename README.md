# juicefs-volume-hook

This mutating webhook transforms all `volumeMounts` of ALL pods which use a `PersistentVolumeClaim` to use `mountPropagation: HostTo'Container` which is needed in order for juiceFS volumes to recover.

## prerequisites

* have `cert-manager` in your cluster installed

## installation

```bash
kubectl create namespace juicefs
kubectl apply -f config/
```
