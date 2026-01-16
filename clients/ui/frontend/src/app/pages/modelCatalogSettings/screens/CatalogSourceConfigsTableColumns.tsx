import * as React from 'react';
import { List, ListItem } from '@patternfly/react-core';
import { kebabTableColumn, SortableData } from 'mod-arch-shared';
import { CatalogSourceConfig } from '~/app/modelCatalogTypes';

export const catalogSourceConfigsColumns: SortableData<CatalogSourceConfig>[] = [
  {
    field: 'name',
    label: 'Source name',
    sortable: (a, b) => a.name.localeCompare(b.name),
    width: 15,
  },
  {
    field: 'allowedOrganization',
    label: 'Organization',
    sortable: false,
    info: {
      popover:
        'Applies only to Hugging Face sources. Shows the organization the source syncs models from (for example, Google). Only models within this organization are included in the catalog.',
    },
    width: 15,
  },
  {
    field: 'filters',
    label: 'Model visibility',
    sortable: false,
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
    label: 'Source type',
    sortable: (a, b) => a.type.localeCompare(b.type),
    width: 15,
  },
  {
    field: 'enabled',
    label: 'Enable',
    sortable: false,
    info: {
      popover:
        'Enable a source to make its models available to users in your organization from the model catalog.',
    },
    width: 10,
  },
  {
    field: 'status',
    label: 'Validation status',
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
