# Considerações sobre Free5GC - Stage 3

## Arquitura de acesso não confiável non-3GPP utilizando o Core Network free5gc stage 3.

Para realizar os procedimento de acesso não confiável via non-3GPP aplicando o padrão **Non 3GPP Interworking Function (N3IWF)**, (1) realizamos a configuração de código Core Network, (2) colocando em produção em estrutura de virtualização baseada em Containers de aplicação. Para esta tarefa foi utilizado a Docker/Docker-compose para composições das estruturas de cada funcionalidade (Function Network), e também as disposições corretas das interfaces de redes com suas respectivas configurações de endereçamento.

![](media/ArquiteturaBasica2.png)

## Componentes do Core que estão em funcionamento:
- [X] AMF  
- [X] SMF -> Verificar Algumas configuração de endereçamento IP 
- [X] NRF 
- [X] PCF
- [X] NSSF
- [X] AUSF
- [X] UDR
- [X] UDM
- [ ] UPF -> Problemas Conainers com tunneis GTP
- [X] N3IWF
- [ ] UE -> Desenvolver funcionalidades do User Equipment (Ike-DAEMON, AN Parametros, UDP Server, IPSec, GRE)

## Considerações de endereçamento nas interfaces de rede NWu,N2,N3:
- A interface NWu, refere-se a interligação lógica entre o UE -> N3IWF
    - Os protocolos relacionados neste contexto são IKEv2, EAP, IPSec, GRE, NAS, PDU.
    - Deve possuir uma interface física, e uma interface virtual.
    - A interface de rede física é destinada para comunicação de acesso a rede non-3GPP com N3IWF, e também para 
estabilização de associação segura (IKE-Daemon) via socket UDP.
    - A interface virtual, deve ser configurada para prosseguimento Tunnel IPSec.
    - Deve existir um Tunnel GRE interligando o UE -> N3IWF, para transporte dos dados relacionados ao plano de dados.

- A interface N2, refere-se a interligação lógica entre o  N3IWF -> AMF 
    - Os protocolos relacionados neste contexto são NAS/NAS-PDU, NGAP, SCTP.
    - Deve possuir uma interface física somente.
    - Contexto para transmissão das mensagens N2 no Plano de Controle .

- A interface N3, refere-se a interligação lógica entre o N3IWF -> UPF 
    - Os protocolos relacioandos neste contexto são GTP, PDU-Sessions
    - Deve possuir uma interface física somente e trambém formação de Tunnel GTP

## Serviços para serem implementados no N3IWF:
- [X] UDP Server Connection Port (500,14500)  
- [X] Ikev2 Daemon 
- [X] Logging 
- [ ] Command Line Test
- [X] GRE Tunnel 
- [X] GTP Tunnel
- [X] IPSec Context
- [X] N2 -SCTP Transport Protocol
- [X] NAS/NAS-PDU/NAS-SMC/NGAP
- [X] Relay UE->N2
- [X] Relay UE->N3

## Serviços para serem implementados no UE:
- [ ] UDP Server Connection Port (500,14500)  
- [ ] IKEv2 Daemon
- [ ] Subscriber Informantion and AN (GUAMI, NSSAI, PLMNid, SUPI/SUCI)
- [ ] EAP (NAS/NAS-PDU)
- [ ] GRE Tunnel 
- [ ] IPSec Context
- [ ] Logging 
- [ ] Command Line Test

## Procedimentos para auxiliar na codificação (Bibliotecas, Documentos Técnicos, Alguns tutoriais):
- Lista de Biblitecas úteis (Golang)
    - Geral
        - "https://github.com/sirupsen/logrus" -> Structured, pluggable logging for Go.
        - "https://github.com/urfave/cli" -> A simple, fast, and fun package for building command line apps in Go
    - (UE -> N3IWF)
        - "https://github.com/vishvananda/netlink" ->  Simple netlink library for go. (iproute, gre, ipsec) 
    
## Criação de Tunnels GRE:
```go
package main

import(
      "log"
      "github.com/vishvananda/netlink"
)

func MakeGRETunnel( name string, localIP string, remoteIP string){
      
         GRETunnel := &netlink.Gretun{
               LinkAttrs := &netlink.Gretun{
                     Name: name
               },
               Local: localIP,
               Remote: remoteIP
         }
         
         err := netlink.LinkAdd(GRETUnnel)
         if err != nil{
            log.Fatal("Error")
         }
}
```
