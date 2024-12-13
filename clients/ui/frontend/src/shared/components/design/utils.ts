import registerModelImg from '~/images/UI_icon-Cubes-RGB.svg';
import modelRegistryEmptyStateImg from '~/images/empty-state-model-registries.svg';
import modelRegistryMissingModelImg from '~/images/no-models-model-registry.svg';
import modelRegistryMissingVersionImg from '~/images/no-versions-model-registry.svg';

import './vars.scss';
/* eslint-disable @typescript-eslint/no-unnecessary-condition */
// These conditions are required for future object types that may be added later.

export enum SectionType {
  setup = 'set-up',
  organize = 'organize',
  training = 'training',
  serving = 'serving',
  general = 'general',
}

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

export const sectionTypeBackgroundColor = (sectionType: SectionType): string => {
  switch (sectionType) {
    case SectionType.setup:
      return 'var(--ai-set-up--BackgroundColor)';
    case SectionType.organize:
      return 'var(--ai-organize--BackgroundColor)';
    case SectionType.training:
      return 'var(--ai-training--BackgroundColor)';
    case SectionType.serving:
      return 'var(--ai-serving--BackgroundColor)';
    case SectionType.general:
      return 'var(--ai-general--BackgroundColor)';
    default:
      return '';
  }
};

export const sectionTypeBorderColor = (sectionType: SectionType): string => {
  switch (sectionType) {
    case SectionType.setup:
      return 'var(--ai-set-up--BorderColor)';
    case SectionType.organize:
      return 'var(--ai-organize--BorderColor)';
    case SectionType.training:
      return 'var(--ai-training--BorderColor)';
    case SectionType.serving:
      return 'var(--ai-serving--BorderColor)';
    case SectionType.general:
      return 'var(--ai-general--BorderColor)';
    default:
      return '';
  }
};
