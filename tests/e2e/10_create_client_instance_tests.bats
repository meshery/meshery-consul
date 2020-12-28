#load 'helpers'

setup() {
    KUBECTL_CONTEXT=$(kubectl config current-context)
    KUBECTL_CONFIG=$(kubectl config view --raw | base64 | tr -d '\n')
    CREATE_MESH_REQ_MSG=$(cat <<EOT
{
  "k8sConfig": "$KUBECTL_CONFIG",
  "contextName": "$KUBECTL_CONTEXT"
}
EOT
)
}

@test "client instance should be created" {
  run bash -c "echo '$CREATE_MESH_REQ_MSG' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.CreateMeshInstance | jq ."
  [ "$status" -eq 0 ]
  # this operation returns an empty JSON map. an error will not return a JSON object.
  [[ $(echo $output | jq length ) = "0" ]]
}
