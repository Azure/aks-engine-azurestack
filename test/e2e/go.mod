module github.com/Azure/aks-engine-azurestack/test/e2e

go 1.20

require (
	github.com/Azure/aks-engine-azurestack v0.43.0
	github.com/Azure/azure-sdk-for-go/profile/p20200901 v0.1.1
	github.com/Azure/go-autorest/autorest v0.11.12
	github.com/influxdata/influxdb v1.7.9
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/onsi/ginkgo/v2 v2.17.1
	github.com/onsi/gomega v1.30.0
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.21.0
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.11.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.5.2 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.10 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.0 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/form3tech-oss/jwt-go v3.2.2+incompatible // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-playground/locales v0.12.1 // indirect
	github.com/go-playground/universal-translator v0.16.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/leonelquinteros/gotext v1.4.0 // indirect
	github.com/magefile/mage v1.10.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/sirupsen/logrus v1.8.0 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	golang.org/x/net v0.23.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.18.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	golang.org/x/tools v0.17.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/go-playground/validator.v9 v9.25.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/Azure/aks-engine-azurestack v0.43.0 => ../..
