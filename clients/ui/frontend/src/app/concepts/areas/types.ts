import { EitherOrBoth } from '~/typeHelpers';
import {
  DashboardCommonConfig,
  DataScienceClusterInitializationKindStatus,
  DataScienceClusterKindStatus,
  DashboardConfigKind,
} from '~/app/k8sTypes';

export type FeatureFlag = keyof DashboardCommonConfig;

export type IsAreaAvailableStatus = {
  /** A single boolean status */
  status: boolean;
  /* Each status portion broken down -- null if no check made */
  devFlags: { [key in string]?: 'on' | 'off' } | null; // simplified. `disableX` flags are weird to read
  featureFlags: { [key in FeatureFlag]?: 'on' | 'off' } | null; // simplified. `disableX` flags are weird to read
  reliantAreas: { [key in SupportedAreaType]?: boolean } | null; // only needs 1 to be true
  requiredComponents: { [key in StackComponent]?: boolean } | null;
  requiredCapabilities: { [key in StackCapability]?: boolean } | null;
  customCondition: (conditionFunc: CustomConditionFunction) => boolean;
};

/** All areas that we need to support in some fashion or another */
export enum SupportedArea {
  HOME = 'home',

  /* Standalone areas */
  // TODO: Jupyter Tile Support? (outside of feature flags today)
  WORKBENCHES = 'workbenches',
  // TODO: Support Applications/Tile area
  // TODO: Support resources area

  /* Pipelines areas */
  DS_PIPELINES = 'ds-pipelines',

  /* Admin areas */
  BYON = 'bring-your-own-notebook',
  CLUSTER_SETTINGS = 'cluster-settings',
  USER_MANAGEMENT = 'user-management',
  ACCELERATOR_PROFILES = 'accelerator-profiles',
  HARDWARE_PROFILES = 'hardware-profiles',
  STORAGE_CLASSES = 'storage-classes',
  ADMIN_CONNECTION_TYPES = 'connection-types',
  FINE_TUNING = 'fine-tuning',

  /* DS Projects specific areas */
  DS_PROJECTS_PERMISSIONS = 'ds-projects-permission',
  DS_PROJECTS_VIEW = 'ds-projects',
  DS_PROJECT_SCOPED = 'ds-project-scoped',

  /* Model Serving areas */
  MODEL_SERVING = 'model-serving-shell',
  CUSTOM_RUNTIMES = 'custom-serving-runtimes',
  K_SERVE = 'kserve',
  K_SERVE_AUTH = 'kserve-auth',
  K_SERVE_METRICS = 'kserve-metrics',
  K_SERVE_RAW = 'kserve-raw',
  MODEL_MESH = 'model-mesh',
  BIAS_METRICS = 'bias-metrics',
  PERFORMANCE_METRICS = 'performance-metrics',
  TRUSTY_AI = 'trusty-ai',
  NIM_MODEL = 'nim-model',
  SERVING_RUNTIME_PARAMS = 'serving-runtime-params',

  /* Distributed Workloads areas */
  DISTRIBUTED_WORKLOADS = 'distributed-workloads',

  /* Model Registry areas */
  MODEL_REGISTRY = 'model-registry',
  MODEL_REGISTRY_SECURE_DB = 'model-registry-secure-db',

  /* Model catalog areas */
  MODEL_CATALOG = 'model-catalog',

  /* Plugins */
  PLUGIN_MODEL_SERVING = 'plugin-model-serving',

  /* RAG & Agentic */
  LLAMA_STACK_CHAT_BOT = 'llama-stack-chat-bot',

  /* LM Eval */
  LM_EVAL = 'lm-eval',
}

