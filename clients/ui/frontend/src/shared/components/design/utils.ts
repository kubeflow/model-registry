import projectImg from '~/shared/images/UI_icon-Red_Hat-Folder-RGB.svg';
import notebookImg from '~/shared/images/UI_icon-Red_Hat-Wrench-RGB.svg';
import pipelineImg from '~/shared/images/UI_icon-Red_Hat-Branch-RGB.svg';
import pipelineRunImg from '~/shared/images/UI_icon-Red_Hat-Double_arrow_right-RGB.svg';
import clusterStorageImg from '~/shared/images/UI_icon-Red_Hat-Storage-RGB.svg';
import modelServerImg from '~/shared/images/UI_icon-Red_Hat-Server-RGB.svg';
import registeredModelsImg from '~/shared/images/Icon-Red_Hat-Layered_A_Black-RGB.svg';
import deployedModelsImg from '~/shared/images/UI_icon-Red_Hat-Cubes-RGB.svg';
import deployingModelsImg from '~/shared/images/UI_icon-Red_Hat-Server_upload-RGB.svg';
import dataConnectionImg from '~/shared/images/UI_icon-Red_Hat-Connected-RGB.svg';
import userImg from '~/shared/images/UI_icon-Red_Hat-User-RGB.svg';
import groupImg from '~/shared/images/UI_icon-Red_Hat-Shared_workspace-RGB.svg';
import projectEmptyStateImg from '~/shared/images/empty-state-project-overview.svg';
import notebookEmptyStateImg from '~/shared/images/empty-state-notebooks.svg';
import pipelineEmptyStateImg from '~/shared/images/empty-state-pipelines.svg';
import clusterStorageEmptyStateImg from '~/shared/images/empty-state-cluster-storage.svg';
import modelServerEmptyStateImg from '~/shared/images/empty-state-model-serving.svg';
import dataConnectionEmptyStateImg from '~/shared/images/empty-state-data-connections.svg';
import modelRegistryEmptyStateImg from '~/shared/images/empty-state-model-registries.svg';
import storageClassesEmptyStateImg from '~/shared/images/empty-state-storage-classes.svg';
import modelRegistryMissingModelImg from '~/shared/images/no-models-model-registry.svg';
import modelRegistryMissingVersionImg from '~/shared/images/no-versions-model-registry.svg';

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
  project = 'project',
  projectContext = 'projectContext',
  notebook = 'notebook',
  notebookImage = 'notebookImage',
  build = 'build',
  pipelineSetup = 'pipeline-setup',
  pipeline = 'pipeline',
  pipelineRun = 'pipeline-run',
  pipelineExperiment = 'pipeline-experiment',
  pipelineExecution = 'pipeline-execution',
  pipelineArtifact = 'pipeline-artifact',
  clusterStorage = 'cluster-storage',
  model = 'model',
  singleModel = 'single-model',
  multiModel = 'multi-model',
  modelServer = 'model-server',
  registeredModels = 'registered-models',
  deployedModels = 'deployed-models',
  deployingModels = 'deploying-models',
  modelRegistrySettings = 'model-registry-settings',
  servingRuntime = 'serving-runtime',
  distributedWorkload = 'distributed-workload',
  dataConnection = 'data-connection',
  connections = 'connections',
  clusterSettings = 'cluster-settings',
  acceleratorProfile = 'accelerator-profile',
  permissions = 'permissions',
  user = 'user',
  group = 'group',
  storageClasses = 'storageClasses',
  enabledApplications = 'enabled-applications',
  exploreApplications = 'explore-applications',
  resources = 'resources',
}

