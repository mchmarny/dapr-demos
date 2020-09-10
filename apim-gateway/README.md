# Dapr & Azure API Management Integration Demo

This repo demonstrates how to expose Dapr API and invoke a service on Kubernetes using [Azure API Management](https://azure.microsoft.com/en-us/services/api-management/) (APIM) self-hosted gateway.

## Prerequisite 

* [Azure account](https://azure.microsoft.com/en-us/free/)
* [Azure CLI](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli?view=azure-cli-latest)
* [Kubernetes cluster with Dapr](https://github.com/dapr/docs/blob/v0.9.0/getting-started/environment-setup.md#installing-dapr-on-a-kubernetes-cluster)

## Setup 

To make this demo easier to reproduce, start by exporting the name of the Azure API Management (APIM) service we will create.

> Note, the name of your API Management service instance name has to be globally unique!

```shell
export APIM_SERVICE_NAME="dapr-apim-demo"
```

In addition to the above name, export also the Azure [Subscription ID](https://docs.bitnami.com/azure/faq/administration/find-subscription-id/) and [Resource Group](https://docs.bitnami.com/azure/faq/administration/find-deployment-resourcegroup-id/) where you would like to create these APIM service.

```shell
export AZ_SUBSCRIPTION_ID="your-subscription-id"
export AZ_RESOURCE_GROUP="your-resource-group"
```

## Azure API Management Configuration 

In this section we will create all the Azure resources. First, create and configure the Azure API Management service.

### Service

Create APIM service instance:

> The `publisher-email` and `publisher-name` below are required to receive system notifications e-mails.

```shell
az apim create --name $APIM_SERVICE_NAME \
               --subscription $AZ_SUBSCRIPTION_ID \
               --resource-group $AZ_RESOURCE_GROUP \
               --publisher-email "you@your-domain.com" \
               --publisher-name "Your Name"
```

> Note, depending on the SKU and resource group configuration, this operation may take 15+ min. While this running, consider quick read on [API Management Concepts](https://docs.microsoft.com/en-us/azure/api-management/api-management-key-concepts#-apis-and-operations)

### API

Each APIM [API](https://docs.microsoft.com/en-us/azure/api-management/api-management-key-concepts#-apis-and-operations) map to back-end service managed by Dapr. This demo will use a simple echo service hosted in Dapr which simply returns posted content. To define that mapping you will need to first update the [api.yaml](./api.yaml) file with the name of the APIM service created above:

```yaml
servers:
  - url: http://YOUR-APIM-SERVICE-NAME.azure-api.net
  - url: https://YOUR-APIM-SERVICE-NAME.azure-api.net
```

When finished, import that OpenApi definition fle into APIM service instance:

```shell
az apim api import --path / \
                   --api-id dapr-echo \
                   --service-name $APIM_SERVICE_NAME \
                   --display-name "Demo Dapr Service API" \
                   --protocols http https \
                   --subscription-required false \
                   --specification-path api.yaml \
                   --specification-format OpenApi
```

### Policy

APIM [Policies](https://docs.microsoft.com/en-us/azure/api-management/api-management-key-concepts#--policies) are defined in XML and sequentially executed on each request and/or response. In this demo we will create a simple `inbound` policy mapping to to the Dapr service method. 

```xml
<policies>
     <inbound>
          <set-backend-service backend-id="dapr" dapr-app-id="echo-service" dapr-method="echo" />
     </inbound>
     <backend>
          <base />
     </backend>
     <outbound>
          <base />
     </outbound>
     <on-error>
          <base />
     </on-error>
</policies>
```

To apply policy we will first need export an Azure management API token: 

```shell
export AZ_API_TOKEN=$(az account get-access-token --resource=https://management.azure.com --query accessToken --output tsv)
```

And then apply the policy:

```shell
curl -X PUT \
     -d @./policy.json \
     -H "Content-Type: application/json" \
     -H "If-Match: *" \
     -H "Authorization: Bearer ${AZ_API_TOKEN}" \
     "https://management.azure.com/subscriptions/${AZ_SUBSCRIPTION_ID}/resourceGroups/${AZ_RESOURCE_GROUP}/providers/Microsoft.ApiManagement/service/${APIM_SERVICE_NAME}/apis/dapr-echo/operations/echo/policies/policy?api-version=2019-12-01"
```

If everything goes well, the API will returned the created policy.

### Gateway

To create a self-hosted gateway which will be then deployed to the Kubernetes cluster, first, we need to create the `demo-apim-gateway` object in APIM:

```shell
curl -v -X PUT -d '{"properties": {"description": "Dapr Gateway","locationData": {"name": "Virtual"}}}' \
     -H "Content-Type: application/json" \
     -H "If-Match: *" \
     -H "Authorization: Bearer ${AZ_API_TOKEN}" \
     "https://management.azure.com/subscriptions/${AZ_SUBSCRIPTION_ID}/resourceGroups/${AZ_RESOURCE_GROUP}/providers/Microsoft.ApiManagement/service/${APIM_SERVICE_NAME}/gateways/demo-apim-gateway?api-version=2019-12-01"
```

And then map the gateway to the previously created API:

```shell
curl -v -X PUT -d '{ "properties": { "provisioningState": "created" } }' \
     -H "Content-Type: application/json" \
     -H "If-Match: *" \
     -H "Authorization: Bearer ${AZ_API_TOKEN}" \
     "https://management.azure.com/subscriptions/${AZ_SUBSCRIPTION_ID}/resourceGroups/${AZ_RESOURCE_GROUP}/providers/Microsoft.ApiManagement/service/${APIM_SERVICE_NAME}/gateways/demo-apim-gateway/apis/dapr-echo?api-version=2019-12-01"
```

If everything goes well, the API returns JSON of the created objects.

## Kubernetes Configuration 

Moving now to your Kubernetes cluster...

### Dapr Service 

To deploy your application as a Dapr service you just need to decorating your Kubernetes deployment template with few Dapr annotations.

```yaml
annotations:
     dapr.io/enabled: "true"
     dapr.io/app-id: "echo-service"
     dapr.io/app-protocol: "http"
     dapr.io/app-port: "8080"
```

> To learn more about Kubernetes sidecar configuration see [Dapr docs](https://github.com/dapr/docs/blob/master/concepts/configuration/README.md#kubernetes-sidecar-configuration).

For this demo we will use a pre-build Docker image of the [http-echo-service](https://github.com/mchmarny/dapr-demos/tree/master/http-echo-service). The Kubernetes deployment file of that service is defined [here](./service.yaml). Deploy it, and check that it is ready:

```shell
kubectl apply -f service.yaml
kubectl get pods -l app=echo-service
```

> Service is ready when its status is `Running` and the ready column is `2/2` (Dapr and our echo service both started)

```shell
NAME                            READY   STATUS    RESTARTS   AGE
echo-service-77d6f5b5bb-crc5q   2/2     Running   0          97s
```

### APIM Gateway 

To connect the self-hosted gateway to APIM service, we will need to create first a Kubernetes secret with the APIM gateway key. First, get the key from APIM API:

> Note, the maximum validity for access tokens is 30. Update the below `expiry` parameter to be withing 30 days from today

```shell
curl -X POST -d '{ "keyType": "primary", "expiry": "2020-10-10T00:00:00Z" }' \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer ${AZ_API_TOKEN}" \
     "https://management.azure.com/subscriptions/${AZ_SUBSCRIPTION_ID}/resourceGroups/${AZ_RESOURCE_GROUP}/providers/Microsoft.ApiManagement/service/${APIM_SERVICE_NAME}/gateways/demo-apim-gateway/generateToken?api-version=2019-12-01"
```

Copy the content of `value` from the response and create a secret:

> Make sure the secret includes the `GatewayKey` + a space ` ` + the value of your token (e.g. `GatewayKey a1b2c3...`)

```shell
kubectl create secret generic demo-apim-gateway-token --type Opaque --from-literal value="GatewayKey YOUR-TOKEN-HERE"
```

Now, create a config map containing the APIM service endpoint that will be used to configure your self-hosted gateway:

```shell
kubectl create configmap demo-apim-gateway-env --from-literal \
     "config.service.endpoint=https://dapr-apim-demo.management.azure-api.net/subscriptions/${AZ_SUBSCRIPTION_ID}/resourceGroups/${AZ_RESOURCE_GROUP}/providers/Microsoft.ApiManagement/service/${APIM_SERVICE_NAME}?api-version=2019-12-01"
```

And finally, deploy the gateway and check that it's ready:

```shell
kubectl apply -f gateway.yaml
kubectl get pods -l app=demo-apim-gateway
```

> Note, the self-hosted gateway is deployed with 2 replicas to ensure availability during upgrades. 

Make sure both instances have status `Running` and container is ready `2/2` (gateway container + Dapr side-car).

```shell
NAME                                 READY   STATUS    RESTARTS   AGE
demo-apim-gateway-6dfb968f5c-cb4t7   2/2     Running   0          26s
demo-apim-gateway-6dfb968f5c-gxrrq   2/2     Running   0          26s
```

## Test

We are ready to test. Start by capturing the cluster load balancer ingress IP:

```shell
export GATEWAY_IP=$(kubectl get svc demo-apim-gateway -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
```

And now, try posting a message to the APIM self-hosted gateway which will be forwarded to the backing Dapr service:

```shell
curl -X POST -d '{ "message": "hello" }' \
     -H "Content-Type: application/json" \
     "http://${GATEWAY_IP}/dapr-echo"
```

If everything is configured correctly, you should see the response from your backing Dapr service: 

```json 
{ "message": "hello" }
```

In addition, you can also check the `echo-service` logs:

```shell
kubectl logs -l app=echo-service -c service
```

This demo illustrates how to setup the APIM service and deploy your self-hosted gateway. Using this gateway can mange access to any number of your Dapr services hosted on Kubernetes. There is a lot more that APIM can do (e.g. Discovery, Access Control, Throttling, Caching, Logging, Traces etc.). You can find out more about APIM [here](https://azure.microsoft.com/en-us/services/api-management/)

## Cleanup 

```shell
kubectl delete -f gateway.yaml
kubectl delete -f service.yaml
kubectl delete secret demo-apim-gateway-token
az apim delete --name daprapimdemo --no-wait --yes
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
