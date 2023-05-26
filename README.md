# swag-ts

[![](https://img.shields.io/github/actions/workflow/status/ditschedev/swag-ts/test.yml?branch=main&longCache=true&label=Test&logo=github%20actions&logoColor=fff)](https://github.com/ditschedev/swag-ts/actions?query=workflow%3ATest)
[![Go Report Card](https://goreportcard.com/badge/github.com/ditschedev/swag-ts)](https://goreportcard.com/report/github.com/ditschedev/swag-ts)

Simply provide a OpenAPI Specification and swag-ts will generate typescript types for you. You can provide json or yaml definitions on your local filesystem or a remote url.

## Motivation
Why another type generator for OpenAPI (Swagger)? Well it's because I could not find a generator that only generates typescript types. 
Most generators also generate runtime code which I don't necessarily need. I just want to have the typescript types to use them in my frontend application.

If thats something for you, feel free to use it. If you need more functionality, feel free to open an issue or a pull request.

## Installation

```bash
go install github.com/ditschedev/swag-ts@latest
```

## Usage

```bash
Usage:
  swag-ts [flags]

Flags:
  -f, --file string     file path or url to the OpenAPI Specification
  -h, --help            help for swag-ts
  -o, --output string   output file for generated definitions (default "./types/swagger.ts")
  -v, --version         shows the version of the cli
```

## Format
This library aims to only provide typescript type definitions from a given OpenAPI Specification. It does not provide any runtime functionality.
All types are exported as `interface`.

For example, the following Schema:
```yaml
LoginResponse:
  required:
    - token
  type: object
  properties:
    token:
      minLength: 1
      type: string
  additionalProperties: false

LoginResponseWrapper:
  required:
    - data
  type: object
  properties:
    data:
      $ref: '#/components/schemas/LoginResponse'
    message:
      type: string
      nullable: true
  additionalProperties: false
```

will be converted to the following typescript definitions:
```typescript
export interface LoginResponse {
  token: string;
}

export interface LoginResponseWrapper {
  data: LoginResponse;
  message?: string | null;
}
```

### Enums
Enums are converted to typescript enums. See the example below:
```yaml
Car:
  required:
    - manufacturer
  type: object
  properties:
    manufacturer:
      $ref: "#/components/schemas/CarManufacturer"
CarManufacturer:
  type: string
  enum:
    - BMW
    - Mercedes
    - Audi  
```

will be converted to the following typescript definitions:
```typescript
export interface Car {
    manufacturer: CarManufacturer;
}

export enum CarManufacturer {
    BMW = "BMW",
    Mercedes = "Mercedes",
    Audi = "Audi",
}
```