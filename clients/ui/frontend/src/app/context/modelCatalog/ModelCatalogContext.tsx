import * as React from 'react';
import { ModelCatalogContextType, ModelCatalogSource } from '~/app/modelCatalogTypes';

// Mock data for initial development
const mockModelCatalogSources: ModelCatalogSource[] = [
  {
    name: 'huggingface',
    displayName: 'Hugging Face Hub',
    description: 'Popular models from Hugging Face Hub',
    provider: 'Hugging Face',
    url: 'https://huggingface.co',
    models: [
      {
        id: 'bert-base-uncased',
        name: 'bert-base-uncased',
        displayName: 'BERT Base Uncased',
        description:
          'Pre-trained BERT model on English language using a masked language modeling objective',
        provider: 'Hugging Face',
        url: 'https://huggingface.co/bert-base-uncased',
        tags: ['transformer', 'bert', 'natural-language-processing'],
        framework: 'PyTorch',
        task: 'text-classification',
        license: 'Apache 2.0',
        metrics: {
          accuracy: 0.92,
          f1: 0.89,
        },
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-06-01T00:00:00Z',
      },
      {
        id: 'gpt2',
        name: 'gpt2',
        displayName: 'GPT-2',
        description: 'OpenAI GPT-2 language model',
        provider: 'OpenAI',
        url: 'https://huggingface.co/gpt2',
        tags: ['transformer', 'gpt', 'language-model'],
        framework: 'PyTorch',
        task: 'text-generation',
        license: 'MIT',
        metrics: {
          perplexity: 18.5,
        },
      },
    ],
  },
  {
    name: 'tensorflow',
    displayName: 'TensorFlow Hub',
    description: 'Pre-trained models from TensorFlow Hub',
    provider: 'Google',
    url: 'https://tfhub.dev',
    models: [
      {
        id: 'efficientnet_b0',
        name: 'efficientnet_b0',
        displayName: 'EfficientNet B0',
        description: 'Lightweight convolutional neural network optimized for mobile devices',
        provider: 'Google',
        url: 'https://tfhub.dev/tensorflow/efficientnet/b0/classification/1',
        tags: ['computer-vision', 'classification', 'mobile'],
        framework: 'TensorFlow',
        task: 'image-classification',
        license: 'Apache 2.0',
        metrics: {
          top1Accuracy: 0.774,
          top6Accuracy: 0.934,
        },
      },
    ],
  },
];

export const ModelCatalogContext = React.createContext<ModelCatalogContextType>({
  sources: [],
  loading: false,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  refreshSources: async () => {},
});

export const ModelCatalogContextProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [sources, setSources] = React.useState(mockModelCatalogSources);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<Error>();

  const refreshSources = React.useCallback(async () => {
    setLoading(true);
    try {
      // TODO: Replace with actual API call
      // eslint-disable-next-line no-promise-executor-return
      await new Promise((resolve) => setTimeout(resolve, 1000)); // Simulate API delay
      setSources(mockModelCatalogSources);
      setError(undefined);
    } catch (e) {
      setError(e instanceof Error ? e : new Error('Failed to fetch catalog sources'));
    } finally {
      setLoading(false);
    }
  }, []);

  React.useEffect(() => {
    refreshSources();
  }, [refreshSources]);

  const value = React.useMemo(
    () => ({
      sources,
      loading,
      error,
      refreshSources,
    }),
    [sources, loading, error, refreshSources],
  );

  return <ModelCatalogContext.Provider value={value}>{children}</ModelCatalogContext.Provider>;
};
