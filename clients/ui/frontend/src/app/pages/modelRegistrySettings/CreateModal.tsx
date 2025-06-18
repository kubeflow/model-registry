import * as React from 'react';
import {
  Alert,
  Bullseye,
  Checkbox,
  HelperText,
  HelperTextItem,
  Spinner,
  TextInput,
  Modal,
  ModalBody,
  ModalFooter,
} from '@patternfly/react-core';
import DashboardModalFooter from '~/app/concepts/dashboard/DashboardModalFooter';
import { ModelRegistryKind } from '~/app/k8sTypes';
import { ModelRegistryModel } from '~/app/api/models';
import {
  createModelRegistryBackend,
  updateModelRegistryBackend,
} from '~/app/services/modelRegistrySettingsService';
import {
  isValidK8sName,
  kindApiVersion,
  translateDisplayNameForK8s,
} from '~/app/concepts/k8s/utils';
import FormSection from '~/app/components/pf-overrides/FormSection';
import { AreaContext } from '~/app/concepts/areas/AreaContext';
import useIsAreaAvailable from '~/app/concepts/areas/useIsAreaAvailable';
import { SupportedArea } from '~/app/concepts/areas/types';
import K8sNameDescriptionField from '~/concepts/k8s/K8sNameDescriptionField/K8sNameDescriptionField';
import { useK8sNameDescriptionFieldData } from '~/app/concepts/k8s/K8sNameDescriptionField/useK8sNameDescriptionField';
import useModelRegistryCertificateNames from '~/app/concepts/modelRegistrySettings/useModelRegistryCertificateNames';
import { RecursivePartial } from '~/typeHelpers';
import { fireFormTrackingEvent } from '~/app/concepts/analyticsTracking/segmentIOUtils';
import { TrackingOutcome } from '~/app/concepts/analyticsTracking/trackingProperties';
import ApplicationsPage from '~/app/pages/ApplicationPage';
import ModelRegistryDatabasePassword from '~/app/pages/settings/ModelRegistryDatabasePassword';
import ThemeAwareFormGroupWrapper from '~/app/pages/settings/components/ThemeAwareFormGroupWrapper';
import RedirectErrorState from '~/app/pages/external/RedirectErrorState';
import { CreateMRSecureDBSection, SecureDBInfo } from './CreateMRSecureDBSection';
import {
  constructRequestBody,
  findConfigMap,
  findSecureDBType,
  isClusterWideCABundleEnabled,
  isOpenshiftCAbundleEnabled,
  CABundleConfig,
} from './utils';
import { ResourceType, SecureDBRType } from './const';

const ODH_TRUSTED_BUNDLE = 'odh-trusted-ca-bundle';
const CA_BUNDLE_CRT = 'ca-bundle.crt';
const ODH_CA_BUNDLE_CRT = 'odh-ca-bundle.crt';

type CreateModalProps = {
  onClose: () => void;
  refresh: () => Promise<unknown>;
  modelRegistry?: ModelRegistryKind;
};

