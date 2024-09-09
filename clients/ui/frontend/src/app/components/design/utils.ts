import registerModelImg from '~/images/UI_icon-Cubes-RGB.svg';
import modelRegistryEmptyStateImg from '~/images/empty-state-model-registries.svg';
import './vars.scss';

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

export const typedEmptyImage = (objectType: ProjectObjectType): string => {
  switch (objectType) {
    case ProjectObjectType.registeredModels:
      return modelRegistryEmptyStateImg;
    default:
      return '';
  }
};
