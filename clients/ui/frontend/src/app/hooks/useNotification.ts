import React, { useContext } from 'react';
import { AlertVariant } from '@patternfly/react-core';
import { NotificationContext, NotificationActionTypes } from 'mod-arch-core';

enum NotificationTypes {
  SUCCESS = 'success',
  ERROR = 'error',
  INFO = 'info',
  WARNING = 'warning',
}

type NotificationProps = (title: string, message?: React.ReactNode) => void;

type NotificationRemoveProps = (id: number | undefined) => void;

type NotificationTypeFunc = {
  [key in NotificationTypes]: NotificationProps;
};

interface NotificationFunc extends NotificationTypeFunc {
  remove: NotificationRemoveProps;
}

export const useNotification = (): NotificationFunc => {
  const { notificationCount, updateNotificationCount, dispatch } = useContext(NotificationContext);

  const success: NotificationProps = React.useCallback(
    (title, message?) => {
      updateNotificationCount(notificationCount + 1);
      dispatch({
        type: NotificationActionTypes.ADD_NOTIFICATION,
        payload: {
          status: AlertVariant.success,
          title,
          timestamp: new Date(),
          message,
          id: notificationCount,
        },
      });
    },
    [dispatch, notificationCount, updateNotificationCount],
  );

  const warning: NotificationProps = React.useCallback(
    (title, message?) => {
      updateNotificationCount(notificationCount + 1);
      dispatch({
        type: NotificationActionTypes.ADD_NOTIFICATION,
        payload: {
          status: AlertVariant.warning,
          title,
          timestamp: new Date(),
          message,
          id: notificationCount,
        },
      });
    },
    [dispatch, notificationCount, updateNotificationCount],
  );

  const error: NotificationProps = React.useCallback(
    (title, message?) => {
      updateNotificationCount(notificationCount + 1);
      dispatch({
        type: NotificationActionTypes.ADD_NOTIFICATION,
        payload: {
          status: AlertVariant.danger,
          title,
          timestamp: new Date(),
          message,
          id: notificationCount,
        },
      });
    },
    [dispatch, notificationCount, updateNotificationCount],
  );

  const info: NotificationProps = React.useCallback(
    (title, message?) => {
      updateNotificationCount(notificationCount + 1);
      dispatch({
        type: NotificationActionTypes.ADD_NOTIFICATION,
        payload: {
          status: AlertVariant.info,
          title,
          timestamp: new Date(),
          message,
          id: notificationCount,
        },
      });
    },
    [dispatch, notificationCount, updateNotificationCount],
  );

  const remove: NotificationRemoveProps = React.useCallback(
    (id) => {
      dispatch({
        type: NotificationActionTypes.DELETE_NOTIFICATION,
        payload: { id },
      });
    },
    [dispatch],
  );

  const notification = React.useMemo(
    () => ({ success, error, info, warning, remove }),
    [success, error, info, warning, remove],
  );

  return notification;
};
