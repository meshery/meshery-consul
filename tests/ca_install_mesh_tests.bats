#load 'helpers'

setup() {
    KUBECTL_CONTEXT=$(kubectl config current-context)
    KUBECTL_CONFIG=$(kubectl config view | base64 -w 0 -)
    CREATE_MESH_REQ_MSG=$(cat <<EOT
{
  "k8sConfig": "$KUBECTL_CONFIG",
  "contextName": "$KUBECTL_CONTEXT"
}
EOT
)
}

@test "client instance should be created" {
  run bash -c "echo '$CREATE_MESH_REQ_MSG' | grpcurl --plaintext -d @ localhost:10002 meshes.MeshService.CreateMeshInstance"
  [ "$status" -eq 0 ]
}

@test "consul_install should be successful" {
  INSTALL_CONSUL=$(cat <<EOT
{
  "opName": "consul_install",
  "namespace": "consul-e2e-tests",
  "username": "",
  "customBody": "",
  "deleteOp": false,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$INSTALL_CONSUL' | grpcurl --plaintext -d @ localhost:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "all pods should be ready" {
  run kubectl wait pod --all -n consul-e2e-tests --for=condition=Ready --timeout=60s
  [ "$status" -eq 0 ]
}