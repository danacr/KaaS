#!/bin/sh

eval `ssh-agent -s`
ssh-keygen -f /root/.ssh/id_rsa -q -N ""

# Get ssh fingerprint
ssh-add /root/.ssh/id_rsa
export TF_VAR_ssh_fingerprint=$(ssh-add -l -E md5| awk '{print $2}'|cut -d':' -f 2-)

# Add ssh pub key to do
export auth="Authorization: Bearer "$TF_VAR_do_token
export payload="{\"name\":\"My SSH Public Key\",\"public_key\":\"$(cat ~/.ssh/id_rsa.pub)\"}"
curl -X POST -H "Content-Type: application/json" -H "$auth" -d "$payload" "https://api.digitalocean.com/v2/account/keys" 

terraform init
terraform plan