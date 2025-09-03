import * as React from 'react';
import NavSidebar from '~/app/standalone/NavSidebar';
import { useNavData } from '~/app/AppRoutes';

const AppNavSidebar: React.FC = () => {
  const navData = useNavData(); // Call useNavData here, safely within context
  return <NavSidebar navData={navData} />;
};

export default AppNavSidebar;
