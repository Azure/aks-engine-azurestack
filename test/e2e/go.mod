module github.com/Azure/aks-engine-azurestack/test/e2e

go 1.23

require (
	github.com/Azure/aks-engine-azurestack v0.43.0
	github.com/Azure/azure-sdk-for-go/profile/p20200901 v0.1.1
	github.com/influxdata/influxdb v1.7.9
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/onsi/ginkgo/v2 v2.17.1
	github.com/onsi/gomega v1.30.0
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.31.0
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.11.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.8.0 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/leonelquinteros/gotext v1.4.0 // indirect
	github.com/magefile/mage v1.10.0 // indirect
	github.com/sirupsen/logrus v1.8.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	gopkg.in/go-playground/validator.v9 v9.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Azure/aks-engine-azurestack v0.43.0 => ../..
