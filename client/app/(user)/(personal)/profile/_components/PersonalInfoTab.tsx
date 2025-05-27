'use client';

import { clientSideFetch } from '@/app/lib/api/apiClient';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { GenericResponse, UserModel } from '@/app/lib/definitions';
import { TextField } from '@/app/components/FormFields';
import { Button, Fieldset } from '@headlessui/react';
import { zodResolver } from '@hookform/resolvers/zod';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { toast } from 'react-toastify';
import { PencilIcon, EnvelopeIcon } from '@heroicons/react/24/outline';
import { z } from 'zod';
import { useRouter } from 'next/navigation';
import { useUser } from '@/app/hooks';

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

export default function PersonalInfoTab() {
  const { user } = useUser();
  const router = useRouter();
  const [editMode, setEditMode] = useState(false);
  const [isSendingVerification, setIsSendingVerification] = useState(false);

  // Profile form
  const {
    register: registerProfile,
    handleSubmit: handleProfileSubmit,
    formState: { errors: profileErrors, isSubmitting: isProfileSubmitting },
    reset: resetProfileForm,
  } = useForm<ProfileFormValues>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      fullname: '',
      email: '',
      phone: '',
    },
  });

  // Set initial profile form values
  useEffect(() => {
    if (user) {
      resetProfileForm({
        fullname: user.fullname,
        email: user.email,
        phone: user.phone || '',
      });
    }
  }, [user, resetProfileForm]);

  // Profile update handler
  const onProfileSubmit = async (data: ProfileFormValues) => {
    try {
      const response = await clientSideFetch<UserModel>(
        PUBLIC_API_PATHS.GET_ME,
        {
          method: 'PATCH',
          body: {
            userId: user?.id,
            fullname: data.fullname,
            email: data.email,
            // phone: data.phone, // Add if your API supports phone updates
          },
        }
      );

      if (response.error) {
        toast.error('Failed to update profile: ' + response.error.details);
        return;
      }
      router.refresh();
      toast.success('Profile updated successfully');
      setEditMode(false);
    } catch (error) {
      toast.error('An error occurred while updating your profile');
      console.error(error);
    }
  };

  // Send verification email handler
  const handleSendVerificationEmail = async () => {
    if (!user || user.verifiedEmail) return;

    setIsSendingVerification(true);
    try {
      const response = await clientSideFetch<GenericResponse<null>>(
        PUBLIC_API_PATHS.GET_ME + '/send-verify-email',
        {
          method: 'POST',
        }
      );

      if (response.error) {
        toast.error(
          'Failed to send verification email: ' + response.error.details
        );
        return;
      }

      toast.success('Verification email sent! Please check your inbox.');
    } catch (error) {
      toast.error('An error occurred while sending the verification email');
      console.error(error);
    } finally {
      setIsSendingVerification(false);
    }
  };

  return (
    <div>
      <div className='flex justify-between items-center mb-6'>
        <h2 className='text-lg font-medium text-gray-900'>
          Personal Information
        </h2>
        <Button
          onClick={() => setEditMode(!editMode)}
          className='inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50'
        >
          <PencilIcon className='h-4 w-4 mr-1' />
          {editMode ? 'Cancel' : 'Edit'}
        </Button>
      </div>

      {editMode ? (
        <Fieldset
          as='form'
          onSubmit={handleProfileSubmit(onProfileSubmit)}
          className='space-y-6'
        >
          <div className='grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6'>
            <div className='sm:col-span-3'>
              <TextField
                {...registerProfile('fullname')}
                label='Full name'
                placeholder='John Doe'
                error={profileErrors.fullname?.message}
              />
            </div>
            <div className='sm:col-span-3'>
              <TextField
                {...registerProfile('phone')}
                label='Phone number'
                placeholder='+1 (555) 987-6543'
                error={profileErrors.phone?.message}
              />
            </div>
            <div className='sm:col-span-6'>
              <TextField
                {...registerProfile('email')}
                label='Email address'
                placeholder='john.doe@example.com'
                error={profileErrors.email?.message}
              />
            </div>
          </div>
          <div className='flex justify-end'>
            <Button
              type='button'
              onClick={() => setEditMode(false)}
              className='mr-3 px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50'
            >
              Cancel
            </Button>
            <Button
              type='submit'
              className='px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700'
              disabled={isProfileSubmitting}
            >
              {isProfileSubmitting ? 'Saving...' : 'Save Changes'}
            </Button>
          </div>
        </Fieldset>
      ) : (
        <div className='bg-gray-50 p-6 rounded-lg'>
          <dl className='grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2'>
            <div>
              <dt className='text-sm font-medium text-gray-500'>Full name</dt>
              <dd className='mt-1 text-sm text-gray-900'>{user?.fullname}</dd>
            </div>
            <div>
              <dt className='text-sm font-medium text-gray-500'>Username</dt>
              <dd className='mt-1 text-sm text-gray-900'>{user?.username}</dd>
            </div>
            <div>
              <dt className='text-sm font-medium text-gray-500'>
                Email address
              </dt>
              <dd className='mt-1 text-sm text-gray-900'>{user?.email}</dd>
            </div>
            <div>
              <dt className='text-sm font-medium text-gray-500'>
                Phone number
              </dt>
              <dd className='mt-1 text-sm text-gray-900'>
                {user?.phone || 'Not provided'}
              </dd>
            </div>
            <div>
              <dt className='text-sm font-medium text-gray-500'>
                Email verified
              </dt>
              <dd className='mt-1 text-sm text-gray-900 flex items-center gap-2'>
                <span
                  className={`inline-flex items-center px-3 py-1 rounded-md text-sm font-medium ${user?.verifiedEmail ? 'bg-green-100 text-green-800' : 'bg-yellow-200 text-orange-800'}`}
                >
                  {user?.verifiedEmail ? 'Verified' : 'Not verified'}
                </span>
                {!user?.verifiedEmail && (
                  <Button
                    onClick={handleSendVerificationEmail}
                    disabled={isSendingVerification}
                    className='inline-flex items-center ml-2 px-3 py-1 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
                  >
                    <EnvelopeIcon className='h-4 w-4 mr-1' />
                    {isSendingVerification ? 'Sending...' : 'Verify Email'}
                  </Button>
                )}
              </dd>
            </div>
            <div>
              <dt className='text-sm font-medium text-gray-500'>
                Phone verified
              </dt>
              <dd className='mt-1 text-sm text-gray-900 flex items-center gap-2'>
                <span
                  className={`inline-flex items-center px-3 py-1 rounded-md text-sm font-medium ${user?.verifiedPhone ? 'bg-green-100 text-green-800' : 'bg-yellow-200 text-orange-800'}`}
                >
                  {user?.verifiedPhone ? 'Verified' : 'Not verified'}
                </span>
              </dd>
            </div>
            <div>
              <dt className='text-sm font-medium text-gray-500'>
                Member since
              </dt>
              <dd className='mt-1 text-sm text-gray-900'>
                {user?.createdAt
                  ? new Date(user.createdAt).toLocaleDateString()
                  : 'Unknown'}
              </dd>
            </div>
          </dl>
        </div>
      )}
    </div>
  );
}
