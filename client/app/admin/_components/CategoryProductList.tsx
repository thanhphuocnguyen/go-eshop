import { ProductCategory } from '@/app/lib/definitions';
import React from 'react';
import Image from 'next/image';
import Link from 'next/link';

interface CategoryProductList {
  products: ProductCategory[];
  onRemove?: (productId: string) => void;
}

const CategoryProductList: React.FC<CategoryProductList> = ({ products, onRemove }) => {
  return (
    <div className='px-10 h-full'>
      <h2 className='text-2xl font-semibold text-gray-600 my-4'>
        Products Linked
      </h2>
      <div className='flex bg-tableRow h-full w-full flex-col space-y-4'>
        {products?.length === 0 ? (
          <div className='text-gray-500 text-center py-8'>No products linked to this category</div>
        ) : (
          products?.map((e) => (
            <div
              className='bg-white flex justify-between gap-x-2 items-center w-full relative p-3 border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow duration-200'
              key={e.id}
            >
              <Link href={`/admin/products/${e.id}`} className="flex gap-4 items-center flex-1 hover:opacity-90 transition-opacity">
                <div className='min-w-[80px] h-[80px] relative'>
                  <Image
                    fill
                    alt={`${e.name} product image`}
                    className='rounded-md object-cover'
                    src={e.imageUrl || '/images/product-placeholder.webp'}
                    sizes="80px"
                  />
                </div>
                <div className='flex flex-col space-y-2'>
                  <h3 className='text-gray-800 text-lg font-bold'>{e.name}</h3>
                  <div className='text-sm text-gray-600 bg-gray-100 inline-block px-3 py-1 rounded-full'>
                    {e.variantCount} {e.variantCount === 1 ? 'Variant' : 'Variants'}
                  </div>
                </div>
              </Link>
              <div className='flex mr-2'>
                <button
                  onClick={() => onRemove && onRemove(e.id)}
                  className='px-3 py-1.5 bg-red-500 hover:bg-red-600 text-white text-sm rounded-md transition-colors duration-200 flex items-center gap-2'
                >
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  Remove
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
};

export default CategoryProductList;
