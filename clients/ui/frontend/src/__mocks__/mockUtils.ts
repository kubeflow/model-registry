import { genRandomChars } from '~/shared/utilities/string';

export const genUID = (name: string): string => `test-uid_${name}_${genRandomChars()}`;
