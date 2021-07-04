load 'helpers'

setup() {
    CONSUL_NAMESPACE=consul-e2e-tests
    CONSUL_OP_NAME=consul_191_demo
}

@test "consul 1.9.1 demo deletion should be successful" {
  DELETE_CONSUL=$(cat <<EOT
{
  "opName": "$CONSUL_OP_NAME",
  "namespace": "$CONSUL_NAMESPACE",
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
  wait_until_namespace_empty "$CONSUL_NAMESPACE"
  [ "$?" -eq 0 ]
}