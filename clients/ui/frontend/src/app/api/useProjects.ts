import { ProjectKind } from '~/app/k8sTypes';
import { listProjects } from '~/app/api/k8s/projects';
import useFetch, { FetchState } from '~/app/utils/useFetch';
import { POLL_INTERVAL } from '~/app/utils/const';
import * as React from 'react';

const useProjects = (): FetchState<ProjectKind[]> => {
    const getProjects = React.useCallback(() => listProjects(), []);
    return useFetch<ProjectKind[]>(getProjects, [], { refreshRate: POLL_INTERVAL });
};

export default useProjects; 