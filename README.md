# go-objectid

[![Go Report Card](https://goreportcard.com/badge/JesseCoretta/go-objectid)](https://goreportcard.com/report/github.com/JesseCoretta/go-objectid) [![Go Reference](https://pkg.go.dev/badge/github.com/JesseCoretta/go-objectid.svg)](https://pkg.go.dev/github.com/JesseCoretta/go-objectid) [![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-objectid/blob/main/LICENSE) [![codecov](https://codecov.io/gh/JesseCoretta/go-objectid/graph/badge.svg?token=RLW4DHLKQP)](https://codecov.io/gh/JesseCoretta/go-objectid) [![contributions welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/JesseCoretta/go-objectid/issues) [![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/JesseCoretta/go-objectid/go.yml)](https://github.com/JesseCoretta/go-objectid/actions/workflows/go.yml) [![Author](https://img.shields.io/badge/author-Jesse_Coretta-darkred?label=%F0%9F%94%BA&labelColor=indigo&color=maroon)](https://www.linkedin.com/in/jessecoretta/) [![Libraries.io dependency status for GitHub repo](https://img.shields.io/librariesio/github/JesseCoretta/go-objectid)](https://github/JesseCoretta/go-objectid) [![Help Animals](https://img.shields.io/badge/donations-yellow?label=%F0%9F%98%BA&labelColor=Yellow)](https://github.com/JesseCoretta/JesseCoretta/blob/main/DONATIONS.md)

<!-- [![GitHub release (with filter)](https://img.shields.io/github/v/release/JesseCoretta/go-objectid)](https://github.com/JesseCoretta/go-objectid/releases) -->

<!-- [![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/JesseCoretta/go-objectid/go.yml)](https://github.com/JesseCoretta/go-objectid/actions/workflows/go.yml) -->

Package objectid offers a convenient ASN.1 Object Identifier type and associated methods.

ASN.1 Object Identifiers encompass information that goes beyond their dotted representation. This tiny package merely facilitates the handling of ASN.1 NameAndNumberForm values and alternate names that may be associated with a given OID in the wild.

## Uint128 Support

Unsigned 128-bit integer support for individual NumberForm values is made possible due to the private incorporation of Luke Champine's awesome Uint128 type, which manifests here through instances of the package-provided NumberForm type.
