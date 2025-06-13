import React from 'react';
import {
  Container,
  Typography,
  Box,
  CircularProgress,
  Breadcrumbs,
  Stack,
} from '@mui/material';
import { Help, Error } from '@mui/icons-material';

type ApplicationsPageProps = {
  title?: React.ReactNode;
  breadcrumb?: React.ReactNode;
  description?: React.ReactNode;
  loaded: boolean;
  empty: boolean;
  loadError?: globalThis.Error;
  loadErrorPage?: React.ReactNode;
  children?: React.ReactNode;
  errorMessage?: string;
  emptyMessage?: string;
  emptyStatePage?: React.ReactNode;
  headerAction?: React.ReactNode;
  headerContent?: React.ReactNode;
  provideChildrenPadding?: boolean;
  removeChildrenTopPadding?: boolean;
  subtext?: React.ReactNode;
  loadingContent?: React.ReactNode;
  noHeader?: boolean;
};

const ApplicationsPage: React.FC<ApplicationsPageProps> = ({
  title,
  breadcrumb,
  description,
  loaded,
  empty,
  loadError,
  loadErrorPage,
  children,
  errorMessage,
  emptyMessage,
  emptyStatePage,
  headerAction,
  headerContent,
  provideChildrenPadding,
  removeChildrenTopPadding,
  subtext,
  loadingContent,
  noHeader,
}) => {
  const renderHeader = () => (
    <Box sx={{ mb: 2 }}>
      <Stack spacing={2}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
          <Box sx={{ flexGrow: 1 }}>
            <Typography variant="h4" component="h1" data-testid="app-page-title">
              {title}
            </Typography>
            <Stack spacing={1}>
              {subtext && <Typography variant="subtitle1">{subtext}</Typography>}
              {description && <Typography variant="body1">{description}</Typography>}
            </Stack>
          </Box>
          <Box>{headerAction}</Box>
        </Box>
        {headerContent && <Box>{headerContent}</Box>}
      </Stack>
    </Box>
  );

  const renderContents = () => {
    if (loadError) {
      return !loadErrorPage ? (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <Error color="error" sx={{ fontSize: 48 }} />
          <Typography variant="h5" component="h1" sx={{ mt: 2 }}>
            {errorMessage !== undefined ? errorMessage : 'Error loading components'}
          </Typography>
          <Typography variant="body1" sx={{ mt: 1 }}>
            {loadError.message}
          </Typography>
        </Box>
      ) : (
        loadErrorPage
      );
    }

    if (!loaded) {
      return (
        loadingContent || (
          <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', py: 4 }}>
            <CircularProgress />
          </Box>
        )
      );
    }

    if (empty) {
      return !emptyStatePage ? (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <Help color="disabled" sx={{ fontSize: 48 }} />
          <Typography variant="h5" component="h1" sx={{ mt: 2 }}>
            {emptyMessage !== undefined ? emptyMessage : 'No Components Found'}
          </Typography>
        </Box>
      ) : (
        emptyStatePage
      );
    }

    if (provideChildrenPadding) {
      return (
        <Container sx={removeChildrenTopPadding ? { pt: 0 } : {}}>
          {children}
        </Container>
      );
    }

    return children;
  };

  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', height: '100%' }}>
      {breadcrumb && <Breadcrumbs sx={{ mb: 2 }}>{breadcrumb}</Breadcrumbs>}
      {!noHeader && renderHeader()}
      <Box sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
        {renderContents()}
      </Box>
    </Box>
  );
};

export default ApplicationsPage; 