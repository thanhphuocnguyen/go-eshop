'use client';
import { useState } from 'react';
import {
  Button,
  Dialog,
  DialogPanel,
  DialogTitle,
  Fieldset,
  Label,
} from '@headlessui/react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { TextField } from '@/app/components/FormFields';
import { AddressModel } from '@/app/lib/definitions/user';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { GenericResponse } from '@/app/lib/definitions';
import { toast } from 'react-toastify';
import { useRouter } from 'next/navigation';
import { useUser } from '@/app/hooks';

// Form schema for address validation
const addressSchema = z.object({
  street: z.string().min(3, { message: 'Street address is required' }),
  city: z.string().min(2, { message: 'City is required' }),
  district: z.string().min(2, { message: 'District is required' }),
  ward: z.string().optional(),
  phone: z
    .string()
    .min(10, { message: 'Valid phone number is required' })
    .max(15),
  isDefault: z.boolean().optional(),
});

type AddressFormValues = z.infer<typeof addressSchema>;

export default function AddressesClient() {
  const { user } = useUser();
  const router = useRouter();
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedAddress, setSelectedAddress] = useState<AddressModel | null>(
    null
  );
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<AddressFormValues>({
    resolver: zodResolver(addressSchema),
    defaultValues: {
      street: '',
      city: '',
      district: '',
      ward: '',
      phone: '',
      isDefault: false,
    },
  });

  const openAddModal = () => {
    reset({
      street: '',
      city: '',
      district: '',
      ward: '',
      phone: '',
      isDefault: false,
    });
    setSelectedAddress(null);
    setIsAddModalOpen(true);
  };

  const openEditModal = (address: AddressModel) => {
    reset({
      street: address.street,
      city: address.city,
      district: address.district,
      ward: address.ward || '',
      phone: address.phone,
      isDefault: address.default,
    });
    setSelectedAddress(address);
    setIsAddModalOpen(true);
  };

  const openDeleteModal = (address: AddressModel) => {
    setSelectedAddress(address);
    setIsDeleteModalOpen(true);
  };

  const handleAddressSubmit = async (data: AddressFormValues) => {
    setIsSubmitting(true);
    try {
      if (selectedAddress) {
        // Edit existing address
        const response = await clientSideFetch<GenericResponse<AddressModel>>(
          PUBLIC_API_PATHS.USER_ADDRESS.replace(
            ':id',
            selectedAddress.id.toString()
          ),
          {
            method: 'PUT',
            body: data,
          }
        );

        if (response.error) {
          toast.error(`Failed to update address: ${response.error.details}`);
          return;
        }

        toast.success('Address updated successfully');
      } else {
        // Add new address
        const response = await clientSideFetch<GenericResponse<AddressModel>>(
          PUBLIC_API_PATHS.USER_ADDRESSES,
          {
            method: 'POST',
            body: data,
          }
        );

        if (response.error) {
          toast.error(`Failed to add address: ${response.error.details}`);
          return;
        }

        toast.success('Address added successfully');
      }

      // Refresh user data to get updated addresses
      router.refresh();
      setIsAddModalOpen(false);
    } catch (error) {
      console.error('Error submitting address:', error);
      toast.error('An unexpected error occurred');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDeleteAddress = async () => {
    if (!selectedAddress) return;

    setIsSubmitting(true);
    try {
      const response = await clientSideFetch<GenericResponse<null>>(
        PUBLIC_API_PATHS.USER_ADDRESS.replace(
          ':id',
          selectedAddress.id.toString()
        ),
        {
          method: 'DELETE',
        }
      );

      if (response.error) {
        toast.error(`Failed to delete address: ${response.error.details}`);
        return;
      }

      toast.success('Address deleted successfully');
      router.refresh();
      setIsDeleteModalOpen(false);
    } catch (error) {
      console.error('Error deleting address:', error);
      toast.error('An unexpected error occurred');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleSetDefaultAddress = async (addressId: string) => {
    try {
      const response = await clientSideFetch<GenericResponse<AddressModel>>(
        PUBLIC_API_PATHS.USER_ADDRESS_DEFAULT.replace(
          ':id',
          addressId.toString()
        ),
        {
          method: 'PATCH',
        }
      );

      if (response.error) {
        toast.error(`Failed to set default address: ${response.error.details}`);
        return;
      }

      toast.success('Default address updated');
      router.refresh();
    } catch (error) {
      console.error('Error setting default address:', error);
      toast.error('An unexpected error occurred');
    }
  };

  return (
    <div className='space-y-6'>
      <h2 className='text-xl font-semibold text-gray-900'>
        Shipping Addresses
      </h2>

      {/* Address list */}
      {user?.addresses && user.addresses.length > 0 ? (
        <div className='grid grid-cols-1 sm:grid-cols-2 gap-4'>
          {user.addresses.map((address) => (
            <div
              key={address.id}
              className='bg-white border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow'
            >
              {address.default && (
                <span className='inline-block px-2 py-1 bg-green-100 text-green-800 text-xs font-medium rounded-md mb-2'>
                  Default
                </span>
              )}
              <p className='text-sm text-gray-500'>{address.street}</p>
              <p className='text-sm text-gray-500'>
                {address.city}, {address.district} {address.ward}
              </p>
              <p className='text-sm text-gray-500'>{address.phone}</p>
              <div className='mt-4 flex space-x-2'>
                <button
                  onClick={() => openEditModal(address)}
                  className='text-sm text-indigo-600 hover:text-indigo-800'
                >
                  Edit
                </button>
                <button
                  onClick={() => openDeleteModal(address)}
                  className='text-sm text-red-600 hover:text-red-800'
                >
                  Delete
                </button>
                {!address.default && (
                  <button
                    onClick={() => handleSetDefaultAddress(address.id)}
                    className='text-sm text-green-600 hover:text-green-800'
                  >
                    Set as Default
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p className='text-gray-500'>
          No addresses found. Add a shipping address to get started.
        </p>
      )}

      <button
        onClick={openAddModal}
        className='bg-indigo-600 text-white px-4 py-2 rounded-md hover:bg-indigo-700 transition-colors'
      >
        Add New Address
      </button>

      {/* Add/Edit Address Modal */}
      <Dialog
        open={isAddModalOpen}
        onClose={() => !isSubmitting && setIsAddModalOpen(false)}
        className='relative z-50'
      >
        <div className='fixed inset-0 bg-black/30' aria-hidden='true' />

        <div className='fixed inset-0 flex items-center justify-center p-4'>
          <Dialog.Panel className='w-full max-w-md transform overflow-hidden rounded-lg bg-white p-6 shadow-xl transition-all'>
            <Dialog.Title className='text-lg font-medium text-gray-900 mb-4'>
              {selectedAddress ? 'Edit Address' : 'Add New Address'}
            </Dialog.Title>

            <Fieldset
              as='form'
              onSubmit={handleSubmit(handleAddressSubmit)}
              className='space-y-4'
            >
              <TextField
                {...register('street')}
                label='Street Address'
                placeholder='123 Main St'
                error={errors.street?.message}
              />

              <div className='grid grid-cols-2 gap-4'>
                <TextField
                  {...register('city')}
                  label='City'
                  placeholder='New York'
                  error={errors.city?.message}
                />

                <TextField
                  {...register('district')}
                  label='District'
                  placeholder='Manhattan'
                  error={errors.district?.message}
                />
              </div>

              <div className='grid grid-cols-2 gap-4'>
                <TextField
                  {...register('ward')}
                  label='Ward/Neighborhood (Optional)'
                  placeholder='Downtown'
                  error={errors.ward?.message}
                />

                <TextField
                  {...register('phone')}
                  label='Phone Number'
                  placeholder='+1 (555) 123-4567'
                  error={errors.phone?.message}
                />
              </div>

              <div className='flex items-center'>
                <input
                  type='checkbox'
                  id='isDefault'
                  {...register('isDefault')}
                  className='h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500'
                />
                <Label
                  htmlFor='isDefault'
                  className='ml-2 block text-sm text-gray-900'
                >
                  Set as default address
                </Label>
              </div>

              <div className='flex justify-end space-x-3 mt-6'>
                <Button
                  type='button'
                  onClick={() => setIsAddModalOpen(false)}
                  className='px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50'
                  disabled={isSubmitting}
                >
                  Cancel
                </Button>
                <Button
                  type='submit'
                  className='px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700'
                  disabled={isSubmitting}
                >
                  {isSubmitting
                    ? 'Saving...'
                    : selectedAddress
                      ? 'Update Address'
                      : 'Add Address'}
                </Button>
              </div>
            </Fieldset>
          </Dialog.Panel>
        </div>
      </Dialog>

      {/* Delete Confirmation Modal */}
      <Dialog
        open={isDeleteModalOpen}
        onClose={() => !isSubmitting && setIsDeleteModalOpen(false)}
        className='relative z-50'
      >
        <div className='fixed inset-0 bg-black/30' aria-hidden='true' />

        <div className='fixed inset-0 flex items-center justify-center p-4'>
          <DialogPanel className='w-full max-w-md transform overflow-hidden rounded-lg bg-white p-6 shadow-xl transition-all'>
            <DialogTitle className='text-lg font-medium text-gray-900 mb-2'>
              Delete Address
            </DialogTitle>

            <p className='text-sm text-gray-500 mb-4'>
              Are you sure you want to delete this address? This action cannot
              be undone.
            </p>

            <div className='flex justify-end space-x-3 mt-6'>
              <Button
                type='button'
                onClick={() => setIsDeleteModalOpen(false)}
                className='px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md shadow-sm hover:bg-gray-50'
                disabled={isSubmitting}
              >
                Cancel
              </Button>
              <Button
                type='button'
                onClick={handleDeleteAddress}
                className='px-4 py-2 text-sm font-medium text-white bg-red-600 border border-transparent rounded-md shadow-sm hover:bg-red-700'
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Deleting...' : 'Delete'}
              </Button>
            </div>
          </DialogPanel>
        </div>
      </Dialog>
    </div>
  );
}
