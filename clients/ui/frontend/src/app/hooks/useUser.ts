import { useContext } from 'react';
import { UserSettings } from '~/shared/types';
import { AppContext } from '~/app/AppContext';

const useUser = (): UserSettings => {
  const { user } = useContext(AppContext);
  return user;
};

export default useUser;
