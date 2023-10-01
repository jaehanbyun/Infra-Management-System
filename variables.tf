# Openstack Public Network
variable "pool" {
  default = "public"
}

variable "floating_subnet_name" {
  default = "public-subnet"
}

# Flavor for master nodes
variable "master_flavor_name" {
  default = "k8s"
}
# Flavor for worker nodes
variable "worker_flavor_name" {
  default = "k8s"
}
# Flavor for bastion nodes
variable "bastion_flavor_name" {
  default = "k8s"
}
# Image for master nodes: Ubuntu 20.04
variable "master_image_uuid" {
  default = "d65a9810-a5b3-426b-9900-e8b326dc3dca"
}
# Image for worker nodes: Ubuntu 20.04
variable "worker_image_uuid" {
  default = "d65a9810-a5b3-426b-9900-e8b326dc3dca"
}
# Image for bastion nodes: Ubuntu 20.04
variable "bastion_image_uuid" {
  default = "d65a9810-a5b3-426b-9900-e8b326dc3dca"
}

variable "ssh_user_name" {
  default = "ubuntu"
}

variable "number_of_master_nodes" {
  type    = number
  default = 1
}

variable "number_of_worker_nodes" {
  type    = number
  default = 1
}

variable "public_key_path" {
  description = "Path to the public key"
  type        = string
}

terraform {
  required_version = ">= 0.14.0"
  required_providers {
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 1.51.1"
    }
  }
  backend "swift" {}
}

variable "cluster_name" {}
variable "os_username" {}
variable "os_project_name" {}
variable "os_password_input" {}
variable "os_auth_url" {}
variable "os_region_name" {}
# variable "ssh_key_file" {}

