load 'helpers'

setup() {
    NAMESPACE=httpbin
}

@test "httpbin installation should be successful" {
  INSTALL_HTTPBIN=$(cat <<EOT
{
  "opName": "httpbin",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": false,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$INSTALL_HTTPBIN' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "deployment/httpbin should be ready" {
  run kubectl rollout status deployment/httpbin -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "httpbin deletion should be successful" {
  DELETE_HTTPBIN=$(cat <<EOT
{
  "opName": "httpbin",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": true,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$DELETE_HTTPBIN' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "no resources should remain in the httpbin namespace" {
  wait_until_namespace_empty "$NAMESPACE"
  [ "$?" -eq 0 ]
}