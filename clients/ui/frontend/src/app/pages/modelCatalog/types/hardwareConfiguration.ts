export enum HardwareType {
  GPU = 'GPU',
  CPU = 'CPU',
  TPU = 'TPU',
}

export enum Framework {
  PYTORCH = 'PyTorch',
  TENSORFLOW = 'TensorFlow',
  JAX = 'JAX',
  ONNX = 'ONNX',
}

export enum Precision {
  FP16 = 'FP16',
  FP32 = 'FP32',
  BF16 = 'BF16',
  INT8 = 'INT8',
}

export type HardwareConfiguration = {
  id: string;
  hardwareType: HardwareType;
  hardwareCount: number;
  latency: number;
  throughput: number;
  guideLLMVersion?: string;
  rhaiisVersion?: string;
  tps: number;
  // Additional fields from the data structure
  modelName?: string;
  modelVersion?: string;
  framework?: Framework;
  precision?: Precision;
  batchSize?: number;
  sequenceLength?: number;
  memoryUsage?: number;
  powerConsumption?: number;
  temperature?: number;
  utilization?: number;
  // Timestamps
  lastUpdated?: string;
  // Other performance metrics
  accuracy?: number;
  f1Score?: number;
  bleuScore?: number;
  perplexity?: number;
};

export type HardwareConfigurationTableProps = {
  configurations: HardwareConfiguration[];
  isLoading?: boolean;
};
