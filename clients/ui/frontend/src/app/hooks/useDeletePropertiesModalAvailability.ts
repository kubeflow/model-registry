import { useBrowserStorage } from 'mod-arch-core';

const useDeletePropertiesModalAvailability = (): [boolean, (v: boolean) => void] =>
  useBrowserStorage<boolean>('delete.properties.modal.preference', false);

export default useDeletePropertiesModalAvailability;
