import ProfileHeader from './_components/ProfileHeader';
import ProfileTabs from './_components/ProfileTabs';

export default async function ProfileLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className='h-full max-w-7xl mx-auto p-4 sm:p-6 lg:p-8'>
      <div className='bg-white shadow rounded-lg'>
        {/* Profile Header */}
        <ProfileHeader />

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
