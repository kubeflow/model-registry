# Dev Setup

## Requirements

This project requires the following tools to be installed on your system:

- [NodeJS and NPM](https://nodejs.org/)
  - Node recommended version -> `20.17.0`
  - NPM recommended version -> `10.8.2`

## Development

1. Clone the repository

      ``` bash
      git clone https://github.com/kubeflow/model-registry
      ```

2. Within the repo context, install project dependencies

     ```bash
     cd model-registry/clients/ui/frontend && npm install
     ```

### Build project

```bash
npm run build
```

### Serve development content

This is the default context for running a local UI.  Make sure you build the project using the instructions above prior to running the command below.

```bash
npm run start:dev
```

For in-depth local run guidance review the [contribution guidelines](../../../ui/CONTRIBUTING.md).

### Testing

Run the mock tests.

  ```bash
  npm run test:cypress-ci
  ```

For in-depth testing guidance review the [testing guidelines](./testing.md)
