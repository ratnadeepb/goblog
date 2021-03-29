#!/bin/bash
function usage {
    echo "$0 [master] [master-ip] or [slave]"
    exit 1
}

node_type=$1
master_ip=$2

if [ "$node_type" = "master" -a "$master_ip" = "" ]
then
    usage
    
elif [ "$node_type" != "master" -a "$node_type" != "slave" ]
then
    usage
fi

sudo apt-get update
sudo apt-get install -y \
default-jre \
default-jdk \
apt-transport-https \
ca-certificates \
curl \
gnupg-agent \
software-properties-common

sudo swapoff -a
cat /etc/fstab | grep -v '^#' | grep -v 'swap' | sudo tee /etc/fstab

# sudo systemctl stop docker
# sudo mv /var/lib/docker /my_mount
# sudo ln -s /my_mount/docker /var/lib/docker
# sudo systemctl restart docker

curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add - && \
echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list && \
sudo apt-get update -q && \
apt-cache madison kubelet

sudo apt-get install -qy kubelet=1.18.5-00 kubectl=1.18.5-00 kubeadm=1.18.5-00

sudo sed -i '/\[Service\]/a\Environment="KUBELET_UNSAFE_SYSCTLS=--allowed-unsafe-sysctls='kernel.shm*,kernel.sem,kernel.msg*,net.core.*'"' /etc/systemd/system/kubelet.service.d/10-kubeadm.con

sudo systemctl daemon-reload
sudo systemctl restart kubelet

if [ "$node_type" = "master" ]
then
    #deploy kubernetes cluster
    sudo kubeadm init --apiserver-advertise-address=$master_ip --pod-network-cidr=10.244.0.0/16
    #for non-root user, make sure that kubernetes config directory has the same permissions as kubernetes config file.
    mkdir -p $HOME/.kube
    sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/config
    sudo chown $(id -u):$(id -g) $HOME/.kube/
    
    #deploy flannel
    sudo sysctl net.bridge.bridge-nf-call-iptables=1
    # https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/
    #after this step, coredns status will be changed to running from pending
    kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
    kubectl get nodes
    kubectl get pods --namespace=kube-system
fi
