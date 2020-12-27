The manifest `consul-1.8.2-demo.yaml` was generated with Helm 3 using
``` 
helm template consul -f consul-values-1.8.2-demo -n consul hashicorp/consul --version 0.24.1 > consul-new.yaml
```

Then, `namespace: consul` was replaced with `namespace: {{.namespace}}` using sed: 
```
sed -E 's/^( +)namespace: +consul *$/\1namespace: {{.namespace}}/g' consul-new.yaml > consul-1.8.2-demo.yaml
```

This makes it possible to deploy Consul to the namespace specified in the Meshery UI.
