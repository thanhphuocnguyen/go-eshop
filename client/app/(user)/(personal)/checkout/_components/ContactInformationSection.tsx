'use client';

import { TextField } from '@/app/components/FormFields';
import { Fieldset } from '@headlessui/react';
import { UseFormRegister } from 'react-hook-form';
import { CheckoutFormValues } from '../_lib/definitions';

interface ContactInformationSectionProps {
  register: UseFormRegister<CheckoutFormValues>;
}

export const ContactInformationSection: React.FC<ContactInformationSectionProps> = ({
  register,
}) => {
  return (
    <div>
      <h3 className='text-lg font-semibold text-gray-600 mb-4'>
        Contact Information
      </h3>
      <Fieldset>
        <TextField
          {...register('email')}
          type='email'
          label='Email address'
        />
      </Fieldset>
    </div>
  );
};
