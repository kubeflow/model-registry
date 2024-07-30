# Model Registry UI Testing

## Cypress Tests

Cypress is used to run tests against the frontend while mocking all network requests.

Single command to run all Cypress tests or a specific test (build frontend, start HTTP server, run Cypress):

```bash
npm run test:cypress-ci

npm run test:cypress-ci -- --spec "**/testfile.cy.ts"
```

Cypress tests require a frontend server to be running.

To best match production, build the frontend and use a lightweight HTTP server to host the files. This method will require manual rebuilds when changes are made to the dashboard frontend code.

```bash
npm run cypress:server:build
npm run cypress:server
```

To run all Cypress tests or a specific test headless

```bash
npm run cypress:run:mock

npm run cypress:run:mock -- --spec "**/testfile.cy.ts"
```

To open the Cypress GUI run

```bash
npm run cypress:open:mock
```
