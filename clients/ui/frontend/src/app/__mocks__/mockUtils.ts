let uid = 0;
export const genUID = (prefix = 'uid-'): string => `${prefix}${uid++}`;
