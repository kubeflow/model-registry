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
