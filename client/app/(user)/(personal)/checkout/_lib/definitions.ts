import { z } from 'zod';

export const CheckoutFormSchema = z.object({
  email: z.string().email().optional(),
  fullname: z.string().optional(),
  addressId: z.string().optional(),
  address: z
    .object({
      street: z.string().min(1, { message: 'Address is required' }),
      city: z.string().min(1, { message: 'City is required' }),
      district: z.string().min(1, { message: 'District is required' }),
      ward: z.string().optional(),
      phone: z.string().min(1, { message: 'Phone number is required' }),
    })
    .optional(),
  paymentMethod: z.enum(['stripe', 'cod']),
  discountCode: z.string().optional(),
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

export type CheckoutDataResponse = {
  orderId: string;
  totalPrice: number;
  paymentId?: string;
  clientSecret?: string;
  paymentIntentId?: string;
};

export type PaymentResponse = {
  id: string;
  gateway: string;
  status: string;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  details: Record<string, any>;
};
