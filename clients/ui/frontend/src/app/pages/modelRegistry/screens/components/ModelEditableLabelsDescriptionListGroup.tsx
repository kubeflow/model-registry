import React from 'react';
import { EditableLabelsDescriptionListGroup } from 'mod-arch-shared';
import { RegisteredModel } from '~/app/types';
import { ModelRegistryContext } from '~/app/context/ModelRegistryContext';
import { getLabels, mergeUpdatedLabels } from '~/app/pages/modelRegistry/screens/utils';

type ModelEditableLabelsDescriptionListGroupProps = {
  isArchiveModel?: boolean;
  rm: RegisteredModel;
  refresh: () => void;
};

const ModelEditableLabelsDescriptionListGroup: React.FC<
  ModelEditableLabelsDescriptionListGroupProps
> = ({ isArchiveModel, rm, refresh }) => {
  const { apiState } = React.useContext(ModelRegistryContext);
  return (
    <EditableLabelsDescriptionListGroup
      labels={getLabels(rm.customProperties)}
      isArchive={isArchiveModel}
      allExistingKeys={Object.keys(rm.customProperties)}
      title="Labels"
      contentWhenEmpty="No labels"
      onLabelsChange={(editedLabels) =>
        apiState.api
          .patchRegisteredModel(
            {},
            {
              customProperties: mergeUpdatedLabels(rm.customProperties, editedLabels),
            },
            rm.id,
          )
          .then(refresh)
      }
    />
  );
};

export default ModelEditableLabelsDescriptionListGroup;
