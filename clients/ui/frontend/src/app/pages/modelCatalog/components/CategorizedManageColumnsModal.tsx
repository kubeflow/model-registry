import * as React from 'react';
import {
  Button,
  Checkbox,
  EmptyState,
  EmptyStateBody,
  EmptyStateVariant,
  ExpandableSection,
  Flex,
  FlexItem,
  Label,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  Popover,
  Stack,
  StackItem,
} from '@patternfly/react-core';
import { HelpIcon, SearchIcon } from '@patternfly/react-icons';
import { DragDropSort } from '@patternfly/react-drag-drop';
import { ManageColumnSearchInput, ManagedColumn, UseManageColumnsResult } from 'mod-arch-shared';
import { reorderColumns } from './categorizedManageColumnsUtils';
import type { ColumnCategory } from './HardwareConfigurationTableColumns';
import './CategorizedManageColumnsModal.scss';

type CategorizedManageColumnsModalProps = {
  manageColumnsResult: Pick<
    UseManageColumnsResult<unknown>,
    | 'managedColumns'
    | 'setVisibleColumnIds'
    | 'defaultVisibleColumnIds'
    | 'isModalOpen'
    | 'closeModal'
  >;
  categories: ColumnCategory[];
  title?: string;
  description?: string;
  searchPlaceholder?: string;
  dataTestId?: string;
};

type CategoryGroup = {
  category: ColumnCategory;
  columns: ManagedColumn[];
};

const normalizeWhitespace = (str: string): string => str.replace(/\s+/g, ' ');

