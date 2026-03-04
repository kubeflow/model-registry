export type McpServerDeploymentMode = 'Remote' | 'Local';

export type McpServerMock = {
  id: string;
  name: string;
  description: string;
  deploymentMode: McpServerDeploymentMode;
  securityVerification: string[];
  category: 'sample' | 'other';
  icon?: string;
};
