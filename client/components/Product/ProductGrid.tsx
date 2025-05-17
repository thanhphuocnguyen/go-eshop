import React from 'react';
import Link from 'next/link';
import { ChevronLeftIcon, ChevronRightIcon } from '@heroicons/react/24/outline';
import ProductCard from './ProductCard';
import { Pagination, ProductDetailModel } from '@/app/lib/definitions';

interface ProductGridProps {
  products: ProductDetailModel[];
  pagination: Pagination;
  basePath: string;
}

export default function ProductGrid({
  products,
  pagination,
  basePath,
}: ProductGridProps) {
  // Calculate total pages
  const totalPages = Math.ceil(pagination.total / pagination.pageSize);
  const currentPage = pagination.page;

  // Function to calculate min and max prices from variants
  const getMinMaxPrice = (product: ProductDetailModel) => {
    if (!product.variants || product.variants.length === 0) {
      return { min: product.price, max: product.price };
    }

    const prices = product.variants.map((variant) => variant.price);
    return {
      min: Math.min(...prices),
      max: Math.max(...prices),
    };
  };

  // Function to get the primary image URL
  const getImageUrl = (product: ProductDetailModel) => {
    if (product.product_images && product.product_images.length > 0) {
      // Try to find a thumbnail image first
      const thumbnail = product.product_images.find(
        (img) => img.role === 'thumbnail'
      );
      if (thumbnail) return thumbnail.url;

      // Otherwise return the first image
      return product.product_images[0].url;
    }
    return null;
  };

  // Calculate average rating
  const calculateRating = (product: ProductDetailModel) => {
    const {
      one_star_count,
      two_star_count,
      three_star_count,
      four_star_count,
      five_star_count,
    } = product;

    const totalRatings =
      one_star_count +
      two_star_count +
      three_star_count +
      four_star_count +
      five_star_count;

    if (totalRatings === 0) return 0;

    const weightedSum =
      one_star_count * 1 +
      two_star_count * 2 +
      three_star_count * 3 +
      four_star_count * 4 +
      five_star_count * 5;

    return weightedSum / totalRatings;
  };

  return (
    <div>
      {/* Product Grid */}
      <div className='grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6'>
        {products.map((product) => {
          const { min, max } = getMinMaxPrice(product);
          const imageUrl = getImageUrl(product);
          const rating = calculateRating(product);

          return (
            <ProductCard
              key={product.id}
              slug={product.slug}
              ID={parseInt(product.id)}
              name={product.name}
              image={imageUrl}
              priceFrom={min}
              priceTo={max}
              rating={rating || 4.5} // Fallback to 4.5 if no ratings
            />
          );
        })}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className='flex items-center justify-center space-x-2 mt-12'>
          <Link
            href={`${basePath}?page=${Math.max(1, currentPage - 1)}`}
            className={`flex items-center px-3 py-2 rounded-md ${
              currentPage === 1
                ? 'text-gray-400 cursor-not-allowed'
                : 'text-gray-700 hover:bg-gray-100'
            }`}
            aria-disabled={currentPage === 1}
          >
            <ChevronLeftIcon className='h-5 w-5 mr-1' />
            Previous
          </Link>

          <div className='flex items-center space-x-1'>
            {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
              <Link
                key={page}
                href={`${basePath}?page=${page}`}
                className={`px-3 py-2 rounded-md ${
                  currentPage === page
                    ? 'bg-blue-600 text-white'
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                {page}
              </Link>
            ))}
          </div>

          <Link
            href={`${basePath}?page=${Math.min(totalPages, currentPage + 1)}`}
            className={`flex items-center px-3 py-2 rounded-md ${
              currentPage === totalPages
                ? 'text-gray-400 cursor-not-allowed'
                : 'text-gray-700 hover:bg-gray-100'
            }`}
            aria-disabled={currentPage === totalPages}
          >
            Next
            <ChevronRightIcon className='h-5 w-5 ml-1' />
          </Link>
        </div>
      )}
    </div>
  );
}
