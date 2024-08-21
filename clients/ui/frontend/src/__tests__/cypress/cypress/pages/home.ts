class Home {
  visit() {
    cy.visit(`/`);
  }

  findButton() {
    return cy.get('button:contains("Primary Action")');
  }
}

export const home = new Home();
