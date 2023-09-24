#### CREATE NETWORKS ###
data "openstack_networking_network_v2" "k8s_network" {
  name = var.pool
}

data "openstack_networking_subnet_v2" "floating_subnet" {
  name = var.floating_subnet_name
}

resource "openstack_networking_network_v2" "k8s_network" {
  name           = "k8s_network"
  admin_state_up = true
}

# Create Subnet
resource "openstack_networking_subnet_v2" "k8s_subnet" {
  name            = "k8s_subnet"
  network_id      = openstack_networking_network_v2.k8s_network.id
  cidr            = "10.0.0.0/24"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# Create Router
resource "openstack_networking_router_v2" "k8s_router" {
  name                = "k8s_router"
  admin_state_up      = true
  external_network_id = data.openstack_networking_network_v2.k8s_network.id
}

# Connect Network Elements
resource "openstack_networking_router_interface_v2" "k8s_network" {
  router_id = openstack_networking_router_v2.k8s_router.id
  subnet_id = openstack_networking_subnet_v2.k8s_subnet.id
}

# Add floating IP to project
resource "openstack_compute_floatingip_v2" "bastion_ip" {
  pool = var.pool
}

