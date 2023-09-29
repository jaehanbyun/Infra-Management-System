#!/bin/bash

# Ensure we get three arguments
if [[ $# -ne 4 ]]; then
  echo "Usage: $0 <CLUSTER_NAME> <BASTION_IP> <MASTER_IPS> <WORKER_IPS>"
  exit 1
fi

CLUSTER_NAME=$1
BASTION_IP=$2
MASTER_IPS=$3
WORKER_IPS=$4

# Create inventory.ini file
echo "[all]" > inventory.ini

# Master nodes
count=1
IFS=',' read -ra ADDR <<< "$MASTER_IPS"
for ip in "${ADDR[@]}"; do
  ip_clean=$(echo $ip | tr -d ' ')
  echo "k8s-${CLUSTER_NAME}-master${count} ansible_host=${ip_clean} ip=${ip_clean} etcd_member_name=etcd${count}" >> inventory.ini
  ((count++))
done

# Worker nodes
count=1
IFS=',' read -ra ADDR <<< "$WORKER_IPS"
for ip in "${ADDR[@]}"; do
  ip_clean=$(echo $ip | tr -d ' ')
  echo "k8s-${CLUSTER_NAME}-worker${count} ansible_host=${ip_clean} ip=${ip_clean}" >> inventory.ini
  ((count++))
done

echo -e "\n[bastion]" >> inventory.ini
echo "k8s-${CLUSTER_NAME}-bastion ansible_host=${BASTION_IP} ansible_user=ubuntu" >> inventory.ini

echo -e "\n[kube_control_plane]" >> inventory.ini
count=1
IFS=',' read -ra ADDR <<< "$MASTER_IPS"
for ip in "${ADDR[@]}"; do
  echo "k8s-${CLUSTER_NAME}-master${count}" >> inventory.ini
  ((count++))
done

echo -e "\n[etcd]" >> inventory.ini
count=1
IFS=',' read -ra ADDR <<< "$MASTER_IPS"
for ip in "${ADDR[@]}"; do
  echo "k8s-${CLUSTER_NAME}-master${count}" >> inventory.ini
  ((count++))
done

echo -e "\n[kube_node]" >> inventory.ini
count=1
IFS=',' read -ra ADDR <<< "$MASTER_IPS"
for ip in "${ADDR[@]}"; do
  echo "k8s-${CLUSTER_NAME}-master${count}" >> inventory.ini
  ((count++))
done
count=1
IFS=',' read -ra ADDR <<< "$WORKER_IPS"
for ip in "${ADDR[@]}"; do
  echo "k8s-${CLUSTER_NAME}-worker${count}" >> inventory.ini
  ((count++))
done

echo -e "\n[k8s_cluster:children]" >> inventory.ini
echo "kube_control_plane" >> inventory.ini
echo "kube_node" >> inventory.ini
