#load 'helpers'

setup() {
    NAMESPACE=bookinfo
}

@test "bookinfo installation should be successful" {
  INSTALL_BOOKINFO=$(cat <<EOT
{
  "opName": "install_book_info",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": false,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$INSTALL_BOOKINFO' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "deployment/details-v1 should be ready" {
  run kubectl rollout status deployment/details-v1 -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/ratings-v1 should be ready" {
  run kubectl rollout status deployment/ratings-v1 -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/reviews-v1 should be ready" {
  run kubectl rollout status deployment/reviews-v1 -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/reviews-v2 should be ready" {
  run kubectl rollout status deployment/reviews-v2 -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/reviews-v3 should be ready" {
  run kubectl rollout status deployment/reviews-v3 -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/productpage-v1 should be ready" {
  run kubectl rollout status deployment/productpage-v1 -n $NAMESPACE
  [ "$status" -eq 0 ]
}

@test "bookinfo deletion should be successful" {
  DELETE_BOOKINFO=$(cat <<EOT
{
  "opName": "install_book_info",
  "namespace": "$NAMESPACE",
  "username": "",
  "customBody": "",
  "deleteOp": true,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$DELETE_BOOKINFO' | grpcurl --plaintext -d @ $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "no resources should remain in the bookinfo namespace" {
  sleep 60
  run bash -c "kubectl get all -n $NAMESPACE -o json | jq -j '.items | length'"
  [ "$status" -eq 0 ]
  [ "$output" = "0" ]
}