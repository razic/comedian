# comedian

`comedian` is a simple HTTP application that presents a joke in plain-text.

## Usage

You can deploy the application to Kubernetes with the following one-liner:

```bash
kubectl apply -f https://raw.githubusercontent.com/razic/comedian/master/contrib/kubernetes.yaml
```

*Note:* Monitoring progress of the Kubernetes deployment, including
load-balancer creation, is out of the scope of this tutorial. Please refer to
the official [Kubernetes] documentation for additional information.

## Design

Under the hood, `comedian` calls out to two external HTTP API's to generate
it's jokes. The external HTTP API's are wrapped into their own [gRPC] services,
which is consumed by `comedian`:

* [uinames](api/services/uinames/uinames.proto)
* [icndb](api/services/icndb/icndb.proto)

[Kubernetes]: https://kubernetes.io/
[gRPC]: http://www.grpc.io
