'use client';

import { useState } from 'react';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import { toast } from 'react-toastify';

interface UserDetails {
  id: string;
  email: string;
  fullname: string;
  phone: string;
  role: string;
  createdAt: string;
  updatedAt: string;
  locked: boolean;
  avatarUrl: string | null;
  avatarImageId: string | null;
}

export default function UserProfileForm({ user }: { user: UserDetails }) {
  const router = useRouter();
  const [isEditing, setIsEditing] = useState(false);
  const [userData, setUserData] = useState({
    fullname: user.fullname,
    email: user.email,
    phone_number: user.phone || '',
    role: user.role,
  });
  const [loading, setLoading] = useState(false);
  const [avatar, setAvatar] = useState<File | null>(null);

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setUserData((prev) => ({ ...prev, [name]: value }));
  };

  const handleAvatarChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setAvatar(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      // If there's an avatar, upload it first
      let avatarImageId = user.avatarImageId;
      let avatarUrl = user.avatarUrl;

      if (avatar) {
        const formData = new FormData();
        formData.append('file', avatar);
        formData.append('type', 'avatar');

        const uploadResponse = await clientSideFetch(
          `${ADMIN_API_PATHS.UPLOADS}`,
          {
            method: 'POST',
            body: formData,
          }
        );

        if (uploadResponse.error) {
          throw new Error(
            uploadResponse.error.details || 'Failed to upload avatar'
          );
        }

        avatarImageId = uploadResponse.data.id;
        avatarUrl = uploadResponse.data.url;
      }

      // Update the user data
      const response = await clientSideFetch(
        `${ADMIN_API_PATHS.USERS}/${user.id}`,
        {
          method: 'PUT',
          body: {
            ...userData,
            avatarImageId: avatarImageId,
            avatarUrl: avatarUrl,
          },
        }
      );

      if (response.error) {
        throw new Error(response.error.details || 'Failed to update user');
      }

      toast.success('User profile updated successfully');
      setIsEditing(false);
      router.refresh();
    } catch (error: any) {
      toast.error(
        error.details || 'An error occurred while updating the profile'
      );
    } finally {
      setLoading(false);
    }
  };

  return (
    <div>
      <div className='flex justify-end mb-4'>
        <button
          type='button'
          onClick={() => setIsEditing(!isEditing)}
          className='inline-flex items-center px-3 py-1.5 border border-gray-300 shadow-sm text-sm font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
        >
          {isEditing ? 'Cancel' : 'Edit Profile'}
        </button>
      </div>

      <form onSubmit={handleSubmit}>
        <div className='flex items-center mb-6'>
          <div className='relative h-24 w-24 rounded-full overflow-hidden bg-gray-100 mr-6'>
            {user.avatarUrl || avatar ? (
              <Image
                src={avatar ? URL.createObjectURL(avatar) : user.avatarUrl!}
                alt='User avatar'
                fill
                className='object-cover'
              />
            ) : (
              <div className='flex items-center justify-center h-full bg-indigo-100 text-indigo-500'>
                <span className='text-2xl font-medium'>
                  {user.fullname
                    .split(' ')
                    .map((e) => e.charAt(0))
                    .join('')}
                </span>
              </div>
            )}
          </div>

          {isEditing && (
            <div>
              <label
                htmlFor='avatar'
                className='block text-sm font-medium text-gray-700 mb-2'
              >
                Profile Picture
              </label>
              <input
                type='file'
                id='avatar'
                name='avatar'
                accept='image/*'
                onChange={handleAvatarChange}
                className='block w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-medium file:bg-indigo-50 file:text-indigo-700 hover:file:bg-indigo-100'
              />
            </div>
          )}
        </div>

        <div className='grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-4'>
          <div>
            <label
              htmlFor='first_name'
              className='block text-sm font-medium text-gray-700'
            >
              First Name
            </label>
            <input
              type='text'
              name='first_name'
              id='first_name'
              value={userData.fullname}
              onChange={handleInputChange}
              disabled={!isEditing}
              className={`mt-1 block w-full shadow-sm sm:text-sm rounded-md ${
                isEditing
                  ? 'border-gray-300 focus:ring-indigo-500 focus:border-indigo-500'
                  : 'border-transparent bg-gray-50'
              }`}
            />
          </div>

          <div>
            <label
              htmlFor='last_name'
              className='block text-sm font-medium text-gray-700'
            >
              Last Name
            </label>
            <input
              type='text'
              name='last_name'
              id='last_name'
              value={userData.fullname}
              onChange={handleInputChange}
              disabled={!isEditing}
              className={`mt-1 block w-full shadow-sm sm:text-sm rounded-md ${
                isEditing
                  ? 'border-gray-300 focus:ring-indigo-500 focus:border-indigo-500'
                  : 'border-transparent bg-gray-50'
              }`}
            />
          </div>

          <div>
            <label
              htmlFor='email'
              className='block text-sm font-medium text-gray-700'
            >
              Email Address
            </label>
            <input
              type='email'
              name='email'
              id='email'
              value={userData.email}
              onChange={handleInputChange}
              disabled={!isEditing}
              className={`mt-1 block w-full shadow-sm sm:text-sm rounded-md ${
                isEditing
                  ? 'border-gray-300 focus:ring-indigo-500 focus:border-indigo-500'
                  : 'border-transparent bg-gray-50'
              }`}
            />
          </div>

          <div>
            <label
              htmlFor='phone_number'
              className='block text-sm font-medium text-gray-700'
            >
              Phone Number
            </label>
            <input
              type='tel'
              name='phone_number'
              id='phone_number'
              value={userData.phone_number}
              onChange={handleInputChange}
              disabled={!isEditing}
              className={`mt-1 block w-full shadow-sm sm:text-sm rounded-md ${
                isEditing
                  ? 'border-gray-300 focus:ring-indigo-500 focus:border-indigo-500'
                  : 'border-transparent bg-gray-50'
              }`}
            />
          </div>

          <div>
            <label
              htmlFor='role'
              className='block text-sm font-medium text-gray-700'
            >
              Role
            </label>
            <select
              id='role'
              name='role'
              value={userData.role}
              onChange={handleInputChange}
              disabled={!isEditing}
              className={`mt-1 block w-full py-2 px-3 border ${
                isEditing ? 'border-gray-300' : 'border-transparent bg-gray-50'
              } bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm`}
            >
              <option value='customer'>Customer</option>
              <option value='staff'>Staff</option>
              <option value='admin'>Admin</option>
            </select>
          </div>
        </div>

        {isEditing && (
          <div className='mt-6'>
            <button
              type='submit'
              disabled={loading}
              className='w-full inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
            >
              {loading ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        )}
      </form>
    </div>
  );
}
