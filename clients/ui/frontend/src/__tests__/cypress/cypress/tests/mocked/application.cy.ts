import { pageNotfound } from "~/src/__tests__/cypress/cypress/pages/pageNoteFound";
import { home } from "~/src/__tests__/cypress/cypress/pages/home";

describe('Application', () => {

    it('Page not found should render', () => {
        pageNotfound.visit()
    });

    it('Home page should have primary button', () => {
        home.visit()
        home.findButton();
    });

});