interface DscStatus {
  components?: {
    modelregistry?: {
      registriesNamespace?: string;
    };
  };
}

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
  const { dscStatus } = React.useContext(AreaContext) as { dscStatus: DscStatus };
  const secureDbEnabled = useIsAreaAvailable(SupportedArea.MODEL_REGISTRY_SECURE_DB).status;
  const configSecretsResult = useModelRegistryCertificateNames(!addSecureDB);
  const configSecrets = configSecretsResult.data;
  const configSecretsLoaded = configSecretsResult.loaded;
  const configSecretsError = configSecretsResult.error;
  const [secureDBInfo, setSecureDBInfo] = React.useState<SecureDBInfo>({
    type: SecureDBRType.CLUSTER_WIDE,
    resourceName: '',
    key: '',
  });
  const modelRegistryNamespace = dscStatus.components?.modelregistry?.registriesNamespace;

  const caBundleConfig: CABundleConfig = React.useMemo(
    () => ({
      trustedBundleName: ODH_TRUSTED_BUNDLE,
      clusterWideKey: CA_BUNDLE_CRT,
      openshiftKey: ODH_CA_BUNDLE_CRT,
    }),
    [],
  );

  React.useEffect(() => {
    if (configSecretsLoaded && !configSecretsError && !mr) {
      setSecureDBInfo((prev) => ({
        ...prev,
        type: isClusterWideCABundleEnabled(configSecrets.configMaps, caBundleConfig)
          ? SecureDBRType.CLUSTER_WIDE
          : isOpenshiftCAbundleEnabled(configSecrets.configMaps, caBundleConfig)
            ? SecureDBRType.OPENSHIFT
            : SecureDBRType.EXISTING,
      }));
    }
  }, [configSecretsLoaded, configSecrets.configMaps, mr, configSecretsError, caBundleConfig]);

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
        const existingInfo: SecureDBInfo = {
          type: findSecureDBType(
            certificateResourceRef.name,
            certificateResourceRef.key,
            caBundleConfig,
          ),
          key: certificateResourceRef.key,
          resourceName: certificateResourceRef.name,
          resourceType: mr.spec.mysql?.sslRootCertificateSecret
            ? ResourceType.Secret
            : ResourceType.ConfigMap,
        };
        setSecureDBInfo(existingInfo);
      }
    }
  }, [mr, caBundleConfig]);

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
          modelRegistry: constructRequestBody(data, secureDBInfo, addSecureDB, caBundleConfig),
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
          name: secureDBInfo.resourceName ?? '',
          key: secureDBInfo.key ?? '',
        };
      } else if (addSecureDB && data.spec.mysql) {
        data.spec.mysql.sslRootCertificateConfigMap = findConfigMap(secureDBInfo, caBundleConfig);
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
    (!addSecureDB || !configSecretsError);

  return (
    <Modal
      isOpen
      title={`${mr ? 'Edit' : 'Create'} model registry`}
      onClose={onCancelClose}
      variant="small"
    >
      <ModalBody>
        <K8sNameDescriptionField dataTestId="mr" data={nameDesc} onDataChange={setNameDesc} />
        <FormSection
          title="Connect to external MySQL database"
          description="This external database is where model data is stored."
        >
          <ThemeAwareFormGroupWrapper label="Host" fieldId="host">
            <TextInput
              isRequired
              type="text"
              id="mr-host"
              name="mr-host"
              value={host}
              onBlur={() => setIsHostTouched(true)}
              onChange={(_e, value) => setHost(value)}
              validated={isHostTouched && !hasContent(host) ? 'error' : 'default'}
            />
          </ThemeAwareFormGroupWrapper>
          <ThemeAwareFormGroupWrapper label="Port" fieldId="port">
            <TextInput
              isRequired
              type="text"
              id="mr-port"
              name="mr-port"
              value={port}
              onBlur={() => setIsPortTouched(true)}
              onChange={(_e, value) => setPort(value)}
              validated={isPortTouched && !hasContent(port) ? 'error' : 'default'}
            />
            {isPortTouched && !hasContent(port) && (
              <HelperText>
                <HelperTextItem variant="error" data-testid="mr-port-error">
                  Port cannot be empty
                </HelperTextItem>
              </HelperText>
            )}
          </ThemeAwareFormGroupWrapper>
          <ThemeAwareFormGroupWrapper label="Username" fieldId="username">
            <TextInput
              isRequired
              type="text"
              id="mr-username"
              name="mr-username"
              value={username}
              onBlur={() => setIsUsernameTouched(true)}
              onChange={(_e, value) => setUsername(value)}
              validated={isUsernameTouched && !hasContent(username) ? 'error' : 'default'}
            />
          </ThemeAwareFormGroupWrapper>
          <ThemeAwareFormGroupWrapper label="Password" fieldId="password">
            <ModelRegistryDatabasePassword
              password={password || ''}
              setPassword={setPassword}
              isPasswordTouched={isPasswordTouched}
              setIsPasswordTouched={setIsPasswordTouched}
              showPassword={showPassword}
            />
          </ThemeAwareFormGroupWrapper>
          <ThemeAwareFormGroupWrapper label="Database" fieldId="database">
            <TextInput
              required
              value={database}
              onBlur={() => setIsDatabaseTouched(true)}
              onChange={(_e, value) => setDatabase(value)}
              validated={isDatabaseTouched && !hasContent(database) ? 'error' : 'default'}
            />
          </ThemeAwareFormGroupWrapper>
          {secureDbEnabled && (
            <>
              <ThemeAwareFormGroupWrapper
                label="Add CA certificate to secure database connection"
                fieldId="addSecureDB"
              >
                <Checkbox
                  id="add-secure-db-checkbox"
                  isChecked={addSecureDB}
                  onChange={(event, checked) => setAddSecureDB(checked)}
                  label="Add CA certificate to secure database connection"
                />
              </ThemeAwareFormGroupWrapper>
              {addSecureDB &&
                (!configSecretsLoaded && !configSecretsError ? (
                  <Bullseye>
                    <Spinner />
                  </Bullseye>
                ) : configSecretsLoaded ? (
                  <CreateMRSecureDBSection
                    secureDBInfo={secureDBInfo}
                    setSecureDBInfo={setSecureDBInfo}
                  />
                ) : (
                  <Alert variant="danger" title="Error fetching config maps and secrets">
                    {configSecretsError?.message}
                  </Alert>
                ))}
            </>
          )}
        </FormSection>
      </ModalBody>
      <ModalFooter>
        <DashboardModalFooter
          onCancel={onCancelClose}
          onSubmit={onSubmit}
          submitLabel={mr ? 'Update' : 'Create'}
          isSubmitLoading={isSubmitting}
          isSubmitDisabled={!canSubmit()}
          error={error}
          alertTitle={`Error ${mr ? 'updating' : 'creating'} model registry`}
        />
      </ModalFooter>
    </Modal>
  );
};

export default CreateModal;
