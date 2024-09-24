/**
 * Verify the relative route to the cypress host
 * e.g. If page is running on `https://localhost:9001/pipelines`
 * calling verifyRelativeURL('/pipelines') will check whether the full URL matches the URL above
 */
export const verifyRelativeURL = (relativeURL: string): Cypress.Chainable<string> => {
  return cy
    .location()
    .then((location) =>
      cy.url().should('eq', `${location.protocol}//${location.host}${relativeURL}`),
    );
};
