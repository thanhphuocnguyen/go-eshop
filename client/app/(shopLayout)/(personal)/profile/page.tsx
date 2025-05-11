import { getUserCache } from '@/lib/cache/user';
import PersonalInfoTab from './_components/PersonalInfoTab';

export default async function ProfilePage() {
  const user = await getUserCache();
  if (!user) {
    return (
      <div className='h-full max-w-7xl mx-auto p-4 sm:p-6 lg:p-8'>
        Please login to view your profile.
      </div>
    );
  }
  return <PersonalInfoTab userData={user} />;
}
