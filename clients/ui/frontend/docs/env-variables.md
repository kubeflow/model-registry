# Environment Variables

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
- `KUBEFLOW_USERNAME="user@example.com"`
  - The default username for the application. **DO NOT CHANGE THIS IF BFF IS SET IN MOCK MODE**.
