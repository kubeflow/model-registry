import * as React from 'react';
import {
  Alert,
  Box,
  Checkbox,
  FormGroup,
  CircularProgress,
  TextField,
  Dialog,
  DialogContent,
  DialogTitle,
  DialogActions,
  FormControlLabel,
} from '@mui/material';
import DashboardModalFooter from '~/app/concepts/dashboard/DashboardModalFooter';
import { ModelRegistryKind } from '~/app/k8sTypes';
import { ModelRegistryModel } from '~/app/api/models';
import {
  createModelRegistryBackend,
  updateModelRegistryBackend,
} from '~/app/services/modelRegistrySettingsService';
import { isValidK8sName, kindApiVersion, translateDisplayNameForK8s } from '~/app/concepts/k8s/utils';
import FormSection from '~/app/components/pf-overrides/FormSection';
import { AreaContext } from '~/app/concepts/areas/AreaContext';
import useIsAreaAvailable from '~/app/concepts/areas/useIsAreaAvailable';
import { SupportedArea } from '~/app/concepts/areas/types';
import K8sNameDescriptionField, {
  useK8sNameDescriptionFieldData,
} from '~/app/concepts/k8s/K8sNameDescriptionField';
import useModelRegistryCertificateNames from '~/app/concepts/modelRegistrySettings/useModelRegistryCertificateNames';
import {
  constructRequestBody,
  findConfigMap,
  findSecureDBType,
  isClusterWideCABundleEnabled,
  isOpenshiftCAbundleEnabled,
} from './utils';
import { RecursivePartial } from '~/typeHelpers';
import { fireFormTrackingEvent } from '~/app/concepts/analyticsTracking/segmentIOUtils';
import { TrackingOutcome } from '~/app/concepts/analyticsTracking/trackingProperties';
import ApplicationsPage from '~/app/pages/ApplicationPage';
import RedirectErrorState from '../external/RedirectErrorState';
import { CreateMRSecureDBSection, SecureDBInfo } from './CreateMRSecureDBSection';
import ModelRegistryDatabasePassword from '~/app/pages/settings/ModelRegistryDatabasePassword';
import { ResourceType, SecureDBRType } from './const';

type CreateModalProps = {
  onClose: () => void;
  refresh: () => Promise<unknown>;
  modelRegistry?: ModelRegistryKind;
};

