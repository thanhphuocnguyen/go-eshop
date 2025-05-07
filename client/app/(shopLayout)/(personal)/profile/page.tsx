'use client';

import { useAppUser } from '@/lib/contexts/AppUserContext';
import PersonalInfoTab from './_components/PersonalInfoTab';

export default function ProfilePage() {
  const { user, mutateUser } = useAppUser();
  
  return <PersonalInfoTab userData={user} mutate={mutateUser} />;
}
