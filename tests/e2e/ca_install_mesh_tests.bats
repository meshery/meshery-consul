#load 'helpers'

setup() {
    NAMESPACE=consul-e2e-tests
    KUBECTL_CONTEXT=$(kubectl config current-context)
    KUBECTL_CONFIG=$(kubectl config view | base64 | tr -d '\n')
    CREATE_MESH_REQ_MSG=$(cat <<EOT
{
  "k8sConfig": "$KUBECTL_CONFIG",
  "contextName": "$KUBECTL_CONTEXT"
}
EOT
)
}

@test "client instance should be created" {
  run bash -c "echo '$CREATE_MESH_REQ_MSG' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.CreateMeshInstance"
  [ "$status" -eq 0 ]
}

@test "consul_install should be successful" {
  INSTALL_CONSUL=$(cat <<EOT
{
  "opName": "consul_install",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": false,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$INSTALL_CONSUL' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "deployment/consul-consul-connect-injector-webhook-deployment should be ready" {
  run kubectl rollout status deployment/consul-consul-connect-injector-webhook-deployment -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "statefulset/consul-consul-server should be ready" {
  run kubectl rollout status statefulset/consul-consul-server -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "daemonset/consul-consul should be ready" {
  run kubectl rollout status daemonset/consul-consul -n $NAMESPACE
  [ "$status" -eq 0 ]
}
