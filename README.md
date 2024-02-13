# common-go-dependency-tree

![Code Coverage](./badges/coverage.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/cjlapao/common-go-dependency-tree)](https://goreportcard.com/report/github.com/cjlapao/common-go-dependency-tree)
[![Lint Codebase](https://github.com/cjlapao/common-go-dependency-tree/actions/workflows/linter.yml/badge.svg)](https://github.com/cjlapao/common-go-dependency-tree/actions/workflows/linter.yml)

This is a simple tool to generate a dependency tree for a given struct in a Go project. you can use it to run modules in sequence depending on their dependencies without a lot of boilerplate code.

It is a generic service and you can have multiple trees for different structs in your project.

