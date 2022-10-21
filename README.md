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

## usage

### default

The controller will convert ALL `volumeMount.mountPropagation` fields which have a `PersistentVolumeClaim` to `HostToContainer`.

### via annotation

Start the controller with the `--pod-annotation` flag and the controller will ONLY process pods which have set the `juicefs.volume.hook/mount-propagation` annotation to `"true"`.

### via storageclass

Start the controller with the `--storage-classes=foobar` flag in order to ONLY process `volumeMounts` that have the given storage classes. Multiple storage classes have to be comma separated.

## development

The `controller-runtime` package always creates a Webhook Server that relies on
TLS, and therefore requires a certificate. As this complicates local
development, you can use the following approach to tunnel the server running
locally with a self-signed certificate through `localtunnel`, exposing it via
public TLS which the Kubernetes API server accepts.

Generate selfsigned *certificates*:

```bash
bash hack/gen-certs.sh
```

Start the local webhook server with the certs configured:

```bash
go run main.go --cert-dir certs --key-name server.key --cert-name server.crt
```

Create a tunnel that exposes your local server:

```bash
npx localtunnel --port 9443 --local-https --local-ca certs/ca.crt --local-cert certs/server.crt --local-key certs/server.key --subdomain juicefs
```

Finally, apply the `MutatingWebhookConfiguration`:

```bash
kubectl apply -f example/webhook.yaml
```

## TODO

* add serviceaccount, role and rolebinding for controller to list PVC

# Nice to Know

I wrote a [Blog Post](https://breuer.dev/blog/kubernetes-webhooks) on how to create a minimal Kubernetes Admission Webhook like this one. Just check it out :)

