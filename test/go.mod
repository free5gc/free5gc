module test

go 1.17

require (
	git.cs.nctu.edu.tw/calee/sctp v1.1.0
	github.com/antihax/optional v1.0.0
	github.com/calee0219/fatal v0.0.1
	github.com/davecgh/go-spew v1.1.1
	github.com/free5gc/amf v0.0.0
	github.com/free5gc/aper v1.0.4
	github.com/free5gc/ausf v0.0.0
	github.com/free5gc/n3iwf v0.0.0
	github.com/free5gc/nas v1.1.1
	github.com/free5gc/ngap v1.0.6
	github.com/free5gc/nrf v0.0.0
	github.com/free5gc/nssf v0.0.0
	github.com/free5gc/openapi v1.0.6
	github.com/free5gc/pcf v0.0.0
	github.com/free5gc/smf v0.0.0
	github.com/free5gc/udm v0.0.0
	github.com/free5gc/udr v0.0.0
	github.com/free5gc/util v1.0.5-0.20230511064842-2e120956883b
	github.com/gin-gonic/gin v1.9.0
	github.com/go-ping/ping v0.0.0-20210506233800-ff8be3320020
	github.com/google/uuid v1.3.0
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/stretchr/testify v1.8.1
	github.com/urfave/cli v1.22.5
	github.com/vishvananda/netlink v1.1.0
	go.mongodb.org/mongo-driver v1.8.4
	golang.org/x/net v0.7.0
	golang.org/x/sys v0.5.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/aead/cmac v0.0.0-20160719120800-7af84192f0b1 // indirect
	github.com/antonfisher/nested-logrus-formatter v1.3.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/bronze1man/radius v0.0.0-20190516032554-afd8baec892d // indirect
	github.com/bytedance/sonic v1.8.0 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/cydev/zero v0.0.0-20160322155811-4a4535dd56e7 // indirect
	github.com/evanphx/json-patch v0.5.2 // indirect
	github.com/free5gc/pfcp v1.0.6 // indirect
	github.com/free5gc/tlv v1.0.2-0.20230131124215-8b6ebd69bf93 // indirect
	github.com/gin-contrib/cors v1.3.1 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.11.2 // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/goccy/go-json v0.10.0 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/h2non/parth v0.0.0-20190131123155-b4df798d6542 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/tim-ywliu/nested-logrus-formatter v1.3.2 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.9 // indirect
	github.com/vishvananda/netns v0.0.0-20211101163701-50045581ed74 // indirect
	github.com/wmnsk/go-gtp v0.8.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8 // indirect
	golang.org/x/sync v0.0.0-20220722155255-886fb9371eb4 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/h2non/gock.v1 v1.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
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
