## install & update Vector
```shell
# add the official vector repository
helm repo add vector https://helm.vector.dev
helm repo update

# Install the Vector to the openobserve namespace (or another monitoring namespace)
helm install vector vector/vector \
  --namespace openobserve \
  -f vector-values.yaml
  
# view pod status
kubectl get pods -n openobserve

# Get the name of the vector pod and look at the logs
kubectl logs -f -n openobserve -l app.kubernetes.io/name=vector  

# load the update configuration
helm upgrade --install vector vector/vector --namespace openobserve -f vector-values.yaml

# remove Vector
helm uninstall vector -n openobserve
```