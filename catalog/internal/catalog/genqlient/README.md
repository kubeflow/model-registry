## Using Genqlient with the Red Hat Ecosystem Catalog

This directory contains the necessary files to generate a type-safe Go GraphQL client for the Red Hat Ecosystem Catalog using [genqlient](https://github.com/Khan/genqlient).

### File Structure

-   `genqlient.yaml`: The configuration file for `genqlient`. It specifies the location of the GraphQL schema, the directory containing the GraphQL queries, and the output file for the generated code.
-   `queries/`: This directory contains the GraphQL schema and query files.
    -   `schema.graphql`: The GraphQL schema for the Red Hat Ecosystem Catalog API.
    -   `*.graphql`: Files containing the GraphQL queries.

### Generating the Client

To regenerate the client, you will first need to ensure the required tools are installed by running `make deps` from the project root. Once the tools are installed, you can generate the client by running the following command from the `catalog` directory:

```bash
make gen/graphql
```

This will generate the `generated.go` file in the current directory.

### Downloading the Schema

The `schema.graphql` file can be updated by downloading the latest version from the Red Hat Ecosystem Catalog API. You can do this by running the following command from the `catalog` directory:

```bash
make download/graphql-schema
```

This will download the schema and save it to the correct location. After updating the schema, you should regenerate the client to ensure it is up to date. 