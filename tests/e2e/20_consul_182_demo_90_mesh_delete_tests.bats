load 'helpers'

setup() {
    NAMESPACE=consul-e2e-tests
}

@test "consul 1.8.2 demo deletion should be successful" {
  DELETE_CONSUL=$(cat <<EOT
{
  "opName": "consul_182_demo",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": true,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$DELETE_CONSUL' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "no resources should remain in the consul namespace" {
  wait_until_namespace_empty "$NAMESPACE"
  [ "$?" -eq 0 ]
}