'use client';

import { TextField } from '@/app/components/FormFields';
import { Button } from '@headlessui/react';
import { UseFormRegister, UseFormSetValue, UseFormWatch } from 'react-hook-form';
import { CheckoutFormValues } from '../_lib/definitions';
import clsx from 'clsx';

type DiscountModel = {
  id: string;
  code: string;
  description: string;
  discountType: 'percentage' | 'fixed';
  discountValue: number;
  startsAt: string;
  expiresAt: string;
  categoryId?: string;
  productId?: string;
};

interface DiscountSectionProps {
  register: UseFormRegister<CheckoutFormValues>;
  setValue: UseFormSetValue<CheckoutFormValues>;
  watch: UseFormWatch<CheckoutFormValues>;
  discounts: DiscountModel[] | undefined;
  discountsLoading: boolean;
  appliedDiscount: DiscountModel | null;
  discountAmount: number;
  onApplyDiscount: () => void;
  onRemoveDiscount: () => void;
}

export const DiscountSection: React.FC<DiscountSectionProps> = ({
  register,
  setValue,
  watch,
  discounts,
  discountsLoading,
  appliedDiscount,
  discountAmount,
  onApplyDiscount,
  onRemoveDiscount,
}) => {
  return (
    <div className='mb-6'>
      <h4 className='text-sm font-medium text-gray-700 mb-3'>
        Apply Discount
      </h4>
      {discountsLoading ? (
        <div className='flex items-center justify-center p-4 bg-gray-50 rounded-md'>
          <div className='animate-spin rounded-full h-5 w-5 border-b-2 border-indigo-600 mr-2'></div>
          <span className='text-sm text-gray-500'>
            Loading discounts...
          </span>
        </div>
      ) : discounts && discounts.length > 0 ? (
        <div className='space-y-3'>
          <div className='flex gap-2'>
            <TextField
              {...register('discountCode')}
              type='text'
              placeholder='Enter discount code'
              className='flex-1'
            />
            <Button
              type='button'
              onClick={onApplyDiscount}
              className='px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 text-sm'
            >
              Apply
            </Button>
          </div>

          {/* Available Discounts List */}
          <div className='border border-gray-200 rounded-md max-h-48 overflow-y-auto'>
            <div className='p-3 bg-gray-50 border-b border-gray-200'>
              <span className='text-xs font-medium text-gray-600 uppercase tracking-wider'>
                Available Discounts
              </span>
            </div>
            {discounts.map((discount) => (
              <div
                key={discount.id}
                className={clsx(
                  'p-3 border-b border-gray-100 last:border-b-0 hover:bg-gray-50 cursor-pointer transition-colors',
                  watch('discountCode') === discount.code
                    ? 'bg-indigo-50'
                    : ''
                )}
                onClick={() => {
                  setValue('discountCode', discount.code);
                }}
              >
                <div className='flex justify-between items-start'>
                  <div className='flex-1'>
                    <div className='flex items-center gap-2'>
                      <span className='font-mono text-sm font-medium text-indigo-600 bg-indigo-50 px-2 py-1 rounded'>
                        {discount.code}
                      </span>
                      <span className='text-xs bg-green-100 text-green-800 px-2 py-1 rounded-full'>
                        {discount.discountType === 'percentage'
                          ? `${discount.discountValue}% OFF`
                          : `$${discount.discountValue} OFF`}
                      </span>
                    </div>
                    {discount.description && (
                      <p className='text-xs text-gray-500 mt-1'>
                        {discount.description}
                      </p>
                    )}
                    <p className='text-xs text-gray-400 mt-1'>
                      Expires:{' '}
                      {new Date(
                        discount.expiresAt
                      ).toLocaleDateString()}
                    </p>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {/* Applied Discount Display */}
          {appliedDiscount && discountAmount > 0 && (
            <div className='flex justify-between text-green-600 p-2 bg-green-50 rounded-md'>
              <div className='flex items-center gap-2'>
                <span>Discount ({appliedDiscount.code})</span>
                <Button
                  type='button'
                  onClick={onRemoveDiscount}
                  className='text-red-500 hover:text-red-700 text-xs'
                >
                  ✕
                </Button>
              </div>
              <span>-${discountAmount.toFixed(2)}</span>
            </div>
          )}
        </div>
      ) : (
        <div className='p-4 bg-gray-50 rounded-md'>
          <p className='text-sm text-gray-500 text-center'>
            No discounts available at the moment
          </p>
        </div>
      )}
    </div>
  );
};
