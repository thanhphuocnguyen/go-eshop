import { z } from 'zod';

export const CheckoutFormSchema = z.object({
  email: z.string().email(),
  fullname: z.string().min(1, { message: 'Full name is required' }),
  address: z.string().min(1, { message: 'Address is required' }),
  city: z.string().min(1, { message: 'City is required' }),
  district: z.string().min(1, { message: 'District is required' }),
  ward: z.string().optional(),
  phone: z.string().min(1, { message: 'Phone number is required' }),
  paymentMethod: z.enum(['credit_card', 'paypal', 'stripe']),
  cardNumber: z.string().optional(),
  cardExpiry: z.string().optional(),
  cardCvc: z.string().optional(),
  paypalEmail: z.string().optional(),
  stripeEmail: z.string().optional(),
  termsAccepted: z.boolean().refine((val) => val === true, {
    message: 'You must accept the terms and conditions',
  }),
});

export type CheckoutFormValues = z.infer<typeof CheckoutFormSchema>;
export type CheckoutFormErrors = Partial<
  Record<keyof CheckoutFormValues, string>
>;
export type CheckoutFormProps = {
  onSubmit: (values: CheckoutFormValues) => void;
  onError: (errors: CheckoutFormErrors) => void;
  initialValues?: Partial<CheckoutFormValues>;
  isLoading?: boolean;
};
