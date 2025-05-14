import { apiFetchServerSide } from '@/app/lib/apis/apiServer';
import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { RatingModel } from '@/app/lib/definitions/rating';
import { StarIcon } from '@heroicons/react/24/solid';
import { StarIcon as StarOutlineIcon } from '@heroicons/react/24/outline';
import Image from 'next/image';
import { formatDate } from '@/app/lib/utils/index';
import React from 'react';
import RatingHelpfulButtons from './RatingHelpfulButtons';

// Client component to render ratings with interactive elements
export async function ReviewsList({ productID }: { productID: string }) {
  const { data: ratings, error } = await apiFetchServerSide<RatingModel[]>(
    PUBLIC_API_PATHS.PRODUCT_RATING.replaceAll(':id', productID)
  );

  if (error) {
    return (
      <div className='mt-20'>
        <h3 className='text-2xl font-semibold text-gray-600'>Reviews</h3>
        <hr className='my-8' />
        <div className='py-8 text-center text-red-500'>
          Error loading reviews: {error.details}
        </div>
      </div>
    );
  }

  if (ratings.length === 0) {
    return (
      <div className='mt-20'>
        <h3 className='text-2xl font-semibold text-gray-600'>Reviews</h3>
        <hr className='my-8' />
        <div className='py-8 text-center text-gray-500'>
          No reviews yet. Be the first to review this product!
        </div>
      </div>
    );
  }

  return (
    <div className='mt-20'>
      <h3 className='text-2xl font-semibold text-gray-600'>
        Reviews ({ratings.length})
      </h3>
      <hr className='my-8' />
      <ul className='flex flex-col gap-y-8 items-start'>
        {ratings.map((rating) => (
          <li
            key={rating.id}
            className='py-6 grid grid-cols-1 md:grid-cols-4 gap-x-6 border-b border-gray-100 items-start'
          >
            <div className='md:col-span-1 mt-2'>
              <div className='font-medium text-gray-800 flex items-center gap-x-2'>
                <span>{rating.name}</span>
                {rating.verified_purchase && (
                  <div className='inline-block bg-green-50 text-green-700 text-xs px-2 py-1 rounded-md border border-green-700'>
                    Verified Purchase
                  </div>
                )}
              </div>
              <div className='text-gray-400 text-base mt-2'>
                {formatDate(new Date(rating.created_at || Date.now()))}
              </div>
            </div>

            <div className='mt-3 flex items-start gap-x-2 md:col-span-1'>
              {Array.from({ length: 5 }, (_, i) => (
                <div key={i}>
                  {i < rating.rating ? (
                    <StarIcon className='size-5 text-yellow-400' />
                  ) : (
                    <StarOutlineIcon className='size-5 text-gray-300' />
                  )}
                </div>
              ))}
              <span className='text-sm font-medium'>{rating.rating}</span>
            </div>

            <div className='md:col-span-2 mt-4 md:mt-0 flex flex-col items-start'>
              <div className='text-lg font-semibold text-gray-800'>
                {rating.review_title}
              </div>
              <p className='mt-2 text-gray-600 whitespace-pre-line'>
                {rating.review_content}
              </p>

              {rating.images && rating.images.length > 0 && (
                <div className='mt-4 self-start w-full'>
                  <p className='text-sm font-medium text-gray-700 mb-2'>
                    Customer Images
                  </p>
                  <div className='flex flex-wrap gap-2 items-start'>
                    {rating.images.map((image) => (
                      <div key={image.id} className='relative h-24 w-24'>
                        <Image
                          src={image.url}
                          alt='Review image'
                          fill
                          className='object-cover rounded-md'
                          sizes='(max-width: 768px) 96px, 96px'
                        />
                      </div>
                    ))}
                  </div>
                </div>
              )}

              <div className='mt-4 self-start'>
                <RatingHelpfulButtons
                  ratingId={rating.id}
                  userId={rating.user_id}
                  helpfulVotes={rating.helpful_votes}
                  unhelpfulVotes={rating.unhelpful_votes}
                />
              </div>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
