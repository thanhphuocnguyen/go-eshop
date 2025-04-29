import { StarIcon } from '@heroicons/react/24/solid';
import React from 'react';

export const ReviewsList = () => {
  return (
    <div className='mt-20'>
      <h3 className='text-2xl font-semibold text-gray-600'>Recent Reviews</h3>
      <hr className='my-8' />
      <ul className='flex flex-col gap-y-4'>
        <li className='py-4 grid grid-cols-4 gap-x-6'>
          <div>
            <div className='font-semibold'>Risako M</div>
            <div className='text-gray-400'>May 16, 2021</div>
          </div>
          <div className='mt-3 flex gap-x-2'>
            {Array.from({ length: 5 }, (_, i) => (
              <div key={i}>
                <StarIcon className='size-6 text-yellow-400' />
              </div>
            ))}
            <span>5</span>
          </div>
          <div className='col-span-2'>
            <div className='text-lg font-semibold'>Can't say enough good things </div>
            <p>
              I was really pleased with the overall shopping experience. My
              order even included a little personal, handwritten note, which
              delighted me! The product quality is amazing, it looks and feel
              even better than I had anticipated.
            </p>
            <p>
              Brilliant stuff! I would gladly recommend this store to my
              friends. And, now that I think of it... I actually have, many
              times!
            </p>
          </div>
        </li>
      </ul>
    </div>
  );
};
