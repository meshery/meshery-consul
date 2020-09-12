@test "consul_install delete should be successful" {
  INSTALL_CONSUL=$(cat <<EOT
{
  "opName": "consul_install",
  "namespace": "consul-e2e-tests",
  "username": "",
  "customBody": "",
  "deleteOp": true,
  "operationId": ""
}
EOT
)
  run bash -c "echo '$INSTALL_CONSUL' | grpcurl --plaintext -d @ localhost:10002 meshes.MeshService.ApplyOperation"
  [ "$status" -eq 0 ]
}

@test "no resources should remain in the test namespace" {
  sleep 30
  run bash -c "kubectl get all -n consul-e2e-tests -o json | jq -j '.items | length'"
  [ "$status" -eq 0 ]
  [ "$output" = "0" ]
}