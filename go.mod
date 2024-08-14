module my5G-RANTester

go 1.21

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/free5gc/aper v1.0.5
	github.com/free5gc/go-gtp5gnl v1.4.5
	github.com/free5gc/nas v1.1.2-0.20230828074825-175b09665828
	github.com/free5gc/ngap v1.0.8
	github.com/free5gc/openapi v1.0.8
	github.com/free5gc/util v1.0.4
	github.com/google/gopacket v1.1.19
	github.com/ishidawataru/sctp v0.0.0-20230406120618-7ff4192f6ff2
	github.com/khirono/go-nl v1.0.4
	github.com/khirono/go-rtnllink v1.1.1
	github.com/mitchellh/mapstructure v1.4.2
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826
	github.com/sirupsen/logrus v1.9.3
	github.com/stretchr/testify v1.8.3
	github.com/tetratelabs/wazero v1.3.0
	github.com/urfave/cli/v2 v2.25.5
	github.com/vishvananda/netlink v1.1.0
	github.com/wmnsk/go-gtp v0.8.6
	golang.org/x/net v0.17.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/aead/cmac v0.0.0-20160719120800-7af84192f0b1 // indirect
	github.com/antonfisher/nested-logrus-formatter v1.3.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/golang-jwt/jwt v3.2.1+incompatible // indirect
	github.com/khirono/go-genl v1.0.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/tim-ywliu/nested-logrus-formatter v1.3.2 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/sys v0.20.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/free5gc/ngap => ../ngap
