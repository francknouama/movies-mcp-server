module github.com/francknouama/movies-mcp-server/godog-server

go 1.24.4

require (
	github.com/cucumber/gherkin/go/v26 v26.2.0
	github.com/cucumber/godog v0.15.0
	github.com/cucumber/messages/go/v21 v21.0.1
	github.com/francknouama/movies-mcp-server/shared-mcp v0.0.0-00010101000000-000000000000
)

require (
	github.com/gofrs/uuid v4.3.1+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-memdb v1.3.4 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
)

replace github.com/francknouama/movies-mcp-server/shared-mcp => ../shared-mcp
