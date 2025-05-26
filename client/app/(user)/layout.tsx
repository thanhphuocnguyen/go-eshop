import Footer from '@/app/components/Footer';
import NavBar from '@/app/components/Nav/NavBar';
import { CartContextProvider } from '@/app/lib/contexts/CartContext';

export default async function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
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
