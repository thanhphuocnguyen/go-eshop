import Footer from '@/components/Footer';
import NavBar from '@/components/Nav/NavBar';
import { getUserCache } from '@/lib/cache/user';
import { CartContextProvider } from '@/lib/contexts/CartContext';

export default async function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const user = await getUserCache();
  return (
    <>
      <CartContextProvider user={user}>
        <div>
          <NavBar />
          <div className='overflow-auto'>{children}</div>
        </div>
      </CartContextProvider>
      <Footer />
    </>
  );
}
