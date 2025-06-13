import {
  projectImg,
  notebookImg,
  pipelineImg,
  pipelineRunImg,
  clusterStorageImg,
  modelServerImg,
  registeredModelsImg,
  deployedModelsImg,
  deployingModelsImg,
  dataConnectionImg,
  userImg,
  groupImg,
  projectEmptyStateImg,
  notebookEmptyStateImg,
  pipelineEmptyStateImg,
  clusterStorageEmptyStateImg,
  modelServerEmptyStateImg,
  dataConnectionEmptyStateImg,
  modelRegistryEmptyStateImg,
  storageClassesEmptyStateImg,
  modelRegistryMissingModelImg,
  modelRegistryMissingVersionImg,
  modelRegistrySelectImg,
} from '../../images';

// import './vars.scss'; // This file does not exist

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
  modelCustomization = 'model-customization',
  labTuning = 'lab-tuning',
  clusterStorage = 'cluster-storage',
  model = 'model',
  singleModel = 'single-model',
  multiModel = 'multi-model',
  modelServer = 'model-server',
  modelCatalog = 'model-catalog',
  registeredModels = 'registered-models',
  modelRegistryContext = 'modelRegistryContext',
  deployedModels = 'deployed-models',
  deployingModels = 'deploying-models',
  deployedModelsList = 'deployed-models-list',
  modelRegistrySettings = 'model-registry-settings',
  modelRegistry = 'model-registry',
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
    case ProjectObjectType.modelCatalog:
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
    case ProjectObjectType.deployedModelsList:
    case ProjectObjectType.modelRegistrySettings:
      return 'var(--ai-set-up--IconColor)';
    case ProjectObjectType.modelRegistryContext:
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
    case ProjectObjectType.modelCatalog:
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
    case ProjectObjectType.modelRegistryContext:
    case ProjectObjectType.deployedModels:
    case ProjectObjectType.deployingModels:
    case ProjectObjectType.deployedModelsList:
    case ProjectObjectType.modelCustomization:
    case ProjectObjectType.labTuning:
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
    case ProjectObjectType.modelCatalog:
      return 'var(--ai-pipeline--Color)';
    case ProjectObjectType.pipelineSetup:
      return 'var(--ai-set-up--Color)';
    case ProjectObjectType.clusterStorage:
      return 'var(--ai-cluster-storage--Color)';
    case ProjectObjectType.modelServer:
    case ProjectObjectType.registeredModels:
    case ProjectObjectType.modelRegistryContext:
    case ProjectObjectType.deployedModels:
    case ProjectObjectType.deployingModels:
    case ProjectObjectType.deployedModelsList:
    case ProjectObjectType.modelCustomization:
    case ProjectObjectType.labTuning:
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
    case ProjectObjectType.modelRegistryContext:
      return modelRegistrySelectImg;
    case ProjectObjectType.deployedModels:
      return deployedModelsImg;
    case ProjectObjectType.deployingModels:
      return deployingModelsImg;
    case ProjectObjectType.deployedModelsList:
      return pipelineRunImg;
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
