module techwave

go 1.25.0

require github.com/gorilla/mux v1.8.0

require (
	github.com/cucumber/godog v0.15.1
	github.com/redis/go-redis/v9 v9.17.2
	github.com/stretchr/testify v1.11.1
	go.uber.org/mock v0.6.0
	golang.org/x/crypto v0.49.0 // VULNERABLE: CVE-2023-48795 (Terrapin Attack)
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cucumber/gherkin/go/v26 v26.2.0 // indirect
	github.com/cucumber/messages/go/v21 v21.0.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	golang.org/x/sys v0.42.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
