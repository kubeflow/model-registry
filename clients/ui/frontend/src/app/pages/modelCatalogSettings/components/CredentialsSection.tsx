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
import { validateOrganization } from '~/app/pages/modelCatalogSettings/utils/validation';
import {
  FORM_LABELS,
  VALIDATION_MESSAGES,
  HELP_TEXT,
  PLACEHOLDERS,
} from '~/app/pages/modelCatalogSettings/constants';
import { TempDevFeature, useTempDevFeatureAvailable } from '~/app/hooks/useTempDevFeatureAvailable';

type CredentialsSectionProps = {
  formData: ManageSourceFormData;
  setData: UpdateObjectAtPropAndValue<ManageSourceFormData>;
  onValidate: () => Promise<void>;
  isValidating: boolean;
  validationError?: Error;
  isValidationSuccess: boolean;
  onClearValidationSuccess: () => void;
};

const CredentialsSection: React.FC<CredentialsSectionProps> = ({
  formData,
  setData,
  onValidate,
  isValidating,
  validationError,
  isValidationSuccess,
  onClearValidationSuccess,
}) => {
  const [isOrganizationTouched, setIsOrganizationTouched] = React.useState(false);

  const isOrganizationValid = validateOrganization(formData.organization);

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

  const organizationFormGroup = (
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
  );

  const accessTokenInput = (
    <PasswordInput
      isRequired
      id="access-token"
      name="access-token"
      data-testid="access-token-input"
      value={formData.accessToken}
      onChange={(_event, value) => setData('accessToken', value)}
      ariaLabelShow="Show access token"
      ariaLabelHide="Hide access token"
    />
  );

  const accessTokenFormGroup = (
    <>
      <FormGroup label={FORM_LABELS.ACCESS_TOKEN} fieldId="access-token">
        <FormHelperText>
          <HelperText>
            <HelperTextItem>{HELP_TEXT.ACCESS_TOKEN}</HelperTextItem>
          </HelperText>
        </FormHelperText>
        <FormFieldset component={accessTokenInput} field="Access token" />
      </FormGroup>
      {validationError && (
        <Alert isInline variant="danger" title="Validation failed" className="pf-v5-u-mt-md">
          {validationError.message}
        </Alert>
      )}
      {isValidationSuccess && (
        <Alert
          isInline
          variant="success"
          className="pf-v5-u-mt-md"
          title="Validation successful"
          actionClose={<AlertActionCloseButton onClose={onClearValidationSuccess} />}
        >
          The organization and accessToken are valid for connection.
        </Alert>
      )}

      <ActionList className="pf-v5-u-mt-md">
        <Button
          isDisabled={!isOrganizationValid || isValidating}
          variant="link"
          onClick={onValidate}
          isLoading={isValidating}
        >
          Validate
        </Button>
      </ActionList>
    </>
  );

  const accessTokenFeatureAvailable = useTempDevFeatureAvailable(
    TempDevFeature.CatalogHuggingFaceApiKey,
  );

  return (
    <FormSection
      title={accessTokenFeatureAvailable ? FORM_LABELS.CREDENTIALS : undefined}
      data-testid="credentials-section"
    >
      {organizationFormGroup}
      {accessTokenFeatureAvailable && accessTokenFormGroup}
    </FormSection>
  );
};

export default CredentialsSection;
