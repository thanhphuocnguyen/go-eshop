import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { GeneralCategoryModel } from '@/app/lib/definitions';
import Link from 'next/link';
import ProductCard from '@/components/Product/ProductCard';
import { ArrowRightIcon } from '@heroicons/react/16/solid';
import { apiFetchServerSide } from '@/app/lib/apis/apiServer';

async function getCategories() {
  // Using apiFetch utility instead of native fetch
  const result = await apiFetchServerSide<GeneralCategoryModel[]>(
    `${PUBLIC_API_PATHS.CATEGORIES}?page=1&page_size=10`
  );

  return result.data || [];
}

export default async function CategoryPage() {
  // Server component with async data fetching
  const categories = await getCategories();
  console.log(categories);
  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='mb-8'>
        <h1 className='text-3xl font-bold'>Categories</h1>
        <p className='text-gray-500 mt-2'>Explore our product categories</p>
      </div>

      <div className='space-y-12'>
        {categories.map((item) => (
          <div
            key={item.id}
            className='pb-8 border-b border-gray-200 last:border-0'
          >
            <div className='mb-4'>
              <div className='flex items-baseline mb-2 justify-between gap-2'>
                <h2 className='text-2xl font-bold text-gray-800'>
                  {item.name}
                </h2>
                <Link
                  href={`/categories/${item.slug}`}
                  className='text-indigo-600 flex items-center text-base font-medium hover:underline'
                >
                  View all
                  <span>
                    <ArrowRightIcon className='size-6' />
                  </span>
                </Link>
              </div>
              <p className='text-gray-600 line-clamp-2'>{item.description}</p>
            </div>

            <>
              {item.products && item.products.length > 0 ? (
                <div className='grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6'>
                  {item.products.slice(0, 5).map((product) => (
                    <ProductCard
                      slug={product.slug}
                      key={product.id}
                      ID={parseInt(product.id)}
                      name={product.name}
                      priceFrom={product.min_price}
                      priceTo={product.max_price}
                      image={product.image_url || ''}
                      rating={4.5}
                    />
                  ))}
                </div>
              ) : (
                <p className='text-gray-500 italic'>
                  No products available in this category.
                </p>
              )}
            </>
          </div>
        ))}

        {categories.length === 0 && (
          <div className='text-center py-12'>
            <p className='text-xl text-gray-600'>No categories found.</p>
          </div>
        )}
      </div>
    </div>
  );
}
