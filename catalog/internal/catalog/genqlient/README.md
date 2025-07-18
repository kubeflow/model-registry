## Using Genqlient with the Red Hat Ecosystem Catalog

This directory contains the necessary files to generate a type-safe Go GraphQL client for the Red Hat Ecosystem Catalog using [genqlient](https://github.com/Khan/genqlient).

### File Structure

-   `genqlient.yaml`: The configuration file for `genqlient`. It specifies the location of the GraphQL schema, the directory containing the GraphQL queries, and the output file for the generated code.
-   `queries/`: This directory contains the GraphQL schema and query files.
    -   `schema.graphql`: The GraphQL schema for the Red Hat Ecosystem Catalog API.
    -   `*.graphql`: Files containing the GraphQL queries.

### Generating the Client

To regenerate the client, you will need to have `genqlient` installed. You can install it by running:

```bash
go install github.com/Khan/genqlient
```

Once installed, you can generate the client by running the following command from the root of the project:

```bash
~/go/bin/genqlient catalog/internal/catalog/genqlient/genqlient.yaml
```

This will generate the `generated.go` file in the current directory.

### Downloading the Schema

The `schema.graphql` file can be updated by downloading the latest version from the Red Hat Ecosystem Catalog API. You can do this using the `get-graphql-schema` tool:

```bash
npx get-graphql-schema https://catalog.redhat.com/api/containers/graphql/ > catalog/internal/catalog/genqlient/queries/schema.graphql
```

This will download the schema and save it to the correct location. After updating the schema, you should regenerate the client to ensure it is up to date. 