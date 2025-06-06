import { genRandomChars } from 'mod-arch-shared';

export const genUID = (name: string): string => `test-uid_${name}_${genRandomChars()}`;
