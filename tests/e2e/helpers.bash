assert_exists() {
  assert [ -e "$1" ]
}

refute_exists() {
  assert [ ! -e "$1" ]
}

trim() {
    # https://stackoverflow.com/questions/369758/how-to-trim-whitespace-from-a-bash-variable
    local var="$*"
    # remove leading whitespace characters
    var="${var#"${var%%[![:space:]]*}"}"
    # remove trailing whitespace characters
    var="${var%"${var##*[![:space:]]}"}"
    echo -n "$var"
}

# Assert that $1 is contained in list $2. list entries are trimmed for leading and trailing whitespace before comparison.
assert_contains() {
  local item
  for item in "${@:2}"; do
    trimmed=$(trim "$item")
    if [ "$trimmed" = "$1" ]; then
      return 0
    fi
  done
  return 1
}

# Assert that $1 is NOT contained in list $2. list entries are trimmed for leading and trailing whitespace before comparison.
assert_not_contains() {
  local item
  for item in "${@:2}"; do
    trimmed=$(trim "$item")
    if [ "$trimmed" = "$1" ]; then
      return 1
    fi
  done
  return 0
}

wait_until_namespace_empty() {
  local max_sleep_count=40
  local sleep_count=0
  local resource_count=1
  until [ "$sleep_count" -ge $max_sleep_count ] || [ "$resource_count" -eq 0 ]; do
    sleep 2
    ((sleep_count=sleep_count+1))
    resource_count=$(kubectl get all -n $1 -o json | jq -j '.items | length')
    echo "# resource_count: ${resource_count}, sleep_count: ${sleep_count}" >&3
  done
  return "$resource_count"
}