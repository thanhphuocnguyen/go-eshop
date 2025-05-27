import { PUBLIC_API_PATHS } from '@/app/lib/constants/api';
import { ProductDetailModel } from '@/app/lib/definitions';
import {
  CurrencyDollarIcon,
  GlobeAsiaAustraliaIcon,
  TagIcon,
} from '@heroicons/react/16/solid';
import React, { cache, Suspense } from 'react';
import {
  AttributesSection,
  ImagesSection,
  RelateProductSection,
  ReviewSection,
  ReviewsList,
  TabsSection,
} from './_components';
import { Metadata } from 'next';
import { serverSideFetch } from '@/app/lib/api/apiServer';

type Props = {
  params: Promise<{ slug: string }>;
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
};
export const getCacheProduct = cache(async (slug: string) => {
  const { data, error } = await serverSideFetch<ProductDetailModel>(
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
  const post = await getCacheProduct(slug);
  return {
    title: post.name,
    description: post.description,
  };
}

async function ProductDetailPage({ params }: Props) {
  const { slug } = await params;
  const productDetail = await getCacheProduct(slug);

  // Calculate average rating
  const {
    oneStarCount,
    twoStarCount,
    threeStarCount,
    fourStarCount,
    fiveStarCount,
  } = productDetail;

  const totalRatings =
    oneStarCount +
    twoStarCount +
    threeStarCount +
    fourStarCount +
    fiveStarCount;
  const avgRating =
    totalRatings > 0
      ? (oneStarCount * 1 +
          twoStarCount * 2 +
          threeStarCount * 3 +
          fourStarCount * 4 +
          fiveStarCount * 5) /
        totalRatings
      : 0;

  // Calculate min and max prices from variants if available
  const prices = productDetail.variants.map((variant) => variant.price);
  const minPrice =
    prices.length > 0 ? Math.min(...prices) : productDetail.price;
  const maxPrice =
    prices.length > 0 ? Math.max(...prices) : productDetail.price;
  return (
    <div className='container mx-auto '>
      {/* Breadcrumbs */}
      <nav className='flex text-sm text-gray-500 my-8' aria-label='Breadcrumb'>
        <ol className='flex items-center space-x-2'>
          <li>
            <a href='/' className='hover:text-gray-700'>
              Home
            </a>
          </li>
          <li>
            <span className='mx-2'>/</span>
            <a href='/products' className='hover:text-gray-700'>
              Products
            </a>
          </li>
          <li>
            <span className='mx-2'>/</span>
            <a
              href={`/categories/${productDetail.category?.slug}`}
              className='hover:text-gray-700'
            >
              {productDetail.category?.name}
            </a>
          </li>
          <li>
            <span className='mx-2'>/</span>
            <span className='text-gray-700 font-medium'>
              {productDetail.name}
            </span>
          </li>
        </ol>
      </nav>

      <div className='lg:grid lg:grid-cols-3 lg:gap-x-12 lg:items-start'>
        {/* Image gallery */}
        <ImagesSection images={productDetail.productImages} />

        {/* Product info */}
        <div className='mt-10 px-4 sm:px-0 sm:mt-16 lg:mt-0'>
          {/* Brand badge */}
          {productDetail.brand && (
            <div className='mb-4'>
              <span className='inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-50 text-blue-700'>
                {productDetail.brand.name}
              </span>
            </div>
          )}

          <div className='flex flex-col mb-6'>
            <h1 className='text-3xl font-semibold tracking-tight text-gray-900 mb-2'>
              {productDetail.name}
            </h1>

            <div className='flex items-center justify-between'>
              <div>
                <h2 className='sr-only'>Product price</h2>
                {productDetail.maxDiscountValue && productDetail.maxDiscountValue > 0 ? (
                  <>
                    {minPrice !== maxPrice ? (
                      <div>
                        <p className='text-sm text-gray-500 line-through'>
                          ${minPrice.toFixed(2)} - ${maxPrice.toFixed(2)}
                        </p>
                        <div className='flex flex-col'>
                          <p className='text-2xl font-bold text-emerald-700'>
                            {productDetail.discountType === 'percentage' 
                              ? `$${(minPrice - (minPrice * productDetail.maxDiscountValue / 100)).toFixed(2)} - $${(maxPrice - (maxPrice * productDetail.maxDiscountValue / 100)).toFixed(2)}`
                              : `$${(minPrice - productDetail.maxDiscountValue).toFixed(2)} - $${(maxPrice - productDetail.maxDiscountValue).toFixed(2)}`}
                          </p>
                          <p className='text-sm font-medium text-emerald-600 mt-1'>
                            Save {productDetail.discountType === 'percentage' 
                              ? `$${((minPrice * productDetail.maxDiscountValue / 100)).toFixed(2)} - $${((maxPrice * productDetail.maxDiscountValue / 100)).toFixed(2)}`
                              : `$${productDetail.maxDiscountValue.toFixed(2)}`}
                          </p>
                        </div>
                      </div>
                    ) : (
                      <div>
                        <p className='text-sm text-gray-500 line-through'>
                          ${minPrice.toFixed(2)}
                        </p>
                        <div className='flex flex-col'>
                          <p className='text-2xl font-bold text-emerald-700'>
                            {productDetail.discountType === 'percentage'
                              ? `$${(minPrice - (minPrice * productDetail.maxDiscountValue / 100)).toFixed(2)}`
                              : `$${(minPrice - productDetail.maxDiscountValue).toFixed(2)}`}
                          </p>
                          <p className='text-sm font-medium text-emerald-600 mt-1'>
                            Save {productDetail.discountType === 'percentage'
                              ? `$${(minPrice * productDetail.maxDiscountValue / 100).toFixed(2)}`
                              : `$${productDetail.maxDiscountValue.toFixed(2)}`}
                          </p>
                        </div>
                      </div>
                    )}
                  </>
                ) : (
                  <>
                    {minPrice !== maxPrice ? (
                      <p className='text-2xl font-bold text-gray-900'>
                        ${minPrice.toFixed(2)} - ${maxPrice.toFixed(2)}
                      </p>
                    ) : (
                      <p className='text-2xl font-bold text-gray-900'>
                        ${minPrice.toFixed(2)}
                      </p>
                    )}
                  </>
                )}
              </div>

              {/* Add to wishlist button can be added here */}
            </div>
          </div>

          <div className='mb-6'>
            {/* Reviews */}
            <ReviewSection
              rating={avgRating ?? 0}
              reviewsCount={productDetail.ratingCount}
            />
          </div>

          {/* Short description */}
          {productDetail.shortDescription && (
            <div className='mt-4 text-gray-500'>
              <p>{productDetail.shortDescription}</p>
            </div>
          )}
          {/* Attribute selection with visual improvements */}
          <div className='mb-8 px-6 rounded-lg bg-gray-50 border border-gray-100'>
            <AttributesSection variants={productDetail.variants} />

            {/* Stock info */}
            <div className='mt-4 flex items-center'>
              <div className='flex items-center'>
                <div className={`w-3 h-3 rounded-full bg-green-500 mr-2`}></div>
                <span className='text-sm text-gray-700'>
                  {productDetail.variants.length > 0
                    ? `Multiple options available`
                    : 'In stock'}
                </span>
              </div>

              {/* SKU info */}
              <div className='ml-6 text-sm text-gray-500'>
                SKU: {productDetail.sku}
              </div>
            </div>
          </div>

          {/* Description with better formatting */}
          <div className='mt-8 mb-6'>
            <h3 className='text-lg font-medium text-gray-900 mb-4'>
              Product Description
            </h3>
            <div
              className='prose prose-sm max-w-none text-gray-700'
              dangerouslySetInnerHTML={{ __html: productDetail.description }}
            />
          </div>

          {/* Product details with tabs */}
          <div className='mt-10 pt-6 border-t border-gray-200'>
            <div className='mb-8'>
              <TabsSection productDetail={productDetail} details={details} />
            </div>
          </div>

          {/* Info Boxes with improved design */}
          <div className='mt-8 grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-1 xl:grid-cols-2'>
            {/* International Delivery */}
            <div className='border border-gray-200 hover:border-indigo-300 hover:bg-indigo-50 transition-colors duration-300 rounded-xl p-5 flex items-start'>
              <div className='bg-indigo-100 rounded-full p-3 mr-4'>
                <GlobeAsiaAustraliaIcon className='size-6 text-indigo-600' />
              </div>
              <div className='text-left'>
                <p className='font-medium text-gray-900 mb-1'>
                  Fast worldwide delivery
                </p>
                <p className='text-sm text-gray-600'>
                  Free shipping on orders over $50
                </p>
              </div>
            </div>

            {/* Loyalty Rewards */}
            <div className='border border-gray-200 hover:border-indigo-300 hover:bg-indigo-50 transition-colors duration-300 rounded-xl p-5 flex items-start'>
              <div className='bg-indigo-100 rounded-full p-3 mr-4'>
                <CurrencyDollarIcon className='size-6 text-indigo-600' />
              </div>
              <div className='text-left'>
                <p className='font-medium text-gray-900 mb-1'>
                  Loyalty rewards
                </p>
                <p className='text-sm text-gray-600'>
                  Earn points with every purchase
                </p>
              </div>
            </div>
          </div>

          {/* Discount information */}
          {productDetail.maxDiscountValue && productDetail.maxDiscountValue > 0 && (
            <div className='mt-6 px-6 py-4 rounded-xl bg-gradient-to-r from-emerald-50 to-teal-50 border border-emerald-100 flex items-center shadow-sm'>
              <div className='bg-gradient-to-r from-emerald-100 to-teal-100 rounded-full p-3 mr-4 shadow-inner'>
                <TagIcon className='size-5 text-emerald-700' />
              </div>
              <div className='flex flex-col'>
                <p className='font-semibold text-emerald-800 text-lg'>
                  Limited Time Offer!
                </p>
                <p className='text-sm font-medium text-emerald-700'>
                  {productDetail.discountType === 'percentage' 
                    ? `${productDetail.maxDiscountValue}% off (You save $${(minPrice * productDetail.maxDiscountValue / 100).toFixed(2)})` 
                    : `$${productDetail.maxDiscountValue.toFixed(2)} off your purchase`}
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Reviews section with better separation */}
      <div className='mt-16 pt-10 border-t border-gray-200'>
        <Suspense
          fallback={
            <div className='text-center py-12'>
              <div className='animate-spin rounded-full h-12 w-12 mx-auto border-t-2 border-b-2 border-indigo-500'></div>
            </div>
          }
        >
          <ReviewsList productID={productDetail.id} />
        </Suspense>
      </div>

      {/* Related products with better heading */}
      <div className='my-8 pt-10 border-t border-gray-200'>
        <h2 className='text-2xl font-bold text-gray-900 mb-8'>
          You might also like
        </h2>
        <RelateProductSection />
      </div>
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
