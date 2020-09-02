# Dapr API in ACI

* Purpose-configured instance of Dapr deployed into Azure Container Instances (ACI) with API token authentication using single command 
* Use of Dapr output binding + Dapr as a microservice (in this case email sending)


## Setup 

The storage account name needs to be globally unique. Set `SNAME` to something 3-24 chars long, containing alphanumerics only, and make sure it's all in lower case.

```shell
export SNAME="demodapr"
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

Create a storage share for config

> For demo purposes share and storage user names are the same 

```shell
az storage share create --name $SNAME --account-name $SNAME
```

Capture storage key 

```shell
export SKEY=$(az storage account keys list --account-name $SNAME --query "[0].value" --output tsv)
```

Create a storage directory for config files  

```shell
az storage directory create --account-name $SNAME --name $SNAME --share-name $SNAME
```

Upload the Dapr component files

> TODO: Make sure you set the Sendgrid API key in the email.yaml

```shell
az storage file upload --account-name $SNAME --share-name $SNAME --source email.yaml
```

## Deployment 

Once the storage is set up, you can deploy. Start by exporting Dapr API Authentication token

```shell
export DTOKEN=$(openssl rand -base64 36)
```

> Note, make sure to save the value exported into `$DTOKEN` variable to ensure you can use it in other terminal sessions. That value will not be recoverable from the ACI service. 

And launch the Dapr container

```shell
az container create \
    --name $SNAME \
    --ports 3500 \
    --protocol TCP \
    --dns-name-label $SNAME \
    --image docker.io/daprio/daprd:0.10.0 \
    --command-line "/daprd --components-path /components --app-protocol http" \
    --secure-environment-variables "DAPR_API_TOKEN=${DTOKEN}" \
    --azure-file-volume-share-name $SNAME \
    --azure-file-volume-account-name $SNAME \
    --azure-file-volume-account-key $SKEY \
    --azure-file-volume-mount-path /components
```

Then check on the status of the deployment 

```shell
az container list -o table
```

The result should look something like this 

```shell
Name      ResourceGroup  Status     Image                          IP:ports           Network  CPU/Memory       OsType    Location
--------  -------------  ---------  -----------------------------  -----------------  -------  ---------------  --------  --------
demodapr  mchmarny       Succeeded  docker.io/daprio/daprd:0.10.0  51.143.49.0:3500   Public   1.0 core/1.5 gb  Linux     westus2
```

If everything went OK, you should be able post to the email output binding below

To restart the service after update of environment variables 

```shell
az container restart --name $SNAME
```

## Use

To use the above deployed instance of Dapr configured with SendGrid output binding, POST to the Dapr API following message using `curl`.

> Note, the from, to, and email subject are configured server side so all you have to submit is a valid output binding message with the `operation` and `data` properties, with the body of the email sent to the user.

```shell
export SREGION=$(az container list --query "[?contains(name, '${SNAME}')].location" --output tsv)
```


```shell
curl -v -X POST -H "Content-Type: application/json" \
    -H "dapr-api-token: ${DTOKEN}" \
    "http://${SNAME}.${SREGION}.azurecontainer.io:3500/v1.0/bindings/email" \
    -d '{ "operation": "create", "data": "<h1>Test Headline</h1><p>Test message</p>"}'
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
