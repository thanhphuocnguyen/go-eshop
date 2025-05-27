'use client';

import { TextField } from '@/app/components/FormFields';
import { Button } from '@headlessui/react';
import { UseFormRegister } from 'react-hook-form';
import { CheckoutFormValues } from '../_lib/definitions';
import { UserModel } from '@/app/lib/definitions';

interface ShippingInformationSectionProps {
  register: UseFormRegister<CheckoutFormValues>;
  user: UserModel | undefined;
  isNewAddress: boolean;
  selectedAddressId: string | null;
  onAddressChange: (e: React.ChangeEvent<HTMLSelectElement>) => void;
  onAddNewAddress: () => void;
}

export const ShippingInformationSection: React.FC<
  ShippingInformationSectionProps
> = ({
  register,
  user,
  isNewAddress,
  selectedAddressId,
  onAddressChange,
  onAddNewAddress,
}) => {
  return (
    <div>
      <hr className='my-8' />
      <h4 className='text-lg font-semibold text-gray-600'>
        Shipping Information
      </h4>

      {user?.addresses && user.addresses.length > 0 && (
        <div className='mt-4 mb-6'>
          <label className='block text-sm font-medium text-gray-700 mb-1'>
            Select Address
          </label>
          <div className='flex gap-4'>
            <select
              className='flex-1 py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 text-sm'
              value={selectedAddressId === null ? -1 : selectedAddressId}
              onChange={onAddressChange}
            >
              {user.addresses.map((address) => (
                <option key={address.id} value={address.id}>
                  {address.street}, {address.district}, {address.city}
                  {address.default ? ' (Default)' : ''}
                </option>
              ))}
              <option value='-1'>+ Add new address</option>
            </select>

            {!isNewAddress && (
              <Button
                type='button'
                onClick={onAddNewAddress}
                className='bg-white border border-gray-300 rounded-md shadow-sm px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50'
              >
                Add New
              </Button>
            )}
          </div>
        </div>
      )}

      <div className='grid grid-cols-2 gap-6 mt-4'>
        <TextField {...register('fullname')} type='text' label='Full name' />
        <TextField
          className=''
          {...register('address.street')}
          type='text'
          label='Street address'
          placeholder='Street address'
          disabled={!isNewAddress}
        />
        <TextField
          {...register('address.city')}
          type='text'
          label='City'
          placeholder='City'
          disabled={!isNewAddress}
        />
        <TextField
          {...register('address.district')}
          type='text'
          label='District'
          placeholder='District'
          disabled={!isNewAddress}
        />
        <TextField
          {...register('address.ward')}
          type='text'
          label='Ward'
          placeholder=''
          disabled={!isNewAddress}
        />
        <TextField
          {...register('address.phone')}
          type='phone'
          className=''
          label='Phone number'
          placeholder='Phone number'
          disabled={!isNewAddress}
        />
      </div>
    </div>
  );
};
