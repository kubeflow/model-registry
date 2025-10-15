import { WorkloadTypeOptionValue } from '~/concepts/modelCatalog/const';

export type WorkloadTypeOption = {
  value: WorkloadTypeOptionValue;
  label: string;
  description: string;
  maxInputTokens: number;
  maxOutputTokens: number;
};

export const WORKLOAD_TYPE_OPTIONS: WorkloadTypeOption[] = [
  {
    value: WorkloadTypeOptionValue.CHAT,
    label: 'Chat (512 input | 256 output tokens)',
    description: 'Conversational AI workload with moderate input/output token lengths',
    maxInputTokens: 512,
    maxOutputTokens: 256,
  },
  {
    value: WorkloadTypeOptionValue.RAG,
    label: 'RAG (4096 input | 512 output tokens)',
    description: 'Retrieval-Augmented Generation with larger context windows',
    maxInputTokens: 4096,
    maxOutputTokens: 512,
  },
  {
    value: WorkloadTypeOptionValue.SUMMARIZATION,
    label: 'Summarization (2048 input | 256 output tokens)',
    description: 'Text summarization tasks with long input documents',
    maxInputTokens: 2048,
    maxOutputTokens: 256,
  },
  {
    value: WorkloadTypeOptionValue.CODE_GENERATION,
    label: 'Code Generation (1024 input | 512 output tokens)',
    description: 'Code generation and completion tasks',
    maxInputTokens: 1024,
    maxOutputTokens: 512,
  },
];

/**
 * Utility function to convert max input/output tokens to workload type
 */
export const maxInputOutputTokensToWorkloadType = (
  maxInputTokens: number | undefined,
  maxOutputTokens: number | undefined,
): WorkloadTypeOptionValue | undefined => {
  if (maxInputTokens === undefined || maxOutputTokens === undefined) {
    return undefined;
  }

  const matchingOption = WORKLOAD_TYPE_OPTIONS.find(
    (option) =>
      option.maxInputTokens === maxInputTokens && option.maxOutputTokens === maxOutputTokens,
  );

  return matchingOption?.value;
};

/**
 * Utility function to convert workload type to max input/output tokens
 */
export const workloadTypeToMaxInputOutputTokens = (
  workloadType: WorkloadTypeOptionValue,
): { maxInputTokens: number; maxOutputTokens: number } | undefined => {
  const matchingOption = WORKLOAD_TYPE_OPTIONS.find((option) => option.value === workloadType);

  return matchingOption
    ? {
        maxInputTokens: matchingOption.maxInputTokens,
        maxOutputTokens: matchingOption.maxOutputTokens,
      }
    : undefined;
};
