import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import ClientToastContainer from '@/app/components/Common/ToastContainer';
import clsx from 'clsx';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc' // ES 2015
import timezone from 'dayjs/plugin/timezone' // ES 2015

dayjs.extend(utc);
dayjs.extend(timezone);

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: {
    template: '%s | Eshop',
    default: 'Eshop', // a default is required when creating a template
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang='en'>
      <body className={clsx(inter.className)}>
        <main className='main'>{children}</main>
        <ClientToastContainer />
      </body>
    </html>
  );
}
