import { SearchType } from '~/shared/components/DashboardSearchField';
import { RegisteredModel } from '~/app/types';

export const asEnumMember = <T extends object>(
  member: T[keyof T] | string | number | undefined | null,
  e: T,
): T[keyof T] | null => (isEnumMember(member, e) ? member : null);

export const isEnumMember = <T extends object>(
  member: T[keyof T] | string | number | undefined | unknown | null,
  e: T,
): member is T[keyof T] => {
  if (member != null) {
    return Object.entries(e)
      .filter(([key]) => Number.isNaN(Number(key)))
      .map(([, value]) => value)
      .includes(member);
  }
  return false;
};

export const filterRegisteredModels = (
  unfilteredRegisteredModels: RegisteredModel[],
  search: string,
  searchType: SearchType,
): RegisteredModel[] =>
  unfilteredRegisteredModels.filter((rm: RegisteredModel) => {
    if (!search) {
      return true;
    }

    switch (searchType) {
      case SearchType.KEYWORD:
        return (
          rm.name.toLowerCase().includes(search.toLowerCase()) ||
          (rm.description && rm.description.toLowerCase().includes(search.toLowerCase()))
        );

      case SearchType.OWNER:
        return rm.owner && rm.owner.toLowerCase().includes(search.toLowerCase());

      default:
        return true;
    }
  });
