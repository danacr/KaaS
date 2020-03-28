module "cluster" {
  source = "git::https://github.com/poseidon/typhoon//digital-ocean/container-linux/kubernetes?ref=v1.17.4"

  # Digital Ocean
  cluster_name = "seventeen"
  region       = "nyc3"
  dns_zone     = "k8stfw.com"
  image        = "coreos-stable"
  # controller_type = "s-4vcpu-8gb"
  # worker_type     = "s-2vcpu-2gb"

  # configuration
  ssh_fingerprints = [var.ssh_fingerprint]

  # optional
  worker_count = 2
}

# Obtain cluster kubeconfig
resource "local_file" "kubeconfig-cluster" {
  content  = module.cluster.kubeconfig-admin
  filename = "./cluster-config"
}
