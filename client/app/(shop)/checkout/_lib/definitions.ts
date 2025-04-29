import { z } from 'zod';

export const CheckoutFormSchema = z.object({
  email: z.string().email().optional(),
  fullname: z.string().optional(),
  address_id: z.number().optional(),
  address: z
    .object({
      street: z.string().min(1, { message: 'Address is required' }),
      city: z.string().min(1, { message: 'City is required' }),
      district: z.string().min(1, { message: 'District is required' }),
      ward: z.string().optional(),
      phone: z.string().min(1, { message: 'Phone number is required' }),
    })
    .optional(),
  payment_method: z.enum(['stripe', 'cod']),
  payment_receipt_email: z.string().optional(),
  terms_accepted: z.boolean().refine((val) => val === true, {
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

export type CheckoutDataResponse = {
  order_id: string;
  payment_id?: string;
  client_secret?: string;
};
