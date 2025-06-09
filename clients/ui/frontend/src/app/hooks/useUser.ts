import { useContext } from 'react';
import { UserSettings } from 'mod-arch-shared';
import { AppContext } from '~/app/AppContext';

const useUser = (): UserSettings => {
  const { user } = useContext(AppContext);
  return user;
};

export default useUser;
