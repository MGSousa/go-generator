# Go Generator

[![Test](https://github.com/MGSousa/go-generator/actions/workflows/tests.yml/badge.svg)](https://github.com/MGSousa/go-generator/actions/workflows/tests.yml)

Static FS Server Generator using [Iris Web Framework](https://github.com/kataras/iris)

## Generate bindata
Converts public assets into byte-codes (for use in Production)
```sh
go-bindata -pkg generator -prefix "assets" -fs ./public/...
```
