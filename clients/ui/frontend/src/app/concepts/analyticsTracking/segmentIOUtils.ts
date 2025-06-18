import { TrackingOutcome } from '~/app/concepts/analyticsTracking/trackingProperties';

export const fireFormTrackingEvent = (
  _eventName: string,
  _properties: { outcome: TrackingOutcome; success?: boolean; error?: string },
): void => {
  // no-op
};
