import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import NavBar from '../components/Nav/NavBar';
import 'react-multi-carousel/lib/styles.css';
import ClientToastContainer from '@/components/Common/ToastContainer';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Eshop',
  description: 'An e-commerce site built with Next.js',
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang='en'>
      <body className={inter.className}>
        <NavBar />
        <div className='main'>
          <div className='gradient' />
        </div>
        <main className='app'>
          {children}
          <ClientToastContainer />
        </main>
      </body>
    </html>
  );
}
