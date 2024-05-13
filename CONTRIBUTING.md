# Welcome to the go-objectid contributing guide <!-- omit in toc -->

First, welcome to the go-objectid repository. This repository offers a package allowing true, unbounded ASN.1 OBJECT IDENTIFIER support in Golang with many convenient features and encoding/decoding support.  The goal of this package is to offer a means of avoiding the shortcomings of the "encoding/asn1" package in terms of OID support.

## Contributor guide

A few things should be reviewed before submitting a contribution to this repository:

 1. Read our [Code of Conduct](./CODE_OF_CONDUCT.md) to keep our community approachable and respectable.
 2. Review the main [![GoDoc](https://pkg.go.dev/github.com/JesseCoretta/go-objectid?status.svg)](https://pkg.go.dev/github.com/JesseCoretta/go-objectid) page, which provides the entire suite of useful documentation rendered in Go's typically slick manner 😎.
 3. Review the [Collaborating with pull requests](https://docs.github.com/en/github/collaborating-with-pull-requests) document, unless you're already familiar with its concepts ...
 4. Keep cyclomatic factors <= 9
 5. Keep test coverage at 100%
 6. All exported constructs -- methods, functions, types, global vars -- must be commented in proper English grammar.

So long as these rules are honored, contributions are welcome.
