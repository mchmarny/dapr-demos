# Dapr integrations demo setup

> Note, this setup assumes Kubernetes cluster created using the [demo cluster setup instructions](../../setup/)

## Components 

Start by create the `order` namespace: 

```shell
kubectl apply -f ./deployment/space.yaml
```

SandGrid secret

> the queue and store secrets were already created during the cluster setup

```shell
kubectl create secret generic email --from-literal=api-key="${SENDGRID_KEY}" -n order
```

Now deploy the queue, state and workflow stores, and email components

```shell
kubectl apply -f ./component
```

## Deployments 

### Auditor 

Deploy Dapr auditor functions and wait for it to be ready 

```shell
kubectl apply -f deployment/auditor.yaml
kubectl rollout status deployment/order-auditor -n order
```

Check logs for errors from both containers

```shell
kubectl logs -l app=order-auditor -c daprd -n order --tail 300
kubectl logs -l app=order-auditor -c auditor -n order
```

### Workflow 

Create the config map to hold the Dapr workflow definition

```shell
kubectl create cm workflows --from-file config/order-cancel.json -n order
```

Create the Azure storage account 

```shell
az storage account create --name daprintdemo --sku Standard_LRS
export AZSAKEY=$(az storage account keys list --account-name daprintdemo --query "[0].value" --output tsv)
```

Create secret to hold the workflow Azure storage account key

> TODO: Remove Azure storage account  dependency 

```shell
kubectl create secret generic dapr-workflows \
  --from-literal=accountName=daprintdemo \
  --from-literal=accountKey=$AZSAKEY \
  -n order
```

Deploy Dapr Workflows host and wait for it to be ready

```shell
kubectl apply -f deployment/workflow.yaml
kubectl rollout status deployment/workflows-host
```

Check logs for errors from both containers

```shell
kubectl logs -l app=workflows-host -c daprd -n order --tail 300
kubectl logs -l app=workflows-host -c host -n order
```

### Viewer

Deploy Dashboard 

```shell
kubectl apply -f deployment/viewer.yaml
kubectl rollout status deployment/order-viewer
```

Create the TLS certs for this domain 

> `demo.dapr.team` is the domain I'm using for this demo

```shell
kubectl create secret tls tls-secret \
    -n order \
    --key ../../setup/certs/demo.dapr.team/cert-pk.pem \
    --cert ../../setup/certs/demo.dapr.team/cert-ca.pem
```

Deploy ingress for `order` 

```shell
kubectl apply -f deployment/ingress.yaml
```


Test it: https://order.demo.dapr.team

## Cleanup 


```shell
kubectl delete -f ./deployment
kubectl delete -f ./component

kubectl delete secret email -n order
kubectl delete secret dapr-workflows -n order
kubectl delete configmap workflows -n order

az storage account delete --name daprintdemo
```


