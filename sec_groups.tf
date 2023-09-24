# Create Security Group For k8s Instances
resource "openstack_networking_secgroup_v2" "k8s_sec_group" {
  name        = "k8s_sec_group"
  description = "Security group for the k8s instances"
}

# Allow port 22 from local network
resource "openstack_networking_secgroup_rule_v2" "k8s_22" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "10.0.0.0/24"
  security_group_id = openstack_networking_secgroup_v2.k8s_sec_group.id
}

# Allow from local network
resource "openstack_networking_secgroup_rule_v2" "k8s_local_network" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 1
  port_range_max    = 35535
  remote_ip_prefix  = "10.0.0.0/24"
  security_group_id = openstack_networking_secgroup_v2.k8s_sec_group.id
}

# Allow port 80
resource "openstack_networking_secgroup_rule_v2" "k8s_80" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 80
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.k8s_sec_group.id
}

# Allow port 443
resource "openstack_networking_secgroup_rule_v2" "k8s_443" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 443
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.k8s_sec_group.id
}

# Allow icmp from local network
resource "openstack_networking_secgroup_rule_v2" "k8s_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "10.0.0.0/24"
  security_group_id = openstack_networking_secgroup_v2.k8s_sec_group.id
}

# Create Security Group For Deployment Instance
resource "openstack_networking_secgroup_v2" "bastion_sec_group" {
  name        = "bastion_sec_group"
  description = "Security group for the bastion instance"
}

# Allow port 22
resource "openstack_networking_secgroup_rule_v2" "bastion_22" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.bastion_sec_group.id
}

# Allow icmp
resource "openstack_networking_secgroup_rule_v2" "bastion_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.bastion_sec_group.id
}

# Create ssh-key
resource "openstack_compute_keypair_v2" "terraform" {
  name       = "terraform"
  public_key = file("~/.ssh/id_rsa.pub")
}

