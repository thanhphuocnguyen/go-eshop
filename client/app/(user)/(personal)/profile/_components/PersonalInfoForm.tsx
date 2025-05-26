'use client';

import { TextField } from '@/app/components/FormFields';
import { Button, Fieldset } from '@headlessui/react';
import { UserModel } from '@/app/lib/definitions/user';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

// Profile form schema
const profileSchema = z.object({
  fullname: z
    .string()
    .min(3, { message: 'Full name must be at least 3 characters' }),
  email: z.string().email({ message: 'Please enter a valid email address' }),
  phone: z
    .string()
    .min(10, { message: 'Phone number must be at least 10 characters' }),
});

type ProfileFormValues = z.infer<typeof profileSchema>;

interface PersonalInfoFormProps {
  userData: UserModel | undefined;
  onUpdate: (data: ProfileFormValues) => Promise<void>;
}

export default function PersonalInfoForm({
  userData,
  onUpdate,
}: PersonalInfoFormProps) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      fullname: userData?.fullname || '',
      email: userData?.email || '',
      phone: userData?.phone || '',
    },
  });

  return (
    <Fieldset as='form' onSubmit={handleSubmit(onUpdate)} className='space-y-6'>
      <div className='grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6'>
        <div className='sm:col-span-3'>
          <TextField
            {...register('fullname')}
            label='Full name'
            placeholder='John Doe'
            error={errors.fullname?.message}
          />
        </div>
        <div className='sm:col-span-3'>
          <TextField
            {...register('phone')}
            label='Phone number'
            placeholder='+1 (555) 987-6543'
            error={errors.phone?.message}
          />
        </div>
        <div className='sm:col-span-6'>
          <TextField
            {...register('email')}
            label='Email address'
            placeholder='john.doe@example.com'
            error={errors.email?.message}
          />
        </div>
      </div>
      <div className='flex justify-end'>
        <Button
          type='submit'
          className='px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700'
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Saving...' : 'Save Changes'}
        </Button>
      </div>
    </Fieldset>
  );
}
