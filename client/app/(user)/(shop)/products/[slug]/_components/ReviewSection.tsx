'use client';
import { StarIcon as SolidStarIcon } from '@heroicons/react/16/solid';
import { StarIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';
import React from 'react';
const renderStars = (rating: number) => {
  const stars = [];
  for (let i = 1; i <= 5; i++) {
    if (i <= Math.floor(rating)) {
      stars.push(<SolidStarIcon key={i} className='size-6 text-yellow-400' />);
    } else if (i === Math.ceil(rating) && !Number.isInteger(rating)) {
      stars.push(<SolidStarIcon key={i} className='size-6 text-yellow-400' />);
    } else {
      stars.push(<StarIcon key={i} className='size-6 text-yellow-400' />);
    }
  }
  return stars;
};
interface ReviewSectionProps {
  rating: number;
  reviewsCount: number;
}

export const ReviewSection: React.FC<ReviewSectionProps> = ({
  rating,
  reviewsCount,
}) => {
  return (
    <div className=''>
      {/* Reviews */}
      <div>
        <h3 className='sr-only'>Reviews</h3>
        <div className='flex items-center'>
          <div className='flex items-center'>{renderStars(rating)}</div>
          <p className='sr-only'>{5} out of 5 stars</p>
          {reviewsCount > 0 ? (
            <Link
              href={`/products/${'product-slug'}/reviews`}
              className='ml-8 text-sm font-medium text-indigo-600 hover:text-indigo-500'
            >
              See all {reviewsCount} reviews
            </Link>
          ) : (
            <span className='ml-8 text-sm text-gray-500'>
              No reviews yet
            </span>
          )}
        </div>
      </div>
    </div>
  );
};
