
key=$(base64 < key.asc)
json="{\"version\":\"v1.17.4\",\"pubkey\":\"$key\"}"

curl --header "Content-Type: application/json" \
  --request POST \
  --data "$json" \
  http://localhost:8080/get