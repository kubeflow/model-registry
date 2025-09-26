export enum ModelCatalogTasks {
  AUDIO_TO_TEXT = 'audio-to-text',
  IMAGE_TEXT_TO_TEXT = 'image-text-to-text',
  IMAGE_TO_TEXT = 'image-to-text',
  TEXT_GENERATION = 'text-generation',
  TEXT_TO_TEXT = 'text-to-text',
  VIDEO_TO_TEXT = 'video-to-text',
}

export const MODEL_CATALOG_TASK_NAME_MAPPING = {
  [ModelCatalogTasks.AUDIO_TO_TEXT]: 'Audio-to-Text',
  [ModelCatalogTasks.IMAGE_TEXT_TO_TEXT]: 'Image-Text-to-Text',
  [ModelCatalogTasks.IMAGE_TO_TEXT]: 'Image-to-Text',
  [ModelCatalogTasks.TEXT_GENERATION]: 'Text Generation',
  [ModelCatalogTasks.TEXT_TO_TEXT]: 'Text-to-Text',
  [ModelCatalogTasks.VIDEO_TO_TEXT]: 'Video-to-Text',
};

export const MODEL_CATALOG_TASK_DESCRIPTION = {
  [ModelCatalogTasks.AUDIO_TO_TEXT]: 'Audio transcription and speech recognition models',
  [ModelCatalogTasks.IMAGE_TEXT_TO_TEXT]: 'Multimodal models that process both images and text',
  [ModelCatalogTasks.IMAGE_TO_TEXT]: 'Image captioning and visual understanding models',
  [ModelCatalogTasks.TEXT_GENERATION]: 'Large language models for text generation',
  [ModelCatalogTasks.TEXT_TO_TEXT]: 'Text transformation and translation models',
  [ModelCatalogTasks.VIDEO_TO_TEXT]: 'Video analysis and description models',
};

export enum ModelCatalogProviders {
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
  [ModelCatalogProviders.ALIBABA_CLOUD]: 'Alibaba Cloud',
  [ModelCatalogProviders.DEEPSEEK]: 'DeepSeek',
  [ModelCatalogProviders.GOOGLE]: 'Google',
  [ModelCatalogProviders.IBM]: 'IBM',
  [ModelCatalogProviders.META]: 'Meta',
  [ModelCatalogProviders.MISTRAL_AI]: 'Mistral AI',
  [ModelCatalogProviders.MOONSHOT_AI]: 'Moonshot AI',
  [ModelCatalogProviders.NEURAL_MAGIC]: 'Neural Magic',
  [ModelCatalogProviders.NVIDIA]: 'NVIDIA',
  [ModelCatalogProviders.NVIDIA_ALTERNATE]: 'NVIDIA',
  [ModelCatalogProviders.RED_HAT]: 'Red Hat',
};

export const MODEL_CATALOG_PROVIDER_NOTABLE_MODELS = {
  [ModelCatalogProviders.ALIBABA_CLOUD]: 'Qwen series models',
  [ModelCatalogProviders.DEEPSEEK]: 'DeepSeek reasoning models',
  [ModelCatalogProviders.GOOGLE]: 'Gemma series models',
  [ModelCatalogProviders.IBM]: 'Granite series models',
  [ModelCatalogProviders.META]: 'Llama series models',
  [ModelCatalogProviders.MISTRAL_AI]: 'Mistral series models',
  [ModelCatalogProviders.MOONSHOT_AI]: 'Kimi series models',
  [ModelCatalogProviders.NEURAL_MAGIC]: 'Quantized model variants',
  [ModelCatalogProviders.NVIDIA]: 'NVIDIA research models',
  [ModelCatalogProviders.NVIDIA_ALTERNATE]: 'NVIDIA research models',
  [ModelCatalogProviders.RED_HAT]: 'Red Hat optimized models',
};

export enum ModelCatalogLicenses {
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
  [ModelCatalogLicenses.APACHE_2_0]: 'Apache 2.0',
  [ModelCatalogLicenses.GEMMA]: 'Gemma',
  [ModelCatalogLicenses.LLLAMA_3_3]: 'Llama 3.3',
  [ModelCatalogLicenses.LLLAMA_3_1]: 'Llama 3.1',
  [ModelCatalogLicenses.LLLAMA_3_3_ALTERNATE]: 'Llama 3.3 (variant)',
  [ModelCatalogLicenses.LLLAMA_4]: 'Llama 4',
  [ModelCatalogLicenses.MIT]: 'MIT',
  [ModelCatalogLicenses.MODIFIED_MIT]: 'Modified MIT',
};

