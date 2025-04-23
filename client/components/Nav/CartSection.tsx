'use client';
import { useAppUser } from '@/app/components/AppUserContext';
import { ShoppingBagIcon } from '@heroicons/react/24/outline';
import clsx from 'clsx';
import Link from 'next/link';

export default function CartSection() {
  const { cart } = useAppUser();
  return (
    <div className='ml-4 flow-root lg:mx-6'>
      <Link href='/checkout' className='group relative -m-2 flex items-center p-2'>
        <ShoppingBagIcon
          aria-hidden='true'
          className='size-8 shrink-0 text-gray-400 group-hover:text-gray-500'
        />
        <span
          className={clsx(
            'top-1 -right-0.5 text-xs font-medium rounded-full px-1.5 py-0.5 absolute',
            cart?.cart_items.length
              ? 'bg-red-500 text-white'
              : 'text-gray-700 group-hover:text-gray-800'
          )}
        >
          {cart?.cart_items.length || 0}
        </span>
        <span className='sr-only'>items in cart, view bag</span>
      </Link>
    </div>
  );
}
