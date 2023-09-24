output "k8s_network_id" {
  value = openstack_networking_network_v2.k8s_network.id
}

output "k8s_subnet_id" {
  value = openstack_networking_subnet_v2.k8s_subnet.id
}

output "bastion_fips" {
  value = openstack_compute_floatingip_v2.bastion_ip[*].address
}

output "master_ips" {
  value = [for i in openstack_compute_instance_v2.masters : i.access_ip_v4]
  description = "IP Addresses of master nodes"
}

output "worker_ips" {
  value = [for i in openstack_compute_instance_v2.workers : i.access_ip_v4]
  description = "IP Addresses of worker nodes"
}

output "floating_network_id" {
  value = data.openstack_networking_network_v2.k8s_network.id
}

output "floating_subnet_id" {
  value = data.openstack_networking_subnet_v2.floating_subnet.id
}

output "os_username_output" {
  value = var.os_username
}

output "os_project_name_output" {
  value = var.os_project_name
}

output "os_password_input_output" {
  value = var.os_password_input
}

output "os_auth_url_output" {
  value = var.os_auth_url
}

output "os_region_name_output" {
  value = var.os_region_name
}

