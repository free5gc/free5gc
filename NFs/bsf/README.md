# BSF - Binding Support Function

5G Core Network BSF (Binding Support Function) implementation for free5GC.

## Features

- 3GPP TS 29.521 compliant BSF implementation
- PCF binding management with dynamic PolicyContext data extraction  
- MongoDB integration for binding persistence
- SBI (Service Based Interface) support
- NRF registration and service discovery
- RESTful API endpoints for binding operations
- Comprehensive metrics and monitoring

## Quick Start

### Build

```bash
go build -o bin/bsf ./cmd/bsf
```

### Run

```bash
./bin/bsf -config config/bsfcfg.yaml
```

## API Endpoints

- `POST /nbsf-management/v1/pcfBindings` - Register PCF binding
- `GET /nbsf-management/v1/pcfBindings/{bindingId}` - Get PCF binding  
- `PATCH /nbsf-management/v1/pcfBindings/{bindingId}` - Update PCF binding
- `DELETE /nbsf-management/v1/pcfBindings/{bindingId}` - Delete PCF binding
- `GET /nbsf-management/v1/pcfBindings` - Query PCF bindings

## Configuration

See `config/bsfcfg.yaml` for configuration options.

## Integration

This BSF is designed to work with:
- free5GC core network
- PCF (Policy Control Function)
- NRF (NF Repository Function)
- MongoDB database

## License

Apache 2.0 License
