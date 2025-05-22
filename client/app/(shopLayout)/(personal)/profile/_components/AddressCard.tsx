'use client';

import { AddressModel } from '@/app/lib/definitions/user';
import { Button } from '@headlessui/react';
import { PencilIcon, TrashIcon } from '@heroicons/react/24/outline';

interface AddressCardProps {
  address: AddressModel;
  onEdit: (address: AddressModel) => void;
  onDelete: (addressId: string) => void;
  onSetDefault: (addressId: string) => void;
}

export default function AddressCard({
  address,
  onEdit,
  onDelete,
  onSetDefault,
}: AddressCardProps) {
  return (
    <div className='bg-gray-50 p-6 rounded-lg border border-gray-200'>
      <div className='flex justify-between items-start'>
        <div className='flex-1'>
          <div className='flex items-center mb-2'>
            <h3 className='text-md font-medium text-gray-900'>
              {address.street || 'Address'}
            </h3>
            {address.default && (
              <span className='ml-2 inline-flex items-center px-2.5 py-0.5 rounded-md text-xs font-medium bg-green-100 text-green-800'>
                Default
              </span>
            )}
          </div>
          <div className='text-sm text-gray-500 flex flex-col gap-4'>
            <p>{address.street}</p>
            <p>{address.city}</p>
            <p className='mt-1'>Phone: {address.phone || 'N/A'}</p>
          </div>
        </div>
        <div className='flex flex-col space-y-2'>
          <Button
            onClick={() => onEdit(address)}
            className='inline-flex items-center px-3 py-1 btn btn-primary'
          >
            <PencilIcon className='h-4 w-4 mr-1' />
            Edit
          </Button>
          <Button
            onClick={() => onDelete(address.id)}
            className='inline-flex btn btn-danger items-center px-3 py-1'
          >
            <TrashIcon className='h-4 w-4 mr-1' />
            Delete
          </Button>
          {!address.default && (
            <Button
              onClick={() => onSetDefault(address.id)}
              className='inline-flex items-center px-3 py-1 btn btn-green'
            >
              Set as default
            </Button>
          )}
        </div>
      </div>
    </div>
  );
}
