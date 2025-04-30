import Footer from '@/components/Footer';
import NavBar from '@/components/Nav/NavBar';
import { JSX } from 'react';
import { CartContextProvider } from '@/lib/contexts/CartContext';

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}): JSX.Element {
  return (
    <>
      <CartContextProvider>
        <div>
          <NavBar />
          <div className='overflow-auto'>{children}</div>
        </div>
      </CartContextProvider>
      <Footer />
    </>
  );
}
