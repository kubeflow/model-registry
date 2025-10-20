import * as React from 'react';
import { LatencyMetric, LatencyPercentile } from '~/concepts/modelCatalog/const';

export type LatencyFilterConfig = {
  metric: LatencyMetric;
  percentile: LatencyPercentile;
};

// Shared state for latency filter configuration
let sharedLatencyConfig: LatencyFilterConfig = {
  metric: LatencyMetric.TTFT,
  percentile: LatencyPercentile.Mean,
};

export const getLatencyFilterConfig = (): LatencyFilterConfig => sharedLatencyConfig;

export const setLatencyFilterConfig = (config: LatencyFilterConfig): void => {
  sharedLatencyConfig = config;
};

// Hook to manage latency filter configuration
export const useLatencyFilterConfig = (): {
  config: LatencyFilterConfig;
  updateConfig: (newConfig: LatencyFilterConfig) => void;
} => {
  const [config, setConfig] = React.useState<LatencyFilterConfig>(sharedLatencyConfig);

  const updateConfig = React.useCallback((newConfig: LatencyFilterConfig) => {
    setConfig(newConfig);
    setLatencyFilterConfig(newConfig);
  }, []);

  return { config, updateConfig };
};
