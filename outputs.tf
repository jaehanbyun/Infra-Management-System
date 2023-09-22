output "k8s_subnet_id" {
  value = "${openstack_networking_subnet_v2.k8s_subnet.id}"
}
output "bastion_fips" {
  value = "${openstack_compute_floatingip_v2.bastion_ip[*].address}"
}

