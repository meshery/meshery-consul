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
  sleep 30
  run bash -c "kubectl get all -n $NAMESPACE -o json | jq -j '.items | length'"
  [ "$status" -eq 0 ]
  [ "$output" = "0" ]
}