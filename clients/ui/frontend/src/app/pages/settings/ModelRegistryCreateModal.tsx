import * as React from 'react';
import {
  Button,
  Form,
  FormGroup,
  HelperText,
  HelperTextItem,
  TextInput,
  Alert,
} from '@patternfly/react-core';
import { Modal } from '@patternfly/react-core/deprecated';
import { useNavigate } from 'react-router';
import { FormSection } from 'mod-arch-shared';
import ModelRegistryCreateModalFooter from '~/app/pages/settings/ModelRegistryCreateModalFooter';
import ModelRegistryDatabasePassword from '~/app/pages/settings/ModelRegistryDatabasePassword';
import K8sNameDescriptionField from '~/concepts/k8s/K8sNameDescriptionField/K8sNameDescriptionField';
import ThemeAwareFormGroupWrapper from './components/ThemeAwareFormGroupWrapper';

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
  const [isHostTouched, setIsHostTouched] = React.useState(false);
  const [isPortTouched, setIsPortTouched] = React.useState(false);
  const [isUsernameTouched, setIsUsernameTouched] = React.useState(false);
  const [isPasswordTouched, setIsPasswordTouched] = React.useState(false);
  const [isDatabaseTouched, setIsDatabaseTouched] = React.useState(false);
  const [showPassword, setShowPassword] = React.useState(false);

  const navigate = useNavigate();

  const onBeforeClose = () => {
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
    hasContent(host) &&
    hasContent(password) &&
    hasContent(port) &&
    hasContent(username) &&
    hasContent(database);

  const onSubmit = () => {
    navigate(`/model-registry-settings`);
    onClose();
  };

  const hostInput = (
    <TextInput
      isRequired
      type="text"
      id="mr-host"
      name="mr-host"
      value={host}
      onBlur={() => setIsHostTouched(true)}
      onChange={(_e, value) => setHost(value)}
    />
  );

  const hostHelperText = isHostTouched && !hasContent(host) && (
    <HelperText>
      <HelperTextItem variant="error" data-testid="mr-host-error">
        Host cannot be empty
      </HelperTextItem>
    </HelperText>
  );

  const portInput = (
    <TextInput
      isRequired
      type="text"
      id="mr-port"
      name="mr-port"
      value={port}
      onBlur={() => setIsPortTouched(true)}
      onChange={(_e, value) => setPort(value)}
    />
  );

  const portHelperText = isPortTouched && !hasContent(port) && (
    <HelperText>
      <HelperTextItem variant="error" data-testid="mr-port-error">
        Port cannot be empty
      </HelperTextItem>
    </HelperText>
  );

  const userNameInput = (
    <TextInput
      isRequired
      type="text"
      id="mr-username"
      name="mr-username"
      value={username}
      onBlur={() => setIsUsernameTouched(true)}
      onChange={(_e, value) => setUsername(value)}
    />
  );

  const usernameHelperText = isUsernameTouched && !hasContent(username) && (
    <HelperText>
      <HelperTextItem variant="error" data-testid="mr-username-error">
        Username cannot be empty
      </HelperTextItem>
    </HelperText>
  );

  const passwordInput = (
    <ModelRegistryDatabasePassword
      password={password || ''}
      setPassword={setPassword}
      isPasswordTouched={isPasswordTouched}
      setIsPasswordTouched={setIsPasswordTouched}
      showPassword={showPassword}
    />
  );

  const passwordHelperText = isPasswordTouched && !hasContent(password) && (
    <HelperText>
      <HelperTextItem variant="error" data-testid="mr-password-error">
        Password cannot be empty
      </HelperTextItem>
    </HelperText>
  );

  const databaseInput = (
    <TextInput
      isRequired
      type="text"
      id="mr-database"
      name="mr-database"
      value={database}
      onBlur={() => setIsDatabaseTouched(true)}
      onChange={(_e, value) => setDatabase(value)}
    />
  );

  const databaseHelperText = isDatabaseTouched && !hasContent(database) && (
    <HelperText>
      <HelperTextItem variant="error" data-testid="mr-database-error">
        Database cannot be empty
      </HelperTextItem>
    </HelperText>
  );

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
          <ThemeAwareFormGroupWrapper
            label="Host"
            fieldId="mr-host"
            isRequired
            helperTextNode={hostHelperText}
          >
            {hostInput}
          </ThemeAwareFormGroupWrapper>

          <ThemeAwareFormGroupWrapper
            label="Port"
            fieldId="mr-port"
            isRequired
            helperTextNode={portHelperText}
          >
            {portInput}
          </ThemeAwareFormGroupWrapper>

          <ThemeAwareFormGroupWrapper
            label="Username"
            fieldId="mr-username"
            isRequired
            helperTextNode={usernameHelperText}
          >
            {userNameInput}
          </ThemeAwareFormGroupWrapper>

          <ThemeAwareFormGroupWrapper
            label="Password"
            fieldId="mr-password"
            isRequired
            helperTextNode={passwordHelperText}
          >
            {passwordInput}
          </ThemeAwareFormGroupWrapper>

          <ThemeAwareFormGroupWrapper
            label="Database"
            fieldId="mr-database"
            isRequired
            helperTextNode={databaseHelperText}
          >
            {databaseInput}
          </ThemeAwareFormGroupWrapper>

          {/* ... Optional TLS section ... */}
        </FormSection>

        {error && (
          <FormGroup>
            <Alert variant="danger" isInline title={error.message} data-testid="mr-error" />
          </FormGroup>
        )}
      </Form>
    </Modal>
  );
};

export default CreateModal;
