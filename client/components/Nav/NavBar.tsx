import { MagnifyingGlassIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import Link from 'next/link';
import { cookies } from 'next/headers';
import CartSection from './CartSection';
import CategorySection from './CategorySection';
import AuthButtons from './AuthButtons';
import { UserModel } from '@/lib/definitions';

export default async function NavBar() {
  const cookieStore = await cookies();
  // const [open, setOpen] = useState(false);
  const role = cookieStore.get('user_role')?.value;
  const fullname = cookieStore.get('user_name')?.value;

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
                    src='/images/logo/logo.webp'
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
                <CartSection />
                <div className='hidden lg:flex lg:flex-1 lg:items-center lg:justify-end lg:space-x-6'>
                  <AuthButtons
                    role={role as UserModel['role']}
                    name={fullname as string}
                  />
                </div>
              </div>
            </div>
          </div>
        </nav>
      </header>
    </div>
  );
}
