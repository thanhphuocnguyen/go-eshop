import { API_PRIVATE_PATHS } from '@/lib/constants/api';
import { GenericResponse } from '@/lib/types';
import { ShoppingBagIcon } from '@heroicons/react/24/outline';
import { cookies } from 'next/headers';
import Link from 'next/link';

export default async function CartSection() {
  let itemCount = 0;
  const cookieStore = await cookies();
  if (cookieStore.get('token')) {
    const res: GenericResponse<number> = await fetch(
      process.env.NEXT_API_URL! + API_PRIVATE_PATHS.CART_ITEM_COUNT,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${cookieStore.get('token')}`,
        },
        next: {
          tags: ['cart'],
        },
      }
    ).then((res) => res.json());
    if (res.data) {
      itemCount = res.data;
    }
  }
  return (
    <div className='ml-4 flow-root lg:ml-6'>
      <Link href='/cart' className='group -m-2 flex items-center p-2'>
        <ShoppingBagIcon
          aria-hidden='true'
          className='size-6 shrink-0 text-gray-400 group-hover:text-gray-500'
        />
        <span className='ml-2 text-sm font-medium text-gray-700 group-hover:text-gray-800'>
          {itemCount}
        </span>
        <span className='sr-only'>items in cart, view bag</span>
      </Link>
    </div>
  );
}
