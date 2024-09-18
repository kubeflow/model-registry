import type { GenericStaticResponse, RouteHandlerController } from 'cypress/types/net-stubbing';
import type {
  ModelArtifact,
  ModelArtifactList,
  ModelRegistry,
  ModelRegistryResponse,
  ModelVersion,
  ModelVersionList,
  RegisteredModel,
  RegisteredModelList,
} from '~/app/types';

type SuccessErrorResponse = {
  success: boolean;
  error?: string;
};

type ApiResponse<V = SuccessErrorResponse> =
  | V
  | GenericStaticResponse<string, V>
  | RouteHandlerController;

type Replacement<R extends string = string> = Record<R, string | undefined>;
type Query<Q extends string = string> = Record<Q, string>;

type Options = { path?: Replacement; query?: Query; times?: number } | null;

/* eslint-disable @typescript-eslint/no-namespace */
declare global {
  namespace Cypress {
    interface Chainable {
      interceptApi: ((
        type: 'GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models',
        options: { path: { modelRegistryName: string; apiVersion: string } },
        response: ApiResponse<ModelRegistryResponse<RegisteredModelList>>,
      ) => Cypress.Chainable<null>) &
        ((
          type: 'POST /api/:apiVersion/model_registry/:modelRegistryName/registered_models',
          options: { path: { modelRegistryName: string; apiVersion: string } },
          response: ApiResponse<RegisteredModel>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId/versions',
          options: {
            path: { modelRegistryName: string; apiVersion: string; registeredModelId: number };
          },
          response: ApiResponse<ModelRegistryResponse<ModelVersionList>>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'POST /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId/versions',
          options: {
            path: { modelRegistryName: string; apiVersion: string; registeredModelId: number };
          },
          response: ApiResponse<ModelVersion>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'GET /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId',
          options: {
            path: { modelRegistryName: string; apiVersion: string; registeredModelId: number };
          },
          response: ApiResponse<RegisteredModel>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'PATCH /api/:apiVersion/model_registry/:modelRegistryName/registered_models/:registeredModelId',
          options: {
            path: { modelRegistryName: string; apiVersion: string; registeredModelId: number };
          },
          response: ApiResponse<RegisteredModel>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId',
          options: {
            path: { modelRegistryName: string; apiVersion: string; modelVersionId: number };
          },
          response: ApiResponse<ModelVersion>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'GET /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId/artifacts',
          options: {
            path: { modelRegistryName: string; apiVersion: string; modelVersionId: number };
          },
          response: ApiResponse<ModelArtifactList>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'POST /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId/artifacts',
          options: {
            path: { modelRegistryName: string; apiVersion: string; modelVersionId: number };
          },
          response: ApiResponse<ModelArtifact>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'PATCH /api/:apiVersion/model_registry/:modelRegistryName/model_versions/:modelVersionId',
          options: {
            path: { modelRegistryName: string; apiVersion: string; modelVersionId: number };
          },
          response: ApiResponse<ModelVersion | undefined>,
        ) => Cypress.Chainable<null>) &
        ((
          type: 'GET /api/:apiVersion/model_registry',
          options: { path: { apiVersion: string } },
          response: ApiResponse<ModelRegistryResponse<ModelRegistry[]>>,
        ) => Cypress.Chainable<null>);
    }
  }
}

Cypress.Commands.add(
  'interceptApi',
  (type: string, ...args: [Options | null, ApiResponse<unknown>] | [ApiResponse<unknown>]) => {
    if (!type) {
      throw new Error('Invalid type parameter.');
    }
    const options = args.length === 2 ? args[0] : null;
    const response = (args.length === 2 ? args[1] : args[0]) ?? '';

    const pathParts = type.match(/:[a-z][a-zA-Z0-9-_]+/g);
    const [method, staticPathname] = type.split(' ');
    let pathname = staticPathname;
    if (pathParts?.length) {
      if (!options || !options.path) {
        throw new Error(`${type}: missing path replacements`);
      }
      const { path: pathReplacements } = options;
      pathParts.forEach((p) => {
        // remove the starting colun from the regex match
        const part = p.substring(1);
        const replacement = pathReplacements[part];
        if (!replacement) {
          throw new Error(`${type} missing path replacement: ${part}`);
        }
        pathname = pathname.replace(new RegExp(`:${part}\\b`), replacement);
      });
    }
    return cy.intercept(
      { method, pathname, query: options?.query, ...(options?.times && { times: options.times }) },
      response,
    );
  },
);
