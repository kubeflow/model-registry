import * as React from 'react';
import { List, ListItem, Panel, PanelMain, PanelMainBody } from '@patternfly/react-core';

type EventLogEvent = {
  timestamp: string;
  type: string;
  reason: string;
  message: string;
};

type EventLogProps = {
  events: EventLogEvent[];
  emptyMessage?: string;
  maxHeight?: string;
  'data-testid'?: string;
};

const getEventFullMessage = (event: EventLogEvent): string =>
  `${event.timestamp} [${event.reason}] [${event.type}] ${event.message}`;

const EventLog: React.FC<EventLogProps> = ({
  events,
  emptyMessage = 'There are no recent events.',
  maxHeight = '300px',
  'data-testid': dataTestId = 'event-log',
}) => {
  const sortedEvents = events.toSorted(
    (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime(),
  );

  return (
    <Panel isScrollable>
      <PanelMain maxHeight={maxHeight}>
        <PanelMainBody>
          {sortedEvents.length > 0 ? (
            <List isPlain data-testid={dataTestId}>
              {sortedEvents.map((event, index) => (
                <ListItem key={`${event.timestamp}-${index}`} data-testid={`${dataTestId}-entry`}>
                  {getEventFullMessage(event)}
                </ListItem>
              ))}
            </List>
          ) : (
            <span className="pf-v6-u-color-200">{emptyMessage}</span>
          )}
        </PanelMainBody>
      </PanelMain>
    </Panel>
  );
};

export default EventLog;
