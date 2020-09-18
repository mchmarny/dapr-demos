# Dapr integrations demo setup

> Note, this setup assumes Kubernetes cluster created using the [demo cluster setup instructions](../../setup/)

## Components 

Start by create the SandGrid secret

> the queue and store secrets were already created during the cluster setup

```shell
kubectl create secret generic email --from-literal=api-key="${SENDGRID_KEY}"
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
kubectl rollout status deployment/order-auditor
```

Check logs for errors from both containers

```shell
kubectl logs -l app=order-auditor -c daprd --tail 300
kubectl logs -l app=order-auditor -c auditor
```

### Workflow 

Create the config map to hold the Dapr workflow definition

```shell
kubectl create configmap workflows --from-file config/order-cancel.json
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
  --from-literal=accountKey=$AZSAKEY
```

Deploy Dapr Workflows host and wait for it to be ready

```shell
kubectl apply -f deployment/workflow.yaml
kubectl rollout status deployment/workflows-host
```

Check logs for errors from both containers

```shell
kubectl logs -l app=workflows-host -c daprd --tail 300
kubectl logs -l app=workflows-host -c host
```

### Viewer

Deploy Dashboard 

```shell
kubectl apply -f deployment/viewer.yaml
kubectl rollout status deployment/order-viewer
```

Patch ingress to add the viewer rule

```shell
kubectl get ing/ingress-rules -o json \
  | jq '.spec.rules += [{"host":"viewer.cloudylabs.dev","http":{"paths":[{"backend": {"serviceName":"order-viewer","servicePort":80},"path":"/"}]}}]' \
  | kubectl apply -f -
```

Test it: https://viewer.cloudylabs.dev

## Cleanup 


```shell
kubectl delete -f ./deployment
kubectl delete -f ./component

kubectl delete secret email
kubectl delete secret dapr-workflows
kubectl delete configmap workflows

az storage account delete --name daprintdemo
```

> TODO: Remove ingress rule 


