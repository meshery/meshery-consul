load 'helpers'

setup() {
    BOOKINFO_NAMESPACE=bookinfo
    BOOKINFO_OP_NAME=bookinfo
}

@test "bookinfo installation should be successful" {
  INSTALL_BOOKINFO=$(cat <<EOT
{
  "opName": "$BOOKINFO_OP_NAME",
  "namespace": "$BOOKINFO_NAMESPACE",
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
  run kubectl rollout status deployment/details-v1 -n $BOOKINFO_NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/ratings-v1 should be ready" {
  run kubectl rollout status deployment/ratings-v1 -n $BOOKINFO_NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/reviews-v1 should be ready" {
  run kubectl rollout status deployment/reviews-v1 -n $BOOKINFO_NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/reviews-v2 should be ready" {
  run kubectl rollout status deployment/reviews-v2 -n $BOOKINFO_NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/reviews-v3 should be ready" {
  run kubectl rollout status deployment/reviews-v3 -n $BOOKINFO_NAMESPACE
  [ "$status" -eq 0 ]
}

@test "deployment/productpage-v1 should be ready" {
  run kubectl rollout status deployment/productpage-v1 -n $BOOKINFO_NAMESPACE
  [ "$status" -eq 0 ]
}

@test "bookinfo deletion should be successful" {
  DELETE_BOOKINFO=$(cat <<EOT
{
  "opName": "$BOOKINFO_OP_NAME",
  "namespace": "$BOOKINFO_NAMESPACE",
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
  wait_until_namespace_empty "$BOOKINFO_NAMESPACE"
  [ "$?" -eq 0 ]
}