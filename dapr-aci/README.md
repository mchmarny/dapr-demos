# Dapr sidecar in ACI

Demo of Dapr sidecar in Azure Container Instances (ACI)

## Setup 

This demo illustrates simple service subscription to pub/sub topic and persistence of event payload into state. To deploy Dapr into ACI however you will need to first setup a SMB volume which will be used to store the Dapr components and mounted in the `daprd` container The storage account name needs to be globally unique so set `SNAME` to something 3-24 chars long, containing alphanumerics only, and make sure it's all in lower case.

```shell
export SNAME="dapraci"
```

> assumes your resource group and location defaults are already set. If not, set them now:

```shell
az account set --subscription <id or name>
az configure --defaults location=<preferred location> group=<preferred resource group>
```


Create a storage account

```shell
az storage account create --name $SNAME --sku Standard_LRS
```

Create a storage share

> For demo purposes share and storage user names are the same 

```shell
az storage share create --name $SNAME --account-name $SNAME
```

Capture storage key 

```shell
export ACCOUNT_KEY=$(az storage account keys list --account-name $SNAME \
                                                  --query "[0].value" \
                                                  --output tsv)
echo $ACCOUNT_KEY
```

Now update `volumes[components].azureFile.storageAccountKey` in `configuration/app.yaml` file so that ACI can mount it.

Upload the Dapr component files

```shell
az storage file upload --account-key $ACCOUNT_KEY \
                       --account-name $SNAME \
                       --share-name $SNAME \
                       --source components/state.yaml

az storage file upload --account-key $ACCOUNT_KEY \
                       --account-name $SNAME \
                       --share-name $SNAME \
                       --source components/pubsub.yaml
```

List files to make sure they are all there

```shell
az storage file list \
    --account-key $ACCOUNT_KEY \
    --share-name $SNAME \
    --account-name $SNAME  \
    --output tsv
```

## Deployment 

Once the storage is set up, you can deploy

```shell
az container create -f deployment/app.yaml
```

When you list the containers:

```shell
az container list -o table
```

The result should look something like this:

```shell
Name      ResourceGroup    Status     Image                                                IP:ports               Network    CPU/Memory       OsType    Location
--------  ---------------  ---------  ---------------------------------------------------  ---------------------  ---------  ---------------  --------  ----------
dapraci   mchmarny         Succeeded  daprio/daprd:0.11.3,ghcr.io/mchmarny/aci-app:v0.2.2  40.xx.xx.xx:3500       Public     1.0 core/1.5 gb  Linux     westus2
```

## Demo

First, capture the IP for ease of access:

```shell
export APP_IP=$(az container show -n dapraci --query "ipAddress.ip" -o tsv)
```

Next, invoke the ping method thru Dapr API:

```shell
curl -i -d '{"message":"ping"}' \
     -H "Content-type: application/json" \
     "http://${APP_IP}:3500/v1.0/invoke/dapraci/method/ping"
```

Response should look something like this:

```json
{ "on": 1604003460965972895, "greeting": "pong" }
```

You can also invoke the PubSub API on Dapr to publish:

```shell
curl -i -d '{"message":"hello"}' \
     -H "Content-type: application/json" \
     "http://${APP_IP}:3500/v1.0/publish/pubsub/messages"
```

Response from post has no body but you should see the headers:

```shell
HTTP/1.1 200 OK
Server: fasthttp
Date: Thu, 29 Oct 2020 21:05:43 GMT
Content-Length: 0
Traceparent: 00-bd0f6f745de1b2cc8b5463f8abaa8656-d2f1d6be6d202273-00
```

This demo also exposes the user container directly. You can disable it but commenting out `- port: 8080` in `ports` section of [configuration/app.yaml](configuration/app.yaml). To invoke the `/ping` route on the deployed app: 

```shell
curl -i -d '{"message":"ping"}' \
     -H "Content-type: application/json" \
     "http://${APP_IP}:8082/ping"
```

### Logs

To view logs from the Dapr container:

```shell
az container logs --name dapraci --container-name daprd
```

> Note, `daprd` is set to log in JSON so you can use `jq` or similar to query the logs and parse out only the messages

```shell
az container logs --name dapraci --container-name daprd | jq ".msg"
```

To query the app container:

```shell
az container logs --name dapraci --container-name app
```

That's it, I hope you found it helpful. 

## Todo

* Dapr API token auth
* Configuration to show ACT
* Secrets to enable Azure Vault 

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
