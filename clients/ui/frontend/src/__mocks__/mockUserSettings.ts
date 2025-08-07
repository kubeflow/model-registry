import { UserSettings } from 'mod-arch-core';

type MockUserSettingsType = {
  userId?: string;
  clusterAdmin?: boolean;
};

export const mockUserSettings = ({
  userId = 'user@example.com',
  clusterAdmin = true,
}: MockUserSettingsType): UserSettings => ({
  userId,
  clusterAdmin,
});
