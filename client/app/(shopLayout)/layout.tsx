import Footer from '@/components/Footer';
import NavBar from '@/components/Nav/NavBar';
import { AppUserProvider } from '../../components/AppUserContext';
import { JSX } from 'react';

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}): JSX.Element {
  return (
    <AppUserProvider>
      <NavBar />
      <div className='overflow-auto'>{children}</div>
      <Footer />
    </AppUserProvider>
  );
}
