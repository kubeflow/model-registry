import { UseCaseOptionValue } from '~/concepts/modelCatalog/const';

export type UseCaseOption = {
  value: UseCaseOptionValue;
  label: string;
  description: string;
  inputTokens: number;
  outputTokens: number;
};

export const USE_CASE_OPTIONS: UseCaseOption[] = [
  {
    value: UseCaseOptionValue.CHATBOT,
    label: 'Chatbot',
    description: 'Conversational AI applications and interactive chat systems',
    inputTokens: 512,
    outputTokens: 256,
  },
  {
    value: UseCaseOptionValue.CODE_FIXING,
    label: 'Code Fixing',
    description: 'Code analysis, debugging, and automated code correction',
    inputTokens: 1024,
    outputTokens: 1024,
  },
  {
    value: UseCaseOptionValue.LONG_RAG,
    label: 'Long RAG',
    description: 'Retrieval-Augmented Generation with extended context windows',
    inputTokens: 10240,
    outputTokens: 1536,
  },
  {
    value: UseCaseOptionValue.RAG,
    label: 'RAG',
    description: 'Retrieval-Augmented Generation with standard context windows',
    inputTokens: 4096,
    outputTokens: 512,
  },
];

/**
 * Utility function to get use case option by value
 */
export const getUseCaseOption = (useCase: UseCaseOptionValue | string): UseCaseOption | undefined =>
  USE_CASE_OPTIONS.find((option) => option.value === useCase);

/**
 * Type guard to check if a string is a valid UseCaseOptionValue
 */
export const isUseCaseOptionValue = (value: string): value is UseCaseOptionValue =>
  USE_CASE_OPTIONS.some((option) => option.value === value);

/**
 * Get display label for a use case value including token information.
 * Format: "Label (inputTokens input | outputTokens output tokens)"
 */
export const getUseCaseDisplayLabel = (value: string): string => {
  const option = getUseCaseOption(value);
  if (option) {
    return `${option.label} (${option.inputTokens} input | ${option.outputTokens} output tokens)`;
  }
  return value;
};

/**
 * Mapping from UseCaseOptionValue to display name for use in filters
 */
export const USE_CASE_NAME_MAPPING: Record<UseCaseOptionValue, string> = {
  [UseCaseOptionValue.CHATBOT]: 'Chatbot',
  [UseCaseOptionValue.CODE_FIXING]: 'Code Fixing',
  [UseCaseOptionValue.LONG_RAG]: 'Long RAG',
  [UseCaseOptionValue.RAG]: 'RAG',
};
