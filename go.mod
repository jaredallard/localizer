module github.com/getoutreach/localizer

go 1.15

require (
	github.com/Microsoft/go-winio v0.5.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/benbjohnson/clock v1.1.0
	github.com/davecgh/go-spew v1.1.1
	github.com/elazarl/goproxy v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/elazarl/goproxy/ext v0.0.0-20210110162100-a92cc753f88e // indirect
	github.com/function61/gokit v0.0.0-20210402130425-341c2c9ecfd0
	github.com/go-logr/logr v0.4.0
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.5
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/metal-stack/go-ipam v1.8.4
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20210503195802-e9a32991a82e
	golang.org/x/net v0.0.0-20210505024714-0287a6fb4125 // indirect
	golang.org/x/oauth2 v0.0.0-20210427180440-81ed05c6b58c // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/term v0.0.0-20210503060354-a79de5458b56 // indirect
	google.golang.org/genproto v0.0.0-20210505142820-a42aa055cf76 // indirect
	google.golang.org/grpc v1.37.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect

	// kubernetes deps
	k8s.io/api v0.21.0
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v0.21.0
	k8s.io/klog/v2 v2.8.0
)

replace k8s.io/client-go => github.com/jaredallard/client-go v0.21.0-jaredallard
