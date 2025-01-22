import * as React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import ModelRegistrySelector from './ModelRegistrySelector';

type ModelRegistrySelectorNavigatorProps = {
  getRedirectPath: (namespace: string) => string;
} & Omit<React.ComponentProps<typeof ModelRegistrySelector>, 'onSelection' | 'modelRegistry'>;

const ModelRegistrySelectorNavigator: React.FC<ModelRegistrySelectorNavigatorProps> = ({
  getRedirectPath,
  ...modelRegistrySelectorProps
}) => {
  const navigate = useNavigate();
  const { modelRegistry: currentModelRegistry } = useParams<{ modelRegistry: string }>();

  return (
    <ModelRegistrySelector
      {...modelRegistrySelectorProps}
      onSelection={(modelRegistryName) => {
        if (modelRegistryName !== currentModelRegistry) {
          navigate(getRedirectPath(modelRegistryName));
        }
      }}
      modelRegistry={currentModelRegistry ?? ''}
    />
  );
};

export default ModelRegistrySelectorNavigator;
