import { FetchState, FetchStateCallbackPromise, NotReadyError, useFetchState } from 'mod-arch-core';
import React from 'react';
import {
  CatalogFilterOptionsList,
  CatalogPerformanceArtifactList,
  ModelCatalogFilterStates,
  PerformanceArtifactsParams,
} from '~/app/modelCatalogTypes';
import { useModelCatalogAPI } from './useModelCatalogAPI';

/**
 * Hook for fetching performance artifacts from the /performance_artifacts endpoint.
 * This endpoint returns only performance artifacts and supports server-side filtering,
 * sorting, and pagination.
 *
 * @param sourceId - The catalog source ID
 * @param modelName - The model name
 * @param params - Performance-specific parameters (targetRPS, latencyProperty, recommendations, pagination)
 * @param filterData - Current filter state for building filterQuery
 * @param filterOptions - Filter options from the API for validation
 * @param enabled - Whether to enable fetching (default: true)
 */
export const useCatalogPerformanceArtifacts = (
  sourceId: string,
  modelName: string,
  params?: PerformanceArtifactsParams,
  filterData?: ModelCatalogFilterStates,
  filterOptions?: CatalogFilterOptionsList | null,
  enabled = true,
): FetchState<CatalogPerformanceArtifactList> => {
  const { api, apiAvailable } = useModelCatalogAPI();

  const performanceParams: PerformanceArtifactsParams | undefined = React.useMemo(
    () => ({
      targetRPS: params?.targetRPS,
      recommendations: params?.recommendations ?? true,
      rpsProperty: params?.rpsProperty,
      latencyProperty: params?.latencyProperty,
      hardwareCountProperty: params?.hardwareCountProperty,
      hardwareTypeProperty: params?.hardwareTypeProperty,
      pageSize: params?.pageSize,
      orderBy: params?.orderBy,
      sortOrder: params?.sortOrder,
      nextPageToken: params?.nextPageToken,
    }),
    [
      params?.targetRPS,
      params?.recommendations,
      params?.rpsProperty,
      params?.latencyProperty,
      params?.hardwareCountProperty,
      params?.hardwareTypeProperty,
      params?.pageSize,
      params?.orderBy,
      params?.sortOrder,
      params?.nextPageToken,
    ],
  );

  const call = React.useCallback<FetchStateCallbackPromise<CatalogPerformanceArtifactList>>(
    (opts) => {
      if (!apiAvailable) {
        return Promise.reject(new Error('API not yet available'));
      }
      if (!sourceId) {
        return Promise.reject(new NotReadyError('No source id'));
      }
      if (!modelName) {
        return Promise.reject(new NotReadyError('No model name'));
      }
      if (!enabled) {
        return Promise.reject(new NotReadyError('Fetching is disabled'));
      }
      return api.getPerformanceArtifacts(
        opts,
        sourceId,
        modelName,
        performanceParams,
        filterData,
        filterOptions,
      );
    },
    [apiAvailable, sourceId, modelName, enabled, api, performanceParams, filterData, filterOptions],
  );

  return useFetchState(
    call,
    { items: [], size: 0, pageSize: 0, nextPageToken: '' },
    {
      initialPromisePurity: true,
    },
  );
};
