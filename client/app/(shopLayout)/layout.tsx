import Footer from '@/components/Footer';
import NavBar from '@/components/Nav/NavBar';
import { CartContextProvider } from '@/lib/contexts/CartContext';

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
