load 'helpers'

@test "adapter should return expected supported operations" {
  run grpcurl --plaintext $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.SupportedOperations
  [ "$status" -eq 0 ]

  [[ $(echo $output | jq '.ops[] | select( .key == "consul_182_demo" )' | jq -j .key ) = "consul_182_demo" ]]
  [[ $(echo $output | jq '.ops[] | select( .key == "bookinfo" )' | jq -j .key ) = "bookinfo" ]]
  [[ $(echo $output | jq '.ops[] | select( .key == "httpbin" )' | jq -j .key ) = "httpbin" ]]
  [[ $(echo $output | jq '.ops[] | select( .key == "imagehub" )' | jq -j .key ) = "imagehub" ]]
  [[ $(echo $output | jq '.ops[] | select( .key == "custom" )' | jq -j .key ) = "custom" ]]
}

@test "adapter should return expected mesh name" {
  run bash -c "grpcurl --plaintext $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.MeshName | jq -j .name"
  [ "$status" -eq 0 ]
  [ "$output" = "Consul" ]
}
