import { RegisteredModelList } from '~/app/types';
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
      labels: [
        'Financial data',
        'Fraud detection',
        'Test label',
        'Machine learning',
        'Next data to be overflow',
      ],
    }),
    mockRegisteredModel({
      name: 'Credit Scoring',
      labels: [
        'Credit Score Predictor',
        'Creditworthiness scoring system',
        'Default Risk Analyzer',
        'Portfolio Management',
        'Risk Assessment',
      ],
    }),
    mockRegisteredModel({
      name: 'Label modal',
      description:
        'A machine learning model trained to detect fraudulent transactions in financial data',
      labels: [
        'Testing label',
        'Financial data',
        'Fraud detection',
        'Long label data to be truncated abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc abc',
        'Machine learning',
        'Next data to be overflow',
        'Label x',
        'Label y',
        'Label z',
      ],
    }),
  ],
}: Partial<RegisteredModelList>): RegisteredModelList => ({
  items,
  nextPageToken: '',
  pageSize: 0,
  size,
});
