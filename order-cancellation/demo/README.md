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

Navigate to https://view.cloudylabs.dev/ to start the order processing dashboard. There won't be any data yet, so this is just to open the WebSocket connection. 

![Initial UI](../img/ui1.png)

### 2. Submit Cancellation 

Submit the order [cancellation.json](data/cancellation.json) file using `curl`

```shell
API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
curl -v \
     -d @data/cancellation.json \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     "https://api.cloudylabs.dev/v1.0/invoke/workflows/method/order-cancel"
```

### 3. Dashboard (updated)


### 4. Email 

Show confirmation email delivered after the processed completed 


## 5. Observability 

### Distributed Traces 

Forward local port to Zipkin


### Logging 

Forward local port to Kibana


### Metrics 

Forward local port to Grafana



## Setup 


## PubSub Queue


```shell
kubectl apply -f config/queue.yaml
```

## State Store 

```shell
kubectl apply -f config/fn-store.yaml \
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

Deploy Dapr Functions

```shell
kubectl apply -f function.yaml
```

Check logs for errors from both containers

```shell
kubectl logs -l app=auditor -c daprd
kubectl logs -l app=auditor -c auditor
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

Deploy Dapr Workflows host

```shell
kubectl apply -f workflow.yaml
```

Check logs for errors from both containers

```shell
kubectl logs -l app=dapr-workflows-host -c daprd
kubectl logs -l app=dapr-workflows-host -c host
```


## Web App 

Deploy Dashboard 

```shell
kubectl apply -f dashboard.yaml
```


Patch ingress

```shell
kubectl patch deployment patch-demo --patch "$(cat config/ingress.yaml)"
```

Test it

https://viewer.cloudylabs.dev



## Cleanup 

> TODO: list cleanup steps 