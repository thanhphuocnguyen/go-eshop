import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import 'react-multi-carousel/lib/styles.css';
import ClientToastContainer from '@/components/Common/ToastContainer';
import clsx from 'clsx';
import { AppUserProvider } from '@/lib/contexts/AppUserContext';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: {
    template: '%s | Eshop',
    default: 'Eshop', // a default is required when creating a template
  },
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang='en'>
      <body className={clsx(inter.className)}>
        <main className='main'>
          <AppUserProvider>{children}</AppUserProvider>
        </main>
        <ClientToastContainer />
      </body>
    </html>
  );
}
