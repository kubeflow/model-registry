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
  hardwareConfiguration: string; // e.g., "1 x A100-80"
  totalHardware: number;
  rpsPerReplica: number;
  totalRps: number;
  // TTFT latency metrics
  ttftLatencyMean: number;
  ttftLatencyP90: number;
  ttftLatencyP95: number;
  ttftLatencyP99: number;
  // E2E latency metrics
  e2eLatencyMean: number;
  e2eLatencyP90: number;
  e2eLatencyP95: number;
  e2eLatencyP99: number;
  // TPS latency metrics
  tpsLatencyMean: number;
  tpsLatencyP90: number;
  tpsLatencyP95: number;
  tpsLatencyP99: number;
  // ITL latency metrics
  itlLatencyMean: number;
  itlLatencyP90: number;
  itlLatencyP95: number;
  itlLatencyP99: number;
  // Token metrics
  maxInputTokens: number;
  maxOutputTokens: number;
  meanInputTokens: number;
  meanOutputTokens: number;
  // Version information
  vllmVersion: string;
  guideLLMVersion: string;
  rhaiisVersion: string;
};

export type HardwareConfigurationTableProps = {
  configurations: HardwareConfiguration[];
  isLoading?: boolean;
};
