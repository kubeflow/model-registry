import 'cypress-axe';

/* eslint-disable @typescript-eslint/no-namespace */
declare global {
  namespace Cypress {
    interface Chainable {
      testA11y: (context?: Parameters<cy['checkA11y']>[0]) => void;
    }
  }
}

Cypress.Commands.add('testA11y', { prevSubject: 'optional' }, (subject, context) => {
  const test = (c: Parameters<typeof cy.checkA11y>[0]) => {
    cy.window({ log: false }).then((win) => {
      // inject on demand
      if (!(win as { axe: unknown }).axe) {
        cy.injectAxe();
      }
      cy.checkA11y(
        c,
        {
          includedImpacts: ['serious', 'critical'],
        },
        (violations) => {
          cy.task(
            'error',
            `${violations.length} accessibility violation${violations.length === 1 ? '' : 's'} ${
              violations.length === 1 ? 'was' : 'were'
            } detected`,
          );
          // pluck specific keys to keep the table readable
          const violationData = violations.map(({ id, impact, description, nodes }) => ({
            id,
            impact,
            description,
            nodes: nodes.length,
          }));

          cy.task('table', violationData);

          cy.task(
            'log',
            violations
              .map(
                ({ nodes }, i) =>
                  `${i}. Affected elements:\n${nodes.map(
                    ({ target, failureSummary, ancestry }) =>
                      `\t${failureSummary} - ${target
                        .map((node) => `"${node}"\n${ancestry}`)
                        .join(', ')}`,
                  )}`,
              )
              .join('\n'),
          );
        },
      );
    });
  };
  if (!context && subject) {
    cy.wrap(subject).each(($el) => {
      Cypress.log({ displayName: 'testA11y', $el });
      test($el[0]);
    });
  } else {
    Cypress.log({ displayName: 'testA11y' });
    test(context);
  }
});
