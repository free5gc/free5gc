module github.com/free5gc/config/cfg-generator

go 1.26.2

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/free5gc/openapi v1.3.1-0.20261022060819-1fa510b3e4a2
	github.com/google/go-cmp v0.7.0
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.12.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.15.0 // indirect
	github.com/bytedance/sonic/loader v0.5.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fclairamb/go-log v0.4.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/free5gc/nas v1.3.0 // indirect
	github.com/free5gc/ngap v1.2.0 // indirect
	github.com/free5gc/sctp v1.2.0 // indirect
	github.com/free5gc/util v1.4.1-0.20261014070359-d865b0ce854b // indirect
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/gin-gonic/gin v1.12.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.1 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/h2non/gock v1.2.0 // indirect
	github.com/h2non/parth v0.0.0-20190131123155-b4df798d6542 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.59.1 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	github.com/tim-ywliu/nested-logrus-formatter v1.3.2 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	github.com/vishvananda/netlink v1.1.0 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	github.com/wmnsk/go-gtp v0.8.11-0.20240705144331-f53bfdd4233b // indirect
	go.mongodb.org/mongo-driver/v2 v2.5.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	go.opentelemetry.io/otel v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.24.0 // indirect
	go.opentelemetry.io/otel/trace v1.24.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/arch v0.23.0 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)

require (
	github.com/free5gc/amf v0.0.0
	github.com/free5gc/ausf v0.0.0
	github.com/free5gc/bsf v0.0.0
	github.com/free5gc/chf v0.0.0
	github.com/free5gc/go-upf v0.0.0
	github.com/free5gc/n3iwf v0.0.0
	github.com/free5gc/nef v0.0.0
	github.com/free5gc/nrf v0.0.0
	github.com/free5gc/nssf v0.0.0
	github.com/free5gc/pcf v0.0.0
	github.com/free5gc/scp v0.0.0
	github.com/free5gc/smf v0.0.0
	github.com/free5gc/tngf v0.0.0
	github.com/free5gc/udm v0.0.0
	github.com/free5gc/udr v0.0.0
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
)

replace github.com/free5gc/amf => ../../NFs/amf

replace github.com/free5gc/smf => ../../NFs/smf

replace github.com/free5gc/ausf => ../../NFs/ausf

replace github.com/free5gc/bsf => ../../NFs/bsf

replace github.com/free5gc/chf => ../../NFs/chf

replace github.com/free5gc/n3iwf => ../../NFs/n3iwf

replace github.com/free5gc/nef => ../../NFs/nef

replace github.com/free5gc/nrf => ../../NFs/nrf

replace github.com/free5gc/nssf => ../../NFs/nssf

replace github.com/free5gc/pcf => ../../NFs/pcf

replace github.com/free5gc/scp => ../../NFs/scp

replace github.com/free5gc/tngf => ../../NFs/tngf

replace github.com/free5gc/udm => ../../NFs/udm

replace github.com/free5gc/udr => ../../NFs/udr

replace github.com/free5gc/go-upf => ../../NFs/upf
