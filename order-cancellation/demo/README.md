# Dapr integrations demo

Dapr integration demo consists of:

1. Starting the order processing dashboard 
2. Submitting cancellation request 
3. Viewing processed request in the dashboard 
4. Querying the state store for cancellation data
5. Showing order cancellation confirmation email 
6. Review of the distributed traces for entire process 

> Note, instructions on how to setup a Kubernetes cluster for this demo are located [here](../setup/README.md)

### 1. Dashboard 

> Note, these instructions assume `cloudylabs.dev` domain setup in the [cluster setup](../setup/README.md) step. You will need to substitute this for your own domain. 

Navigate to https://viewer.cloudylabs.dev to start the order processing dashboard. There won't be any data yet, so this is just to open the WebSocket connection. 

![Initial UI](../img/ui1.png)

### 2. Submit Cancellation 

Submit the order [cancellation.json](data/cancellation.json) file using `curl`

```shell
API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
curl -i \
     -d @data/cancellation.json \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/invoke/workflows/method/order-cancel"
```

### 3. Dashboard (updated)

View the dashboard again at https://viewer.cloudylabs.dev to see the orders

### 4. Email 

Show confirmation email delivered after the processed completed 

> Make sure to check junk mail!

## 5. Observability 

### Distributed Traces 

Forward local port to Zipkin

http://localhost:9411

### Logging 

Forward local port to Kibana

http://localhost:5601

### Metrics 

Forward local port to Grafana

http://localhost:8888


## Setup 


## PubSub Queue


```shell
kubectl apply -f config/queue.yaml
```

## State Store 

```shell
kubectl apply -f config/audit-store.yaml \
              -f config/workflow-store.yaml
```

## Email 

Create the SandGrid secret

```shell
kubectl create secret generic email --from-literal=api-key="<YOUR-SECRET>"
```

And the email component

```shell
kubectl apply -f config/email.yaml
```


## Auditor 

Deploy Dapr auditor functions and wait for it to be ready 

```shell
kubectl apply -f auditor.yaml
kubectl rollout status deployment/order-auditor
```

Check logs for errors from both containers

```shell
kubectl logs -l app=order-auditor -c daprd --tail 300
kubectl logs -l app=order-auditor -c auditor
```

## Workflow 

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
kubectl apply -f workflow.yaml
kubectl rollout status deployment/workflows-host
```

Check logs for errors from both containers

```shell
kubectl logs -l app=workflows-host -c daprd --tail 300
kubectl logs -l app=workflows-host -c host
```


## Web App 

Deploy Dashboard 

```shell
kubectl apply -f viewer.yaml
```

Patch ingress to add the viewer rule

```shell
kubectl get ing/ingress-rules -o json \
  | jq '.spec.rules += [{"host":"viewer.cloudylabs.dev","http":{"paths":[{"backend": {"serviceName":"order-viewer","servicePort":80},"path":"/"}]}}]' \
  | kubectl apply -f -
```

Test it

https://viewer.cloudylabs.dev

## Restart Gateway

```shell
kubectl rollout restart deployment/nginx-ingress-nginx-controller
kubectl rollout status deployment/nginx-ingress-nginx-controller
```


## Cleanup 

> TODO: list cleanup steps 