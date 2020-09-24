load 'helpers'

@test "adapter should return expected supported operations" {
  run grpcurl --plaintext $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.SupportedOperations
  [ "$status" -eq 0 ]

  [[ $(echo $output | jq '.ops[] | select( .key == "consul_install" )' | jq -j .key ) = "consul_install" ]]
  [[ $(echo $output | jq '.ops[] | select( .key == "install_book_info" )' | jq -j .key ) = "install_book_info" ]]
  [[ $(echo $output | jq '.ops[] | select( .key == "install_http_bin" )' | jq -j .key ) = "install_http_bin" ]]
  [[ $(echo $output | jq '.ops[] | select( .key == "install_image_hub" )' | jq -j .key ) = "install_image_hub" ]]
  [[ $(echo $output | jq '.ops[] | select( .key == "custom" )' | jq -j .key ) = "custom" ]]
}

@test "adapter should return expected mesh name" {
  run bash -c "grpcurl --plaintext $MESHERY_ADAPTER_ADDR:10002 meshes.MeshService.MeshName | jq -j .name"
  [ "$status" -eq 0 ]
  [ "$output" = "Consul" ]
}
