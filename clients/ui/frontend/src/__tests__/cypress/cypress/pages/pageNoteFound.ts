class PageNotFound {
  visit() {
    cy.visit(`/force-not-found-page`, { failOnStatusCode: false });
    this.wait();
  }

  private wait() {
    this.findPage();
    cy.testA11y();
  }

  findPage() {
    return cy.get('h1:contains("404 Page not found")');
  }
}

export const pageNotfound = new PageNotFound();