const createEventName = 'Model Registry Created';
const updateEventName = 'Model Registry Updated';
const CreateModal: React.FC<CreateModalProps> = ({ onClose, refresh, modelRegistry: mr }) => {
  const [isSubmitting, setIsSubmitting] = React.useState(false);
  const [error, setError] = React.useState<Error>();
  const { data: nameDesc, onDataChange: setNameDesc } = useK8sNameDescriptionFieldData({
    initialData: mr,
  });
  const [host, setHost] = React.useState('');
  const [port, setPort] = React.useState('');
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [database, setDatabase] = React.useState('');
  const [addSecureDB, setAddSecureDB] = React.useState(false);
  const [isHostTouched, setIsHostTouched] = React.useState(false);
  const [isPortTouched, setIsPortTouched] = React.useState(false);
  const [isUsernameTouched, setIsUsernameTouched] = React.useState(false);
  const [isPasswordTouched, setIsPasswordTouched] = React.useState(false);
  const [isDatabaseTouched, setIsDatabaseTouched] = React.useState(false);
  const [showPassword, setShowPassword] = React.useState(false);
  const { dscStatus } = React.useContext(AreaContext);
  const secureDbEnabled = useIsAreaAvailable(SupportedArea.MODEL_REGISTRY_SECURE_DB).status;
  const configSecretsResult = useModelRegistryCertificateNames(!addSecureDB);
  const configSecrets = configSecretsResult.data;
  const configSecretsLoaded = configSecretsResult.loaded;
  const configSecretsError = configSecretsResult.error;
  const [secureDBInfo, setSecureDBInfo] = React.useState<SecureDBInfo>({
    type: SecureDBRType.CLUSTER_WIDE,
    nameSpace: '',
    resourceName: '',
    certificate: '',
    key: '',
    isValid: true,
  });
  const modelRegistryNamespace = dscStatus?.components?.modelregistry?.registriesNamespace;

  React.useEffect(() => {
    if (configSecretsLoaded && !configSecretsError && !mr) {
      setSecureDBInfo((prev: SecureDBInfo) => ({
        ...prev,
        type: isClusterWideCABundleEnabled(configSecrets.configMaps)
          ? SecureDBRType.CLUSTER_WIDE
          : isOpenshiftCAbundleEnabled(configSecrets.configMaps)
          ? SecureDBRType.OPENSHIFT
          : SecureDBRType.EXISTING,
        isValid: !!(
          isClusterWideCABundleEnabled(configSecrets.configMaps) ||
          isOpenshiftCAbundleEnabled(configSecrets.configMaps)
        ),
      }));
    }
  }, [configSecretsLoaded, configSecrets.configMaps, mr, configSecretsError]);

  React.useEffect(() => {
    if (mr) {
      const dbSpec = mr.spec.mysql || mr.spec.postgres;
      setHost(dbSpec?.host || 'Unknown');
      setPort(dbSpec?.port?.toString() || 'Unknown');
      setUsername(dbSpec?.username || 'Unknown');
      setDatabase(dbSpec?.database || 'Unknown');
      const certificateResourceRef =
        mr.spec.mysql?.sslRootCertificateConfigMap || mr.spec.mysql?.sslRootCertificateSecret;
      if (certificateResourceRef) {
        setAddSecureDB(true);
        const existingInfo = {
          type: findSecureDBType(certificateResourceRef.name, certificateResourceRef.key),
          nameSpace: '',
          key: certificateResourceRef.key,
          resourceName: certificateResourceRef.name,
          resourceType: mr.spec.mysql?.sslRootCertificateSecret
            ? ResourceType.Secret
            : ResourceType.ConfigMap,
          certificate: '',
        };
        setSecureDBInfo({ ...existingInfo, isValid: true });
      }
    }
  }, [mr]);

  if (!modelRegistryNamespace) {
    return (
      <ApplicationsPage loaded empty={false}>
        <RedirectErrorState
          title="Could not load component state"
          errorMessage="No registries namespace could be found"
        />
      </ApplicationsPage>
    );
  }

  const onCancelClose = () => {
    fireFormTrackingEvent(mr ? updateEventName : createEventName, {
      outcome: TrackingOutcome.cancel,
    });
    onBeforeClose();
  };

  const onBeforeClose = () => {
    setIsSubmitting(false);
    setError(undefined);

    setHost('');
    setPort('');
    setUsername('');
    setPassword('');
    setDatabase('');
    setIsHostTouched(false);
    setIsPortTouched(false);
    setIsUsernameTouched(false);
    setIsPasswordTouched(false);
    setIsDatabaseTouched(false);
    setShowPassword(false);
    onClose();
  };

  const onSubmit = async () => {
    setIsSubmitting(true);
    setError(undefined);

    const newDatabaseCACertificate =
      addSecureDB && secureDBInfo.type === SecureDBRType.NEW ? secureDBInfo.certificate : undefined;

    if (mr) {
      const data: RecursivePartial<ModelRegistryKind> = {
        metadata: {
          annotations: {
            'openshift.io/description': nameDesc.description,
            'openshift.io/display-name': nameDesc.name.trim(),
          },
        },
        spec: {
          oauthProxy: {},
          mysql: {
            host,
            port: Number(port),
            database,
            username,
          },
        },
      };

      try {
        await updateModelRegistryBackend(mr.metadata.name, {
          modelRegistry: constructRequestBody(data, secureDBInfo, addSecureDB),
          databasePassword: password,
          newDatabaseCACertificate,
        });
        await refresh();
        fireFormTrackingEvent(updateEventName, {
          outcome: TrackingOutcome.submit,
          success: true,
        });
        onBeforeClose();
      } catch (e) {
        if (e instanceof Error) {
          setError(e);
          fireFormTrackingEvent(updateEventName, {
            outcome: TrackingOutcome.submit,
            success: false,
            error: e.message,
          });
        }
        setIsSubmitting(false);
      }
    } else {
      const data: ModelRegistryKind = {
        apiVersion: kindApiVersion(ModelRegistryModel),
        kind: 'ModelRegistry',
        metadata: {
          name: nameDesc.k8sName.value || translateDisplayNameForK8s(nameDesc.name),
          namespace: modelRegistryNamespace,
          annotations: {
            'openshift.io/description': nameDesc.description,
            'openshift.io/display-name': nameDesc.name.trim(),
          },
        },
        spec: {
          oauthProxy: {},
          grpc: {},
          rest: {},
          mysql: {
            host,
            port: Number(port),
            database,
            username,
            skipDBCreation: false,
          },
        },
      };

      if (addSecureDB && secureDBInfo.resourceType === ResourceType.Secret && data.spec.mysql) {
        data.spec.mysql.sslRootCertificateSecret = {
          name: secureDBInfo.resourceName,
          key: secureDBInfo.key,
        };
      } else if (addSecureDB && data.spec.mysql) {
        data.spec.mysql.sslRootCertificateConfigMap = findConfigMap(secureDBInfo);
      }

      try {
        await createModelRegistryBackend({
          modelRegistry: data,
          databasePassword: password,
          newDatabaseCACertificate,
        });
        fireFormTrackingEvent(createEventName, {
          outcome: TrackingOutcome.submit,
          success: true,
        });
        await refresh();
        onBeforeClose();
      } catch (e) {
        if (e instanceof Error) {
          setError(e);
          fireFormTrackingEvent(createEventName, {
            outcome: TrackingOutcome.submit,
            success: false,
            error: e.message,
          });
        }
        setIsSubmitting(false);
      }
    }
  };

  const hasContent = (value: string): boolean => !!value.trim().length;

  const canSubmit = () =>
    !isSubmitting &&
    isValidK8sName(nameDesc.k8sName.value || translateDisplayNameForK8s(nameDesc.name)) &&
    hasContent(host) &&
    hasContent(password) &&
    hasContent(port) &&
    hasContent(username) &&
    hasContent(database) &&
    (!addSecureDB || (secureDBInfo.isValid && !configSecretsError));

  return (
    <Dialog open onClose={onCancelClose}>
      <DialogTitle>{`${mr ? 'Edit' : 'Create'} model registry`}</DialogTitle>
      <DialogContent>
        <K8sNameDescriptionField dataTestId="mr" data={nameDesc} onDataChange={setNameDesc} />
        <FormSection
          title="Connect to external MySQL database"
          description="This external database is where model data is stored."
        >
          <FormGroup>
            <TextField
              label="Host"
              required
              value={host}
              onBlur={() => setIsHostTouched(true)}
              onChange={(e) => setHost(e.target.value)}
              error={isHostTouched && !hasContent(host)}
              helperText={isHostTouched && !hasContent(host) && "Host cannot be empty"}
            />
          </FormGroup>
          <FormGroup>
            <TextField
              label="Port"
              required
              value={port}
              onBlur={() => setIsPortTouched(true)}
              onChange={(e) => setPort(e.target.value)}
              error={isPortTouched && !hasContent(port)}
              helperText={isPortTouched && !hasContent(port) && "Port cannot be empty"}
            />
          </FormGroup>
          <FormGroup>
            <TextField
              label="Username"
              required
              value={username}
              onBlur={() => setIsUsernameTouched(true)}
              onChange={(e) => setUsername(e.target.value)}
              error={isUsernameTouched && !hasContent(username)}
              helperText={isUsernameTouched && !hasContent(username) && "Username cannot be empty"}
            />
          </FormGroup>
          <FormGroup>
            <ModelRegistryDatabasePassword
              password={password || ''}
              setPassword={setPassword}
              isPasswordTouched={isPasswordTouched}
              setIsPasswordTouched={setIsPasswordTouched}
              showPassword={showPassword}
            />
          </FormGroup>
          <FormGroup>
            <TextField
              label="Database"
              required
              value={database}
              onBlur={() => setIsDatabaseTouched(true)}
              onChange={(e) => setDatabase(e.target.value)}
              error={isDatabaseTouched && !hasContent(database)}
              helperText={isDatabaseTouched && !hasContent(database) && "Database cannot be empty"}
            />
          </FormGroup>
          {secureDbEnabled && (
            <>
              <FormGroup>
                <FormControlLabel control={<Checkbox
                  checked={addSecureDB}
                  onChange={(e) => setAddSecureDB(e.target.checked)}
                />} label="Add CA certificate to secure database connection" />
              </FormGroup>
              {addSecureDB &&
                (!configSecretsLoaded && !configSecretsError ? (
                  <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                    <CircularProgress />
                  </Box>
                ) : configSecretsLoaded ? (
                  <CreateMRSecureDBSection
                    secureDBInfo={secureDBInfo}
                    modelRegistryNamespace={modelRegistryNamespace}
                    k8sName={nameDesc.k8sName.value}
                    existingCertConfigMaps={configSecrets.configMaps}
                    existingCertSecrets={configSecrets.secrets}
                    setSecureDBInfo={setSecureDBInfo}
                  />
                ) : (
                  <Alert
                    severity="error"
                    title="Error fetching config maps and secrets"
                  >
                    {configSecretsError?.message}
                  </Alert>
                ))}
            </>
          )}
        </FormSection>
      </DialogContent>
      <DialogActions>
        <DashboardModalFooter
          onCancel={onCancelClose}
          onSubmit={onSubmit}
          submitLabel={mr ? 'Update' : 'Create'}
          isSubmitLoading={isSubmitting}
          isSubmitDisabled={!canSubmit()}
          error={error}
          alertTitle={`Error ${mr ? 'updating' : 'creating'} model registry`}
        />
      </DialogActions>
    </Dialog>
  );
};

export default CreateModal; 