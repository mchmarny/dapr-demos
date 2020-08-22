# Dapr integrations demo

> WIP: This demo is being currently updated to Dapr v0.10

## Use-case

Order cancellation demo, loosely based on [Shopify API](https://shopify.dev/docs/admin-api/rest/reference/orders/order?api[version]=2020-04) order cancellation use-case to showcase multiple Dapr service integrations in a single solution: 

![Use-case](img/usecase.png)

## Component Overview 

* **Dapr API** endpoint published with JWT token auth in **Daprized Ngnx** ingress
* **Dapr Workflows** to orchestrate cancellation process using **Logic Apps** runtime
* **Dapr Functions** extensions to create and persist audit state into **Mongo DB**
* **Dapr Eventing** using **Redis Queue** for order message queue
* **Daprized Web App** as order processing dashboard
* **Dapr Binding** to send confirmation emails using **SendGrid**
* **Dapr Distributed Tracing** to capture and forward traces to **Zipkin**

![Draft Demo Flow Diagram](img/diagram.png)

## Demo 

For script through this demo see [demo](./demo).



