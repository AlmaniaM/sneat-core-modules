module github.com/AlmaniaM/sneat-core-modules

go 1.21.4

//replace github.com/sneat-co/sneat-go-core => ../sneat-go-core
//replace github.com/bots-go-framework/bots-fw-store => ../../bots-go-framework/bots-fw-store

require (
	github.com/bots-go-framework/bots-fw-store v0.1.2
	github.com/crediterra/money v0.2.1
	github.com/dal-go/dalgo v0.12.0
	github.com/dal-go/mocks4dalgo v0.1.17
	github.com/golang/mock v1.6.0
	github.com/gosimple/slug v1.13.1
	github.com/sneat-co/sneat-go-core v0.18.0
	github.com/stretchr/testify v1.8.4
	github.com/strongo/random v0.0.1
	github.com/strongo/slice v0.1.4
	github.com/strongo/validation v0.0.6
)

require (
	github.com/alexsergivan/transliterator v1.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/strongo/decimal v0.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
