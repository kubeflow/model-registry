import * as React from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  IconButton,
  Button,
} from '@mui/material';
import { Edit, Delete } from '@mui/icons-material';
import { ModelRegistryKind, RoleBindingKind } from '~/app/k8sTypes';
import { FetchState } from '~/app/utils/useFetch';

type ModelRegistriesTableProps = {
  modelRegistries: ModelRegistryKind[];
  roleBindings: FetchState<RoleBindingKind[]>;
  refresh: () => void;
  onCreateModelRegistryClick: () => void;
};

const ModelRegistriesTable: React.FC<ModelRegistriesTableProps> = ({
  modelRegistries,
  roleBindings,
  refresh,
  onCreateModelRegistryClick,
}) => (
  <TableContainer component={Paper}>
    <Table>
      <TableHead>
        <TableRow>
          <TableCell>Name</TableCell>
          <TableCell>Owner</TableCell>
          <TableCell>Created</TableCell>
          <TableCell>Actions</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {modelRegistries.map((mr) => (
          <TableRow key={mr.metadata.name}>
            <TableCell>{mr.metadata.annotations?.['openshift.io/display-name'] || mr.metadata.name}</TableCell>
            <TableCell>{mr.metadata.annotations?.['opendatahub.io/username'] || 'Unknown'}</TableCell>
            <TableCell>{new Date(mr.metadata.creationTimestamp || '').toLocaleString()}</TableCell>
            <TableCell>
              <IconButton size="small">
                <Edit />
              </IconButton>
              <IconButton size="small">
                <Delete />
              </IconButton>
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
    <Button onClick={onCreateModelRegistryClick}>Create model registry</Button>
  </TableContainer>
);

export default ModelRegistriesTable; 