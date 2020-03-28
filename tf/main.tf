module "cluster" {
  source = "git::https://github.com/poseidon/typhoon//digital-ocean/container-linux/kubernetes?ref=CLUSTER_VERSION"

  # Digital Ocean
  cluster_name = var.cluster_name
  region       = "nyc3"
  dns_zone     = "k8stfw.com"
  image        = "coreos-stable"

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
