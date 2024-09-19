export class Contextual<E extends HTMLElement> {
  constructor(private parentSelector: () => Cypress.Chainable<JQuery<E>>) {}

  find(): Cypress.Chainable<JQuery<E>> {
    return this.parentSelector();
  }
}
