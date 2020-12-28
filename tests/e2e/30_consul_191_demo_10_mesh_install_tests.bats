#load 'helpers'

setup() {
    CONSUL_NAMESPACE=consul-e2e-tests
    CONSUL_OP_NAME=consul_191_demo
}

@test "consul 1.9.1 demo installation should be successful" {
  INSTALL_CONSUL=$(cat <<EOT
{
  "opName": "$CONSUL_OP_NAME",
  "namespace": "$CONSUL_NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": false,
  "operationId": "testid_consul_191_demo"
}
EOT
)
  run bash -c "echo '$INSTALL_CONSUL' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
  # This operation returns a JSON map if successful. An error doesn't return a JSON object, unless the implementation in
  # api/grpc/handlers.go:ApplyOperation is changed (see TODO there).
  [[ $(echo $output | jq -j ".operationId") = "testid_consul_191_demo" ]]
}

@test "deployment/consul-consul-connect-injector-webhook-deployment should be ready" {
  run kubectl rollout status deployment/consul-consul-connect-injector-webhook-deployment -n $CONSUL_NAMESPACE
  [ "$status" -eq 0 ]
}

@test "statefulset/consul-consul-server should be ready" {
  run kubectl rollout status statefulset/consul-consul-server -n $CONSUL_NAMESPACE
  [ "$status" -eq 0 ]
}

@test "daemonset/consul-consul should be ready" {
  run kubectl rollout status daemonset/consul-consul -n $CONSUL_NAMESPACE
  [ "$status" -eq 0 ]
}
