The manifest consul.yaml was generated with Helm 3 using
``` 
helm template consul -f consul-values.yaml -n consul hashicorp/consul --version 0.24.1 > consul-new.yaml
```

Then, `namespace: consul` was replaced with `namespace: {{.namespace}}` using sed: 
```
sed -E 's/^( +)namespace: +consul *$/\1namespace: {{.namespace}}/g' consul-new.yaml > consul.yaml
```

This makes it possible to deploy Consul to the namespace specified in the Meshery UI.

Note: Helm support in this adapter is planned.