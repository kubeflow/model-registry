import path from 'path';
import { env } from 'process';
import dotenv from 'dotenv';

[
  `.env.cypress${env.CY_MOCK ? '.mock' : ''}.local`,
  `.env.cypress${env.CY_MOCK ? '.mock' : ''}`,
  '.env.test',
  '.env.local',
  '.env',
].forEach((file) =>
  dotenv.config({
    path: path.resolve(__dirname, '../../../../../', file),
  }),
);

export const BASE_URL = env.BASE_URL || '';

// re-export the updated process env
export { env };
