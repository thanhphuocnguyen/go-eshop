'use client';

import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { CategoryListResponse, GenericResponse } from '@/lib/definitions';
import Link from 'next/link';
import { useEffect, useState } from 'react';
import CategoryProductSkeleton from '@/components/Product/CategoryProductSkeleton';
import ProductCard from '@/components/Product/ProductCard';
import { ArrowRightIcon } from '@heroicons/react/16/solid';

export default function CategoryPage() {
  const [categories, setCategories] = useState<CategoryListResponse[]>([]);
  const [loading, setLoading] = useState(true);

  // Fetch all categories
  useEffect(() => {
    async function fetchCategories() {
      try {
        const { data } = await apiFetch<
          GenericResponse<CategoryListResponse[]>
        >(`${PUBLIC_API_PATHS.CATEGORIES}?page=1&page_size=10`);

        if (data) {
          setCategories(data);
        }
      } catch (error) {
        console.error('Failed to fetch categories:', error);
      } finally {
        setLoading(false);
      }
    }

    fetchCategories();
  }, []);

  if (loading) {
    return (
      <div className='container mx-auto px-4 py-8'>
        <h1 className='text-3xl font-bold mb-8'>Categories</h1>
        <div className='space-y-12'>
          {[1, 2, 3].map((i) => (
            <div key={i} className='animate-pulse'>
              <div className='h-8 bg-gray-200 rounded w-1/4 mb-6'></div>
              <CategoryProductSkeleton />
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div className='container mx-auto px-4 py-8'>
      <div className='mb-8'>
        <h1 className='text-3xl font-bold'>Categories</h1>
        <p className='text-gray-500 mt-2'>Explore our product categories</p>
      </div>

      <div className='space-y-12'>
        {categories.map((item) => (
          <div
            key={item.category.id}
            className='pb-8 border-b border-gray-200 last:border-0'
          >
            <div className='mb-4'>
              <div className='flex items-baseline mb-2 justify-between gap-2'>
                <h2 className='text-2xl font-bold text-gray-800'>
                  {item.category.name}
                </h2>
                <Link
                  href={`/categories/${item.category.slug}`}
                  className='text-indigo-600 flex items-center text-base font-medium hover:underline'
                >
                  View all
                  <span>
                    <ArrowRightIcon className='size-6' />
                  </span>
                </Link>
              </div>
              <p className='text-gray-600 line-clamp-2'>
                {item.category.description}
              </p>
            </div>

            <>
              {item.products.length > 0 ? (
                <div className='grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-6'>
                  {item.products.slice(0, 5).map((product) => (
                    <ProductCard
                      key={product.id}
                      ID={product.id}
                      name={product.name}
                      image={product.image_url}
                      priceFrom={product.price_from}
                      priceTo={product.price_to}
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
            {/* )} */}
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
