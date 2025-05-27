'use client';

import { CartModel } from '@/app/lib/definitions';

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

interface OrderTotalSectionProps {
  cart: CartModel | undefined;
  appliedDiscount: DiscountModel | null;
  discountAmount: number;
  shipping?: number;
  taxes?: number;
}

export const OrderTotalSection: React.FC<OrderTotalSectionProps> = ({
  cart,
  appliedDiscount,
  discountAmount,
  shipping = 0,
  taxes = 0.2,
}) => {
  const subtotal = cart?.totalPrice || 0;
  const finalTotal = subtotal + shipping + taxes - discountAmount;

  return (
    <div className='flex flex-col gap-6'>
      <div className='flex justify-between'>
        <span>Subtotal</span>
        <span>${subtotal.toFixed(2)}</span>
      </div>
      <div className='flex justify-between'>
        <span>Shipping</span>
        <span>${shipping.toFixed(2)}</span>
      </div>
      <div className='flex justify-between'>
        <span>Taxes</span>
        <span>${taxes.toFixed(2)}</span>
      </div>
      {appliedDiscount && discountAmount > 0 && (
        <div className='flex justify-between text-green-600'>
          <span>Discount ({appliedDiscount.code})</span>
          <span>-${discountAmount.toFixed(2)}</span>
        </div>
      )}
      <hr className='my-2' />
      <div className='flex justify-between font-semibold text-lg'>
        <span>Total</span>
        <span>${finalTotal.toFixed(2)}</span>
      </div>
    </div>
  );
};
