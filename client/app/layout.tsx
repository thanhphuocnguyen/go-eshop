import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import NavBar from '../components/NavBar';
import { SWRConfig } from 'swr';
import 'react-multi-carousel/lib/styles.css';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Eshop',
  description: 'An e-commerce site built with Next.js',
};

export default function RootLayout({
  children,
}: Readonly<{
  auth: React.ReactNode;
  children: React.ReactNode;
}>) {
  return (
    <html lang='en'>
      <body className={inter.className}>
        <NavBar />
        <div className='main'>
          <div className='gradient' />
        </div>
        <main className='app'>{children}</main>
      </body>
    </html>
  );
}
