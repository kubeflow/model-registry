class Home {
  visit() {
    cy.visit(`/`);
  }

  findTitle() {
    cy.get(`h1`).should(`have.text`, `Model registry`);
  }
}

export const home = new Home();