export const typedIconColor = (objectType: ProjectObjectType): string => {
  switch (objectType) {
    case ProjectObjectType.project:
      return 'var(--ai-project--IconColor)';
    case ProjectObjectType.projectContext:
      return 'var(--ai-project-context--IconColor)';
    case ProjectObjectType.notebook:
      return 'var(--ai-notebook--IconColor)';
    case ProjectObjectType.notebookImage:
      return 'var(--ai-set-up--IconColor)';
    case ProjectObjectType.pipeline:
    case ProjectObjectType.pipelineRun:
    case ProjectObjectType.pipelineExperiment:
    case ProjectObjectType.pipelineExecution:
    case ProjectObjectType.pipelineArtifact:
      return 'var(--ai-pipeline--IconColor)';
    case ProjectObjectType.pipelineSetup:
      return 'var(--ai-set-up--IconColor)';
    case ProjectObjectType.clusterStorage:
    case ProjectObjectType.storageClasses:
      return 'var(--ai-cluster-storage--IconColor)';
    case ProjectObjectType.model:
    case ProjectObjectType.singleModel:
    case ProjectObjectType.multiModel:
    case ProjectObjectType.modelServer:
    case ProjectObjectType.registeredModels:
    case ProjectObjectType.deployedModels:
    case ProjectObjectType.deployingModels:
      return 'var(--ai-model-server--IconColor)';
    case ProjectObjectType.modelRegistrySettings:
      return 'var(--ai-set-up--IconColor)';
    case ProjectObjectType.dataConnection:
    case ProjectObjectType.connections:
      return 'var(--ai-data-connection--IconColor)';
    case ProjectObjectType.user:
      return 'var(--ai-user--IconColor)';
    case ProjectObjectType.group:
      return 'var(--ai-group--IconColor)';
    case ProjectObjectType.permissions:
      return 'var(--ai-set-up--IconColor)';
    case ProjectObjectType.enabledApplications:
    case ProjectObjectType.exploreApplications:
      return 'var(--ai-config--IconColor)';
    case ProjectObjectType.resources:
      return 'var(--ai-general--IconColor)';
    case ProjectObjectType.distributedWorkload:
      return 'var(--ai-serving--IconColor)';
    case ProjectObjectType.clusterSettings:
    case ProjectObjectType.acceleratorProfile:
      return 'var(--ai-set-up--IconColor)';
    case ProjectObjectType.servingRuntime:
      return 'var(--ai-set-up--IconColor)';
    default:
      return '';
  }
};

export const typedBackgroundColor = (objectType: ProjectObjectType): string => {
  switch (objectType) {
    case ProjectObjectType.project:
      return 'var(--ai-project--BackgroundColor)';
    case ProjectObjectType.projectContext:
      return 'var(--ai-project-context--BackgroundColor)';
    case ProjectObjectType.notebook:
      return 'var(--ai-notebook--BackgroundColor)';
    case ProjectObjectType.notebookImage:
      return 'var(--ai-set-up--BackgroundColor)';
    case ProjectObjectType.pipeline:
    case ProjectObjectType.pipelineRun:
    case ProjectObjectType.pipelineExperiment:
    case ProjectObjectType.pipelineExecution:
    case ProjectObjectType.pipelineArtifact:
      return 'var(--ai-pipeline--BackgroundColor)';
    case ProjectObjectType.pipelineSetup:
      return 'var(--ai-set-up--BackgroundColor)';
    case ProjectObjectType.clusterStorage:
    case ProjectObjectType.storageClasses:
      return 'var(--ai-cluster-storage--BackgroundColor)';
    case ProjectObjectType.model:
    case ProjectObjectType.singleModel:
    case ProjectObjectType.multiModel:
    case ProjectObjectType.modelServer:
    case ProjectObjectType.registeredModels:
    case ProjectObjectType.deployedModels:
    case ProjectObjectType.deployingModels:
      return 'var(--ai-model-server--BackgroundColor)';
    case ProjectObjectType.modelRegistrySettings:
      return 'var(--ai-set-up--BackgroundColor)';
    case ProjectObjectType.dataConnection:
    case ProjectObjectType.connections:
      return 'var(--ai-data-connection--BackgroundColor)';
    case ProjectObjectType.user:
      return 'var(--ai-user--BackgroundColor)';
    case ProjectObjectType.group:
      return 'var(--ai-group--BackgroundColor)';
    case ProjectObjectType.permissions:
      return 'var(--ai-set-up--BackgroundColor)';
    case ProjectObjectType.enabledApplications:
    case ProjectObjectType.exploreApplications:
      return 'var(--ai-config--BackgroundColor)';
    case ProjectObjectType.resources:
      return 'var(--ai-general--BackgroundColor)';
    case ProjectObjectType.distributedWorkload:
      return 'var(--ai-serving--BackgroundColor)';
    case ProjectObjectType.clusterSettings:
    case ProjectObjectType.acceleratorProfile:
      return 'var(--ai-set-up--BackgroundColor)';
    case ProjectObjectType.servingRuntime:
      return 'var(--ai-set-up--BackgroundColor)';
    default:
      return '';
  }
};