export const MODEL_CATALOG_LICENSE_DETAILS = {
  [ModelCatalogLicenses.APACHE_2_0]: {
    name: 'Apache 2.0',
    type: 'Open Source',
    description: 'Permissive Apache License 2.0',
  },
  [ModelCatalogLicenses.GEMMA]: {
    name: 'Gemma',
    type: 'Custom',
    description: 'Google Gemma model license',
  },
  [ModelCatalogLicenses.LLLAMA_3_3]: {
    name: 'Llama 3.3',
    type: 'Custom',
    description: 'Meta Llama 3.3 license',
  },
  [ModelCatalogLicenses.LLLAMA_3_1]: {
    name: 'Llama 3.1',
    type: 'Custom',
    description: 'Meta Llama 3.1 license',
  },
  [ModelCatalogLicenses.LLLAMA_3_3_ALTERNATE]: {
    name: 'Llama 3.3 (variant)',
    type: 'Custom',
    description: 'Meta Llama 3.3 license (variant)',
  },
  [ModelCatalogLicenses.LLLAMA_4]: {
    name: 'Llama 4',
    type: 'Custom',
    description: 'Meta Llama 4 license',
  },
  [ModelCatalogLicenses.MIT]: {
    name: 'MIT',
    type: 'Open Source',
    description: 'Permissive MIT license',
  },
  [ModelCatalogLicenses.MODIFIED_MIT]: {
    name: 'Modified MIT',
    type: 'Open Source',
    description: 'Modified MIT license',
  },
};

export enum EuropeanLanguagesCodes {
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
  [EuropeanLanguagesCodes.BG]: 'Bulgarian',
  [EuropeanLanguagesCodes.CA]: 'Catalan',
  [EuropeanLanguagesCodes.CS]: 'Czech',
  [EuropeanLanguagesCodes.DA]: 'Danish',
  [EuropeanLanguagesCodes.DE]: 'German',
  [EuropeanLanguagesCodes.EL]: 'Greek',
  [EuropeanLanguagesCodes.EN]: 'English',
  [EuropeanLanguagesCodes.ES]: 'Spanish',
  [EuropeanLanguagesCodes.FI]: 'Finnish',
  [EuropeanLanguagesCodes.FR]: 'French',
  [EuropeanLanguagesCodes.HR]: 'Croatian',
  [EuropeanLanguagesCodes.HU]: 'Hungarian',
  [EuropeanLanguagesCodes.IS]: 'Icelandic',
  [EuropeanLanguagesCodes.IT]: 'Italian',
  [EuropeanLanguagesCodes.NL]: 'Dutch',
  [EuropeanLanguagesCodes.NLD]: 'Dutch - variant',
  [EuropeanLanguagesCodes.NO]: 'Norwegian',
  [EuropeanLanguagesCodes.PL]: 'Polish',
  [EuropeanLanguagesCodes.PT]: 'Portuguese',
  [EuropeanLanguagesCodes.RO]: 'Romanian',
  [EuropeanLanguagesCodes.RU]: 'Russian',
  [EuropeanLanguagesCodes.SK]: 'Slovak',
  [EuropeanLanguagesCodes.SL]: 'Slovenian',
  [EuropeanLanguagesCodes.SR]: 'Serbian',
  [EuropeanLanguagesCodes.SV]: 'Swedish',
  [EuropeanLanguagesCodes.UK]: 'Ukrainian',
};

export enum AsianLanguagesCodes {
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
  [AsianLanguagesCodes.JA]: 'Japanese',
  [AsianLanguagesCodes.KO]: 'Korean',
  [AsianLanguagesCodes.ZH]: 'Chinese',
  [AsianLanguagesCodes.HI]: 'Hindi',
  [AsianLanguagesCodes.TH]: 'Thai',
  [AsianLanguagesCodes.VI]: 'Vietnamese',
  [AsianLanguagesCodes.ID]: 'Indonesian',
  [AsianLanguagesCodes.MS]: 'Malay',
  [AsianLanguagesCodes.ZSM]: 'Standard Malay',
};

export enum MiddleEasternAndOtherLanguagesCodes {
  AR = 'ar',
  FA = 'fa',
  HE = 'he',
  TR = 'tr',
  UR = 'ur',
  TL = 'tl',
}

export const MODEL_CATALOG_MIDDLE_EASTERN_AND_OTHER_LANGUAGES_DETAILS = {
  [MiddleEasternAndOtherLanguagesCodes.AR]: 'Arabic',
  [MiddleEasternAndOtherLanguagesCodes.FA]: 'Persian',
  [MiddleEasternAndOtherLanguagesCodes.HE]: 'Hebrew',
  [MiddleEasternAndOtherLanguagesCodes.TR]: 'Turkish',
  [MiddleEasternAndOtherLanguagesCodes.UR]: 'Urdu',
  [MiddleEasternAndOtherLanguagesCodes.TL]: 'Tagalog',
};
