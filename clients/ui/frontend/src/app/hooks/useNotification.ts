import React, { useContext, useRef } from 'react';
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
  const { updateNotificationCount, dispatch } = useContext(NotificationContext);
  const nextIdRef = useRef(0);

  const addNotification = React.useCallback(
    (status: AlertVariant, title: string, message?: React.ReactNode) => {
      const id = nextIdRef.current++;
      updateNotificationCount(id);
      dispatch({
        type: NotificationActionTypes.ADD_NOTIFICATION,
        payload: {
          status,
          title,
          timestamp: new Date(),
          message,
          id,
        },
      });
    },
    [dispatch, updateNotificationCount],
  );

  const success: NotificationProps = React.useCallback(
    (title, message?) => addNotification(AlertVariant.success, title, message),
    [addNotification],
  );

  const warning: NotificationProps = React.useCallback(
    (title, message?) => addNotification(AlertVariant.warning, title, message),
    [addNotification],
  );

  const error: NotificationProps = React.useCallback(
    (title, message?) => addNotification(AlertVariant.danger, title, message),
    [addNotification],
  );

  const info: NotificationProps = React.useCallback(
    (title, message?) => addNotification(AlertVariant.info, title, message),
    [addNotification],
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
