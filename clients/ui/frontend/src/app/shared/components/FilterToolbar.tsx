import * as React from 'react';
import {
  ToolbarFilter,
  ToolbarGroup,
  ToolbarItem,
  ToolbarToggleGroup,
  Dropdown,
  DropdownItem,
  MenuToggle,
  DropdownList,
} from '@patternfly/react-core';
import { FilterIcon } from '@patternfly/react-icons';

type FilterOptionRenders = {
  onChange: (value?: string, label?: string) => void;
  value?: string;
  label?: string;
};

export type ToolbarFilterProps<T extends string> = React.ComponentProps<typeof ToolbarGroup> & {
  children?: React.ReactNode;
  filterOptions: { [key in T]?: string };
  filterOptionRenders: Record<T, (props: FilterOptionRenders) => React.ReactNode>;
  filterData: Record<T, string | { label: string; value: string } | undefined>;
  onFilterUpdate: (filterType: T, value?: string | { label: string; value: string }) => void;
  testId?: string;
};

function FilterToolbar<T extends string>({
  filterOptions,
  filterOptionRenders,
  filterData,
  onFilterUpdate,
  children,
  testId = 'filter-toolbar',
  ...toolbarGroupProps
}: ToolbarFilterProps<T>): React.JSX.Element {
  // eslint-disable-next-line @typescript-eslint/consistent-type-assertions
  const keys = Object.keys(filterOptions) as Array<T>;
  const [open, setOpen] = React.useState(false);
  const [currentFilterType, setCurrentFilterType] = React.useState<T>(keys[0]);
  const filterItem = filterData[currentFilterType];

  return (
    <>
      <ToolbarToggleGroup breakpoint="md" toggleIcon={<FilterIcon />}>
        <ToolbarGroup variant="filter-group" data-testid={testId} {...toolbarGroupProps}>
          <ToolbarItem>
            <Dropdown
              onOpenChange={(isOpenChange) => setOpen(isOpenChange)}
              shouldFocusToggleOnSelect
              toggle={(toggleRef) => (
                <MenuToggle
                  data-testid={`${testId}-dropdown`}
                  id={`${testId}-toggle-button`}
                  ref={toggleRef}
                  aria-label="Filter toggle"
                  onClick={() => setOpen(!open)}
                  isExpanded={open}
                  icon={<FilterIcon />}
                >
                  {filterOptions[currentFilterType]}
                </MenuToggle>
              )}
              isOpen={open}
              popperProps={{ appendTo: 'inline' }}
            >
              <DropdownList>
                {keys.map((filterKey) => (
                  <DropdownItem
                    key={filterKey}
                    id={filterKey}
                    onClick={() => {
                      setOpen(false);
                      setCurrentFilterType(filterKey);
                    }}
                  >
                    {filterOptions[filterKey]}
                  </DropdownItem>
                ))}
              </DropdownList>
            </Dropdown>
          </ToolbarItem>
          {keys.map((filterKey) => {
            const optionValue = filterOptions[filterKey];
            const data = filterData[filterKey];
            const dataValue: { label: string; value: string } | undefined =
              typeof data === 'string' ? { label: data, value: data } : data;
            return optionValue ? (
              <ToolbarFilter
                key={filterKey}
                categoryName={optionValue}
                data-testid={`${testId}-text-field`}
                labels={
                  data && dataValue
                    ? [
                        {
                          key: filterKey,
                          node: (
                            <span data-testid={`${filterKey}-filter-chip`}>{dataValue.label}</span>
                          ),
                        },
                      ]
                    : []
                }
                deleteLabel={() => {
                  onFilterUpdate(filterKey, '');
                }}
                showToolbarItem={currentFilterType === filterKey}
              >
                {filterOptionRenders[filterKey]({
                  onChange: (value, label) =>
                    onFilterUpdate(filterKey, label && value ? { label, value } : value),
                  ...(typeof filterItem === 'string' ? { value: filterItem } : filterItem),
                })}
              </ToolbarFilter>
            ) : null;
          })}
        </ToolbarGroup>
      </ToolbarToggleGroup>
      {children}
    </>
  );
}

export default FilterToolbar;
