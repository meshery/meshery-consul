@test "gRPCurl should be installed" {
  run grpcurl -version
  [ "$status" -eq 0 ]
}

@test "jq should be installed" {
  run jq --help
  [ "$status" -eq 0 ]
}

@test "base64 should be installed" {
  run which base64
  [ "$status" -eq 0 ]
}

@test "kubectl should be installed" {
  run kubectl config current-context
  [ "$status" -eq 0 ]
  echo '# current kubectl context is:' "$output" >&3
}

@test "cluster should be reachable" {
  run kubectl get ns kube-system
  [ "$status" -eq 0 ]
}

@test "all pods should be ready" {
  run kubectl wait pod --all --all-namespaces --for=condition=Ready --timeout=20s
  [ "$status" -eq 0 ]
}