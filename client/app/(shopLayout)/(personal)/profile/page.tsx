import { Metadata } from 'next';
import PersonalInfoTab from './_components/PersonalInfoTab';

export const metadata: Metadata = {
  title: 'Profile | E-Shop',
  description: 'Manage your profile information and settings.',
};

export default async function ProfilePage() {
  return <PersonalInfoTab />;
}
