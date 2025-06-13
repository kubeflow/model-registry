import * as React from 'react';
import { NavDataItem, NavSidebar } from 'mod-arch-shared';
import { useNavData } from './AppRoutes';

const AppNavSidebar: React.FC = () => {
  const navData = useNavData(); // Call useNavData here, safely within context
  return <NavSidebar navData={navData as NavDataItem[]} />;
};

export default AppNavSidebar;
