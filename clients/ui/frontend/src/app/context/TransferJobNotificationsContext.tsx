import React from 'react';
import { useNotification } from '~/app/hooks/useNotification';
import { getListModelTransferJobs } from '~/app/api/service';
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
  jobId: string;
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

    const pollJobs = async () => {
      const jobs = watchedJobsRef.current;
      if (jobs.length === 0) {
        setPolling(false);
        return;
      }

      const registryGroups = new Map<string, WatchedJob[]>();
      for (const job of jobs) {
        const group = registryGroups.get(job.registryName) ?? [];
        group.push(job);
        registryGroups.set(job.registryName, group);
      }

      const resolvedIds: string[] = [];

      await Promise.all(
        Array.from(registryGroups, async ([registryName, registryJobs]) => {
          const hostPath = `${URL_PREFIX}/api/${BFF_API_VERSION}/model_registry/${registryName}`;
          const listJobs = getListModelTransferJobs(hostPath);

          try {
            const result = await listJobs({});

            for (const watched of registryJobs) {
              const found = result.items.find((j) => j.id === watched.jobId);
              if (!found || !TERMINAL_STATUSES.has(found.status)) {
                continue;
              }

              resolvedIds.push(watched.jobId);

              if (found.status === ModelTransferJobStatus.COMPLETED) {
                notificationRef.current.success(
                  REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_SUCCEEDED,
                  getRegisterAndStoreToastMessage(watched.displayParams),
                );
              } else if (found.status === ModelTransferJobStatus.FAILED) {
                notificationRef.current.error(
                  REGISTRATION_TOAST_TITLES.REGISTER_AND_STORE_ERROR,
                  getRegisterAndStoreToastMessage(watched.displayParams),
                );
              }
            }
          } catch {
            // API errors are transient; keep polling
          }
        }),
      );

      if (resolvedIds.length > 0) {
        watchedJobsRef.current = watchedJobsRef.current.filter(
          (j) => !resolvedIds.includes(j.jobId),
        );
        if (watchedJobsRef.current.length === 0) {
          setPolling(false);
        }
      }
    };

    pollJobs();
    const interval = setInterval(pollJobs, POLL_INTERVAL);
    return () => clearInterval(interval);
  }, [polling]);

  const contextValue = React.useMemo(() => ({ watchJob }), [watchJob]);

  return (
    <TransferJobNotificationsContext.Provider value={contextValue}>
      {children}
    </TransferJobNotificationsContext.Provider>
  );
};
