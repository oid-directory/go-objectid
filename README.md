# go-objectid

[![Go Report Card](https://goreportcard.com/badge/JesseCoretta/go-objectid)](https://goreportcard.com/report/github.com/JesseCoretta/go-objectid) [![GoDoc](https://godoc.org/github.com/JesseCoretta/go-objectid?status.svg)](https://godoc.org/github.com/JesseCoretta/go-objectid) ![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)

Package objectid offers a convenient ASN.1 Object Identifier type and associated methods.

ASN.1 Object Identifiers encompass information that goes beyond their dotted representation. This tiny package merely facilitates the handling of ASN.1 NameAndNumberForm values and alternate names that may be associated with a given OID in the wild.

## Uint128 Support

Unsigned 128-bit integer support for individual NumberForm values is made possible due to the private incorporation of Luke Champine's awesome Uint128 type, which manifests here through instances of the package-provided NumberForm type.
