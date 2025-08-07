import { useContext } from 'react';
import { UserSettings } from 'mod-arch-core';
import { AppContext } from '~/app/context/AppContext';

const useUser = (): UserSettings => {
  const { user } = useContext(AppContext);
  return user;
};

export default useUser;
