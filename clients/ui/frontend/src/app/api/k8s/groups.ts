import * as React from 'react';
import { GroupKind } from '~/app/k8sTypes';

export const useGroups = (): [GroupKind[], boolean, Error | undefined] => {
    return [[], true, undefined];
}; 