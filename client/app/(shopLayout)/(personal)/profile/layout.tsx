'use client';

import { useAppUser } from '@/lib/contexts/AppUserContext';
import ProfileHeader from './_components/ProfileHeader';
import ProfileTabs from './_components/ProfileTabs';

export default function ProfileLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const { user, isLoading } = useAppUser();

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
        <ProfileHeader userData={user} />

        {/* Navigation Tabs */}
        <ProfileTabs />

        <div className='p-6'>
          {/* Content for each tab will be rendered here */}
          {children}
        </div>
      </div>
    </div>
  );
}