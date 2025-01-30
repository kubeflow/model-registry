import { SubComponentBase } from '~/__tests__/cypress/cypress/pages/components/subComponents/SubComponentBase';

export class SearchSelector extends SubComponentBase {
  constructor(
    private selectorId: string,
    contextSelectorId?: string,
  ) {
    super(contextSelectorId);
  }

  private findContextualItem(suffix: string): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.findScope().document().findByTestId(`${this.selectorId}-${suffix}`);
  }

  findItem(name: string, useMenuList: boolean): Cypress.Chainable<JQuery<HTMLElement>> {
    const list = useMenuList ? this.findMenuList() : this.findResultTableList();
    return list.contains(name).should('exist');
  }

  selectItem(name: string, useMenuList = false): void {
    this.findItem(name, useMenuList).click();
  }

  findSearchInput(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.findContextualItem('search');
  }

  findToggleButton(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.findContextualItem('toggle');
  }

  findResultTableList(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.findContextualItem('table-list');
  }

  findSearchHelpText(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.findContextualItem('searchHelpText');
  }

  findMenu(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.findContextualItem('menu');
  }

  findMenuList(): Cypress.Chainable<JQuery<HTMLElement>> {
    return this.findContextualItem('menuList');
  }

  // Search for an item by typing into the search input
  searchItem(name: string): void {
    this.findSearchInput().clear().type(name);
  }

  // Perform the entire process: open, search, and select
  openAndSelectItem(name: string, useMenuList = false): void {
    this.findToggleButton().click();
    this.searchItem(name);
    this.selectItem(name, useMenuList);
  }
}
