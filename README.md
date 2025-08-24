# CPF/CNPJ Validation for Go

[![Go Version](https://img.shields.io/badge/go-1.24.2+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/n0vdd/cpf_cnpj)](https://goreportcard.com/report/github.com/n0vdd/cpf_cnpj)

A modern Go package for validating and formatting Brazilian taxpayer identification numbers (CPF and CNPJ). Supports both traditional numeric CNPJs and the new **CNPJ Alfanumérico** format.

**For international users**: CPF is for individuals, CNPJ is for companies - both are mandatory tax identification numbers in Brazil.

## Installation

```bash
go get github.com/n0vdd/cpf_cnpj
```

## Quick Start

```go
import "github.com/n0vdd/cpf_cnpj"

// CPF validation - accepts both formatted and clean input
cpf, err := cpfcnpj.NewCpf("716.566.867-59")  // formatted
cpf, err := cpfcnpj.NewCpf("71656686759")     // clean
if err != nil {
    return err
}
fmt.Println(cpf.String()) // "716.566.867-59"

// CNPJ validation - accepts both formatted and clean input
cnpj, err := cpfcnpj.NewCnpj("22.796.729/0001-59")  // formatted
cnpj, err := cpfcnpj.NewCnpj("22796729000159")      // clean
if err != nil {
    return err
}
fmt.Println(cnpj.String()) // "22.796.729/0001-59"
```

## Usage Examples

### CPF Validation

```go
package main

import (
    "errors"
    "fmt"
    "github.com/n0vdd/cpf_cnpj"
)

func main() {
    // Both formatted and clean inputs work
    examples := []string{
        "716.566.867-59",  // formatted
        "71656686759",     // clean
    }
    
    for _, input := range examples {
        cpf, err := cpfcnpj.NewCpf(input)
        if err != nil {
            if errors.Is(err, cpfcnpj.ErrCPFInvalidLength) {
                fmt.Printf("CPF %s has invalid length\n", input)
            }
            continue
        }
        
        fmt.Printf("Valid CPF: %s\n", cpf.String())
        // Output: Valid CPF: 716.566.867-59
    }
}
```

### CNPJ Validation (Numeric)

```go
// Both formatted and clean inputs work
examples := []string{
    "22.796.729/0001-59",  // formatted
    "22796729000159",      // clean
}

for _, input := range examples {
    cnpj, err := cpfcnpj.NewCnpj(input)
    if err != nil {
        if errors.Is(err, cpfcnpj.ErrCNPJInvalidChecksum) {
            fmt.Printf("CNPJ %s has invalid checksum\n", input)
        }
        continue
    }
    
    fmt.Printf("Valid CNPJ: %s\n", cnpj.String())
    // Output: Valid CNPJ: 22.796.729/0001-59
}
```

### Document Cleaning

```go
// Clean removes formatting and normalizes input
fmt.Println(cpfcnpj.Clean("716.566.867-59"))      // "71656686759"
fmt.Println(cpfcnpj.Clean("22.796.729/0001-59"))  // "22796729000159"
fmt.Println(cpfcnpj.Clean("12.abc.345/01de-35"))  // "12ABC34501DE35"
```

## CNPJ Alfanumérico

This package supports the new Brazilian **CNPJ Alfanumérico** format introduced by [Instrução Normativa RFB nº 2.119/2022](https://www.in.gov.br/en/web/dou/-/instrucao-normativa-rfb-n-2.119-de-21-de-dezembro-de-2022-454078082).

```go
// Both formatted and clean alphanumeric CNPJs work
examples := []string{
    "12.ABC.345/01DE-35",  // formatted
    "12ABC34501DE35",      // clean
}

for _, input := range examples {
    cnpj, err := cpfcnpj.NewCnpj(input)
    if err != nil {
        fmt.Printf("Invalid CNPJ: %v\n", err)
        continue
    }
    
    fmt.Printf("Valid Alphanumeric CNPJ: %s\n", cnpj.String())
    // Output: Valid Alphanumeric CNPJ: 12.ABC.345/01DE-35
}
```

### Official Documentation

- [Receita Federal - CNPJ Alfanumérico](https://www.gov.br/receitafederal/pt-br/assuntos/orientacao-tributaria/cadastros/cnpj/cnpj-alfanumerico)
- [Instrução Normativa RFB nº 2.119/2022](https://www.in.gov.br/en/web/dou/-/instrucao-normativa-rfb-n-2.119-de-21-de-dezembro-de-2022-454078082)

## API Reference

### Types

```go
type CPF string
type CNPJ string
```

### Constructors

```go
// NewCpf validates and creates a CPF instance
func NewCpf(s string) (CPF, error)

// NewCnpj validates and creates a CNPJ instance (supports alphanumeric)
func NewCnpj(s string) (CNPJ, error)
```

### Utilities

```go
// Clean removes formatting and normalizes input
func Clean(s string) string
```

### Methods

```go
// String returns formatted document
func (c *CPF) String() string   // Returns: "716.566.867-59"
func (c *CNPJ) String() string  // Returns: "22.796.729/0001-59" or "12.ABC.345/01DE-35"
```

### Error Types

```go
var (
    // Generic errors
    ErrAllSameDigits    = errors.New("document cannot have all same digits")
    ErrInvalidCharacter = errors.New("document contains invalid character")
    
    // CPF errors
    ErrCPFInvalidLength   = errors.New("CPF must have exactly 11 digits")
    ErrCPFInvalidChecksum = errors.New("CPF checksum validation failed")
    
    // CNPJ errors
    ErrCNPJInvalidLength       = errors.New("CNPJ must have exactly 14 characters")
    ErrCNPJInvalidChecksum     = errors.New("CNPJ checksum validation failed")
    ErrCNPJInvalidAlphanumeric = errors.New("CNPJ alphanumeric format invalid")
)
```

## Error Handling

```go
cpf, err := cpfcnpj.NewCpf("invalid-input")
if err != nil {
    if errors.Is(err, cpfcnpj.ErrCPFInvalidLength) {
        fmt.Println("CPF has wrong length")
    } else if errors.Is(err, cpfcnpj.ErrCPFInvalidChecksum) {
        fmt.Println("CPF checksum is invalid")
    }
    return err
}

// cpf is guaranteed to be valid
fmt.Println("Valid CPF:", cpf.String())
```

## Input Flexibility

This package accepts both formatted and clean inputs for maximum convenience:

| Document Type | Formatted Input | Clean Input | Output |
|---------------|-----------------|-------------|--------|
| CPF | `"716.566.867-59"` | `"71656686759"` | `"716.566.867-59"` |
| CNPJ Numeric | `"22.796.729/0001-59"` | `"22796729000159"` | `"22.796.729/0001-59"` |
| CNPJ Alphanumeric | `"12.ABC.345/01DE-35"` | `"12ABC34501DE35"` | `"12.ABC.345/01DE-35"` |

## License

[MIT](LICENSE)