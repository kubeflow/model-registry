/**
 * A SubComponent is a component that doesn't make up a full page and will be consumed in other page
 * objects. This could be a complex field, a group of fields, or some section. Typically not large
 * enough to warrant its own standalone page object.
 *
 * Primary use-case example:
 *   class Foo extends SubComponentBase {
 *     constructor(private myTestId: string, scopedTestId?: string) {
 *       super(scopedTestId);
 *     }
 *
 *     private find(suffix: string) {
 *       return this.findScope().getByTestId(`${this.myTestId}-${suffix}`);
 *     }
 *
 *     selectItem(name: string) {
 *       // "list" would be an internal suffix for your component to know where the "items" are
 *       return this.find('list').findDropdownItem(name);
 *     }
 *   }
 *
 * Search uses of this component to see further examples
 */
export class SubComponentBase {
  constructor(private scopedBaseTestId?: string) {}

  /** Allows for extended classes to make use of a simple one-check for their `find()` calls */
  protected findScope(): (Cypress.cy & CyEventEmitter) | Cypress.Chainable<JQuery<HTMLElement>> {
    if (this.scopedBaseTestId) {
      return cy.findByTestId(this.scopedBaseTestId);
    }

    return cy;
  }
}
