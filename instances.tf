resource "openstack_compute_flavor_v2" "k8s_flavor" {
  name      = "k8s-${var.cluster_name}-flavor"
  ram       = var.flavor_ram
  vcpus     = var.flavor_vcpu
  disk      = var.flavor_disk
  is_public = true
}

resource "openstack_compute_instance_v2" "bastion" {
  name            = "k8s-${var.cluster_name}-bastion"
  flavor_name     = openstack_compute_flavor_v2.k8s_flavor.name
  key_pair        = openstack_compute_keypair_v2.terraform.name
  security_groups = [openstack_networking_secgroup_v2.bastion_sec_group.name]

  block_device {
    uuid                  = var.node_image_uuid
    source_type           = "image"
    volume_size           = 15
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    uuid = openstack_networking_network_v2.k8s_network.id
  }

  depends_on = [openstack_networking_subnet_v2.k8s_subnet, openstack_compute_flavor_v2.k8s_flavor]
}

# Associate Floating IP
resource "openstack_compute_floatingip_associate_v2" "floatip" {
  floating_ip = openstack_compute_floatingip_v2.bastion_ip.address
  instance_id = openstack_compute_instance_v2.bastion.id
}


resource "openstack_compute_instance_v2" "masters" {
  count           = var.number_of_master_nodes
  name            = "k8s-${var.cluster_name}-master${count.index + 1}"
  flavor_name     = openstack_compute_flavor_v2.k8s_flavor.name
  key_pair        = openstack_compute_keypair_v2.terraform.name
  security_groups = [openstack_networking_secgroup_v2.k8s_sec_group.name]
  user_data       = <<-EOF
              #cloud-config
              runcmd:
                - [ sudo, apt, update ]
              EOF

  block_device {
    uuid                  = var.node_image_uuid
    source_type           = "image"
    volume_size           = 10
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    uuid = openstack_networking_network_v2.k8s_network.id
  }

  depends_on = [openstack_networking_subnet_v2.k8s_subnet, openstack_compute_flavor_v2.k8s_flavor]
}

resource "openstack_compute_instance_v2" "workers" {
  count           = var.number_of_worker_nodes
  name            = "k8s-${var.cluster_name}-worker${count.index + 1}"
  flavor_name     = openstack_compute_flavor_v2.k8s_flavor.name
  key_pair        = openstack_compute_keypair_v2.terraform.name
  security_groups = [openstack_networking_secgroup_v2.k8s_sec_group.name]
  user_data       = <<-EOF
              #cloud-config
              runcmd:
                - [ sudo, apt, update ]
              EOF  

  block_device {
    uuid                  = var.node_image_uuid
    source_type           = "image"
    volume_size           = 20
    boot_index            = 0
    destination_type      = "volume"
    delete_on_termination = true
  }

  network {
    uuid = openstack_networking_network_v2.k8s_network.id
  }

  depends_on = [openstack_networking_subnet_v2.k8s_subnet, openstack_compute_flavor_v2.k8s_flavor]
}

