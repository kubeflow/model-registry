import React, { createContext } from 'react';
import { Notification, NotificationActionTypes, NotificationAction } from '~/app/types';

type NotificationContextProps = {
  notifications: Notification[];
  notificationCount: number;
  updateNotificationCount: React.Dispatch<React.SetStateAction<number>>;
  dispatch: React.Dispatch<NotificationAction>;
};

export const NotificationContext = createContext<NotificationContextProps>({
  notifications: [],
  notificationCount: 0,
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  updateNotificationCount: () => {},
  // eslint-disable-next-line @typescript-eslint/no-empty-function
  dispatch: () => {},
});

const notificationReducer: React.Reducer<Notification[], NotificationAction> = (
  notifications,
  action,
) => {
  switch (action.type) {
    case NotificationActionTypes.ADD_NOTIFICATION: {
      return [
        ...notifications,
        {
          status: action.payload.status,
          title: action.payload.title,
          timestamp: action.payload.timestamp,
          message: action.payload.message,
          id: action.payload.id,
        },
      ];
    }
    case NotificationActionTypes.DELETE_NOTIFICATION: {
      return notifications.filter((t) => t.id !== action.payload.id);
    }
    default: {
      return notifications;
    }
  }
};

type NotificationContextProviderProps = {
  children: React.ReactNode;
};

export const NotificationContextProvider: React.FC<NotificationContextProviderProps> = ({
  children,
}) => {
  const [notifications, dispatch] = React.useReducer(notificationReducer, []);
  const [notificationCount, setNotificationCount] = React.useState(0);

  return (
    <NotificationContext.Provider
      value={React.useMemo(
        () => ({
          notifications,
          notificationCount,
          updateNotificationCount: setNotificationCount,
          dispatch,
        }),
        [notifications, notificationCount, setNotificationCount, dispatch],
      )}
    >
      {children}
    </NotificationContext.Provider>
  );
};
