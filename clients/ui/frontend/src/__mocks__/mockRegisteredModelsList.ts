import { ModelRegistryMetadataType, RegisteredModelList } from '~/app/types';
import { mockRegisteredModel } from './mockRegisteredModel';

export const mockRegisteredModelList = ({
  size = 5,
  items = [
    mockRegisteredModel({ name: 'test-1' }),
    mockRegisteredModel({ name: 'test-2' }),
    mockRegisteredModel({
      name: 'Fraud detection model',
      description:
        'A machine learning model trained to detect fraudulent transactions in financial data',
      customProperties: {
        'Financial data': {
          metadataType: ModelRegistryMetadataType.STRING,
          // eslint-disable-next-line camelcase
          string_value: '',
        },
      },
    }),
    mockRegisteredModel({
      name: 'Credit Scoring',
      customProperties: {
        'Credit Score Predictor': {
          metadataType: ModelRegistryMetadataType.STRING,
          // eslint-disable-next-line camelcase
          string_value: '',
        },
      },
    }),
    mockRegisteredModel({
      name: 'Label modal',
      description:
        'A machine learning model trained to detect fraudulent transactions in financial data',
      customProperties: {
        'Testing label': {
          metadataType: ModelRegistryMetadataType.STRING,
          // eslint-disable-next-line camelcase
          string_value: '',
        },
      },
    }),
  ],
}: Partial<RegisteredModelList>): RegisteredModelList => ({
  items,
  nextPageToken: '',
  pageSize: 0,
  size,
});
