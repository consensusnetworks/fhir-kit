# FHIR Kit

> Simple to spin-up FHIR development kit.

## Setup

1. Clone the repo

```bash
git clone https://github.com/consensusnetworks/fhir-kit
```

2. Start the server in a container

```bash
make start
```

To stop the server either press `Ctrl+C` or run `make stop` in a separate terminal.

3. Make sure the FHIR server is running

```bash
curl http://localhost:9090/ping
```

## FHIR Server

The FHIR Server is a RESTful API that uses the [FHIR R4](http://hl7.org/fhir/R4/) specification and supports US Core profiles.

Base URL: `http://localhost:9090`

To access the supported resources just append the resource name (PascalCase) to the base URL.

Example: `http://localhost:9090/Patient`

## Supported Resources

| Resource                                                               |   Endpoint   |
|------------------------------------------------------------------------|:------------:|
| [CapabilityStatement](http://hl7.org/fhir/R4/capabilitystatement.html) |  /metadata   |
| [Patient](http://hl7.org/fhir/R4/patient.html)                         |   /Patient   |
| [Procedure](http://hl7.org/fhir/R4/procedure.html)                     |  /Procedure  |
 | [Observation](http://hl7.org/fhir/R4/observation.html)                 | /Observation |

## Examples

Get the FHIR server's capability statement. The capability statement describes
the server's supported resources and operations so a client
may use it as an interface definition when interacting with the server.

```bash
curl -X GET http://localhost:9090/metadata
```

Create patient

```bash
curl -X POST -H "Content-Type: application/json" -d '{"resourceType": "Patient", "name": [{"given": ["John"], "family": "Doe"}]}' http://localhost:9090/Patient
```

Create procedure

```bash
curl -X POST -H "Content-Type: application/json" -d '{"subject":{"reference":"25oYHe8zCfx52wp9S8RKEVjEyTw"}}' http://localhost:9090/Procedure
```
