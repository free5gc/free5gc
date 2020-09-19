Criando um cluster kubernetes
========================

# Install docker

```
curl cat -fsSL https://get.docker.com | bash
```
```
cat > /etc/docker/daemon.json <<EOF
{
"exec-opts": ["native.cgroupdriver=systemd"],
"log-driver": "json-file",
"log-opts": {
"max-size": "100m"
},
"storage-driver": "overlay2"
}
EOF
```
```
mkdir -p /etc/systemd/system/docker.service.d
```
```
systemctl daemon-reload
```
```
systemctl restart docker
```

- A saída desse comando deve ser **Cgroup Driver: systemd**.

```
docker info | grep -i cgroup
```

# Instalação Kubernets

- Baixe o pacote e use a versão 1.15.3 que é 

```
apt-get update && apt-get install -y apt-transport-https gnupg2
```
```
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
```
```
echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" > /etc/apt/sources.list.d/kubernetes.list
```
```
apt-get update
```
```
apt-get install -y kubelet=1.15.3-00 kubeadm=1.15.3-00 kubectl=1.15.3-00 --allow-downgrades
```
```
swapoff -a
```

# Iniciando o Cluster Kubernetes

- Os comandos a seguir devem ser rodados apenas no Master

`kubeadm config images pull`

- Ao executar o comando a seguir será retornado um token iniciado com **kubeadm join**. Copie esse token completo (inclusive com a *kubeadm join*) e cole na(s) máquina(s) Worker(s).

`kubeadm init`

`mkdir -p $HOME/.kube`

`sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config`

`sudo chown $(id -u):$(id -g) $HOME/.kube/config`

- Instale o Pod Network para criar uma rede entre os Pods.

`kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"`

`kubectl get pods -n kube-system`
