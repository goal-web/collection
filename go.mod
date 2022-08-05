module github.com/goal-web/collection

go 1.18

require (
	github.com/goal-web/contracts v0.1.13
	github.com/goal-web/supports v0.1.10
	github.com/shopspring/decimal v1.3.1
	github.com/spf13/cast v1.5.0
	github.com/stretchr/testify v1.7.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200605160147-a5ece683394c // indirect
)

replace github.com/goal-web/contracts => ../contracts