const CategorizedManageColumnsModal: React.FC<CategorizedManageColumnsModalProps> = ({
  manageColumnsResult,
  categories,
  title = 'Customize columns',
  description,
  searchPlaceholder = 'Filter by column name',
  dataTestId = 'manage-columns-modal',
}) => {
  const {
    managedColumns: initialColumns,
    setVisibleColumnIds,
    defaultVisibleColumnIds,
    isModalOpen,
    closeModal,
  } = manageColumnsResult;

  const [columns, setColumns] = React.useState<ManagedColumn[]>(initialColumns);
  const [searchValue, setSearchValue] = React.useState('');
  const [expandedCategories, setExpandedCategories] = React.useState<Set<string>>(
    () => new Set(categories.map((c) => c.id)),
  );

  const headingId = React.useId();
  const descriptionId = React.useId();

  React.useEffect(() => {
    if (isModalOpen) {
      setColumns(initialColumns);
      setSearchValue('');
      setExpandedCategories(new Set(categories.map((c) => c.id)));
    }
  }, [isModalOpen, initialColumns, categories]);

  const columnsMatchingSearch = React.useMemo(
    () =>
      searchValue
        ? columns.filter((col) =>
            normalizeWhitespace(col.label.toLowerCase()).includes(
              normalizeWhitespace(searchValue.toLowerCase()),
            ),
          )
        : columns,
    [columns, searchValue],
  );

  const categoryGroups = React.useMemo((): CategoryGroup[] => {
    const matchingIds = new Set(columnsMatchingSearch.map((c) => c.id));

    return categories
      .map((category) => {
        const categoryColumnIds = new Set<string>(category.columnIds);
        const categoryColumns = columns.filter(
          (col) => categoryColumnIds.has(col.id) && matchingIds.has(col.id),
        );
        return { category, columns: categoryColumns };
      })
      .filter((group) => group.columns.length > 0);
  }, [categories, columns, columnsMatchingSearch]);

  const selectedCount = columns.filter((col) => col.isVisible).length;

  const isDefaultState = React.useMemo(() => {
    if (!defaultVisibleColumnIds) {
      return false;
    }
    const defaultSet = new Set(defaultVisibleColumnIds);
    const currentVisibleIds = new Set(columns.filter((c) => c.isVisible).map((c) => c.id));
    if (defaultSet.size !== currentVisibleIds.size) {
      return false;
    }
    return [...defaultSet].every((id) => currentVisibleIds.has(id));
  }, [columns, defaultVisibleColumnIds]);

  const handleUpdate = React.useCallback(() => {
    const visibleColumnIds = columns.filter((col) => col.isVisible).map((col) => col.id);
    setVisibleColumnIds(visibleColumnIds);
    closeModal();
  }, [columns, setVisibleColumnIds, closeModal]);

  const handleToggleColumn = React.useCallback((columnId: string, isChecked: boolean) => {
    setColumns((prev) =>
      prev.map((col) => (col.id === columnId ? { ...col, isVisible: isChecked } : col)),
    );
  }, []);

  const handleSectionDrop = React.useCallback(
    (category: ColumnCategory, _: unknown, newItems: { id: string | number }[]) => {
      const reorderedIds = newItems.map((item) => String(item.id));
      const matchingIds = new Set<string>(category.columnIds);
      setColumns((prev) => reorderColumns(prev, matchingIds, reorderedIds));
    },
    [],
  );

  const handleRestoreDefaults = React.useCallback(() => {
    if (!defaultVisibleColumnIds) {
      return;
    }
    setColumns((prev) => {
      const defaultSet = new Set(defaultVisibleColumnIds);
      const defaultColumns = defaultVisibleColumnIds
        .map((id) => prev.find((col) => col.id === id))
        .filter((col): col is ManagedColumn => col !== undefined)
        .map((col) => ({ ...col, isVisible: true }));
      const nonDefaultColumns = prev
        .filter((col) => !defaultSet.has(col.id))
        .map((col) => ({ ...col, isVisible: false }));
      return [...defaultColumns, ...nonDefaultColumns];
    });
  }, [defaultVisibleColumnIds]);

  const handleToggleExpanded = React.useCallback((categoryId: string) => {
    setExpandedCategories((prev) => {
      const next = new Set(prev);
      if (next.has(categoryId)) {
        next.delete(categoryId);
      } else {
        next.add(categoryId);
      }
      return next;
    });
  }, []);

  if (!isModalOpen) {
    return null;
  }

  const buildDraggableItems = (sectionColumns: ManagedColumn[]) =>
    sectionColumns.map((col) => {
      const currentCol = columns.find((c) => c.id === col.id) ?? col;
      const checkbox = (
        <Checkbox
          id={`${dataTestId}-checkbox-${col.id}`}
          aria-label={col.label}
          isChecked={currentCol.isVisible}
          onChange={(_, checked) => handleToggleColumn(col.id, checked)}
        />
      );
      return {
        id: col.id,
        content: (
          <div className="pf-v6-u-display-inline-block">
            <Flex
              alignItems={{ default: 'alignItemsCenter' }}
              flexWrap={{ default: 'nowrap' }}
              spaceItems={{ default: 'spaceItemsSm' }}
            >
              <FlexItem>{checkbox}</FlexItem>
              <FlexItem>{col.label}</FlexItem>
            </Flex>
          </div>
        ),
        props: {
          checked: currentCol.isVisible,
          className: 'categorized-manage-columns__draggable',
        },
      };
    });

  const hasNoResults = columnsMatchingSearch.length === 0 && searchValue;

  return (
    <Modal
      data-testid={dataTestId}
      isOpen
      variant="small"
      onClose={closeModal}
      aria-labelledby={headingId}
      aria-describedby={description ? descriptionId : undefined}
    >
      <ModalHeader
        title={title}
        labelId={headingId}
        description={
          <div id={descriptionId}>
            <Stack hasGutter>
              {description && <StackItem>{description}</StackItem>}
              <StackItem>
                <ManageColumnSearchInput
                  value={searchValue}
                  placeholder={searchPlaceholder}
                  onSearch={setSearchValue}
                  dataTestId={`${dataTestId}-search`}
                />
              </StackItem>
              <StackItem>
                <Flex
                  justifyContent={{ default: 'justifyContentSpaceBetween' }}
                  alignItems={{ default: 'alignItemsCenter' }}
                >
                  <FlexItem>
                    <Label data-testid={`${dataTestId}-selected-count`}>
                      {selectedCount} / total {columns.length} selected
                    </Label>
                  </FlexItem>
                  {defaultVisibleColumnIds && (
                    <FlexItem>
                      <Button
                        variant="link"
                        isInline
                        isDisabled={isDefaultState}
                        onClick={handleRestoreDefaults}
                        data-testid={`${dataTestId}-restore-defaults`}
                      >
                        Restore default columns
                      </Button>
                    </FlexItem>
                  )}
                </Flex>
              </StackItem>
            </Stack>
          </div>
        }
        data-testid="generic-modal-header"
      />
      <ModalBody aria-label="Column names">
        {hasNoResults ? (
          <EmptyState
            headingLevel="h4"
            icon={SearchIcon}
            titleText="No results found"
            variant={EmptyStateVariant.sm}
          >
            <EmptyStateBody>
              No columns match your search. Try adjusting your search terms.
            </EmptyStateBody>
          </EmptyState>
        ) : (
          <Stack>
            {categoryGroups.map(({ category, columns: sectionColumns }) => (
              <StackItem key={category.id} className="categorized-manage-columns__section">
                <ExpandableSection
                  toggleContent={
                    <span className="categorized-manage-columns__section-header">
                      <strong>{category.label}</strong>
                      {category.description && (
                        <Popover bodyContent={category.description}>
                          <Button
                            variant="plain"
                            aria-label={`More info for ${category.label}`}
                            className="pf-v6-u-p-0"
                            icon={<HelpIcon />}
                            onClick={(e) => e.stopPropagation()}
                          />
                        </Popover>
                      )}
                    </span>
                  }
                  isExpanded={expandedCategories.has(category.id)}
                  onToggle={() => handleToggleExpanded(category.id)}
                  data-testid={`${dataTestId}-section-${category.id}`}
                >
                  <DragDropSort
                    items={buildDraggableItems(sectionColumns)}
                    variant="default"
                    onDrop={(_, newItems) => handleSectionDrop(category, _, newItems)}
                  />
                </ExpandableSection>
              </StackItem>
            ))}
          </Stack>
        )}
      </ModalBody>
      <ModalFooter>
        <Button
          variant="primary"
          onClick={handleUpdate}
          data-testid={`${dataTestId}-update-button`}
        >
          Update
        </Button>
        <Button variant="link" onClick={closeModal} data-testid={`${dataTestId}-cancel-button`}>
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
};

export default CategorizedManageColumnsModal;
