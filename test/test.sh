
key=$(base64 < key.asc)
json="{\"version\":\"v1.17.4\",\"pubkey\":\"$key\"}"

curl --header "Content-Type: application/json" \
  --request POST \
  --data "$json" \
  https://k8stfw.com/get