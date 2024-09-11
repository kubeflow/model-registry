import * as React from 'react';
import { intersection, xor } from 'lodash-es';
import type Table from './Table';

export type UseCheckboxTableBaseProps<DataType> = {
  selections: DataType[];
  tableProps: Required<Pick<React.ComponentProps<typeof Table>, 'selectAll'>>;
  toggleSelection: (selection: DataType) => void;
  isSelected: (selection: DataType) => boolean;
  disableCheck: (item: DataType, enabled: boolean) => void;
  setSelections: React.Dispatch<React.SetStateAction<DataType[]>>;
};

const useCheckboxTableBase = <T>(
  data: T[],
  selectedData: T[],
  setSelectedData: React.Dispatch<React.SetStateAction<T[]>>,
  dataMappingHelper: (selectData: T) => string,
  options?: { selectAll?: { selected?: boolean; disabled?: boolean }; persistSelections?: boolean },
): UseCheckboxTableBaseProps<T> => {
  const dataIds = React.useMemo(() => data.map(dataMappingHelper), [data, dataMappingHelper]);

  const [disabledData, setDisabledData] = React.useState<T[]>([]);

  const selectedDataIds = React.useMemo(
    () => selectedData.map(dataMappingHelper),
    [selectedData, dataMappingHelper],
  );

  // remove selected ids that are no longer present in the provided dataIds
  React.useEffect(() => {
    if (options?.persistSelections) {
      return;
    }

    const newSelectedIds = intersection(selectedDataIds, dataIds);
    const newSelectedData = newSelectedIds
      .map((id) => data.find((d) => dataMappingHelper(d) === id))
      .filter((v): v is T => !!v);
    if (selectedData.length !== newSelectedData.length) {
      setSelectedData(newSelectedData);
    }
  }, [
    data,
    dataIds,
    dataMappingHelper,
    options?.persistSelections,
    selectedData,
    selectedDataIds,
    setSelectedData,
  ]);

  const disableCheck = React.useCallback<UseCheckboxTableBaseProps<T>['disableCheck']>(
    (item, disabled) =>
      setDisabledData((prevData) =>
        disabled
          ? prevData.some((d) => dataMappingHelper(d) === dataMappingHelper(item))
            ? prevData
            : [...prevData, item]
          : prevData.filter((d) => dataMappingHelper(d) !== dataMappingHelper(item)),
      ),
    [dataMappingHelper],
  );

  return React.useMemo(() => {
    // Header is selected if all selections and all ids are equal
    // This will allow for checking of the header to "reset" to provided ids during a trim/filter
    const checkable = data.filter(
      (d) => !disabledData.some((item) => dataMappingHelper(item) === dataMappingHelper(d)),
    );

    const headerSelected =
      selectedDataIds.length > 0 &&
      xor(selectedDataIds, checkable.map(dataMappingHelper)).length === 0;

    const allDisabled = selectedData.length === 0 && disabledData.length === data.length;

    return {
      selections: selectedData,
      setSelections: setSelectedData,
      tableProps: {
        selectAll: {
          disabled: allDisabled,
          tooltip: allDisabled ? 'No selectable rows' : undefined,
          onSelect: (value) => {
            setSelectedData(value ? checkable : []);
          },
          selected: headerSelected,
          ...options?.selectAll,
        },
      },
      disableCheck,
      isSelected: (selection) => selectedDataIds.includes(dataMappingHelper(selection)),
      toggleSelection: (selection) => {
        const id = dataMappingHelper(selection);
        setSelectedData((prevData) =>
          prevData.map(dataMappingHelper).includes(id)
            ? prevData.filter(
                (currentSelectedData) => dataMappingHelper(currentSelectedData) !== id,
              )
            : [...prevData, selection],
        );
      },
    };
  }, [
    data,
    selectedDataIds,
    dataMappingHelper,
    selectedData,
    disabledData,
    setSelectedData,
    options?.selectAll,
    disableCheck,
  ]);
};

export default useCheckboxTableBase;
