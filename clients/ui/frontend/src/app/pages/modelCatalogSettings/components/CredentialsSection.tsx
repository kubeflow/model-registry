import * as React from 'react';
import {
  FormSection,
  FormGroup,
  TextInput,
  FormHelperText,
  HelperText,
  HelperTextItem,
  InputGroup,
  InputGroupItem,
  Button,
} from '@patternfly/react-core';
import { EyeIcon, EyeSlashIcon } from '@patternfly/react-icons';
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
  touched: Record<string, boolean>;
  onDataChange: (key: keyof ManageSourceFormData, value: string) => void;
  onFieldBlur: (field: string) => void;
};

const CredentialsSection: React.FC<CredentialsSectionProps> = ({
  formData,
  touched,
  onDataChange,
  onFieldBlur,
}) => {
  const [isPasswordVisible, setIsPasswordVisible] = React.useState(false);
  const isOrganizationValid = validateOrganization(formData.organization);
  const isAccessTokenValid = validateAccessToken(formData.accessToken);

  return (
    <FormSection title={FORM_LABELS.CREDENTIALS} data-testid="credentials-section">
      <FormGroup label={FORM_LABELS.ORGANIZATION} isRequired fieldId="organization">
        <FormHelperText>
          <HelperText>
            <HelperTextItem>{HELP_TEXT.ORGANIZATION}</HelperTextItem>
          </HelperText>
        </FormHelperText>
        <TextInput
          isRequired
          type="text"
          id="organization"
          name="organization"
          data-testid="organization-input"
          placeholder={PLACEHOLDERS.ORGANIZATION}
          value={formData.organization}
          onChange={(_event, value) => onDataChange('organization', value)}
          onBlur={() => onFieldBlur('organization')}
          validated={touched.organization && !isOrganizationValid ? 'error' : 'default'}
        />
        {touched.organization && !isOrganizationValid && (
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
        <InputGroup>
          <InputGroupItem isFill>
            <TextInput
              isRequired
              type={isPasswordVisible ? 'text' : 'password'}
              id="access-token"
              name="access-token"
              data-testid="access-token-input"
              value={formData.accessToken}
              onChange={(_event, value) => onDataChange('accessToken', value)}
              onBlur={() => onFieldBlur('accessToken')}
              validated={touched.accessToken && !isAccessTokenValid ? 'error' : 'default'}
            />
          </InputGroupItem>
          <InputGroupItem>
            <Button
              variant="control"
              onClick={() => setIsPasswordVisible(!isPasswordVisible)}
              aria-label={isPasswordVisible ? 'Hide access token' : 'Show access token'}
              data-testid="access-token-toggle-visibility"
            >
              {isPasswordVisible ? <EyeSlashIcon /> : <EyeIcon />}
            </Button>
          </InputGroupItem>
        </InputGroup>
        {touched.accessToken && !isAccessTokenValid && (
          <FormHelperText>
            <HelperText>
              <HelperTextItem variant="error" data-testid="access-token-error">
                {VALIDATION_MESSAGES.ACCESS_TOKEN_REQUIRED}
              </HelperTextItem>
            </HelperText>
          </FormHelperText>
        )}
      </FormGroup>
    </FormSection>
  );
};

export default CredentialsSection;
