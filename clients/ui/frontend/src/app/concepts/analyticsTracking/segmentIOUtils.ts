import { TrackingOutcome } from '~/app/concepts/analyticsTracking/trackingProperties';

export const fireFormTrackingEvent = (eventName: string, properties: { outcome: TrackingOutcome, success?: boolean, error?: string }) => {
    // no-op
}; 