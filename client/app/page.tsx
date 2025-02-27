import CategoryProductSkeleton from '@/components/Product/CategoryProductSkeleton';
import ProductCarousel from '@/components/Product/ProductCarousel';
import { API_PUBLIC_PATHS } from '@/lib/constants/api';
import { Category, GenericListResponse } from '@/lib/types';
import Image from 'next/image';
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
  const categories = categoriesResp.data;
  return (
    <div className='bg-white'>
      {categories.map((category) => (
        <div
          key={category.category_id}
          className='mx-auto max-w-2xl px-4 py-16 sm:px-6 sm:py-24 lg:max-w-7xl lg:px-8'
        >
          <div className='flex justify-between'>
            <h2 className='text-2xl font-bold tracking-tight text-gray-900'>
              {category.name}
            </h2>
          </div>
          <div className='grid grid-cols-1 gap-x-6 gap-y-10 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 xl:gap-x-8'>
            {category.products.map((product) => (
              <Link
                href='/products/[id]'
                as={`/products/${product.id}`}
                className='group'
              >
                <Image
                  src={product.image_url || '/images/product-placeholder.webp'}
                  width={300}
                  height={220}
                  alt='Tall slender porcelain bottle with natural clay textured body and cork stopper.'
                  className='aspect-square w-full rounded-lg bg-gray-200 object-cover group-hover:opacity-75 xl:aspect-7/8'
                />
                <h3 className='mt-4 text-sm text-gray-700'>{product.name}</h3>
                <p className='mt-1 text-lg font-medium text-gray-900'>${product.price_from} - {product.price_to}</p>
              </Link>
            ))}
            {/* <!-- More products... --> */}
          </div>
        </div>
      ))}
    </div>
  );
}
