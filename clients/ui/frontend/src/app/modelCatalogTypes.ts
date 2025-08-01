export type ModelCatalogSource = {
  name: string;
  displayName: string;
  description?: string;
  provider?: string;
  url?: string;
  models?: ModelCatalogItem[];
};

export type ModelCatalogItem = {
  id: string;
  name: string;
  displayName: string;
  description?: string;
  provider?: string;
  url?: string;
  logo?: string;
  tags?: string[];
  framework?: string;
  task?: string;
  license?: string;
  metrics?: {
    [key: string]: string | number;
  };
  createdAt?: string;
  updatedAt?: string;
};

export type ModelCatalogContextType = {
  sources: ModelCatalogSource[];
  loading: boolean;
  error?: Error;
  refreshSources: () => Promise<void>;
};
