## FN-KUBE-STATE

A small utility to show the state of k8s cluster


#### Local testing steps (For Mac):

Download kind

```bash
brew install kind
```

Create cluster
```bash
kind create cluster --name fn-kubestate-cl
```

Set and get the kubectl context
```
kubectl config set-context kind-fn-kubestate-cl
kubectl cluster-info --context kind-fn-kubestate-cl
```

Quick check to see the nodes
```
kubectl get nodes
```

Add test deployments
```
kubectl apply -f services.yaml
```

Check the pods
```
kubectl get pods
```

Run the go file

```
go mod tidy
go run main.go
``

