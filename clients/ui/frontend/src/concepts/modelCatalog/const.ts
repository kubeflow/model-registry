export enum ModelCatalogStringFilterKey {
  TASK = 'tasks',
  PROVIDER = 'provider',
  LICENSE = 'license',
  LANGUAGE = 'language',
  HARDWARE_TYPE = 'hardware_type',
}

export enum ModelCatalogNumberFilterKey {
  MIN_RPS = 'rps_mean',
  MAX_LATENCY = 'ttft_mean',
  MAX_INPUT_TOKENS = 'max_input_tokens',
  MAX_OUTPUT_TOKENS = 'max_output_tokens',
}

// Separate enum for all latency metrics used in validation and filtering
export enum LatencyMetricKey {
  // TTFT (Time To First Token)
  TTFT_MEAN = 'ttft_mean',
  TTFT_P90 = 'ttft_p90',
  TTFT_P95 = 'ttft_p95',
  TTFT_P99 = 'ttft_p99',
  // E2E (End-to-End)
  E2E_MEAN = 'e2e_mean',
  E2E_P90 = 'e2e_p90',
  E2E_P95 = 'e2e_p95',
  E2E_P99 = 'e2e_p99',
  // TPS (Tokens Per Second)
  TPS_MEAN = 'tps_mean',
  TPS_P90 = 'tps_p90',
  TPS_P95 = 'tps_p95',
  TPS_P99 = 'tps_p99',
  // ITL (Inter-Token Latency)
  ITL_MEAN = 'itl_mean',
  ITL_P90 = 'itl_p90',
  ITL_P95 = 'itl_p95',
  ITL_P99 = 'itl_p99',
}

export enum WorkloadTypeOptionValue {
  CHAT = 'chat',
  RAG = 'rag',
  SUMMARIZATION = 'summarization',
  CODE_GENERATION = 'code_generation',
}

export enum ModelCatalogTask {
  AUDIO_TO_TEXT = 'audio-to-text',
  IMAGE_TEXT_TO_TEXT = 'image-text-to-text',
  IMAGE_TO_TEXT = 'image-to-text',
  TEXT_GENERATION = 'text-generation',
  TEXT_TO_TEXT = 'text-to-text',
  VIDEO_TO_TEXT = 'video-to-text',
}

export const MODEL_CATALOG_TASK_NAME_MAPPING = {
  [ModelCatalogTask.AUDIO_TO_TEXT]: 'Audio-to-text',
  [ModelCatalogTask.IMAGE_TEXT_TO_TEXT]: 'Image-text-to-text',
  [ModelCatalogTask.IMAGE_TO_TEXT]: 'Image-to-text',
  [ModelCatalogTask.TEXT_GENERATION]: 'Text generation',
  [ModelCatalogTask.TEXT_TO_TEXT]: 'Text-to-text',
  [ModelCatalogTask.VIDEO_TO_TEXT]: 'Video-to-text',
};

export const MODEL_CATALOG_TASK_DESCRIPTION = {
  [ModelCatalogTask.AUDIO_TO_TEXT]: 'Audio transcription and speech recognition models',
  [ModelCatalogTask.IMAGE_TEXT_TO_TEXT]: 'Multimodal models that process both images and text',
  [ModelCatalogTask.IMAGE_TO_TEXT]: 'Image captioning and visual understanding models',
  [ModelCatalogTask.TEXT_GENERATION]: 'Large language models for text generation',
  [ModelCatalogTask.TEXT_TO_TEXT]: 'Text transformation and translation models',
  [ModelCatalogTask.VIDEO_TO_TEXT]: 'Video analysis and description models',
};

export enum ModelCatalogProvider {
  ALIBABA_CLOUD = 'Alibaba Cloud',
  DEEPSEEK = 'DeepSeek',
  GOOGLE = 'Google',
  IBM = 'IBM',
  META = 'Meta',
  MISTRAL_AI = 'Mistral AI',
  MOONSHOT_AI = 'Moonshot AI',
  NEURAL_MAGIC = 'Neural Magic',
  NVIDIA = 'NVIDIA',
  NVIDIA_ALTERNATE = 'Nvidia', // alternate casing
  RED_HAT = 'Red Hat',
}

export const MODEL_CATALOG_PROVIDER_NAME_MAPPING = {
  [ModelCatalogProvider.ALIBABA_CLOUD]: 'Alibaba Cloud',
  [ModelCatalogProvider.DEEPSEEK]: 'DeepSeek',
  [ModelCatalogProvider.GOOGLE]: 'Google',
  [ModelCatalogProvider.IBM]: 'IBM',
  [ModelCatalogProvider.META]: 'Meta',
  [ModelCatalogProvider.MISTRAL_AI]: 'Mistral AI',
  [ModelCatalogProvider.MOONSHOT_AI]: 'Moonshot AI',
  [ModelCatalogProvider.NEURAL_MAGIC]: 'Neural Magic',
  [ModelCatalogProvider.NVIDIA]: 'NVIDIA',
  [ModelCatalogProvider.NVIDIA_ALTERNATE]: 'NVIDIA',
  [ModelCatalogProvider.RED_HAT]: 'Red Hat',
};

