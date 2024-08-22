import * as React from 'react';
import {CubesIcon} from '@patternfly/react-icons';
import {
    Button,
    Content,
    ContentVariants,
    EmptyState,
    EmptyStateBody,
    EmptyStateFooter,
    EmptyStateVariant,
    PageSection,
} from '@patternfly/react-core';

const Admin: React.FunctionComponent = () => (
    <PageSection>
        <EmptyState variant={EmptyStateVariant.full} titleText="Empty State (Stub Admin Module)" icon={CubesIcon}>
            <EmptyStateBody>
                <Content component={ContentVariants.p}>
                    This represents an the empty state pattern in Patternfly 6. Hopefully it&apos;s simple enough to use
                    but
                    flexible enough to meet a variety of needs.
                </Content>
            </EmptyStateBody><EmptyStateFooter>
            <Button variant="primary">Primary Action</Button>

        </EmptyStateFooter></EmptyState>
    </PageSection>
);

export {Admin};