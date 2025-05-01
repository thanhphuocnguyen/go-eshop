'use client';

import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { AddressModel, GenericResponse } from '@/lib/definitions';
import { TextField } from '@/components/FormFields';
import { Button, Fieldset, Switch } from '@headlessui/react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import { z } from 'zod';
import { PlusCircleIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import { getCookie } from 'cookies-next';
import AddressCard from './_components/AddressCard';
import OrderHistoryTab from './_components/OrderHistoryTab';
import PersonalInfoTab from './_components/PersonalInfoTab';
import { useAppUser } from '@/lib/contexts/AppUserContext';

// Address form schema
const addressSchema = z.object({
  id: z.number().optional(),
  street: z.string().min(3, { message: 'Street address is required' }),
  district: z.string().min(2, { message: 'District is required' }),
  city: z.string().min(2, { message: 'City is required' }),
  ward: z.string().optional(),
  phone: z
    .string()
    .min(10, { message: 'Phone number must be at least 10 characters' }),
});

type AddressFormValues = z.infer<typeof addressSchema>;

export default function PersonalInfoPage() {
  const [activeTab, setActiveTab] = useState<
    'profile' | 'addresses' | 'orders' | 'security'
  >('profile');

  const [addAddressMode, setAddAddressMode] = useState(false);
  const [editAddressId, setEditAddressId] = useState<number | null>(null);
  const router = useRouter();
  const { user, mutateUser, isLoading } = useAppUser();

  // Address form
  const {
    register: registerAddress,
    handleSubmit: handleAddressSubmit,
    formState: { errors: addressErrors, isSubmitting: isAddressSubmitting },
    reset: resetAddressForm,
  } = useForm<AddressFormValues>({
    resolver: zodResolver(addressSchema),
    defaultValues: {
      street: '',
      district: '',
      city: '',
      ward: '',
      phone: '',
    },
  });

  // Address submit handler
  const onAddressSubmit = async (body: AddressFormValues) => {
    const url = editAddressId
      ? PUBLIC_API_PATHS.USER_ADDRESS.replace(':id', editAddressId.toString())
      : PUBLIC_API_PATHS.USER_ADDRESSES;

    const method = editAddressId ? 'PATCH' : 'POST';

    const { data, error } = await apiFetch<GenericResponse<AddressModel>>(url, {
      method,
      body: {
        user_id: user?.id,
        ...body,
      },
    });

    if (error) {
      toast.error(
        `Failed to ${editAddressId ? 'update' : 'add'} address: ` +
          error.details
      );
      return;
    }

    if (data) {
      await mutateUser();
      toast.success(
        `Address ${editAddressId ? 'updated' : 'added'} successfully`
      );
      resetAddressForm();
      setAddAddressMode(false);
      setEditAddressId(null);
    }
  };

  // Delete address handler
  const handleDeleteAddress = async (addressId: number) => {
    if (!confirm('Are you sure you want to delete this address?')) {
      return;
    }

    try {
      const response = await apiFetch(
        PUBLIC_API_PATHS.USER_ADDRESS.replace(':id', addressId.toString()),
        {
          method: 'DELETE',
          headers: {
            Authorization: `Bearer ${getCookie('access_token')}`,
          },
        }
      );

      if (response.error) {
        toast.error('Failed to delete address: ' + response.error.details);
        return;
      }

      await mutateUser();
      toast.success('Address deleted successfully');
    } catch (error) {
      toast.error('An error occurred while deleting your address');
      console.error(error);
    }
  };

  // Edit address handler
  const handleEditAddress = (address: AddressModel) => {
    setEditAddressId(address.id);
    setAddAddressMode(true);
    resetAddressForm({
      street: address.street || '',
      district: address.district || '',
      city: address.city || '',
      ward: address.ward || '',
      phone: address.phone || '',
    });
  };

  // Set address as default handler
  const handleSetDefaultAddress = async (addressId: number) => {
    const response = await apiFetch<GenericResponse<boolean>>(
      PUBLIC_API_PATHS.USER_ADDRESS_DEFAULT.replace(
        ':id',
        addressId.toString()
      ),
      {
        method: 'PATCH',
        headers: {
          Authorization: `Bearer ${getCookie('access_token')}`,
        },
      }
    );

    if (response.error) {
      toast.error('Failed to set default address: ' + response.error.details);
      return;
    }

    if (response.data) {
      await mutateUser();
      toast.success('Default address updated successfully');
    }
  };

  if (isLoading) {
    return (
      <div className='flex items-center justify-center h-screen'>
        <div className='animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-500'></div>
      </div>
    );
  }

  return (
    <div className='h-full max-w-7xl mx-auto p-4 sm:p-6 lg:p-8'>
      <div className='bg-white shadow rounded-lg'>
        {/* Profile Header */}
        <div className='px-6 py-8 border-b border-gray-200 sm:px-8'>
          <div className='flex items-center justify-between'>
            <div className='flex items-center'>
              <div className='h-16 w-16 rounded-full bg-indigo-600 flex items-center justify-center text-white text-2xl font-bold'>
                {user?.fullname?.charAt(0) || user?.username?.charAt(0) || '?'}
              </div>
              <div className='ml-4'>
                <h1 className='text-2xl font-bold text-gray-900'>
                  {user?.fullname}
                </h1>
                <p className='text-sm text-gray-500'>{user?.email}</p>
                <p className='text-xs mt-1'>
                  <span className='inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-indigo-100 text-indigo-800'>
                    {user?.role}
                  </span>
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Navigation Tabs */}
        <div className='border-b border-gray-200'>
          <nav className='flex -mb-px'>
            <button
              onClick={() => setActiveTab('profile')}
              className={clsx(
                'w-1/4 py-4 px-1 text-center border-b-2 font-medium text-sm',
                activeTab === 'profile'
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              )}
            >
              Personal Info
            </button>
            <button
              onClick={() => setActiveTab('addresses')}
              className={clsx(
                'w-1/4 py-4 px-1 text-center border-b-2 font-medium text-sm',
                activeTab === 'addresses'
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              )}
            >
              Addresses
            </button>
            <button
              onClick={() => setActiveTab('orders')}
              className={clsx(
                'w-1/4 py-4 px-1 text-center border-b-2 font-medium text-sm',
                activeTab === 'orders'
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              )}
            >
              Orders
            </button>
            <button
              onClick={() => setActiveTab('security')}
              className={clsx(
                'w-1/4 py-4 px-1 text-center border-b-2 font-medium text-sm',
                activeTab === 'security'
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              )}
            >
              Security
            </button>
          </nav>
        </div>

        <div className='p-6'>
          {/* Profile Info Tab */}
          {activeTab === 'profile' && (
            <PersonalInfoTab userData={user} mutate={mutateUser} />
          )}

          {/* Addresses Tab */}
          {activeTab === 'addresses' && (
            <div>
              <div className='flex justify-between items-center mb-6'>
                <h2 className='text-lg font-medium text-gray-900'>
                  Your Addresses
                </h2>
                <Button
                  onClick={() => {
                    setAddAddressMode(!addAddressMode);
                    setEditAddressId(null);
                    resetAddressForm();
                  }}
                  className='inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50'
                >
                  <PlusCircleIcon className='h-4 w-4 mr-1' />
                  {addAddressMode ? 'Cancel' : 'Add Address'}
                </Button>
              </div>

              {addAddressMode && (
                <Fieldset
                  as='form'
                  onSubmit={handleAddressSubmit(onAddressSubmit)}
                  className='mb-8 border border-gray-200 rounded-md p-6'
                >
                  <h3 className='text-md font-medium mb-4'>
                    {editAddressId ? 'Edit Address' : 'New Address'}
                  </h3>
                  <div className='grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6'>
                    <div className='sm:col-span-6'>
                      <TextField
                        {...registerAddress('street')}
                        label='Street address'
                        placeholder='123 Main St'
                        error={addressErrors.street?.message}
                      />
                    </div>
                    <div className='sm:col-span-2'>
                      <TextField
                        {...registerAddress('district')}
                        label='District'
                        placeholder='District'
                        error={addressErrors.district?.message}
                      />
                    </div>
                    <div className='sm:col-span-2'>
                      <TextField
                        {...registerAddress('ward')}
                        label='Ward (optional)'
                        placeholder='Ward'
                        error={addressErrors.ward?.message}
                      />
                    </div>
                    <div className='sm:col-span-2'>
                      <TextField
                        {...registerAddress('city')}
                        label='City'
                        placeholder='City'
                        error={addressErrors.city?.message}
                      />
                    </div>
                    <div className='sm:col-span-3'>
                      <TextField
                        {...registerAddress('phone')}
                        label='Phone number'
                        placeholder='+1 (555) 987-6543'
                        error={addressErrors.phone?.message}
                      />
                    </div>
                  </div>
                  <div className='flex justify-end mt-6'>
                    <Button
                      type='button'
                      onClick={() => {
                        setAddAddressMode(false);
                        setEditAddressId(null);
                      }}
                      className='mr-3 px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50'
                    >
                      Cancel
                    </Button>
                    <Button
                      type='submit'
                      className='px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700'
                      disabled={isAddressSubmitting}
                    >
                      {isAddressSubmitting
                        ? 'Saving...'
                        : editAddressId
                          ? 'Update Address'
                          : 'Add Address'}
                    </Button>
                  </div>
                </Fieldset>
              )}

              <div className='space-y-4'>
                {user?.addresses?.length ? (
                  user.addresses.map((address) => (
                    <AddressCard
                      key={address.id}
                      address={address}
                      onEdit={handleEditAddress}
                      onDelete={handleDeleteAddress}
                      onSetDefault={handleSetDefaultAddress}
                    />
                  ))
                ) : (
                  <div className='text-center py-12 border-2 border-dashed border-gray-300 rounded-lg'>
                    <p className='text-gray-500'>
                      You don&apos;t have any saved addresses yet.
                    </p>
                    <Button
                      onClick={() => setAddAddressMode(true)}
                      className='mt-4 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700'
                    >
                      <PlusCircleIcon className='h-5 w-5 mr-2' />
                      Add Your First Address
                    </Button>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Orders Tab */}
          {activeTab === 'orders' && (
            <div>
              <div className='flex justify-between items-center mb-6'>
                <h2 className='text-lg font-medium text-gray-900'>
                  Order History
                </h2>
              </div>
              <OrderHistoryTab />
            </div>
          )}

          {/* Security Tab */}
          {activeTab === 'security' && (
            <div>
              <h2 className='text-lg font-medium text-gray-900 mb-6'>
                Security Settings
              </h2>
              <div className='space-y-6'>
                <div className='bg-gray-50 p-6 rounded-lg border border-gray-200'>
                  <div className='flex items-center justify-between'>
                    <div>
                      <h3 className='text-md font-medium text-gray-900'>
                        Change Password
                      </h3>
                      <p className='text-sm text-gray-500 mt-1'>
                        Update your password to keep your account secure
                      </p>
                    </div>
                    <Button
                      onClick={() => router.push('/change-password')}
                      className='px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50'
                    >
                      Change Password
                    </Button>
                  </div>
                </div>

                <div className='bg-gray-50 p-6 rounded-lg border border-gray-200'>
                  <div className='flex items-center justify-between'>
                    <div>
                      <h3 className='text-md font-medium text-gray-900'>
                        Email Notifications
                      </h3>
                      <p className='text-sm text-gray-500 mt-1'>
                        Receive email updates about your orders and account
                        activity
                      </p>
                    </div>
                    <Switch defaultChecked={true} />
                  </div>
                </div>

                <div className='bg-gray-50 p-6 rounded-lg border border-gray-200'>
                  <div className='flex items-center justify-between'>
                    <div>
                      <h3 className='text-md font-medium text-gray-900'>
                        Two-Factor Authentication
                      </h3>
                      <p className='text-sm text-gray-500 mt-1'>
                        Add an extra layer of security to your account
                      </p>
                    </div>
                    <Button className='px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50'>
                      Setup 2FA
                    </Button>
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
