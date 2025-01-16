[dev setup documentation]: docs/dev-setup.md#requirements

# Contributing

Contributing encompasses repository specific requirements.

## Requirements

To review the requirements, please refer to the [dev setup documentation].

## Writing code

### Running locally

The recommended flow for development is to have two sessions, one for the "frontend":

```bash
npm run start:dev
```

And one for the "backend":

```bash
cd ../bff
go run ./cmd/main.go  --port=4000 --static-assets-dir=./static --mock-k8s-client=true --mock-mr-client=true --dev-mode=true --dev-mode-port=8080  --standalone-mode=true
```

Once you have both services ready, you can open the dashboard locally at: `http://localhost:4010`. The dev server will reload automatically when you make changes.

You can also run an automated command to run both services:

```bash
cd .. && make dev-start
```

## Debugging and Testing

See [frontend testing guidelines](docs/testing.md) for more information.

### Environment Variables

The following environment variables are used to configure the deployment and development environment for the Model Registry UI. These variables should be defined in a `.env.local` file in the `clients/ui` directory of the project. **These values will affect the build and push commands**.

- `LOGO=logo-light-theme.svg`
  - The file name for the logo used in the light theme.
- `LOGO_DARK=logo-dark-theme.svg`
  - The file name for the logo used in the dark theme.
- `FAVICON=favicon.ico`
  - The file name for the favicon of the application.
- `PRODUCT_NAME="Model Registry"`
  - The name of the product displayed in the UI.
- `STYLE_THEME=mui-theme`
  - The style theme used for the UI, in this case, Material-UI theme.
