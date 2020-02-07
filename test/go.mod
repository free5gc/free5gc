module test

go 1.14

require (
	github.com/free5gc/CommonConsumerTestData v1.0.0
	github.com/free5gc/MongoDBLibrary v1.0.0
	github.com/free5gc/UeauCommon v1.0.0
	github.com/free5gc/amf v1.3.0
	github.com/free5gc/aper v1.0.0
	github.com/free5gc/ausf v1.3.0
	github.com/free5gc/http2_util v1.0.0
	github.com/free5gc/logger_util v1.0.0
	github.com/free5gc/milenage v1.0.0
	github.com/free5gc/n3iwf v1.3.0
	github.com/free5gc/nas v1.0.0
	github.com/free5gc/ngap v1.0.0
	github.com/free5gc/nrf v1.3.0
	github.com/free5gc/nssf v1.3.0
	github.com/free5gc/openapi v1.0.0
	github.com/free5gc/path_util v1.0.0
	github.com/free5gc/pcf v1.3.0
	github.com/free5gc/smf v1.3.0
	github.com/free5gc/udm v1.3.0
	github.com/free5gc/udr v1.3.0
	git.cs.nctu.edu.tw/calee/sctp v1.1.0
	github.com/Djarvur/go-err113 v0.1.0 // indirect
	github.com/aws/aws-sdk-go v1.36.24 // indirect
	github.com/calee0219/fatal v0.0.1
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ping/ping v0.0.0-20200918120429-e8ae07c3cec8
	github.com/golangci/golangci-lint v1.35.2
	github.com/golangci/misspell v0.3.5 // indirect
	github.com/golangci/revgrep v0.0.0-20180812185044-276a5c0a1039 // indirect
	github.com/gostaticanalysis/analysisutil v0.6.1 // indirect
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/klauspost/compress v1.11.6 // indirect
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/matoous/godox v0.0.0-20200801072554-4fb83dc2941e // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/nbutton23/zxcvbn-go v0.0.0-20201221231540-e56b841a3c88 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/quasilyte/go-ruleguard v0.2.1 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20200805063351-8f842688393c // indirect
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.6.1
	github.com/tdakkota/asciicheck v0.0.0-20200416200610-e657995f937b // indirect
	github.com/timakin/bodyclose v0.0.0-20200424151742-cb6215831a94 // indirect
	github.com/tomarrell/wrapcheck v0.0.0-20201130113247-1683564d9756 // indirect
	github.com/ugorji/go v1.2.3 // indirect
	github.com/urfave/cli v1.22.5
	github.com/vishvananda/netlink v1.1.0
	go.mongodb.org/mongo-driver v1.4.4
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20201224014010-6772e930b67b
	golang.org/x/sys v0.0.0-20210112091331-59c308dcf3cc
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.0.0-20210111221946-d33bae441459 // indirect
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	honnef.co/go/tools v0.1.0 // indirect
	mvdan.cc/unparam v0.0.0-20210104141923-aac4ce9116a7 // indirect
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
