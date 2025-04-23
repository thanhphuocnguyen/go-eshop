import { z } from 'zod';

export const CheckoutFormSchema = z.object({
  email: z.string().email(),
  firstName: z.string().min(1, { message: 'First name is required' }),
  lastName: z.string().min(1, { message: 'Last name is required' }),
  address: z.string().min(1, { message: 'Address is required' }),
  city: z.string().min(1, { message: 'City is required' }),
  state: z.string().min(1, { message: 'State is required' }),
  zip: z.string().min(1, { message: 'Zip code is required' }),
  country: z.string().min(1, { message: 'Country is required' }),
  phone: z.string().min(1, { message: 'Phone number is required' }),
  paymentMethod: z.enum(['credit_card', 'paypal']),
  cardNumber: z.string().optional(),
  cardExpiry: z.string().optional(),
  cardCvc: z.string().optional(),
  paypalEmail: z.string().optional(),
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
