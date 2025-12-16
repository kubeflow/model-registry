export const FORM_LABELS = {
  NAME: 'Name',
  SOURCE_TYPE: 'Source type',
  ORGANIZATION: 'Organization',
  ACCESS_TOKEN: 'Access token',
  YAML_CONTENT: 'Upload a YAML file',
  MODEL_VISIBILITY: 'Model visibility',
  ALLOWED_MODELS: 'Included models',
  EXCLUDED_MODELS: 'Excluded models',
  ENABLE_SOURCE: 'Enable source',
  CREDENTIALS: 'Credentials',
} as const;

export const BUTTON_LABELS = {
  ADD: 'Add',
  SAVE: 'Save',
  PREVIEW: 'Preview',
  CANCEL: 'Cancel',
} as const;

export const SOURCE_TYPE_LABELS = {
  HUGGING_FACE: 'Hugging Face repository',
  YAML: 'YAML file',
} as const;

export const VALIDATION_MESSAGES = {
  NAME_REQUIRED: 'Name is required',
  ORGANIZATION_REQUIRED: 'Organization is required',
  YAML_CONTENT_REQUIRED: 'YAML content is required',
} as const;

export const HELP_TEXT = {
  ACCESS_TOKEN:
    'Enter your fine-grained Hugging Face access token. Public models can be pulled into catalog without an access token. For private/gated models, a token is recommended to ensure full metadata is displayed, otherwise only limited metadata may be available. The token must have the following permissions: read repos in your namespace, read public repos that you can access.',
  ORGANIZATION:
    'Limiting each Hugging Face source to a single organization helps prevent performance issues when loading large model sets.',
  YAML: 'Upload or paste a YAML string.',
} as const;

export const PLACEHOLDERS = {
  ORGANIZATION: 'Example: Google',
  ALLOWED_MODELS_HF: 'Enter model names, one per line (e.g., gemma-7b*)',
  ALLOWED_MODELS_GENERIC: 'Enter model names, one per line',
  EXCLUDED_MODELS_HF: 'Enter model names, one per line (e.g., gemma-7b-test*)',
  EXCLUDED_MODELS_GENERIC: 'Enter model names, one per line',
} as const;

export const DESCRIPTIONS = {
  ENABLE_SOURCE:
    'Enable users in your organization to view models from this source in the model catalog.',
  FILTER_INFO_GENERIC:
    'Optionally filter which models from your source appear in the model catalog. If no filters are set, all models from the source will be visible.',
} as const;

export const PAGE_TITLES = {
  MODEL_CATALOG_PREVIEW: 'Model catalog preview',
  PREVIEW_MODELS: 'Preview models',
} as const;

export const getFilterInfoWithOrg = (organization: string): string =>
  `Optionally filter which ${organization} models from your source appear in the model catalog. If no filters are set, all ${organization} models from the source will be visible.`;

export const getAllowedModelsHelp = (organization?: string): string =>
  organization
    ? `Enter the names of ${organization} models to include from this source. These models will appear in the model catalog.`
    : 'Enter the names of models to include from this source. These models will appear in the model catalog.';

export const getExcludedModelsHelp = (organization?: string): string =>
  organization
    ? `Enter the names of ${organization} models to exclude from this source. These models will not appear in the model catalog.`
    : 'Enter the names of models to exclude from this source. These models will not appear in the model catalog.';

export const FIELD_HELPER_TEXT = {
  INCLUDED_MODELS:
    'Separate model names using commas. To include all models with a specific prefix, enter the prefix followed by an asterisk. Example, Llama*',
  EXCLUDED_MODELS:
    'Separate model names using commas. To exclude all models with a specific prefix, enter the prefix followed by an asterisk. Example, Llama*',
} as const;
