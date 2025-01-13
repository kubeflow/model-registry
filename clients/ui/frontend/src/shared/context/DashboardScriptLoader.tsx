import React, { useEffect, useState } from 'react';
import { Bullseye, Spinner } from '@patternfly/react-core';
import { isIntegrated } from '~/shared/utilities/const';

type DashboardScriptLoaderProps = {
  children: React.ReactNode;
};

const loadScript = (src: string, onLoad: () => void, onError: () => void) => {
  const script = document.createElement('script');
  script.src = src;
  script.async = true;
  script.onload = onLoad;
  script.onerror = onError;
  document.head.appendChild(script);
};

/* eslint-disable no-console */
const DashboardScriptLoader: React.FC<DashboardScriptLoaderProps> = ({ children }) => {
  const [scriptLoaded, setScriptLoaded] = useState(false);

  useEffect(() => {
    const scriptUrl = '/dashboard_lib.bundle.js';

    if (!isIntegrated()) {
      console.warn(
        'DashboardScriptLoader: Script not loaded because deployment mode is not integrated',
      );
      setScriptLoaded(true);
      return;
    }

    fetch(scriptUrl, { method: 'HEAD' })
      .then((response) => {
        if (response.ok) {
          loadScript(
            scriptUrl,
            () => setScriptLoaded(true),
            () => console.error('Failed to load the script'),
          );
        } else {
          console.warn('Script not found');
        }
      })
      .catch((error) => console.error('Error checking script existence', error));
  }, []);

  return !scriptLoaded ? (
    <Bullseye>
      <Spinner />
    </Bullseye>
  ) : (
    <>{children}</>
  );
};

export default DashboardScriptLoader;