export const MODEL_CATALOG_PROVIDER_NOTABLE_MODELS = {
  [ModelCatalogProvider.ALIBABA_CLOUD]: 'Qwen series models',
  [ModelCatalogProvider.DEEPSEEK]: 'DeepSeek reasoning models',
  [ModelCatalogProvider.GOOGLE]: 'Gemma series models',
  [ModelCatalogProvider.IBM]: 'Granite series models',
  [ModelCatalogProvider.META]: 'Llama series models',
  [ModelCatalogProvider.MISTRAL_AI]: 'Mistral series models',
  [ModelCatalogProvider.MOONSHOT_AI]: 'Kimi series models',
  [ModelCatalogProvider.NEURAL_MAGIC]: 'Quantized model variants',
  [ModelCatalogProvider.NVIDIA]: 'NVIDIA research models',
  [ModelCatalogProvider.NVIDIA_ALTERNATE]: 'NVIDIA research models',
  [ModelCatalogProvider.RED_HAT]: 'Red Hat optimized models',
};

export enum ModelCatalogLicense {
  APACHE_2_0 = 'apache-2.0',
  GEMMA = 'gemma',
  LLLAMA_3_3 = 'llama-3.3',
  LLLAMA_3_1 = 'llama3.1',
  LLLAMA_3_3_ALTERNATE = 'llama3.3',
  LLLAMA_4 = 'llama4',
  MIT = 'mit',
  MODIFIED_MIT = 'modified-mit',
}

export const MODEL_CATALOG_LICENSE_NAME_MAPPING = {
  [ModelCatalogLicense.APACHE_2_0]: 'Apache 2.0',
  [ModelCatalogLicense.GEMMA]: 'Gemma',
  [ModelCatalogLicense.LLLAMA_3_3]: 'Llama 3.3',
  [ModelCatalogLicense.LLLAMA_3_1]: 'Llama 3.1',
  [ModelCatalogLicense.LLLAMA_3_3_ALTERNATE]: 'Llama 3.3 (variant)',
  [ModelCatalogLicense.LLLAMA_4]: 'Llama 4',
  [ModelCatalogLicense.MIT]: 'MIT',
  [ModelCatalogLicense.MODIFIED_MIT]: 'Modified MIT',
};

export const MODEL_CATALOG_LICENSE_DETAILS = {
  [ModelCatalogLicense.APACHE_2_0]: {
    name: 'Apache 2.0',
    type: 'Open Source',
    description: 'Permissive Apache License 2.0',
  },
  [ModelCatalogLicense.GEMMA]: {
    name: 'Gemma',
    type: 'Custom',
    description: 'Google Gemma model license',
  },
  [ModelCatalogLicense.LLLAMA_3_3]: {
    name: 'Llama 3.3',
    type: 'Custom',
    description: 'Meta Llama 3.3 license',
  },
  [ModelCatalogLicense.LLLAMA_3_1]: {
    name: 'Llama 3.1',
    type: 'Custom',
    description: 'Meta Llama 3.1 license',
  },
  [ModelCatalogLicense.LLLAMA_3_3_ALTERNATE]: {
    name: 'Llama 3.3 (variant)',
    type: 'Custom',
    description: 'Meta Llama 3.3 license (variant)',
  },
  [ModelCatalogLicense.LLLAMA_4]: {
    name: 'Llama 4',
    type: 'Custom',
    description: 'Meta Llama 4 license',
  },
  [ModelCatalogLicense.MIT]: {
    name: 'MIT',
    type: 'Open Source',
    description: 'Permissive MIT license',
  },
  [ModelCatalogLicense.MODIFIED_MIT]: {
    name: 'Modified MIT',
    type: 'Open Source',
    description: 'Modified MIT license',
  },
};

export enum EuropeanLanguagesCode {
  BG = 'bg',
  CA = 'ca',
  CS = 'cs',
  DA = 'da',
  DE = 'de',
  EL = 'el',
  EN = 'en',
  ES = 'es',
  FI = 'fi',
  FR = 'fr',
  HR = 'hr',
  HU = 'hu',
  IS = 'is',
  IT = 'it',
  NL = 'nl',
  NLD = 'nld',
  NO = 'no',
  PL = 'pl',
  PT = 'pt',
  RO = 'ro',
  RU = 'ru',
  SK = 'sk',
  SL = 'sl',
  SR = 'sr',
  SV = 'sv',
  UK = 'uk',
}

