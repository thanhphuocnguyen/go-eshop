import Link from 'next/link';
import Image from 'next/image';
import { ProductListModel } from '@/app/lib/definitions';
import { StarIcon } from '@heroicons/react/24/solid';
import { formatCurrency } from '@/app/utils';

export default function CategoryProducts({
  products,
  loadingProducts,
  minPrice,
  maxPrice,
  selectedRating,
}: {
  products: ProductListModel[];
  loadingProducts: boolean;
  minPrice?: number;
  maxPrice?: number;
  selectedRating?: number | null;
}) {
  // Check if we need to display filter information
  const hasFilters =
    minPrice !== undefined || maxPrice !== undefined || selectedRating !== null;

  // If loading, show skeleton
  if (loadingProducts) {
    return (
      <div className='flex-1'>
        <div className='grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6'>
          {[...Array(6)].map((_, index) => (
            <div key={index} className='animate-pulse border rounded-lg p-4'>
              <div className='h-48 bg-gray-300 rounded-md mb-4'></div>
              <div className='h-4 bg-gray-300 rounded w-3/4 mb-2'></div>
              <div className='h-4 bg-gray-300 rounded w-1/2 mb-4'></div>
              <div className='h-6 bg-gray-300 rounded w-1/4'></div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  // No products found
  if (products.length === 0) {
    return (
      <div className='flex-1'>
        <div className='text-center p-8 bg-gray-50 rounded-lg'>
          <h3 className='text-lg font-medium text-gray-900 mb-1'>
            No products found
          </h3>
          <p className='text-gray-500 mb-4'>
            {hasFilters
              ? 'Try adjusting your filters or browse all products in this category.'
              : "This category doesn't have any products yet."}
          </p>

          {hasFilters && (
            <Link
              href='?'
              className='inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700'
            >
              Clear all filters
            </Link>
          )}
        </div>
      </div>
    );
  }

  // Show products
  return (
    <div className='flex-1'>
      <div className='grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6'>
        {products.map((product) => (
          <ProductCard key={product.id} product={product} />
        ))}
      </div>
    </div>
  );
}

export function ProductCard({ product }: { product: ProductListModel }) {
  return (
    <Link
      href={`/products/${product.slug}`}
      className='group border rounded-lg p-4 hover:shadow-md transition-shadow duration-200 flex flex-col h-full'
    >
      <div className='relative h-48 mb-4 bg-gray-100 rounded-md overflow-hidden'>
        {product.imageUrl ? (
          <Image
            src={product.imageUrl}
            alt={product.name}
            fill
            sizes='(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw'
            className='object-contain group-hover:scale-105 transition-transform duration-200'
          />
        ) : (
          <div className='flex items-center justify-center h-full bg-gray-200'>
            <span className='text-gray-500'>No image</span>
          </div>
        )}
      </div>

      <h3 className='text-sm font-medium text-gray-900 mb-1 line-clamp-2'>
        {product.name}
      </h3>

      {product.reviewCount > 0 ? (
        <div className='flex items-center mb-2'>
          {[...Array(5)].map((_, i) => (
            <StarIcon
              key={i}
              className={`h-4 w-4 ${
                i < Math.round(product.avgRating)
                  ? 'text-yellow-400'
                  : 'text-gray-300'
              }`}
            />
          ))}
          <span className='text-xs text-gray-500 ml-1'>
            ({product.reviewCount || 0})
          </span>
        </div>
      ) : (
        <div>
          <span className='text-xs text-gray-500'>No reviews yet</span>
        </div>
      )}

      <div className='mt-auto'>
        <span className='text-sm font-medium text-gray-900'>
          {product.minPrice === product.maxPrice
            ? formatCurrency(product.minPrice)
            : `${formatCurrency(product.minPrice)} - ${formatCurrency(product.maxPrice)}`}
        </span>
      </div>
    </Link>
  );
}
