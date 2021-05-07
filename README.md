# authj

[![GoDoc](https://godoc.org/github.com/things-go/authj?status.svg)](https://godoc.org/github.com/things-go/authj)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/authj?tab=doc)
[![Build Status](https://travis-ci.org/things-go/authj.svg)](https://travis-ci.org/things-go/authj)
[![codecov](https://codecov.io/gh/things-go/authj/branch/master/graph/badge.svg)](https://codecov.io/gh/things-go/authj)
![Action Status](https://github.com/things-go/authj/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/authj)](https://goreportcard.com/report/github.com/things-go/authj)
[![Licence](https://img.shields.io/github/license/things-go/authj)](https://raw.githubusercontent.com/things-go/authj/master/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/authj)](https://github.com/things-go/authj/tags)

authj is an authorization middleware for [Gin](https://github.com/gin-gonic/gin), it's based on
 [casbin](https://github.com/casbin/casbin).

## Installation

```bash
go get github.com/things-go/authj
```

## Example

```Go
package main

import (
    "github.com/casbin/casbin/v2"
    "github.com/gin-gonic/gin"

    "github.com/things-go/authj"
)

func main() {
    // load the casbin model and policy from files, database is also supported.
    e ,err := casbin.NewEnforcer("authj_model.conf", "authj_policy.csv")
    if err!= nil{
        panic(err)    
    }   

    // define your router, and use the Casbin authj middleware.
    // the access that is denied by authj will return HTTP 403 error.
    router := gin.New()
    router.Use(func(c *gin.Context) {
        // got subject
        authj.ContextWithSubject(c, "admin")
    })
    router.Use(authj.NewAuthorizer(e))
}
```

## Documentation

The authorization determines a request based on ``{subject, object, action}``, which means what ``subject`` can perform what ``action`` on what ``object``. In this plugin, the meanings are:

1. ``subject``: the logged-on user name
2. ``object``: the URL path for the web resource like "dataset1/item1"
3. ``action``: HTTP method like GET, POST, PUT, DELETE, or the high-level actions you defined like "read-file", "write-blog"

For how to write authorization policy and other details, please refer to [the Casbin's documentation](https://github.com/casbin/casbin).

## Getting Help

- [Casbin](https://github.com/casbin/casbin)
- [Gin](https://github.com/gin-gonic/gin)
- [Gin-authz](https://github.com/gin-contrib/authz)

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
