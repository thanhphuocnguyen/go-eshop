import { z } from 'zod';

export const BaseOptionSchema = z.object({
  id: z.number().or(z.string()),
  name: z.string(),
});
