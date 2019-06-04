1. Create a signed cert/key pair and store it in a Kubernetes `secret` that will be consumed by sidecar deployment
./webhook-create-signed-cert.sh \
    --service sidecar-injector-webhook-svc \
    --secret sidecar-injector-webhook-certs \
    --namespace default
```

2. Patch the `MutatingWebhookConfiguration` by set `caBundle` with correct value from Kubernetes cluster
```
cat mutatingwebhook.yaml | \
    webhook-patch-ca-bundle.sh > \
    mutatingwebhook-ca-bundle.yaml
```

3. Deploy resources
```
kubectl create -f deployment/nginxconfigmap.yaml
kubectl create -f deployment/configmap.yaml
kubectl create -f deployment/deployment.yaml
kubectl create -f deployment/service.yaml
kubectl create -f deployment/mutatingwebhook-ca-bundle.yaml
