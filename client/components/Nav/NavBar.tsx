import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import Link from 'next/link';
import CartSection from './CartSection';
import AuthButtons from './AuthButtons';
import MainNavigation from './MainNavigation';

export default async function NavBar() {
  return (
    <div className='bg-white sticky top-0 z-20'>
      <header className='relative bg-white'>
        <p className='flex h-10 items-center justify-center bg-indigo-600 px-4 text-sm font-medium text-white sm:px-6 lg:px-8'>
          Get free delivery on orders over $100
        </p>

        <nav
          aria-label='Top'
          className='mx-auto max-w-8xl px-4 sm:px-6 lg:px-8'
        >
          <div className='border-b border-gray-200'>
            <div className='flex h-16 items-center'>
              <div className='ml-4 flex lg:ml-0'>
                <Link href='/'>
                  <span className='sr-only'>Eshop</span>
                  <Image
                    alt=''
                    src='/images/logos/logo.webp'
                    className='h-8 w-auto rounded-md'
                    width={50}
                    height={50}
                  />
                </Link>
              </div>

              {/* New navigation component */}
              <MainNavigation />

              <div className='ml-auto flex items-center'>
                {/* Search */}
                <div className='flex lg:ml-6'>
                  <Link href="/search" className='p-2 text-gray-400 hover:text-gray-500'>
                    <span className='sr-only'>Search</span>
                    <MagnifyingGlassIcon
                      aria-hidden='true'
                      className='size-6'
                    />
                  </Link>
                </div>

                <CartSection />
                <div className='hidden lg:flex lg:flex-1 lg:items-center lg:justify-end lg:space-x-6'>
                  <AuthButtons />
                </div>
              </div>
            </div>
          </div>
        </nav>
      </header>
    </div>
  );
}
