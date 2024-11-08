export class ToastNotification {
  constructor(private title: string) {}

  find(): Cypress.Chainable<JQuery<HTMLElement>> {
    return cy.findByText(this.title);
  }
}
