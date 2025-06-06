// TODO: Move this code to shared library once the migration completes.
import { genRandomChars } from 'mod-arch-shared';

export const genUID = (name: string): string => `test-uid_${name}_${genRandomChars()}`;
