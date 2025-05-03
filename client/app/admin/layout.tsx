import './admin.css';
import { Metadata } from 'next';
import { AdminNavbar, AdminSidebar } from './_components';

export const metadata: Metadata = {
  title: {
    template: '%s | Admin Dashboard',
    default: 'Admin Dashboard',
  },
  description: 'E-commerce admin dashboard for managing your online store',
};

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className='flex h-screen'>
      <AdminSidebar />
      <main className='w-5/6 block bg-white'>
        <AdminNavbar />
        <section className='p-4 m-5 border-2 border-gray-200 shadow-sm flex flex-col h-[90%] rounded-md'>
          {children}
        </section>
      </main>
    </div>
  );
}
