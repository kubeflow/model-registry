import * as React from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { ModelRegistrySelectorContext } from '~/app/context/ModelRegistrySelectorContext';
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
  const { modelRegistries } = React.useContext(ModelRegistrySelectorContext);
  const selection = modelRegistries.find((mr) => mr.name === (currentModelRegistry ?? ''));
  // When parent passes hasError (e.g. CoreLoader for unavailable page), use it; otherwise derive from selection.
  const hasError =
    modelRegistrySelectorProps.hasError ?? (!!selection && selection.isAvailable === false);

  return (
    <ModelRegistrySelector
      {...modelRegistrySelectorProps}
      hasError={hasError}
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
