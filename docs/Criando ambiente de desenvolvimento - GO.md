Criando ambiente de desenvolvimento - GO
========================

## Instalação - Linux

### Baixe o pacote da linguagem no site 

- `golang.org/dl/`

### Execute os comandos a seguir no local do download

- `tar -C /usr/local -xzf go1.14.1.linux-amd64.tar.gz`
- `export PATH=$PATH:/usr/local/go/bin`

- Adicione variaveis de ambiente no arquivo .profile

- `sudo gedit ~/.profile`
- `export GOPATH=$HOME/go`
- `export PATH=$PATH:$GOPATH/bin:/usr/local/go/bin`

- Salve e Carregue as variáveis de ambiente.
- `. /.profile`

- Você pode ver suas variáveis utilizando
- `echo $PATH`

### Agora teste sua instalação.

- Para testar a instalação basta criar um arquivo com o seguinte código 

```
 package main

import "fmt"

func main() {
	fmt.Printf("hello, world\n")
} 
```
- Agora builde.

`go build hello.go`

- Execute.

`./hello`

- Com estas configurações você poderá executar arquivos .go em qualquer lugar dentro da sua máquina.

### Agora instale o Visual Studio Code para desenvolvimento. (à gosto)

- Baixe o .deb no site e instale.

`https://code.visualstudio.com/`

- Instale a extensão para Go
    - Para isso clique em extensões e procure por Go Lang.
    - Instale a primeira que provavelmente será a distribuida pela microsoft.