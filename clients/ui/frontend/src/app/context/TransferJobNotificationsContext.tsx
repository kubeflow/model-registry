import React from 'react';
import { useQueryParamNamespaces } from 'mod-arch-core';
import { useNotification } from '~/app/hooks/useNotification';
import { getModelTransferJob } from '~/app/api/service';
import { ModelTransferJobStatus } from '~/app/types';
import {
  BFF_API_VERSION,
  POLL_INTERVAL,
  URL_PREFIX,
  REGISTRATION_TOAST_TITLES,
} from '~/app/utilities/const';
import {
  getRegisterAndStoreToastMessage,
  type RegistrationToastMessagesParams,
} from '~/app/pages/modelRegistry/screens/RegisterModel/registrationToastMessages';

type WatchedJob = {
  jobName: string;
  registryName: string;
  displayParams: RegistrationToastMessagesParams;
};

type TransferJobNotificationsContextType = {
  watchJob: (job: WatchedJob) => void;
};

export const TransferJobNotificationsContext =
  React.createContext<TransferJobNotificationsContextType>({
    watchJob: () => undefined,
  });

const TERMINAL_STATUSES = new Set<ModelTransferJobStatus>([
  ModelTransferJobStatus.COMPLETED,
  ModelTransferJobStatus.FAILED,
  ModelTransferJobStatus.CANCELLED,
]);

export const TransferJobNotificationsProvider: React.FC<React.PropsWithChildren> = ({
  children,
}) => {
  const notification = useNotification();
  const notificationRef = React.useRef(notification);
  notificationRef.current = notification;

  const queryParams = useQueryParamNamespaces();

  const watchedJobsRef = React.useRef<WatchedJob[]>([]);
  const [polling, setPolling] = React.useState(false);

  const watchJob = React.useCallback((job: WatchedJob) => {
    watchedJobsRef.current = [...watchedJobsRef.current, job];
    setPolling(true);
  }, []);

  React.useEffect(() => {
    if (!polling) {
      return undefined;
    }

    let cancelled = false;

    const pollJobs = async () => {
      const jobs = watchedJobsRef.current;
      if (jobs.length === 0) {
        setPolling(false);
        return;
      }

      const resolvedNames: string[] = [];

      await Promise.all(
        jobs.map(async (watched) => {
          const hostPath = `${URL_PREFIX}/api/${BFF_API_VERSION}/model_registry/${watched.registryName}`;
          const fetchJob = getModelTransferJob(hostPath, queryParams);

          try {
            const job = await fetchJob({}, watched.jobName);
            if (!TERMINAL_STATUSES.has(job.status)) {
              return;
            }

            resolvedNames.push(watched.jobName);

            if (job.status === ModelTransferJobStatus.COMPLETED) {
              notificationRef.current.success(
                REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_SUCCEEDED,
                getRegisterAndStoreToastMessage(watched.displayParams),
              );
            } else if (job.status === ModelTransferJobStatus.FAILED) {
              notificationRef.current.error(
                REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_ERROR,
                getRegisterAndStoreToastMessage(watched.displayParams),
              );
            }
          } catch {
            // API errors are transient; keep polling
          }
        }),
      );

      if (cancelled) {
        return;
      }

      if (resolvedNames.length > 0) {
        watchedJobsRef.current = watchedJobsRef.current.filter(
          (j) => !resolvedNames.includes(j.jobName),
        );
        if (watchedJobsRef.current.length === 0) {
          setPolling(false);
        }
      }
    };

    pollJobs();
    const interval = setInterval(pollJobs, POLL_INTERVAL);
    return () => {
      cancelled = true;
      clearInterval(interval);
    };
  }, [polling, queryParams]);

  const contextValue = React.useMemo(() => ({ watchJob }), [watchJob]);

  return (
    <TransferJobNotificationsContext.Provider value={contextValue}>
      {children}
    </TransferJobNotificationsContext.Provider>
  );
};