export const MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS = {
  [EuropeanLanguagesCode.BG]: 'Bulgarian',
  [EuropeanLanguagesCode.CA]: 'Catalan',
  [EuropeanLanguagesCode.CS]: 'Czech',
  [EuropeanLanguagesCode.DA]: 'Danish',
  [EuropeanLanguagesCode.DE]: 'German',
  [EuropeanLanguagesCode.EL]: 'Greek',
  [EuropeanLanguagesCode.EN]: 'English',
  [EuropeanLanguagesCode.ES]: 'Spanish',
  [EuropeanLanguagesCode.FI]: 'Finnish',
  [EuropeanLanguagesCode.FR]: 'French',
  [EuropeanLanguagesCode.HR]: 'Croatian',
  [EuropeanLanguagesCode.HU]: 'Hungarian',
  [EuropeanLanguagesCode.IS]: 'Icelandic',
  [EuropeanLanguagesCode.IT]: 'Italian',
  [EuropeanLanguagesCode.NL]: 'Dutch',
  [EuropeanLanguagesCode.NLD]: 'Dutch - variant',
  [EuropeanLanguagesCode.NO]: 'Norwegian',
  [EuropeanLanguagesCode.PL]: 'Polish',
  [EuropeanLanguagesCode.PT]: 'Portuguese',
  [EuropeanLanguagesCode.RO]: 'Romanian',
  [EuropeanLanguagesCode.RU]: 'Russian',
  [EuropeanLanguagesCode.SK]: 'Slovak',
  [EuropeanLanguagesCode.SL]: 'Slovenian',
  [EuropeanLanguagesCode.SR]: 'Serbian',
  [EuropeanLanguagesCode.SV]: 'Swedish',
  [EuropeanLanguagesCode.UK]: 'Ukrainian',
};

export enum AsianLanguagesCode {
  JA = 'ja',
  KO = 'ko',
  ZH = 'zh',
  HI = 'hi',
  TH = 'th',
  VI = 'vi',
  ID = 'id',
  MS = 'ms',
  ZSM = 'zsm',
}

export const MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS = {
  [AsianLanguagesCode.JA]: 'Japanese',
  [AsianLanguagesCode.KO]: 'Korean',
  [AsianLanguagesCode.ZH]: 'Chinese',
  [AsianLanguagesCode.HI]: 'Hindi',
  [AsianLanguagesCode.TH]: 'Thai',
  [AsianLanguagesCode.VI]: 'Vietnamese',
  [AsianLanguagesCode.ID]: 'Indonesian',
  [AsianLanguagesCode.MS]: 'Malay',
  [AsianLanguagesCode.ZSM]: 'Standard Malay',
};

export enum MiddleEasternAndOtherLanguagesCode {
  AR = 'ar',
  FA = 'fa',
  HE = 'he',
  TR = 'tr',
  UR = 'ur',
  TL = 'tl',
}

export const MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS = {
  [MiddleEasternAndOtherLanguagesCode.AR]: 'Arabic',
  [MiddleEasternAndOtherLanguagesCode.FA]: 'Persian',
  [MiddleEasternAndOtherLanguagesCode.HE]: 'Hebrew',
  [MiddleEasternAndOtherLanguagesCode.TR]: 'Turkish',
  [MiddleEasternAndOtherLanguagesCode.UR]: 'Urdu',
  [MiddleEasternAndOtherLanguagesCode.TL]: 'Tagalog',
};

export const AllLanguageCodesMap = {
  ...MODEL_CATALOG_EUROPEAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_ASIAN_LANGUAGES_DETAILS,
  ...MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS,
};

export enum AllLanguageCode {
  BG = 'bg',
  CA = 'ca',
  CS = 'cs',
  DA = 'da',
  DE = 'de',
  EL = 'el',
  EN = 'en',
  ES = 'es',
  FI = 'fi',
  FR = 'fr',
  HR = 'hr',
  HU = 'hu',
  IS = 'is',
  IT = 'it',
  NL = 'nl',
  NLD = 'nld',
  NO = 'no',
  PL = 'pl',
  PT = 'pt',
  RO = 'ro',
  RU = 'ru',
  SK = 'sk',
  SL = 'sl',
  SR = 'sr',
  SV = 'sv',
  UK = 'uk',
  JA = 'ja',
  KO = 'ko',
  ZH = 'zh',
  HI = 'hi',
  TH = 'th',
  VI = 'vi',
  ID = 'id',
  MS = 'ms',
  ZSM = 'zsm',
  AR = 'ar',
  FA = 'fa',
  HE = 'he',
  TR = 'tr',
  UR = 'ur',
  TL = 'tl',
}
