export const POLL_INTERVAL = process.env.POLL_INTERVAL ? parseInt(process.env.POLL_INTERVAL) : 30000;
export const FAST_POLL_INTERVAL = process.env.FAST_POLL_INTERVAL
  ? parseInt(process.env.FAST_POLL_INTERVAL)
  : 3000;
export const SERVER_TIMEOUT = process.env.SERVER_TIMEOUT ? parseInt(process.env.SERVER_TIMEOUT) : 300000; // 5 minutes 