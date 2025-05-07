'use client';

import ProductCard from '@/components/Product/ProductCard';
import CategoryProductSkeleton from '@/components/Product/CategoryProductSkeleton';
import { ProductListModel } from '@/lib/definitions';
import { useMemo, useState } from 'react';

interface CategoryProductsProps {
  products: ProductListModel[];
  loadingProducts: boolean;
  resetFilters: () => void;
  minPrice: number;
  maxPrice: number;
  selectedRating: number | null;
}

export default function CategoryProducts({
  products,
  loadingProducts,
  resetFilters,
  minPrice,
  maxPrice,
  selectedRating,
}: CategoryProductsProps) {
  const [sortOption, setSortOption] = useState<string>('default');

  // Filter products based on selected filters
  const filteredProducts = useMemo(() => {
    if (!products.length) return [];

    return products.filter((product) => {
      // Filter by price range
      const isPriceInRange =
        (product.min_price >= minPrice || product.max_price >= minPrice) &&
        (product.min_price <= maxPrice || product.max_price <= maxPrice);

      // Return true if all conditions are met
      return isPriceInRange;
    });
  }, [products, minPrice, maxPrice]);

  // Sort products based on selected sorting option
  const sortedProducts = useMemo(() => {
    if (!filteredProducts.length) return [];

    const productsCopy = [...filteredProducts];

    switch (sortOption) {
      case 'price-low-high':
        return productsCopy.sort((a, b) => a.min_price - b.min_price);
      case 'price-high-low':
        return productsCopy.sort((a, b) => b.min_price - a.min_price);
      case 'newest':
        return productsCopy.sort(
          (a, b) =>
            new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        );
      default:
        return productsCopy;
    }
  }, [filteredProducts, sortOption]);

  if (loadingProducts) {
    return <CategoryProductSkeleton />;
  }

  return (
    <div className='flex-1'>
      <div className='mb-6 flex justify-between items-center'>
        <h2 className='text-xl font-semibold text-gray-800'>
          Products in this category
        </h2>

        {/* Sort dropdown */}
        <div className='flex items-center'>
          <label htmlFor='sort' className='mr-2 text-sm text-gray-600'>
            Sort by:
          </label>
          <select
            id='sort'
            value={sortOption}
            onChange={(e) => setSortOption(e.target.value)}
            className='p-2 border border-gray-300 rounded-md bg-white'
          >
            <option value='default'>Featured</option>
            <option value='price-low-high'>Price: Low to High</option>
            <option value='price-high-low'>Price: High to Low</option>
            <option value='newest'>Newest Arrivals</option>
          </select>
        </div>
      </div>

      {sortedProducts.length > 0 ? (
        <>
          <p className='text-sm text-gray-500 mb-4'>
            {sortedProducts.length}{' '}
            {sortedProducts.length === 1 ? 'product' : 'products'} found
          </p>
          <div className='grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6'>
            {sortedProducts.map((product) => (
              <ProductCard
                key={product.id}
                ID={product.id}
                name={product.name}
                image={product.image_url}
                priceFrom={product.min_price}
                priceTo={product.max_price}
                rating={4.5}
              />
            ))}
          </div>
        </>
      ) : (
        <div className='text-center py-12 bg-gray-50 rounded-lg'>
          <p className='text-lg text-gray-600'>
            No products available with the selected filters.
          </p>
          <button
            onClick={resetFilters}
            className='text-indigo-600 hover:underline mt-2 inline-block'
          >
            Clear all filters
          </button>
        </div>
      )}
    </div>
  );
}
