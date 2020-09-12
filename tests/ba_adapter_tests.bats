load 'helpers'

@test "adapter should return expected supported operations" {
  run grpcurl --plaintext localhost:10002 meshes.MeshService.SupportedOperations
  [ "$status" -eq 0 ]
  local r=$(echo $output | jq -r '.ops[] | .key')
  assert_contains "consul_install" $r
  assert_contains "install_book_info" $r
  assert_contains "install_http_bin" $r
  assert_contains "install_image_hub" $r
  assert_contains "custom" $r
}

@test "adapter should return expected mesh name" {
  run bash -c "grpcurl --plaintext localhost:10002 meshes.MeshService.MeshName | jq -j .name"
  [ "$status" -eq 0 ]
  [ "$output" = "Consul" ]
}
