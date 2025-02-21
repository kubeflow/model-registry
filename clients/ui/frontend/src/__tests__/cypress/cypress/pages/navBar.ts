class NavBar {
  findBrand() {
    return cy.get('.pf-v5-c-brand');
  }

  findNavToggleButton() {
    return cy.get('#page-nav-toggle');
  }

  findNamespaceSelector() {
    return cy.get('.kubeflow-u-namespace-select');
  }

  selectNamespace(name: string) {
    this.findNamespaceSelector().findByRole('button').click();
    cy.findByRole('option', { name }).click();
  }

  findUsername() {
    return cy.findByTestId('user-menu-toggle-button');
  }

  openUserMenu() {
    this.findUsername().click();
  }

  shouldNamespaceSelectorHaveNoItems() {
    this.findNamespaceSelector().click();
    cy.findByRole('option').should('not.exist');
  }
}

export const navBar = new NavBar();
