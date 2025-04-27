import React from 'react';

const Footer = () => {
  return (
    <footer className='pt-16 pb-12 px-20 text-white bg-primary flex flex-col content-between border-t border-gray-400 shadow-sm'>
      <div className='grid grid-cols-6 gap-4'>
        <div className=''>
          <div>Shop</div>
          <div></div>
        </div>
        <div className=''>
          <div>Company</div>
        </div>
        <div className=''>
          <div>Account</div>
        </div>
        <div className='flex flex-col gap-5'>
          <div className='font-semibold'>Connect</div>
          <div>Contact us</div>
          <div>Facebook</div>
          <div>Twitter</div>
          <div>Instagram</div>
        </div>
        <div className='col-span-2'>
          <div className='mb-8'>Sign up for our newsletter</div>
          <div>The latest deals and savings, sent to your inbox weekly.</div>
          <div className='flex gap-2 mt-4'>
            <input
              type='text'
              placeholder='Enter your email'
              className='w-72 px-3 py-2 rounded-md'
            />
            <button className='bg-button-tertiary text-green-700 rounded-md px-4'>
              Subscribe
            </button>
          </div>
        </div>
      </div>
      <div className='mt-40'>Copyright © 2021 Your Company, Inc.</div>
    </footer>
  );
};

export default Footer;
