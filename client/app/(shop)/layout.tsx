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
    <div className='h-full'>
      <AppUserProvider>
        <NavBar />
        <div className=''>{children}</div>
        <Footer />
      </AppUserProvider>
    </div>
  );
}
