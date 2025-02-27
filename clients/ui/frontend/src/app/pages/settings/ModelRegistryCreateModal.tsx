import * as React from 'react';
import {
  Button,
  Form,
  FormGroup,
  HelperText,
  HelperTextItem,
  TextInput,
} from '@patternfly/react-core';
import { Modal } from '@patternfly/react-core/deprecated';
import { useNavigate } from 'react-router';
import ModelRegistryCreateModalFooter from '~/app/pages/settings/ModelRegistryCreateModalFooter';
import FormSection from '~/shared/components/pf-overrides/FormSection';

import ModelRegistryDatabasePassword from '~/app/pages/settings/ModelRegistryDatabasePassword';
import K8sNameDescriptionField from '~/concepts/k8s/K8sNameDescriptionField/K8sNameDescriptionField';

type CreateModalProps = {
  onClose: () => void;
  // refresh: () => Promise<unknown>;
  // modelRegistry: ModelRegistry;
};

const CreateModal: React.FC<CreateModalProps> = ({
  onClose,
  // refresh,
  // modelRegistry,
}) => {
  const [error, setError] = React.useState<Error>();

  const [host, setHost] = React.useState('');
  const [port, setPort] = React.useState('');
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [database, setDatabase] = React.useState('');
  //   const [addSecureDB, setAddSecureDB] = React.useState(false);
  const [isHostTouched, setIsHostTouched] = React.useState(false);
  const [isPortTouched, setIsPortTouched] = React.useState(false);
  const [isUsernameTouched, setIsUsernameTouched] = React.useState(false);
  const [isPasswordTouched, setIsPasswordTouched] = React.useState(false);
  const [isDatabaseTouched, setIsDatabaseTouched] = React.useState(false);
  const [showPassword, setShowPassword] = React.useState(false);

  const navigate = useNavigate();

  const onBeforeClose = () => {
    // setIsSubmitting(false);
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

  const hasContent = (value: string): boolean => !!value.trim().length;

  const canSubmit = () =>
    // TODO: implement once we have the endpoint
    // !isSubmitting &&
    // isValidK8sName(nameDesc.k8sName.value || translateDisplayNameForK8s(nameDesc.name))
    // &&
    hasContent(host) &&
    hasContent(password) &&
    hasContent(port) &&
    hasContent(username) &&
    hasContent(database);
  // &&
  // (!addSecureDB || (secureDBInfo.isValid && !configSecretsError))

  const onSubmit = () => {
    navigate(`/model-registry-settings`);
    onClose();
  };

  return (
    <Modal
      isOpen
      title="Create model registry"
      onClose={onBeforeClose}
      actions={[
        <Button key="create-button" variant="primary" isDisabled={!canSubmit()} onClick={onSubmit}>
          Create
        </Button>,
        <Button key="cancel-button" variant="secondary" onClick={onBeforeClose}>
          Cancel
        </Button>,
      ]}
      variant="medium"
      footer={
        <ModelRegistryCreateModalFooter
          onCancel={onBeforeClose}
          onSubmit={onSubmit}
          submitLabel="Create"
          // isSubmitLoading={isSubmitting}
          isSubmitDisabled={!canSubmit()}
          error={error}
          alertTitle={`Error ${'creating'} model registry`}
        />
      }
    >
      <Form>
        <K8sNameDescriptionField
          dataTestId="mr"
          // data={nameDesc}
          //  onDataChange={setNameDesc}
        />
        <FormSection
          title="Connect to external MySQL database"
          description="This external database is where model data is stored."
        >
          <FormGroup label="Host" isRequired fieldId="mr-host">
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
            {isHostTouched && !hasContent(host) && (
              <HelperText>
                <HelperTextItem variant="error" data-testid="mr-host-error">
                  Host cannot be empty
                </HelperTextItem>
              </HelperText>
            )}
          </FormGroup>
          <FormGroup label="Port" isRequired fieldId="mr-port">
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
          </FormGroup>
          <FormGroup label="Username" isRequired fieldId="mr-username">
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
            {isUsernameTouched && !hasContent(username) && (
              <HelperText>
                <HelperTextItem variant="error" data-testid="mr-username-error">
                  Username cannot be empty
                </HelperTextItem>
              </HelperText>
            )}
          </FormGroup>
          <FormGroup label="Password" isRequired fieldId="mr-password">
            <ModelRegistryDatabasePassword
              password={password || ''}
              setPassword={setPassword}
              isPasswordTouched={isPasswordTouched}
              setIsPasswordTouched={setIsPasswordTouched}
              showPassword={showPassword}
              //   editRegistry={mr}
            />
          </FormGroup>
          <FormGroup label="Database" isRequired fieldId="mr-database">
            <TextInput
              isRequired
              type="text"
              id="mr-database"
              name="mr-database"
              value={database}
              onBlur={() => setIsDatabaseTouched(true)}
              onChange={(_e, value) => setDatabase(value)}
              validated={isDatabaseTouched && !hasContent(database) ? 'error' : 'default'}
            />
            {isDatabaseTouched && !hasContent(database) && (
              <HelperText>
                <HelperTextItem variant="error" data-testid="mr-database-error">
                  Database cannot be empty
                </HelperTextItem>
              </HelperText>
            )}
          </FormGroup>
          {/* {secureDbEnabled && (
            <>
              <FormGroup>
                <Checkbox
                  label="Add CA certificate to secure database connection"
                  isChecked={addSecureDB}
                  onChange={(_e, value) => setAddSecureDB(value)}
                  id="add-secure-db"
                  data-testid="add-secure-db-mr-checkbox"
                  name="add-secure-db"
                />
              </FormGroup>
              {addSecureDB &&
                (!configSecretsLoaded && !configSecretsError ? (
                  <EmptyState icon={Spinner} />
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
                    isInline
                    variant="danger"
                    title="Error fetching config maps and secrets"
                    data-testid="error-fetching-resource-alert"
                  >
                    {configSecretsError?.message}
                  </Alert>
                ))}
            </>
          )} */}
        </FormSection>
      </Form>
    </Modal>
  );
};

export default CreateModal;
