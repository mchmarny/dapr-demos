# hardened demo 

In addition to support for Kubernetes namespace isolation and Role-Based Access Control (RBAC) authorization, Dapr also provides additional, more granular, controls to harden applications deployment in Kubernetes. Some security related features, like in-transit encryption for all sidecar-to-sidecar communication using mutual TLS, are enabled by default. Others, require opt-in. This demo will overview: 

* Token-based authentication on Dapr APIs exposed on cluster ingress
* Cross-namespace service invocation with [SPIFFE](https://spiffe.io/) identity verification 
* Trust relationship management using logical domain groups (can span multiple namespaces)
* Per operation verb access control settings (e.g. deny all except `POST` from `app2` on `/op1`)
* Component scoping (i.e. which app should be able to access a given component)
* Pub/Sub topic scoping (i.e. which app should be able to publish or subscriber to a given topic)
* Secret access control per application (i.e. which secrets the app should be able to access)

![](img/overview.png)

> You can replicate this demo on any Kubernetes cluster configured with Dapr. To demo the cross-namespace service invocation with external API gateway you will need "dapr'ized' cluster ingress (ingress with Dapr sidecar). You can setup fully configured Dapr cluster with all these dependencies using included [Dapr cluster setup](../setup#dapr-cluster-setup).

## Setup 

In Kubernetes, [namespaces](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/) provide a way to divide cluster resources between multiple users or applications. To isolate all the microservices in this demo, first, create a namespace on your cluster.

> For purposes of this demo, the namespace will be called `hardened` but you can choose your own name.

```shell
kubectl create namespace hardened
```

Also, to illustrate Dapr component scoping (e.g. PubSub and State), this demo will use in-cluster Redis deployment (see [Redis setup](../setup#usage)). To showcase the declarative access control for applications over secrets this demo will use `redis-secret` defined in the `hardened` namespace.

```shell
kubectl create secret generic redis-secret \
    --from-literal=password="${REDIS_PASS}" \
    -n hardened 
```

> If this is Redis on your cluster you can look it up using `kubectl get svc nginx-ingress-nginx-controller -o jsonpath='{.status.loadBalancer.ingress[0].ip}'` and define the `REDIS_PASS` environment variable with that secret. 

Also, create another secret to demonstrate later how Dapr controls application's access to secrets.

```shell
kubectl create secret generic test-secret --from-literal=test="test" -n hardened 
```

## Deploy

With the namespace configured and the Redis password created, it's time to deploy:

* [app1.yaml](./k8s/app1.yaml), [app2.yaml](./k8s/app1.yaml), and [app2.yaml](./k8s/app1.yaml) are the Kubernetes deployments with their Dapr configuration.
* [pubsub.yaml](./k8s/pubsub.yaml) and [state.yaml](./k8s/state.yaml) are the configuration files for PubSub and State components using Redis
* [role.yaml](./k8s/role.yaml) defines the Role and RoleBinding required for Dapr application access the Kubernetes secrets in the `hardened` namespace.

> This demo uses [prebuilt application images](https://github.com/mchmarny?tab=packages&q=hardened-app). You can review the code for these 3 applications in the [src](./src) directory.

Now, apply the demo resources to the cluster.

```shell
kubectl apply -f k8s/ -n hardened
```

The response from the above command should confirm that all the resources were configured.

```shell
deployment.apps/app1 configured
configuration.dapr.io/app1-config configured
deployment.apps/app2 configured
configuration.dapr.io/app2-config configured
deployment.apps/app3 configured
configuration.dapr.io/app3-config configured
component.dapr.io/pubsub configured
role.rbac.authorization.k8s.io/secret-reader configured
rolebinding.rbac.authorization.k8s.io/dapr-secret-reader configured
component.dapr.io/state configured
```

## Verify 

To ensure the rest of the demo goes smoothly, check that everything was deployed correctly.

```shell
kubectl get pods -n hardened
```

If everything went well, the response should include `app1`, `app2`, and `app3` pods with the status `Running` and the ready state of `2/2` indicating that the Dapr sidecar has been injected and components successfully loaded.

```shell
NAME                    READY   STATUS    RESTARTS   AGE
app1-6df587fb45-k46sz   2/2     Running   0          40s
app2-685fd94f69-5vkwl   2/2     Running   0          40s
app3-6d57778cbd-mxn2k   2/2     Running   0          40s
```

## Demo 

The Dapr API exposed on the cluster ingress is protected with [token authentication](https://github.com/dapr/docs/tree/master/howto/enable-dapr-api-token-based-authentication#enable-dapr-apis-token-based-authentication). Start by exporting that token from the cluster secret to allow for API invocation in this demo.

```shell
export API_TOKEN=$(kubectl get secret dapr-api-token -o jsonpath="{.data.token}" | base64 --decode)
```

### Service Invocation

The app identity and its access control within Dapr as controlled using [policies](https://github.com/dapr/docs/blob/master/howto/allowlists-serviceinvocation/README.md) which are defined in the app configuration. To "attach" [configuration](https://github.com/dapr/docs/blob/master/howto/allowlists-serviceinvocation/README.md), the app deployment template has to be annotated with the name of the configuration:

```yaml
annotations:
  dapr.io/config: "app1-config"
```

In this demo, to allow only the Dapr'ized NGNX ingress to invoke the `/ping` method on [app1.yaml](./k8s/app1.yaml), the default action is set to `deny` and an explicit policy created for `nginx-ingress` in the `default` namespace which also, first denies access to all methods on that app, and only then allows access on the `/ping` method (aka operation) when the HTTP verb is `POST`. 

```yaml
accessControl:
  defaultAction: deny
  trustDomain: "hardened"
  policies:
  - appId: nginx-ingress
    defaultAction: deny 
    trustDomain: "public"
    namespace: "default"
    operations:
    - name: /ping
      httpVerb: ["POST"] 
      action: allow
```

To demo this now, invoke the `ping` method on `app1` in the `hardened` namespace using the Dapr API exposed on the NGNX ingress.

> The [Dapr cluster setup](../setup#dapr-cluster-setup) includes custom domain and TLS certificate support. This demo users `thingz.io` domain and a wildcard certificates for al (`*`) subdomains.

```shell
curl -i -d '{ "message": "hello" }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     https://api.thingz.io/v1.0/invoke/app1.hardened/method/ping
```

Dapr should respond with HTTP status code `200` as well as parent trace ID for this invocation (`traceparent`) in the header, and a JSON payload with the number of API invocations and nano epoch timestamp.

> The count of API invocations is persisted in the Dapr sate store configured in [State component](./k8s/state.yaml)

```shell
HTTP/2 200
date: Sun, 25 Oct 2020 12:05:56 GMT
content-type: text/plain; charset=utf-8
content-length: 39
traceparent: 00-ecbbc473826b3e328ea00f5ac0ce222b-0824d3896092d8ce-01
strict-transport-security: max-age=15724800; includeSubDomains

{ "on": 1603627556200126373, "count": 8 }
```

To demo the active access policy, try also to invoke the `counter` method on `app2` in the `hardened` namespace.

```shell
curl -i -d '{ "on": 1603627556200126373, "count": 2 }' \
     -H "Content-type: application/json" \
     -H "dapr-api-token: ${API_TOKEN}" \
     https://api.thingz.io/v1.0/invoke/app2.hardened/method/counter
```

That invocation will result in an error. The response will include `PermissionDenied` message:

```json
{
  "errorCode": "ERR_DIRECT_INVOKE",
  "message": "rpc error: code = PermissionDenied desc = access control policy has denied access to appid: app2 operation: ping verb: POST"
}
```

The access control defined above applies also to in-cluster invocation ([app2.yaml](./k8s/app2.yaml)). Where the additional `trustDomain` setting on `app2` configuration is used to only allow access to invoke the `/counter` method when the calling app is `app1`:

```yaml
policies:
  - appId: app1
    defaultAction: deny 
    trustDomain: "hardened"
    namespace: "hardened"
    operations:
    - name: /counter
      httpVerb: ["POST"] 
      action: allow
```

To demo this, forward local port to any other Dapr sidecar besides `app1` in that cluster.

```shell
kubectl port-forward deployment/app2 3500 -n hardened
```

And then try to invoke the `/ping` method on the `app1`. That too will result in `PermissionDenied` message. 

```shell
curl -i -d '{ "message": "hello" }' \
     -H "Content-type: application/json" \
     http://localhost:3500/v1.0/invoke/app1/method/ping
```

> In this configuration, all invocations that are not explicitly permitted in Dapr access policy will be denied!

### Topic Publishing and Subscription 

Just like in case of invocation, access to components in Dapr is also driven by configuration. The [pubsub](./k8s/pubsub.yaml) component in this demo is scoped to only be accessible by `app2` and `app3`:

```yaml
scopes:
- app2
- app3
```

The topic access of the PubSub component is further defined by the `publishingScopes` and `subscriptionScopes` lists. In this case `app2` can only publish, and the `app3` can only subscribe to the `messages` topic:

```yaml
- name: publishingScopes
  value: "app2=messages"
- name: subscriptionScopes
  value: "app3=messages"
```

To demo this, while still forwarding local port to the `app2` pod, try publish to any other topic besides `messages`.

```shell
curl -i -d '{ "message": "test" }' \
     -H "Content-type: application/json" \
     http://localhost:3500/v1.0/publish/pubsub/test
```

The above publish will result in error:

```json
{
  "errorCode": "ERR_PUBSUB_PUBLISH_MESSAGE",
  "message": "topic test is not allowed for app id app2"
}
```

You can also try to subscribe to the `messages` topic or even forward port to `app3` and try to publish to the valid topic there, and still receive the same error, because that application is only allowed to subscribe to the `messages` topic, not publish to it.

### Secrets 

Application access to secrets within Dapr is also driven by configuration. In this demo, the `app2` for example, has its secrets configuration defined as follow: `deny` this application's access to all secrets except `redis-secret`: 

```yaml
secrets:
  scopes:
    - storeName: kubernetes
      defaultAccess: deny
      allowedSecrets: ["redis-secret"]
```        

To demo this, while still forwarding local port to the `app2` pod, try access the other secret we created in the `hardened` namespace during setup: `test-secret`.

```shell
curl -i http://localhost:3500/v1.0/secrets/kubernetes/test-secret
```

The above query will result in `403 Forbidden` as the `test-secret` secret is not listed in the `allowedSecrets` list and the `defaultAccess` is set to `deny`.

```json
{
  "errorCode": "ERR_PERMISSION_DENIED", 
  "message": "Access denied by policy to get test-secret from kubernetes"
}
```

## Summary 

This demo illustrated just a few of the options that Dapr provides to harden application deployments. For more security-related information (including network, threat model, and latest security audit) see the [Security section](https://github.com/dapr/docs/blob/master/concepts/security/README.md) in Dapr documentation. 

## Restarts

If you update components you may have to restart the deployments.

```shell
kubectl rollout restart deployment/app1 -n hardened
kubectl rollout restart deployment/app2 -n hardened
kubectl rollout restart deployment/app3 -n hardened
kubectl rollout status deployment/app1 -n hardened
kubectl rollout status deployment/app2 -n hardened
kubectl rollout status deployment/app3 -n hardened
```

## Cleanup

```shell
kubectl delete -f k8s/ -n hardened
kubectl delete secret redis-secret -n hardened
kubectl delete ns hardened
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.

## License

This software is released under the [MIT](../LICENSE)
