import { apiFetch } from '@/lib/apis/api';
import { PUBLIC_API_PATHS } from '@/lib/constants/api';
import { GenericResponse, ProductDetailModel } from '@/lib/definitions';
import {
  CurrencyDollarIcon,
  GlobeAsiaAustraliaIcon,
} from '@heroicons/react/16/solid';
import React, { cache } from 'react';
import {
  AttributesSection,
  ImagesSection,
  RelateProductSection,
  ReviewSection,
  ReviewsList,
} from './_components';
import { Metadata } from 'next';

type Props = {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
};
export const getProduct = cache(async (slug: string) => {
  const { data, error } = await apiFetch<GenericResponse<ProductDetailModel>>(
    PUBLIC_API_PATHS.PRODUCT_DETAIL.replace(':id', slug),
    {
      nextOptions: {
        next: {
          tags: ['product'],
        },
      },
    }
  );
  if (error) {
    throw new Error(error.details, {
      cause: error,
    });
  }
  return data;
});

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  const post = await getProduct(slug);
  return {
    title: post.name,
    description: post.description,
  };
}

async function ProductDetailPage({ params }: Props) {
  const { slug } = await params;
  const data = await getProduct(slug);

  return (
    <div className='container mx-auto px-8 py-8'>
      <div className='lg:grid lg:grid-cols-3 lg:gap-x-8 lg:items-start'>
        {/* Image gallery */}
        <ImagesSection images={data.product_images} />

        {/* Product info */}
        <div className='mt-10 px-4 sm:px-0 sm:mt-16 lg:mt-0'>
          <div className='flex items-center justify-between mb-6'>
            <h1 className='text-3xl font-semibold tracking-tight text-gray-900'>
              {data.name}
            </h1>

            <div className=''>
              <h2 className='sr-only'>Product information</h2>
              <p className='text-2xl text-gray-900'>${data.price}</p>
            </div>
          </div>

          <div>
            {/* Reviews */}
            <ReviewSection rating={5} reviewsCount={30} />
          </div>
          <AttributesSection variants={data.variants} />
          <div className='mt-6'>
            <h3 className='sr-only'>Description</h3>
            <div
              className='text-base text-gray-700 space-y-6 list-inside list-disc'
              dangerouslySetInnerHTML={{ __html: data.description }} // Use only if description is trusted HTML
              // Or just: <p className="text-base text-gray-700">{data.description}</p>
            />
          </div>

          <div
            className='mt-6 text-base text-gray-700 space-y-4'
            dangerouslySetInnerHTML={{
              __html: data.description,
            }}
          />

          <div className='mt-10 pt-10 border-t border-gray-200'>
            <h3 className='text-sm font-medium text-gray-900'>Fabric & Care</h3>
            <div className='mt-4 prose prose-sm text-gray-500'>
              <ul role='list'>
                {details.map((item) => (
                  <li
                    className='text-sm text-gray-500 list-inside list-disc'
                    key={item}
                  >
                    {item}
                  </li>
                ))}
              </ul>
            </div>
          </div>

          {/* Info Boxes */}
          <div className='mt-8 grid grid-cols-1 gap-y-6 sm:grid-cols-2 sm:gap-x-6 lg:grid-cols-1 xl:grid-cols-2'>
            {/* International Delivery */}
            <div className='border border-gray-200 rounded-lg p-6 text-center'>
              <div className='flex items-center justify-center text-gray-400 mb-2'>
                <GlobeAsiaAustraliaIcon className='size-6' />
              </div>
              <p className='text-sm font-medium text-gray-900'>
                International delivery
              </p>
              <p className='mt-1 text-sm text-gray-500'>
                Get your order in 2 years
              </p>{' '}
              {/* Update text as needed */}
            </div>
            {/* Loyalty Rewards */}
            <div className='border border-gray-200 rounded-lg p-6 text-center'>
              <div className='flex items-center justify-center text-gray-400 mb-2'>
                <CurrencyDollarIcon className='size-6' />
              </div>
              <p className='text-sm font-medium text-gray-900'>
                Loyalty rewards
              </p>
              <p className='mt-1 text-sm text-gray-500'>
                Don&apos;t look at other tees
              </p>{' '}
              {/* Update text as needed */}
            </div>
          </div>
        </div>
      </div>
      <ReviewsList />
      <RelateProductSection />
    </div>
  );
}

export default ProductDetailPage;

const details = [
  'Only the best materials',
  'Ethically and locally made',
  'Pre-washed and pre-shrunk',
  'Machine wash cold with similar colors',
];
