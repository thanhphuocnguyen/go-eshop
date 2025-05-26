import { Metadata } from 'next';
import Image from 'next/image';
import { CartItemSection } from './CartItemSection';

export const metadata: Metadata = {
  title: 'Cart Review',
  description: 'Please review your cart before checkout',
};

export default function CartPage() {
  return (
    <div className='mx-auto flex flex-col'>
      <div className='w-[768px] m-auto py-14'>
        <div className='w-full'>
          <h1 className='text-3xl font-bold text-center mb-12'>
            Shopping Cart
          </h1>
          <hr />
          <CartItemSection />
        </div>
      </div>
      <div className='border-t border-gray-300 m-auto w-full py-32 bg-gray-100 shadow-sm px-20 flex gap-16'>
        <div className='flex flex-col items-center'>
          {/* Image */}
          <Image
            src={'/images/logos/icon-returns-light.svg'}
            alt='Returns'
            width={100}
            height={100}
          />
          <div className='text-lg mb-3 font-bold'>Free returns</div>
          <div className='text-center'>
            Not what you expected? Place it back in the parcel and attach the
            pre-paid postage stamp.
          </div>
        </div>
        <div className='flex flex-col items-center'>
          {/* Image */}
          <Image
            src={'/images/logos/icon-calendar-light.svg'}
            alt='Returns'
            width={100}
            height={100}
          />
          <div className='text-lg mb-3 font-bold'>Same day delivery</div>
          <div className='text-center'>
            We offer a delivery service that has never been done before.
            Checkout today and receive your products within hours.
          </div>
        </div>
        <div className='flex flex-col items-center'>
          {/* Image */}
          <Image
            src={'/images/logos/icon-gift-card-light.svg'}
            alt='Returns'
            width={100}
            height={100}
          />
          <div className='text-lg mb-3 font-bold'>All year discount</div>
          <div className='text-center'>
            Looking for a deal? You can use the code &quot;ALLYEAR &quot; at
            checkout and get money off all year round.
          </div>
        </div>
        <div className='flex mb-3 flex-col items-center'>
          {/* Image */}
          <Image
            src={'/images/logos/icon-planet-light.svg'}
            alt='Returns'
            width={100}
            height={100}
          />
          <div className='text-lg font-bold'>For the planet</div>
          <div className='text-center'>
            We’ve pledged 1% of sales to the preservation and restoration of the
            natural environment.
          </div>
        </div>
      </div>
    </div>
  );
}
