# Filler

[![go](https://img.shields.io/badge/go-v1.13.X-5C4EE5.svg)](https://golang.org/)
[![make](https://img.shields.io/badge/make-v3.8.X-yellow.svg)](https://linux.die.net/man/1/make)

> `Filler` is a golang package pretty power you, when you want populate a structure with dynamic tags and implementation (we call these providers)
>

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quickstart](#quickstart)

## Prerequisites

You will need the following things properly installed on your computer:

* [Git](http://git-scm.com/)
* [Make](https://linux.die.net/man/1/make)
* [Go](https://golang.org/)

## Quickstart

### Provider

Filler is a library that you can to use for populate your structure from different data source, we have two implementation right now, but you can develop more

- The first one `envconfig` it has the power to get values from environment variables e.g :
```go
    provider.NewEnvProvider()
```
- The second one `ssm` it has the power to get values from aws parameter store
```go
    ssmProvider, err := provider.NewSsmProvider()
```

Now we need associated a provider to a tag, we can use suture `config.ProviderByTag` e.g :

````go
    config.ProviderByTag{
		Provider: provider.NewEnvProvider(),
		Tag:      "envconfig",
	}
````

### Read config

By read config and populate a struts just it`s necessary the next step,

1) Create struts with tags

````go
type (
	DB struct {
		URL          string        `ssm:"/base-path/second-path/value"`
		Username     string        `envconfig:"DB_USER,optional" ssm:"/base-path/second-path/valueTwo"`
		Password     string        `yml:"password" envconfig:"DB_PASS,optional"`
	}
	Server struct {
		Port    string        `envconfig:"SERVER_PORT,optional"`
	}
	AppConfig struct {
		DB
		Server
		IsDebug     bool   `envconfig:"DEBUG,optional"`
		Environment string `envconfig:"ENVIRONMENT,optional"`
	}
)
````
2) Initialization config
````go
    appConfig := new(AppConfig)
	err = config.Marshall(appConfig)
````
### Custom Provider

You can add a custom provider just need create a structured that implement an interface

```go
type Provider interface {
	Execute(*Context) (string, error)
}
```

