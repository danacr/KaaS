variable "do_token" {}
variable "ssh_fingerprint" {}

provider "digitalocean" {
  version = "1.14.0"
  token   = var.do_token
}

provider "ct" {
  version = "0.4.0"
}
