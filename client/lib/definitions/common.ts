import { z } from 'zod';

export const BaseOptionSchema = z.object({
  id: z.string(),
  name: z.string(),
});
