import {
  MagnifyingGlassIcon,
  ShoppingBagIcon,
} from '@heroicons/react/24/outline';
import Image from 'next/image';
import Link from 'next/link';
import CartSection from './CartSection';
import CategorySection from './CategorySection';
import AuthButtons from './AuthButtons';
import { getUserCache } from '@/lib/cache/user';

export default async function NavBar() {
  const user = await getUserCache();
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
              {/* <Button
                type='button'
                onClick={() => setOpen(true)}
                className='relative rounded-md bg-white p-2 text-gray-400 lg:hidden'
              >
                <span className='absolute -inset-0.5' />
                <span className='sr-only'>Open menu</span>
                <Bars3Icon aria-hidden='true' className='size-6' />
              </Button> */}

              {/* Logo */}
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

              <CategorySection />

              <div className='ml-auto flex items-center'>
                {/* <div className='hidden lg:ml-8 lg:flex'>
                  <Link
                    href='/'
                    className='flex items-center text-gray-700 hover:text-gray-800'
                  >
                    <Image
                      width={40}
                      height={40}
                      alt=''
                      src='https://tailwindui.com/plus/img/flags/flag-canada.svg'
                      className='block h-auto w-5 shrink-0'
                    />
                    <span className='ml-3 block text-sm font-medium'>CAD</span>
                    <span className='sr-only'>, change currency</span>
                  </Link>
                </div> */}

                {/* Search */}
                <div className='flex lg:ml-6'>
                  <a href='#' className='p-2 text-gray-400 hover:text-gray-500'>
                    <span className='sr-only'>Search</span>
                    <MagnifyingGlassIcon
                      aria-hidden='true'
                      className='size-6'
                    />
                  </a>
                </div>

                {/* Orders - only shown when logged in */}
                {!!user && (
                  <div className='hidden lg:flex lg:ml-6'>
                    <Link
                      href='/orders'
                      className='flex items-center p-2 text-gray-700 hover:text-indigo-600'
                    >
                      <ShoppingBagIcon className='size-6 mr-1' />
                      <span className='text-sm font-medium'>Orders</span>
                    </Link>
                  </div>
                )}

                <CartSection />
                <div className='hidden lg:flex lg:flex-1 lg:items-center lg:justify-end lg:space-x-6'>
                  <AuthButtons user={user} />
                </div>
              </div>
            </div>
          </div>
        </nav>
      </header>
    </div>
  );
}
