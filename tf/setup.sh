#!/bin/bash

# Get supported versions from typhoon
echo "Supported versions"
curl https://api.github.com/repos/poseidon/typhoon/tags|jq  '.[].name'

# Figure out which os we are running on
case "$OSTYPE" in
  darwin*)  os="darwin" ;; 
  linux*)   os="linux" ;;
  *)        echo "unknown: $OSTYPE" && exit 1;;
esac

# Download the typhoon ct provider
wget https://github.com/poseidon/terraform-provider-ct/releases/download/v0.4.0/terraform-provider-ct-v0.4.0-${os}-amd64.tar.gz
tar xzf terraform-provider-ct-v0.4.0-${os}-amd64.tar.gz
mkdir -p  ~/.terraform.d/plugins/
mv terraform-provider-ct-v0.4.0-${os}-amd64/terraform-provider-ct ~/.terraform.d/plugins/terraform-provider-ct_v0.4.0
rm -r terraform-provider-ct*

# Get ssh fingerprint
ssh-add ~/.ssh/id_rsa
ssh-add -L
export TF_VAR_ssh_fingerprint=$(ssh-add -l -E md5| awk '{print $2}'|cut -d':' -f 2-)

# Add ssh pub key to do
export auth="Authorization: Bearer "$TF_VAR_do_token
export payload="{\"name\":\"My SSH Public Key\",\"public_key\":\"$(cat ~/.ssh/id_rsa.pub)\"}"
curl -X POST -H "Content-Type: application/json" -H "$auth" -d "$payload" "https://api.digitalocean.com/v2/account/keys" 

terraform init
terraform plan