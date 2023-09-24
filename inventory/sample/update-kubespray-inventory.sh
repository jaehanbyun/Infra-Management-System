#!/bin/bash

# Parsing Terraform outputs
BASTION_IP=$(terraform output bastion_fips | tr -d '[]," \n')
MASTER_IPS=$(terraform output master_ips | tr -d '[]' | sed 's/\"//g' | tr -d '\n')
WORKER_IPS=$(terraform output worker_ips | tr -d '[]' | sed 's/\"//g' | tr -d '\n')

# Create inventory.ini file
echo "[all]" > inventory.ini

# Master nodes
count=1
IFS=',' read -ra ADDR <<< "$MASTER_IPS"
for ip in "${ADDR[@]}"; do
  ip_clean=$(echo $ip | tr -d ' ')
  echo "k8s-master${count} ansible_host=${ip_clean} ip=${ip_clean} etcd_member_name=etcd${count}" >> inventory.ini
  ((count++))
done

# Worker nodes
count=1
IFS=',' read -ra ADDR <<< "$WORKER_IPS"
for ip in "${ADDR[@]}"; do
  ip_clean=$(echo $ip | tr -d ' ')
  echo "k8s-worker${count} ansible_host=${ip_clean} ip=${ip_clean}" >> inventory.ini
  ((count++))
done

echo -e "\n[bastion]" >> inventory.ini
echo "k8s-bastion ansible_host=${BASTION_IP} ansible_user=ubuntu" >> inventory.ini

echo -e "\n[kube_control_plane]" >> inventory.ini
count=1
IFS=',' read -ra ADDR <<< "$MASTER_IPS"
for ip in "${ADDR[@]}"; do
  echo "k8s-master${count}" >> inventory.ini
  ((count++))
done

echo -e "\n[etcd]" >> inventory.ini
count=1
IFS=',' read -ra ADDR <<< "$MASTER_IPS"
for ip in "${ADDR[@]}"; do
  echo "k8s-master${count}" >> inventory.ini
  ((count++))
done

echo -e "\n[kube_node]" >> inventory.ini
count=1
IFS=',' read -ra ADDR <<< "$MASTER_IPS"
for ip in "${ADDR[@]}"; do
  echo "k8s-master${count}" >> inventory.ini
  ((count++))
done
count=1
IFS=',' read -ra ADDR <<< "$WORKER_IPS"
for ip in "${ADDR[@]}"; do
  echo "k8s-worker${count}" >> inventory.ini
  ((count++))
done

echo -e "\n[k8s_cluster:children]" >> inventory.ini
echo "kube_control_plane" >> inventory.ini
echo "kube_node" >> inventory.ini
