# Openstack Public Network
variable "pool" {
  default = "public"
}

variable "floating_subnet_name" {
  default = "public-subnet"
}

variable "node_image_uuid" {
  default = "c586196e-c122-4541-90f6-60c49b9f91c6"
}

variable "flavor_vcpu" {
  default = "5"
}

variable "flavor_ram" {
  type    = number
  default = 4096
}

variable "flavor_disk" {
  type    = number
  default = 40
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

# Image for master nodes: Ubuntu 20.04
# variable "master_image_uuid" {
#   default = "d65a9810-a5b3-426b-9900-e8b326dc3dca"
# }
# # Image for worker nodes: Ubuntu 20.04
# variable "worker_image_uuid" {
#   default = "d65a9810-a5b3-426b-9900-e8b326dc3dca"
# }
# # Image for bastion nodes: Ubuntu 20.04
# variable "bastion_image_uuid" {
#   default = "d65a9810-a5b3-426b-9900-e8b326dc3dca"
# }
