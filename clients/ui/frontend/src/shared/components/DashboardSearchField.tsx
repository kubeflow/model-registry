import * as React from 'react';
import { InputGroup, InputGroupItem, TextInput } from '@patternfly/react-core';
import SimpleSelect from '~/shared/components/SimpleSelect';
import { asEnumMember } from '~/app/utils';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';

// List all the possible search fields here
export enum SearchType {
  NAME = 'Name',
  DESCRIPTION = 'Description',
  USER = 'User',
  PROJECT = 'Project',
  METRIC = 'Metric',
  PROTECTED_ATTRIBUTE = 'Protected attribute',
  PRIVILEGED_VALUE = 'Privileged value',
  UNPRIVILEGED_VALUE = 'Unprivileged value',
  OUTPUT = 'Output',
  OUTPUT_VALUE = 'Output value',
  PROVIDER = 'Provider',
  IDENTIFIER = 'Identifier',
  KEYWORD = 'Keyword',
  AUTHOR = 'Author',
  OWNER = 'Owner',
}

type DashboardSearchFieldProps = {
  types: SearchType[];
  searchType: SearchType;
  onSearchTypeChange: (searchType: SearchType) => void;
  searchValue: string;
  onSearchValueChange: (searchValue: string) => void;
  icon?: React.ReactNode;
};

const DashboardSearchField: React.FC<DashboardSearchFieldProps> = ({
  types,
  searchValue,
  searchType,
  onSearchValueChange,
  onSearchTypeChange,
  icon,
}) => (
  <InputGroup data-testid="dashboard-table-toolbar">
    <InputGroupItem>
      <SimpleSelect
        aria-label="Filter type"
        dataTestId="filter-dropdown-select"
        options={types.map((key) => ({
          key,
          label: key,
        }))}
        value={searchType}
        onChange={(key) => {
          const enumMember = asEnumMember(key, SearchType);
          if (enumMember !== null) {
            onSearchTypeChange(enumMember);
          }
        }}
        icon={icon}
      />
    </InputGroupItem>
    <InputGroupItem>
      <FormFieldset
        className="toolbar-fieldset-wrapper"
        component={
          <TextInput
            value={searchValue}
            type="text"
            onChange={(_, newSearch) => {
              onSearchValueChange(newSearch);
            }}
            style={{ minWidth: '200px' }}
            aria-label="Search"
          />
        }
        field={`Find by ${searchType.toLowerCase()}`}
      />
    </InputGroupItem>
  </InputGroup>
);

export default DashboardSearchField;
