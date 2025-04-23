import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import 'react-multi-carousel/lib/styles.css';
import ClientToastContainer from '@/components/Common/ToastContainer';
import clsx from 'clsx';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Simple Life',
  description: 'An e-commerce site built with Next.js',
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang='en'>
      <body className={clsx(inter.className, 'overflow-hidden')}>
        <main className='main h-full'>{children}</main>
        <ClientToastContainer />
      </body>
    </html>
  );
}
