import Footer from '@/components/Footer';
import NavBar from '@/components/Nav/NavBar';
import { JSX } from 'react';

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}): JSX.Element {
  return (
    <div className='h-full overflow-auto'>
      <NavBar />
      {children}
      <Footer />
    </div>
  );
}
