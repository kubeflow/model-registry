import * as React from 'react';
import {
  Alert,
  Button,
  Form,
  FormGroup,
  FormHelperText,
  HelperText,
  HelperTextItem,
  MenuToggle,
  Modal,
  ModalBody,
  ModalFooter,
  ModalHeader,
  ModalVariant,
  Radio,
  Select,
  SelectList,
  SelectOption,
  TextInput,
} from '@patternfly/react-core';
import { useNavigate } from 'react-router';
import { FormSection } from 'mod-arch-shared';
import { createModelRegistrySettings } from '~/app/api/k8s';
import ModelRegistryDatabasePassword from '~/app/pages/settings/ModelRegistryDatabasePassword';
import K8sNameDescriptionField from '~/concepts/k8s/K8sNameDescriptionField/K8sNameDescriptionField';
import ThemeAwareFormGroupWrapper from '~/app/pages/settings/components/ThemeAwareFormGroupWrapper';
import { DatabaseType, DatabaseMode, ModelRegistryPayload } from '~/app/types';

type NameDescType = {
  name: string;
  description: string;
};

type CreateModalProps = {
  onClose: () => void;
  refresh: () => void;
};

const CreateModal: React.FC<CreateModalProps> = ({ onClose, refresh }) => {
  const [error, setError] = React.useState<Error>();
  const [nameDesc, setNameDesc] = React.useState<NameDescType>({
    name: '',
    description: '',
  });
  const [databaseMode, setDatabaseMode] = React.useState<DatabaseMode>(DatabaseMode.External);
  const [databaseType, setDatabaseType] = React.useState<DatabaseType>(DatabaseType.MySQL);
  const [databaseTypeSelectOpen, setDatabaseTypeSelectOpen] = React.useState(false);

  // Common fields
  const [database, setDatabase] = React.useState('');
  const [isDatabaseTouched, setIsDatabaseTouched] = React.useState(false);

  // External database fields
  const [host, setHost] = React.useState('');
  const [port, setPort] = React.useState('3306'); // Default MySQL port
  const [username, setUsername] = React.useState('');
  const [password, setPassword] = React.useState('');
  const [isHostTouched, setIsHostTouched] = React.useState(false);
  const [isPortTouched, setIsPortTouched] = React.useState(false);
  const [isUsernameTouched, setIsUsernameTouched] = React.useState(false);
  const [isPasswordTouched, setIsPasswordTouched] = React.useState(false);
  const [showPassword, setShowPassword] = React.useState(false);

  const navigate = useNavigate();

  // Auto-fill port when database type changes (only if user hasn't manually changed it)
  React.useEffect(() => {
    if (databaseMode === DatabaseMode.External && !isPortTouched) {
      setPort(databaseType === DatabaseType.MySQL ? '3306' : '5432');
    }
  }, [databaseType, databaseMode, isPortTouched]);

  const onBeforeClose = () => {
    setError(undefined);
    setNameDesc({ name: '', description: '' });
    setDatabaseMode(DatabaseMode.External);
    setDatabaseType(DatabaseType.MySQL);
    setHost('');
    setPort('3306'); // Reset to MySQL default
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

  const canSubmit = () => {
    // Name is always required
    if (!hasContent(nameDesc.name)) {
      return false;
    }

    // For external databases, all connection fields are required including database name
    if (databaseMode === DatabaseMode.External) {
      return (
        hasContent(host) &&
        hasContent(port) &&
        hasContent(username) &&
        hasContent(password) &&
        hasContent(database)
      );
    }

    // For default database, only name is required
    return true;
  };

  const onSubmit = async () => {
    setError(undefined);

    let payload: ModelRegistryPayload;

    if (databaseMode === DatabaseMode.Default) {
      // Default database with generateDeployment
      payload = {
        modelRegistry: {
          metadata: {
            name: nameDesc.name,
            annotations: {
              'openshift.io/display-name': nameDesc.name,
              'openshift.io/description': nameDesc.description,
            },
          },
          spec: {
            postgres: {
              database: 'model_registry',
              generateDeployment: true,
            },
          },
        },
      };
    } else if (databaseType === DatabaseType.MySQL) {
      // External MySQL database
      payload = {
        modelRegistry: {
          metadata: {
            name: nameDesc.name,
            annotations: {
              'openshift.io/display-name': nameDesc.name,
              'openshift.io/description': nameDesc.description,
            },
          },
          spec: {
            mysql: {
              host,
              port: port ? Number(port) : undefined,
              username,
              database,
            },
          },
        },
        databasePassword: password,
      };
    } else {
      // External PostgreSQL database
      payload = {
        modelRegistry: {
          metadata: {
            name: nameDesc.name,
            annotations: {
              'openshift.io/display-name': nameDesc.name,
              'openshift.io/description': nameDesc.description,
            },
          },
          spec: {
            postgres: {
              host,
              port: port ? Number(port) : undefined,
              username,
              database,
            },
          },
        },
        databasePassword: password,
      };
    }

    try {
      await createModelRegistrySettings(window.location.origin, {
        namespace: 'model-registry',
      })({}, payload);
      refresh();
      navigate(`/model-registry-settings`);
      onClose();
    } catch (e) {
      if (e instanceof Error) {
        setError(e);
      }
    }
  };

  const databaseTypeSelectToggle = (toggleRef: React.Ref<HTMLButtonElement>) => (
    <MenuToggle
      ref={toggleRef}
      onClick={() => setDatabaseTypeSelectOpen(!databaseTypeSelectOpen)}
      isExpanded={databaseTypeSelectOpen}
    >
      {databaseType === DatabaseType.MySQL ? 'MySQL' : 'PostgreSQL'}
    </MenuToggle>
  );

  const isExternalMode = databaseMode === DatabaseMode.External;

  return (
    <Modal
      isOpen
      variant={ModalVariant.medium}
      onClose={onBeforeClose}
      data-testid="create-model-registry-modal"
    >
      <ModalHeader title="Create model registry" />
      <ModalBody>
        <Form>
          <K8sNameDescriptionField dataTestId="mr" data={nameDesc} onDataChange={setNameDesc} />

          {/* Database Mode Selection */}
          <FormSection
            title="Database"
            description="Choose where to store model data."
            data-testid="mr-database-section"
          >
            <FormGroup role="radiogroup">
              <Radio
                id="mr-database-mode-default"
                name="mr-database-mode"
                label="Default database (non-production)"
                description="PostgreSQL database enabled by default on the cluster."
                isChecked={databaseMode === DatabaseMode.Default}
                onChange={() => setDatabaseMode(DatabaseMode.Default)}
                data-testid="mr-database-mode-default"
                body={
                  databaseMode === DatabaseMode.Default && (
                    <Alert
                      variant="info"
                      isInline
                      isPlain
                      title="This default database is for development and testing purposes only. It is not supported by Red Hat for production use cases."
                      data-testid="mr-database-mode-default-warning"
                    />
                  )
                }
              />
              <Radio
                id="mr-database-mode-external"
                name="mr-database-mode"
                label="External database"
                description="Connect a MySQL or PostgreSQL database."
                isChecked={databaseMode === DatabaseMode.External}
                onChange={() => setDatabaseMode(DatabaseMode.External)}
                data-testid="mr-database-mode-external"
              />
            </FormGroup>
          </FormSection>

          {/* External Database Configuration */}
          {isExternalMode && (
            <FormSection
              title={`Connect to external ${databaseType === DatabaseType.MySQL ? 'MySQL' : 'PostgreSQL'} database`}
              description="This external database is where model data is stored."
              data-testid="mr-external-database-section"
            >
              {/* Database Type Selection */}
              <ThemeAwareFormGroupWrapper
                label="Database type"
                fieldId="mr-database-type"
                isRequired
                data-testid="mr-database-type-group"
              >
                <Select
                  id="mr-database-type"
                  isOpen={databaseTypeSelectOpen}
                  selected={databaseType}
                  onSelect={(_event, value) => {
                    if (value === DatabaseType.MySQL) {
                      setDatabaseType(DatabaseType.MySQL);
                    } else if (value === DatabaseType.PostgreSQL) {
                      setDatabaseType(DatabaseType.PostgreSQL);
                    } else if (value === 'mysql') {
                      setDatabaseType(DatabaseType.MySQL);
                    } else if (value === 'postgres') {
                      setDatabaseType(DatabaseType.PostgreSQL);
                    }
                    setDatabaseTypeSelectOpen(false);
                  }}
                  onOpenChange={setDatabaseTypeSelectOpen}
                  toggle={databaseTypeSelectToggle}
                  data-testid="mr-database-type-select"
                >
                  <SelectList>
                    <SelectOption value={DatabaseType.MySQL} data-testid="mr-database-type-mysql">
                      MySQL
                    </SelectOption>
                    <SelectOption
                      value={DatabaseType.PostgreSQL}
                      data-testid="mr-database-type-postgresql"
                    >
                      PostgreSQL
                    </SelectOption>
                  </SelectList>
                </Select>
              </ThemeAwareFormGroupWrapper>

              {/* Host */}
              <ThemeAwareFormGroupWrapper
                label="Host"
                fieldId="mr-host"
                isRequired
                helperTextNode={
                  isHostTouched &&
                  !hasContent(host) && (
                    <FormHelperText>
                      <HelperText>
                        <HelperTextItem variant="error" data-testid="mr-host-error">
                          Host cannot be empty
                        </HelperTextItem>
                      </HelperText>
                    </FormHelperText>
                  )
                }
              >
                <TextInput
                  isRequired
                  type="text"
                  id="mr-host"
                  name="mr-host"
                  value={host}
                  onBlur={() => setIsHostTouched(true)}
                  onChange={(_e, value) => setHost(value)}
                  validated={isHostTouched && !hasContent(host) ? 'error' : 'default'}
                  data-testid="mr-host-input"
                />
              </ThemeAwareFormGroupWrapper>

              {/* Port */}
              <ThemeAwareFormGroupWrapper
                label="Port"
                fieldId="mr-port"
                isRequired
                helperTextNode={
                  isPortTouched &&
                  !hasContent(port) && (
                    <FormHelperText>
                      <HelperText>
                        <HelperTextItem variant="error" data-testid="mr-port-error">
                          Port cannot be empty
                        </HelperTextItem>
                      </HelperText>
                    </FormHelperText>
                  )
                }
              >
                <TextInput
                  isRequired
                  type="text"
                  id="mr-port"
                  name="mr-port"
                  value={port}
                  onBlur={() => setIsPortTouched(true)}
                  onChange={(_e, value) => setPort(value)}
                  validated={isPortTouched && !hasContent(port) ? 'error' : 'default'}
                  data-testid="mr-port-input"
                />
              </ThemeAwareFormGroupWrapper>

              {/* Username */}
              <ThemeAwareFormGroupWrapper
                label="Username"
                fieldId="mr-username"
                isRequired
                helperTextNode={
                  isUsernameTouched &&
                  !hasContent(username) && (
                    <FormHelperText>
                      <HelperText>
                        <HelperTextItem variant="error" data-testid="mr-username-error">
                          Username cannot be empty
                        </HelperTextItem>
                      </HelperText>
                    </FormHelperText>
                  )
                }
              >
                <TextInput
                  isRequired
                  type="text"
                  id="mr-username"
                  name="mr-username"
                  value={username}
                  onBlur={() => setIsUsernameTouched(true)}
                  onChange={(_e, value) => setUsername(value)}
                  validated={isUsernameTouched && !hasContent(username) ? 'error' : 'default'}
                  data-testid="mr-username-input"
                />
              </ThemeAwareFormGroupWrapper>

              {/* Password */}
              <ThemeAwareFormGroupWrapper
                label="Password"
                fieldId="mr-password"
                isRequired
                helperTextNode={
                  isPasswordTouched &&
                  !hasContent(password) && (
                    <FormHelperText>
                      <HelperText>
                        <HelperTextItem variant="error" data-testid="mr-password-error">
                          Password cannot be empty
                        </HelperTextItem>
                      </HelperText>
                    </FormHelperText>
                  )
                }
              >
                <ModelRegistryDatabasePassword
                  password={password || ''}
                  setPassword={setPassword}
                  isPasswordTouched={isPasswordTouched}
                  setIsPasswordTouched={setIsPasswordTouched}
                  showPassword={showPassword}
                />
              </ThemeAwareFormGroupWrapper>
            </FormSection>
          )}

          {/* Database Name - Only for external mode */}
          {isExternalMode && (
            <FormSection
              title="Database settings"
              description="Configure database-specific settings."
              data-testid="mr-database-settings-section"
            >
              <ThemeAwareFormGroupWrapper
                label="Database"
                fieldId="mr-database"
                isRequired
                helperTextNode={
                  isDatabaseTouched &&
                  !hasContent(database) && (
                    <FormHelperText>
                      <HelperText>
                        <HelperTextItem variant="error" data-testid="mr-database-error">
                          Database cannot be empty
                        </HelperTextItem>
                      </HelperText>
                    </FormHelperText>
                  )
                }
              >
                <TextInput
                  isRequired
                  type="text"
                  id="mr-database"
                  name="mr-database"
                  value={database}
                  onBlur={() => setIsDatabaseTouched(true)}
                  onChange={(_e, value) => setDatabase(value)}
                  validated={isDatabaseTouched && !hasContent(database) ? 'error' : 'default'}
                  data-testid="mr-database-input"
                />
              </ThemeAwareFormGroupWrapper>
            </FormSection>
          )}

          {error && (
            <FormGroup>
              <Alert variant="danger" isInline title={error.message} data-testid="mr-error" />
            </FormGroup>
          )}
        </Form>
      </ModalBody>
      <ModalFooter>
        <Button
          key="create-button"
          variant="primary"
          isDisabled={!canSubmit()}
          onClick={onSubmit}
          data-testid="mr-create-button"
        >
          Create
        </Button>
        <Button
          key="cancel-button"
          variant="secondary"
          onClick={onBeforeClose}
          data-testid="mr-cancel-button"
        >
          Cancel
        </Button>
      </ModalFooter>
    </Modal>
  );
};

export default CreateModal;
