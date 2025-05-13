import { z } from 'zod';

export const BaseOptionSchema = z.object({
  id: z.string(),
  name: z.string(),
});

export enum OrderStatus {
  Pending = 'pending',
  Confirmed = 'confirmed',
  Delivering = 'delivering',
  Delivered = 'delivered',
  Cancelled = 'cancelled',
  Refunded = 'refunded',
  Completed = 'completed',
}
