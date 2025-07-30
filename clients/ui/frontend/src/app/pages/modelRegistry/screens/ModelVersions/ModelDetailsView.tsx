import * as React from 'react';
import { Grid, GridItem, Stack } from '@patternfly/react-core';
import { RegisteredModel } from '~/app/types';
import ModelDetailsCard from './ModelDetailsCard';

type ModelDetailsViewProps = {
  registeredModel: RegisteredModel;
  refresh: () => void;
  isArchiveModel?: boolean;
};

const ModelDetailsView: React.FC<ModelDetailsViewProps> = ({
  registeredModel: rm,
  refresh,
  isArchiveModel,
}) => (
  <Grid hasGutter>
    <GridItem span={12} lg={8}>
      <ModelDetailsCard registeredModel={rm} refresh={refresh} isArchiveModel={isArchiveModel} />
    </GridItem>
    <GridItem span={12} lg={4}>
      <Stack hasGutter>
        {/* TODO: Add latest model versions card here */}
        {/* TODO: Add latest deployments card here (as an extension)*/}
      </Stack>
    </GridItem>
  </Grid>
);

export default ModelDetailsView;
