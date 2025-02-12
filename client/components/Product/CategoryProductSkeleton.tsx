import React from 'react';
import ProductSkeleton from './ProductSkeleton';

const CategoryProductSkeleton = () => {
  return (
    <div className='grid gap-4 md:grid-cols-2 lg:grid-cols-4'>
      {new Array(4).fill(0).map((_, index) => (
        <ProductSkeleton key={index} />
      ))}
    </div>
  );
};

export default CategoryProductSkeleton;