export type SupportedAreaType = SupportedArea | string;
/** Components deployed by the Operator. Part of the DSC Status. */
export enum StackComponent {
  CODE_FLARE = 'codeflare',
  DS_PIPELINES = 'data-science-pipelines-operator',
  K_SERVE = 'kserve',
  MODEL_MESH = 'model-mesh',
  // Bug: https://github.com/opendatahub-io/opendatahub-operator/issues/641
  DASHBOARD = 'odh-dashboard',
  RAY = 'ray',
  WORKBENCHES = 'workbenches',
  TRUSTY_AI = 'trustyai',
  KUEUE = 'kueue',
  MODEL_REGISTRY = 'model-registry-operator',
}

/** The possible component names that are used as keys in the `components` object of the DSC Status.
 * Each component's key (e.g., 'codeflare', 'dashboard', etc.) maps to a specific component status.
 **/
export enum DataScienceStackComponent {
  CODE_FLARE = 'codeflare',
  DASHBOARD = 'dashboard',
  DS_PIPELINES = 'datasciencepipelines',
  K_SERVE = 'kserve',
  KUEUE = 'kueue',
  MODEL_MESH_SERVING = 'modelmeshserving',
  MODEL_REGISTRY = 'modelregistry',
  FEAST_OPERATOR = 'feastoperator',
  RAY = 'ray',
  TRAINING_OPERATOR = 'trainingoperator',
  TRUSTY_AI = 'trustyai',
  WORKBENCHES = 'workbenches',
}

/** Capabilities of the Operator. Part of the DSCI Status. */
export enum StackCapability {
  SERVICE_MESH = 'CapabilityServiceMesh',
  SERVICE_MESH_AUTHZ = 'CapabilityServiceMeshAuthorization',
}

/**
 * Optional function to check for a condition that is not covered by other checks.
 *
 * Example, checking there exists a specific condition in the DSC status.
 *
 * @param state.dashboardConfigSpec The dashboard config spec
 * @param state.dscStatus The data science cluster status
 * @param state.dsciStatus The data science cluster initialization status
 * @returns True if the condition is met, false otherwise
 */
export type CustomConditionFunction = (state: {
  dashboardConfigSpec: DashboardConfigKind['spec'];
  dscStatus: DataScienceClusterKindStatus | null;
  dsciStatus: DataScienceClusterInitializationKindStatus | null;
}) => boolean;

// TODO: Support extra operators, like the pipelines operator -- maybe as a "external dependency need?"
export type SupportedComponentFlagValue = {
  /**
   * An area can be reliant on another area being enabled. The list is "OR"-ed together.
   *
   * Example, Model Serving is a shell for either KServe or ModelMesh. It has no value on its own.
   * It can also be a chain of reliance... example, Custom Runtimes is a Model Serving feature.
   *
   * TODO: support AND -- maybe double array?
   */
  reliantAreas?: SupportedAreaType[];
  /**
   * Required capabilities supported by the Operator. The list is "AND"-ed together.
   * If the Operator does not support the capability, the area is not available.
   * The capabilities are retrieved from the DSCI status.
   */
  requiredCapabilities?: StackCapability[];
} & EitherOrBoth<
  {
    /**
     * Flags that are only available to developers.
     */
    devFlags?: string[];
  },
  EitherOrBoth<
    {
      /**
       * Refers to OdhDashboardConfig's feature flags, any number of them to be "enabled", the result
       * is AND-ed. Omit to not be related to any feature flag.
       *
       * Note: "disable<FlagName>" methodology is confusing and needs to be removed
       * Note: "Enabled" will mean "disable<FlagName>" is false
       * @see https://github.com/opendatahub-io/odh-dashboard/issues/1108
       */
      featureFlags: FeatureFlag[];
    },
    {
      /**
       * Refers to the related stack component names. If a backend component is not installed, this
       * can prevent the feature flag from enabling the item. Omit to not be reliant on a backend
       * component.
       */
      requiredComponents: StackComponent[];
    }
  >
>;

/**
 * Relationships between areas and the state of the cluster.
 */
export type SupportedAreasState = {
  [key in SupportedAreaType]: SupportedComponentFlagValue;
}; 