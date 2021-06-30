load 'helpers'

setup() {
    NAMESPACE=httpbin-custom
    CUSTOM_BODY='apiVersion: v1\nkind: Service\nmetadata:\n  name: httpbin\n  labels:\n    app: httpbin\nspec:\n  type: LoadBalancer\n  ports:\n  - name: http\n    port: 8000\n    targetPort: 80\n  selector:\n    app: httpbin\n---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: httpbin\nspec:\n  replicas: 1\n  selector:\n    matchLabels:\n      app: httpbin\n  template:\n    metadata:\n      labels:\n        app: httpbin\n        version: v1\n      annotations:\n        \"consul.hashicorp.com/connect-inject\": \"true\"\n    spec:\n      containers:\n      - image: docker.io/kennethreitz/httpbin\n        imagePullPolicy: IfNotPresent\n        name: httpbin\n        ports:\n        - containerPort: 80\n\n'
}

@test "custom-operation: httpbin installation should be successful" {
  INSTALL_HTTBIN_CUSTOM=$(cat <<EOT
{
  "opName": "custom",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "$CUSTOM_BODY",
  "deleteOp": false,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$INSTALL_HTTBIN_CUSTOM' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "custom-operation: deployment/httpbin should be ready" {
  run kubectl rollout status deployment/httpbin -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "custom-operation: httpbin deletion should be successful" {
  DELETE_HTTPBIN=$(cat <<EOT
{
  "opName": "custom",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "$CUSTOM_BODY",
  "deleteOp": true,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$DELETE_HTTPBIN' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "custom-operation: no resources should remain in the httpbin namespace" {
  wait_until_namespace_empty "$NAMESPACE"
  [ "$?" -eq 0 ]
}