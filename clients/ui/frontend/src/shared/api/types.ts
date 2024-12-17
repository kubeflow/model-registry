export type APIOptions = {
  dryRun?: boolean;
  signal?: AbortSignal;
  parseJSON?: boolean;
  headers?: Record<string, string>;
};

export type APIError = {
  error: {
    code: string;
    message: string;
  };
};

export type APIState<T> = {
  /** If API will successfully call */
  apiAvailable: boolean;
  /** The available API functions */
  api: T;
};
