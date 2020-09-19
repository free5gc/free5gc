# Passos para testar o ping diretamente no Host

## 1. Baixando o repo
```shell
git clone https://github.com/LABORA-INF-UFG/non-3GPP ~/go/src/free5gc

cd ~/go/src/free5gc

# baixa os repos utilizados no release 3.0.1 (As libs e cada NF virou um repo separando, com um total de 55 repos diferentes)
git submodule update
```


## 2. Cadastrar UE via WebConsole
```shell
cd ~/go/src/free5gc/webconsole
go run server.go

# Acesse o WebConsole pelo navegador no endereço: http://localhost:5000
# Credenciais de acesso: admin free5gc
# Vá em Subscribers -> New Subscriber
# Informe os seguintes dados:
# plmnid: 20893
# Supi (imsi): 2089300007487
# Authentication Method: 5G_AKA
# k: 5122250214c33e723a5dd523fc145fc0
# Operator Code Type: OP
# Op: c9e8763286b5b9ffbdf56e1297d0887b
```


## 3. Instalando dependências, compilandos NFs e preparando o ambiente de execução


```shell
# instala as dependências utilizadas pelo free5gc
./install_env.sh

# compila os NFs, incluindo o UPF
./build.sh

# prepara o ambiente, com todas as veths, namespaces e módulos necessários
./setup_dev.sh up
```

## 4. Preparar o monitoramento (Opcional)

```shell
# Monitorar o tráfego de rede

# Exemplo de filtro para utilizar no wireshark para evitar pacotes lixos
# (!mdns and !ndps and !ssdp and !arp and !dns and  !icmpv6 and !(ip.src == 127.0.0.1 and ip.dst == 127.0.0.1)) and (!icmp or (ip.addr == 192.168.127.0/24 or ip.addr == 10.200.200.0/24 or ip.addr == 10.0.0.1/24))
sudo ./wireshark_ue.sh &>/dev/null &

# Exemplo de filtro para utilizar no wireshark para evitar pacotes lixos
# !dns and !arp and !icmpv6 and !ssdp and !mdns and !(ip.src == 127.0.0.1 and ip.dst == 127.0.0.1)
sudo ./wireshark_upf.sh &>/dev/null &

# Exemplo de filtro para utilizar no wireshark para evitar pacotes lixos
# !dns and !arp and !icmpv6 and !ssdp and !mdns
wireshark -kni any &>/dev/null &

# Monitorar os túneis GTP no UPF (São utilizados para implementar as regras de tratamento de pacotes informadas pelo SMF durante o PFCP Session Establishment)

# Rodar em um terminal diferente (Vai ficar a cada 1s atualizando a lista de túneis disponíveis no UPF)
sudo ./upf_gtp_info.sh pdr

# Rodar em um terminal difererente (Vai ficar a cada 1s atualizando a lista de túneis disponíveis no UPF)
sudo ./upf_gtp_info.sh far

# Rodar em um terminal difererente (Vai ficar a cada 1s atualizando a tabela de rotas no UPF)
sudo ip netns exec UPFns ip route

# Rodar em um terminal difererente (Vai ficar a cada 1s atualizando a tabela de rotas no UE)
# Para finalizar CTRL + C
sudo ip netns exec UEns ip route
```



## 5. Executando Core e UE

```shell
# Roda o Core, UPF e n3iwf (P/ finalizar CTRL + C)
# Para finalizar CTRL + C
sudo ./run.sh

# Roda o UE (Interessante fazer em um terminal novo, e deixar o atual mostrando os logs do Core, UPF e n3iwf)
# O script abaixo já contém os parâmetros necessários para executar o UE como OpCode, IMSI, PLMNID, ...
sudo ./run_ue.sh
```

## 6. Limpando o ambiente após a execução
```shell
# Remove as veths, namespaces e módulo gtp5g
sudo ./setup_dev.sh down

# Se quiser remover o database do free5gc
mongo free5gc --eval "db.dropDatabase()"
```
