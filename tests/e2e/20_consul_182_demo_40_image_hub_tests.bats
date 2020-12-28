#load 'helpers'

setup() {
    NAMESPACE=image-hub
}

@test "image-hub installation should be successful" {
  INSTALL_IMAGE_HUB=$(cat <<EOT
{
  "opName": "imagehub",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": false,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$INSTALL_IMAGE_HUB' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "deployment/web-deployment should be ready" {
  run kubectl rollout status deployment/web-deployment -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/api-deployment-v1 should be ready" {
  run kubectl rollout status deployment/api-deployment-v1 -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/envoy should be ready" {
  run kubectl rollout status deployment/envoy -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "image-hub deletion should be successful" {
  DELETE_IMAGE_HUB=$(cat <<EOT
{
  "opName": "imagehub",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": true,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$DELETE_IMAGE_HUB' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "no resources should remain in the image-hub namespace" {
  sleep 50
  run bash -c "kubectl get all -n $NAMESPACE -o json | jq -j '.items | length'"
  [ "$status" -eq 0 ]
  [ "$output" = "0" ]
}