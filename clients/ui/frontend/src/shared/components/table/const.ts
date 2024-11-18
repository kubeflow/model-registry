import { SortableData } from './types';

export const CHECKBOX_FIELD_ID = 'checkbox';
export const KEBAB_FIELD_ID = 'kebab';
export const EXPAND_FIELD_ID = 'expand';

export const checkboxTableColumn = (): SortableData<unknown> => ({
  label: '',
  field: CHECKBOX_FIELD_ID,
  sortable: false,
});

export const kebabTableColumn = (): SortableData<unknown> => ({
  label: '',
  field: KEBAB_FIELD_ID,
  sortable: false,
});

export const expandTableColumn = (): SortableData<unknown> => ({
  label: '',
  field: EXPAND_FIELD_ID,
  sortable: false,
});
