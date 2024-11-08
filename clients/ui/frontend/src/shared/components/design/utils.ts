import registerModelImg from '~/images/UI_icon-Cubes-RGB.svg';
import modelRegistryEmptyStateImg from '~/images/empty-state-model-registries.svg';
import modelRegistryMissingModelImg from '~/images/no-models-model-registry.svg';
import modelRegistryMissingVersionImg from '~/images/no-versions-model-registry.svg';

import './vars.scss';
/* eslint-disable @typescript-eslint/no-unnecessary-condition */
// These conditions are required for future object types that may be added later.

export enum ProjectObjectType {
  registeredModels = 'registered-models',
}

export const typedBackgroundColor = (objectType: ProjectObjectType): string => {
  switch (objectType) {
    case ProjectObjectType.registeredModels:
      return 'var(--ai-model-server--BackgroundColor)';
    default:
      return '';
  }
};

export const typedObjectImage = (objectType: ProjectObjectType): string => {
  switch (objectType) {
    case ProjectObjectType.registeredModels:
      return registerModelImg;
    default:
      return '';
  }
};

export const typedEmptyImage = (objectType: ProjectObjectType, option?: string): string => {
  switch (objectType) {
    case ProjectObjectType.registeredModels:
      switch (option) {
        case 'MissingModel':
          return modelRegistryMissingModelImg;
        case 'MissingVersion':
          return modelRegistryMissingVersionImg;
        default:
          return modelRegistryEmptyStateImg;
      }
    default:
      return '';
  }
};
