'use client';

import {
  CategoryProductModel,
  GeneralCategoryModel,
  GenericResponse,
  ProductListModel,
} from '@/lib/definitions';
import { useState, useEffect } from 'react';
import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import CategoryFilters from './CategoryFilters';
import CategoryProducts from './CategoryProducts';
import CategoryProductSkeleton from '@/components/Product/CategoryProductSkeleton';

interface CategoryDetailClientProps {
  category: GeneralCategoryModel;
  initialProducts?: CategoryProductModel[];
}

export default function CategoryDetailClient({
  category,
}: CategoryDetailClientProps) {
  const [products, setProducts] = useState<ProductListModel[]>([]);
  const [loading, setLoading] = useState(!category);
  const [loadingProducts, setLoadingProducts] = useState(!category);

  // Filter states
  const [priceRange, setPriceRange] = useState<{ min: number; max: number }>({
    min: 0,
    max: 10000,
  });
  const [minPrice, setMinPrice] = useState<number>(0);
  const [maxPrice, setMaxPrice] = useState<number>(10000);
  const [selectedRating, setSelectedRating] = useState<number | null>(null);
  const [filterOpen, setFilterOpen] = useState(false);

  // For mobile responsiveness
  const toggleFilters = () => {
    setFilterOpen(!filterOpen);
  };

  useEffect(() => {
    // Only fetch if we don't have initial data
    fetchProducts();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [category]);

  async function fetchProducts() {
    setLoading(true);
    setLoadingProducts(true);

    // First fetch the category to get its ID
    const { data, error } = await apiFetch<GenericResponse<ProductListModel[]>>(
      `${PUBLIC_API_PATHS.PRODUCTS}?category_id=${category.id}&page=1&page_size=100`
    );

    if (error) {
      console.error('Failed to fetch products:', error);
      setLoading(false);
      setLoadingProducts(false);
      return;
    }
    if (data) {
      // Then fetch products for this category

      setProducts(data);

      // Set price range based on products
      if (data.length > 0) {
        const minProductPrice = Math.min(...data.map((p) => p.min_price));
        const maxProductPrice = Math.max(...data.map((p) => p.max_price));
        setPriceRange({ min: minProductPrice, max: maxProductPrice });
        setMinPrice(minProductPrice);
        setMaxPrice(maxProductPrice);
      }
    }
    setLoading(false);
    setLoadingProducts(false);
  }

  const resetFilters = () => {
    setMinPrice(priceRange.min);
    setMaxPrice(priceRange.max);
    setSelectedRating(null);
  };

  if (loading) {
    return (
      <div className='container mx-auto px-4 py-8'>
        <div className='animate-pulse'>
          <div className='h-10 bg-gray-200 rounded w-1/3 mb-6'></div>
          <div className='h-6 bg-gray-200 rounded w-1/2 mb-12'></div>
          <CategoryProductSkeleton />
        </div>
      </div>
    );
  }

  if (!category) {
    return (
      <div className='container mx-auto px-4 py-8'>
        <Link
          href='/categories'
          className='inline-flex items-center text-indigo-600 hover:underline mb-6'
        >
          <ArrowLeftIcon className='h-4 w-4 mr-2' />
          Back to all categories
        </Link>
        <div className='text-center py-12'>
          <h1 className='text-2xl font-bold text-gray-800 mb-2'>
            Category Not Found
          </h1>
          <p className='text-gray-600'>
            The category you are looking for does not exist or has been removed.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className='container mx-auto px-4 py-8'>
      <Link
        href='/categories'
        className='inline-flex items-center text-indigo-600 hover:underline mb-6'
      >
        <ArrowLeftIcon className='h-4 w-4 mr-2' />
        Back to all categories
      </Link>

      <div className='flex items-center mb-8'>
        {category.image_url && (
          <div className='h-24 w-24 mr-6 relative overflow-hidden rounded-lg shadow-md'>
            <Image
              src={category.image_url}
              alt={category.name}
              fill
              className='object-cover'
            />
          </div>
        )}
        <div>
          <h1 className='text-3xl font-bold text-gray-800'>{category.name}</h1>
          {category.description && (
            <p className='text-gray-600 mt-2 max-w-2xl'>
              {category.description}
            </p>
          )}
        </div>
      </div>

      <div className='flex flex-col lg:flex-row gap-8'>
        <CategoryFilters
          priceRange={priceRange}
          minPrice={minPrice}
          maxPrice={maxPrice}
          selectedRating={selectedRating}
          setMinPrice={setMinPrice}
          setMaxPrice={setMaxPrice}
          setSelectedRating={setSelectedRating}
          resetFilters={resetFilters}
          filterOpen={filterOpen}
          toggleFilters={toggleFilters}
        />

        <CategoryProducts
          products={products}
          loadingProducts={loadingProducts}
          resetFilters={resetFilters}
          minPrice={minPrice}
          maxPrice={maxPrice}
          selectedRating={selectedRating}
        />
      </div>
    </div>
  );
}
