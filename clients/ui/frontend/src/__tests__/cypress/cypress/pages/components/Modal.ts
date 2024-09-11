import type { ByRoleOptions } from '@testing-library/react';

export class Modal {
  constructor(private title: ByRoleOptions['name']) {}

  shouldBeOpen(open = true): void {
    if (open) {
      this.find().testA11y();
    } else {
      this.find().should('not.exist');
    }
  }

  find(): Cypress.Chainable<JQuery<HTMLElement>> {
    return cy.findByRole('dialog', { name: this.title });
  }

  findCloseButton(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.find().findByRole('button', { name: 'Close' });
  }

  findCancelButton(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.findFooter().findByRole('button', { name: 'Cancel' });
  }

  findFooter(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.find().find('footer');
  }
}
