import { UseCaseOptionValue } from '~/concepts/modelCatalog/const';

export type UseCaseOption = {
  value: UseCaseOptionValue;
  label: string;
  description: string;
};

export const USE_CASE_OPTIONS: UseCaseOption[] = [
  {
    value: UseCaseOptionValue.CHATBOT,
    label: 'Chatbot',
    description: 'Conversational AI applications and interactive chat systems',
  },
  {
    value: UseCaseOptionValue.CODE_FIXING,
    label: 'Code Fixing',
    description: 'Code analysis, debugging, and automated code correction',
  },
  {
    value: UseCaseOptionValue.LONG_RAG,
    label: 'Long RAG',
    description: 'Retrieval-Augmented Generation with extended context windows',
  },
  {
    value: UseCaseOptionValue.RAG,
    label: 'RAG',
    description: 'Retrieval-Augmented Generation with standard context windows',
  },
];

/**
 * Utility function to get use case option by value
 */
export const getUseCaseOption = (useCase: UseCaseOptionValue): UseCaseOption | undefined =>
  USE_CASE_OPTIONS.find((option) => option.value === useCase);

/**
 * Mapping from UseCaseOptionValue to display name for use in filters
 */
export const USE_CASE_NAME_MAPPING: Record<UseCaseOptionValue, string> = {
  [UseCaseOptionValue.CHATBOT]: 'Chatbot',
  [UseCaseOptionValue.CODE_FIXING]: 'Code Fixing',
  [UseCaseOptionValue.LONG_RAG]: 'Long RAG',
  [UseCaseOptionValue.RAG]: 'RAG',
};
