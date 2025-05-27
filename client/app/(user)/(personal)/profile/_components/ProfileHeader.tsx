'use client';

import { useUser } from '@/app/hooks';


export default function ProfileHeader() {
  const { user } = useUser();
  return (
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
  );
}
