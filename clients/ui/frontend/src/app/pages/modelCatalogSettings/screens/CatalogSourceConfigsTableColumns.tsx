import * as React from 'react';
import { List, ListItem } from '@patternfly/react-core';
import { kebabTableColumn, SortableData } from 'mod-arch-shared';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';
import {
  TABLE_COLUMN_LABELS,
  TABLE_COLUMN_POPOVERS,
} from '~/app/pages/modelCatalogSettings/constants';

export const catalogSourceConfigsColumns: SortableData<CatalogSourceConfig>[] = [
  {
    field: 'name',
    label: TABLE_COLUMN_LABELS.SOURCE_NAME,
    sortable: (a, b) => a.name.localeCompare(b.name),
    width: 15,
  },
  {
    field: 'allowedOrganization',
    label: TABLE_COLUMN_LABELS.ORGANIZATION,
    sortable: (a, b) =>
      ('allowedOrganization' in a ? (a.allowedOrganization ?? '') : '').localeCompare(
        'allowedOrganization' in b ? (b.allowedOrganization ?? '') : '',
      ),
    info: {
      popover: TABLE_COLUMN_POPOVERS.ORGANIZATION,
    },
    width: 15,
  },
  {
    field: 'filters',
    label: TABLE_COLUMN_LABELS.MODEL_VISIBILITY,
    sortable: (a: CatalogSourceConfig, b: CatalogSourceConfig): number => {
      const aFiltered = (a.includedModels?.length ?? 0) + (a.excludedModels?.length ?? 0);
      const bFiltered = (b.includedModels?.length ?? 0) + (b.excludedModels?.length ?? 0);
      return aFiltered - bFiltered;
    },
    info: {
      popover: (
        <div>
          <p>
            Shows whether all models from a source appear in the model catalog or if visibility is
            filtered.
          </p>
          <List>
            <ListItem>
              <strong>All models:</strong> Every model from the source appears in the catalog.
            </ListItem>
            <ListItem>
              <strong>Filtered:</strong> Only specific models appear, based on the visibility
              settings for that source.
            </ListItem>
          </List>
        </div>
      ),
    },
    width: 15,
  },
  {
    field: 'type',
    label: TABLE_COLUMN_LABELS.SOURCE_TYPE,
    sortable: (a, b) => a.type.localeCompare(b.type),
    width: 15,
  },
  {
    field: 'enabled',
    label: TABLE_COLUMN_LABELS.ENABLE,
    sortable: (a, b) => Number(a.enabled ?? true) - Number(b.enabled ?? true),
    info: {
      popover: TABLE_COLUMN_POPOVERS.ENABLE,
    },
    width: 10,
  },
  {
    field: 'status',
    label: TABLE_COLUMN_LABELS.VALIDATION_STATUS,
    sortable: false,
    width: 10,
  },
  {
    field: 'actions',
    label: '',
    sortable: false,
    width: 15,
  },
  kebabTableColumn(),
];