export const typedColor = (objectType: ProjectObjectType): string => {
  switch (objectType) {
    case ProjectObjectType.project:
      return 'var(--ai-project--Color)';
    case ProjectObjectType.projectContext:
      return 'var(--ai-project-context--Color)';
    case ProjectObjectType.notebook:
    case ProjectObjectType.notebookImage:
      return 'var(--ai-training--BackgroundColor)';
    case ProjectObjectType.build:
      return 'var(--ai-model-server--Color)';
    case ProjectObjectType.pipeline:
    case ProjectObjectType.pipelineRun:
    case ProjectObjectType.pipelineExecution:
    case ProjectObjectType.pipelineArtifact:
      return 'var(--ai-pipeline--Color)';
    case ProjectObjectType.pipelineSetup:
      return 'var(--ai-set-up--Color)';
    case ProjectObjectType.clusterStorage:
      return 'var(--ai-cluster-storage--Color)';
    case ProjectObjectType.modelServer:
    case ProjectObjectType.registeredModels:
    case ProjectObjectType.deployedModels:
    case ProjectObjectType.deployingModels:
      return 'var(--ai-model-server--Color)';
    case ProjectObjectType.modelRegistrySettings:
      return 'var(--ai-set-up--Color)';
    case ProjectObjectType.dataConnection:
    case ProjectObjectType.connections:
      return 'var(--ai-data-connection--Color)';
    case ProjectObjectType.user:
      return 'var(--ai-user--Color)';
    case ProjectObjectType.group:
      return 'var(--ai-group--Color)';
    default:
      return '';
  }
};

export const typedObjectImage = (objectType: ProjectObjectType): string => {
  switch (objectType) {
    case ProjectObjectType.project:
    case ProjectObjectType.projectContext:
      return projectImg;
    case ProjectObjectType.notebook:
      return notebookImg;
    case ProjectObjectType.pipeline:
    case ProjectObjectType.pipelineSetup:
      return pipelineImg;
    case ProjectObjectType.pipelineRun:
      return pipelineRunImg;
    case ProjectObjectType.clusterStorage:
      return clusterStorageImg;
    case ProjectObjectType.modelServer:
      return modelServerImg;
    case ProjectObjectType.registeredModels:
      return registeredModelsImg;
    case ProjectObjectType.deployedModels:
      return deployedModelsImg;
    case ProjectObjectType.deployingModels:
      return deployingModelsImg;
    case ProjectObjectType.dataConnection:
    case ProjectObjectType.connections:
      return dataConnectionImg;
    case ProjectObjectType.user:
      return userImg;
    case ProjectObjectType.group:
      return groupImg;
    default:
      return '';
  }
};

export const typedEmptyImage = (objectType: ProjectObjectType, option?: string): string => {
  switch (objectType) {
    case ProjectObjectType.project:
    case ProjectObjectType.projectContext:
      return projectEmptyStateImg;
    case ProjectObjectType.notebook:
      return notebookEmptyStateImg;
    case ProjectObjectType.pipeline:
    case ProjectObjectType.pipelineRun:
    case ProjectObjectType.pipelineSetup:
      return pipelineEmptyStateImg;
    case ProjectObjectType.clusterStorage:
      return clusterStorageEmptyStateImg;
    case ProjectObjectType.modelServer:
      return modelServerEmptyStateImg;
    case ProjectObjectType.registeredModels:
      switch (option) {
        case 'MissingModel':
          return modelRegistryMissingModelImg;
        case 'MissingVersion':
          return modelRegistryMissingVersionImg;
        case 'MissingDeployment':
          return modelServerEmptyStateImg;
        default:
          return modelRegistryEmptyStateImg;
      }
    case ProjectObjectType.storageClasses:
      return storageClassesEmptyStateImg;
    case ProjectObjectType.dataConnection:
    case ProjectObjectType.connections:
      return dataConnectionEmptyStateImg;
    default:
      return '';
  }
};

export const sectionTypeIconColor = (sectionType: SectionType): string => {
  switch (sectionType) {
    case SectionType.setup:
      return 'var(--ai-set-up--IconColor)';
    case SectionType.organize:
      return 'var(--ai-organize--IconColor)';
    case SectionType.training:
      return 'var(--ai-training--IconColor)';
    case SectionType.serving:
      return 'var(--ai-serving--IconColor)';
    case SectionType.general:
      return 'var(--ai-general--IconColor)';
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
