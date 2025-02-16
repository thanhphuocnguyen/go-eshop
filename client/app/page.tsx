import CategoryProductSkeleton from '@/components/Product/CategoryProductSkeleton';
import ProductCarousel from '@/components/Product/ProductCarousel';
import { API_PUBLIC_PATHS } from '@/lib/constants/api';
import { Category, GenericListResponse } from '@/lib/types';
import Link from 'next/link';
import { Suspense } from 'react';

export default async function Home() {
  const categoriesResp: GenericListResponse<Category> = await fetch(
    process.env.NEXT_API_URL + API_PUBLIC_PATHS.CATEGORIES,
    {
      next: {
        tags: ['categories'],
      },
    }
  ).then((res) => res.json());

  return (
    <div className='w-full py-3 h-full flex flex-col gap-3'>
      {categoriesResp.data.map((category) => (
        <div key={category.category_id}>
          <Link
            href={`/category/${category.category_id}`}
            className='text-xl font-bold hover:text-green-800 text-green-700 cursor-pointer'
            key={category.category_id}
          >
            {category.name}
          </Link>
          <div className='w-full py-2'>
            <Suspense fallback={<CategoryProductSkeleton />}>
              <ProductCarousel categoryID={category.category_id} />
            </Suspense>
          </div>
        </div>
      ))}
    </div>
  );
}
