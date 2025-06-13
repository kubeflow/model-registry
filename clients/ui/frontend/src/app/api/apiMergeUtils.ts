import { K8sAPIOptions } from '~/app/k8sTypes';

export const applyK8sAPIOptions = <T extends { queryParams?: { [key: string]: string } }>(
    options: T,
    opts?: K8sAPIOptions,
): T => {
    if (!opts) {
        return options;
    }
    return {
        ...options,
        ...(opts.dryRun && { queryParams: { ...options.queryParams, dryRun: 'All' } }),
        signal: opts.signal,
    };
}; 