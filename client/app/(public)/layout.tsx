import Footer from '@/components/Footer';
import NavBar from '@/components/Nav/NavBar';
import { AppUserProvider } from '../components/AppUserContext';
import { JSX } from 'react';

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}): JSX.Element {
  return (
    <AppUserProvider>
      <div className=''>
        <NavBar />
        <div className='h-full'>{children}</div>
        <Footer />
      </div>
    </AppUserProvider>
  );
}
