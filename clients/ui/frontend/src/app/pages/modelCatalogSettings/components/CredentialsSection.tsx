import * as React from 'react';
import {
  FormGroup,
  TextInput,
  FormHelperText,
  HelperText,
  HelperTextItem,
  Button,
  ActionList,
  Alert,
  AlertActionCloseButton,
} from '@patternfly/react-core';
import { UpdateObjectAtPropAndValue } from 'mod-arch-shared';
import PasswordInput from '~/app/shared/components/PasswordInput';
import FormFieldset from '~/app/pages/modelRegistry/screens/components/FormFieldset';
import FormSection from '~/app/pages/modelRegistry/components/pf-overrides/FormSection';
import { ManageSourceFormData } from '~/app/pages/modelCatalogSettings/useManageSourceData';
import {
  validateOrganization,
  validateAccessToken,
} from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  FORM_LABELS,
  VALIDATION_MESSAGES,
  HELP_TEXT,
  PLACEHOLDERS,
} from '~/app/pages/modelCatalogSettings/constants';

type CredentialsSectionProps = {
  formData: ManageSourceFormData;
  setData: UpdateObjectAtPropAndValue<ManageSourceFormData>;
};

const CredentialsSection: React.FC<CredentialsSectionProps> = ({ formData, setData }) => {
  const [isOrganizationTouched, setIsOrganizationTouched] = React.useState(false);
  const [isAccessTokenTouched, setIsAccessTokenTouched] = React.useState(false);

  const isOrganizationValid = validateOrganization(formData.organization);
  const isAccessTokenValid = validateAccessToken(formData.accessToken);
  const [validationError, setValidationError] = React.useState<Error | undefined>(undefined);
  const [isValidating, setIsValidating] = React.useState(false);
  const [isValidationSuccess, setIsValidationSuccess] = React.useState(false);

  const handleValidate = async () => {
    // setIsValidating(true);
    // setValidationError(undefined);

    // TODO: Implement validation logic
    // setShowAlert(true);

    // if success
    setValidationError(undefined);
    setIsValidationSuccess(true);
    setIsValidating(false);

    //if fails
    // setValidationError(new Error('error'));
    // setIsValidationSuccess(false);
  };

  const organizationInput = (
    <TextInput
      isRequired
      type="text"
      id="organization"
      name="organization"
      data-testid="organization-input"
      placeholder={PLACEHOLDERS.ORGANIZATION}
      value={formData.organization}
      onChange={(_event, value) => setData('organization', value)}
      onBlur={() => setIsOrganizationTouched(true)}
      validated={isOrganizationTouched && !isOrganizationValid ? 'error' : 'default'}
    />
  );

  const accessTokenInput = (
    <PasswordInput
      isRequired
      id="access-token"
      name="access-token"
      data-testid="access-token-input"
      value={formData.accessToken}
      onChange={(_event, value) => setData('accessToken', value)}
      onBlur={() => setIsAccessTokenTouched(true)}
      validated={isAccessTokenTouched && !isAccessTokenValid ? 'error' : 'default'}
      ariaLabelShow="Show access token"
      ariaLabelHide="Hide access token"
    />
  );

  return (
    <FormSection title={FORM_LABELS.CREDENTIALS} data-testid="credentials-section">
      <FormGroup label={FORM_LABELS.ORGANIZATION} isRequired fieldId="organization">
        <FormHelperText>
          <HelperText>
            <HelperTextItem>{HELP_TEXT.ORGANIZATION}</HelperTextItem>
          </HelperText>
        </FormHelperText>
        <FormFieldset component={organizationInput} field="Allowed organization" />
        {isOrganizationTouched && !isOrganizationValid && (
          <FormHelperText>
            <HelperText>
              <HelperTextItem variant="error" data-testid="organization-error">
                {VALIDATION_MESSAGES.ORGANIZATION_REQUIRED}
              </HelperTextItem>
            </HelperText>
          </FormHelperText>
        )}
      </FormGroup>

      <FormGroup label={FORM_LABELS.ACCESS_TOKEN} isRequired fieldId="access-token">
        <FormHelperText>
          <HelperText>
            <HelperTextItem>{HELP_TEXT.ACCESS_TOKEN}</HelperTextItem>
          </HelperText>
        </FormHelperText>
        <FormFieldset component={accessTokenInput} field="Access token" />
        {isAccessTokenTouched && !isAccessTokenValid && (
          <FormHelperText>
            <HelperText>
              <HelperTextItem variant="error" data-testid="access-token-error">
                {VALIDATION_MESSAGES.ACCESS_TOKEN_REQUIRED}
              </HelperTextItem>
            </HelperText>
          </FormHelperText>
        )}
      </FormGroup>
      {validationError && (
        <Alert isInline variant="danger" title="Validation failed" className="pf-v5-u-mt-md">
          The system cannot establish a connection to the source. Ensure that the organization and
          access token are accurate, then try again.
        </Alert>
      )}
      {isValidationSuccess && (
        <Alert
          isInline
          variant="success"
          className="pf-v5-u-mt-md"
          title="Validation successful"
          actionClose={<AlertActionCloseButton onClose={() => setIsValidationSuccess(false)} />}
        >
          The organization and accessToken are valid for connection.
        </Alert>
      )}

      <ActionList className="pf-v5-u-mt-md">
        <Button
          isDisabled={!isAccessTokenValid}
          variant="link"
          onClick={handleValidate}
          isLoading={isValidating}
        >
          Validate
        </Button>
      </ActionList>
    </FormSection>
  );
};

export default CredentialsSection;
