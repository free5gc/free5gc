module test

go 1.14

require (
	git.cs.nctu.edu.tw/calee/sctp v1.1.0
	github.com/antihax/optional v1.0.0
	github.com/calee0219/fatal v0.0.1
	github.com/davecgh/go-spew v1.1.1
	github.com/free5gc/amf v0.0.0-00010101000000-000000000000
	github.com/free5gc/aper v1.0.4
	github.com/free5gc/ausf v0.0.0-00010101000000-000000000000
	github.com/free5gc/n3iwf v0.0.0-00010101000000-000000000000
	github.com/free5gc/nas v1.0.7
	github.com/free5gc/ngap v1.0.6
	github.com/free5gc/nrf v0.0.0-00010101000000-000000000000
	github.com/free5gc/nssf v0.0.0-00010101000000-000000000000
	github.com/free5gc/openapi v1.0.5
	github.com/free5gc/pcf v0.0.0-00010101000000-000000000000
	github.com/free5gc/smf v0.0.0-00010101000000-000000000000
	github.com/free5gc/udm v0.0.0-00010101000000-000000000000
	github.com/free5gc/udr v0.0.0-00010101000000-000000000000
	github.com/free5gc/util v1.0.3
	github.com/go-ping/ping v0.0.0-20210506233800-ff8be3320020
	github.com/google/uuid v1.3.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli v1.22.5
	github.com/vishvananda/netlink v1.1.0
	go.mongodb.org/mongo-driver v1.8.4
	golang.org/x/net v0.0.0-20211020060615-d418f374d309
	golang.org/x/sys v0.0.0-20211007075335-d3039528d8ac
	gopkg.in/yaml.v2 v2.4.0
)

replace (
	github.com/free5gc/amf => ../NFs/amf
	github.com/free5gc/ausf => ../NFs/ausf
	github.com/free5gc/n3iwf => ../NFs/n3iwf
	github.com/free5gc/nrf => ../NFs/nrf
	github.com/free5gc/nssf => ../NFs/nssf
	github.com/free5gc/pcf => ../NFs/pcf
	github.com/free5gc/smf => ../NFs/smf
	github.com/free5gc/udm => ../NFs/udm
	github.com/free5gc/udr => ../NFs/udr
)
