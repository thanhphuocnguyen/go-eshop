'use client';

import { useState, useEffect, useCallback } from 'react';
import { CartItem } from '@/app/lib/definitions/cart';
import { UseFormSetValue } from 'react-hook-form';
import { CheckoutFormValues } from '../../_lib/definitions';

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

interface UseDiscountLogicProps {
  selectedDiscountCode: string | undefined;
  discounts: DiscountModel[] | undefined;
  cart: any;
  setValue: UseFormSetValue<CheckoutFormValues>;
}

export const useDiscountLogic = ({
  selectedDiscountCode,
  discounts,
  cart,
  setValue,
}: UseDiscountLogicProps) => {
  const [appliedDiscount, setAppliedDiscount] = useState<DiscountModel | null>(
    null
  );
  const [discountAmount, setDiscountAmount] = useState(0);

  // Calculate discount amount for cart items
  const calculateDiscountAmount = useCallback(
    (discount: DiscountModel, cartItems: CartItem[]) => {
      if (!cartItems || cartItems.length === 0) return 0;

      let totalDiscountAmount = 0;

      // Loop through each cart item to check for product-specific discounts
      for (const item of cartItems) {
        let itemSubtotal = item.price * item.quantity;
        let itemDiscountAmount = 0;

        // Check if discount applies to this specific product
        if (discount.productId && discount.productId === item.productId) {
          // Product-specific discount
          if (discount.discountType === 'percentage') {
            itemDiscountAmount = (itemSubtotal * discount.discountValue) / 100;
          } else {
            itemDiscountAmount = Math.min(discount.discountValue, itemSubtotal);
          }
        }
        // Check if discount applies to this product's category
        else if (
          discount.categoryId &&
          discount.categoryId === item.categoryId
        ) {
          // Category-specific discount
          if (discount.discountType === 'percentage') {
            itemDiscountAmount = (itemSubtotal * discount.discountValue) / 100;
          } else {
            itemDiscountAmount = Math.min(discount.discountValue, itemSubtotal);
          }
        }

        totalDiscountAmount += itemDiscountAmount;
      }

      return totalDiscountAmount;
    },
    []
  );

  // Apply discount function
  const handleApplyDiscount = useCallback(() => {
    const discountCode = selectedDiscountCode?.trim();
    if (!discountCode || !discounts) return;

    const foundDiscount = discounts.find(
      (discount) => discount.code.toLowerCase() === discountCode.toLowerCase()
    );

    if (foundDiscount) {
      const cartItems = cart?.cartItems || [];
      const amount = calculateDiscountAmount(foundDiscount, cartItems);
      setAppliedDiscount(foundDiscount);
      setDiscountAmount(amount);
      return {
        success: true,
        message: `Discount "${foundDiscount.code}" applied successfully!`,
      };
    } else {
      setAppliedDiscount(null);
      setDiscountAmount(0);
      return { success: false, message: 'Invalid discount code' };
    }
  }, [
    selectedDiscountCode,
    discounts,
    cart?.cartItems,
    calculateDiscountAmount,
  ]);

  // Remove discount function
  const handleRemoveDiscount = useCallback(() => {
    setAppliedDiscount(null);
    setDiscountAmount(0);
    setValue('discountCode', '');
    return { success: true, message: 'Discount removed' };
  }, [setValue]);

  // Auto-apply discount when discount code changes
  useEffect(() => {
    if (selectedDiscountCode && discounts) {
      const foundDiscount = discounts.find(
        (discount) =>
          discount.code.toLowerCase() === selectedDiscountCode.toLowerCase()
      );

      if (foundDiscount && appliedDiscount?.code !== foundDiscount.code) {
        const cartItems = cart?.cartItems || [];
        const amount = calculateDiscountAmount(foundDiscount, cartItems);
        setAppliedDiscount(foundDiscount);
        setDiscountAmount(amount);
      } else if (!foundDiscount && appliedDiscount) {
        setAppliedDiscount(null);
        setDiscountAmount(0);
      }
    } else if (!selectedDiscountCode && appliedDiscount) {
      setAppliedDiscount(null);
      setDiscountAmount(0);
    }
  }, [
    selectedDiscountCode,
    discounts,
    cart?.totalPrice,
    appliedDiscount,
    calculateDiscountAmount,
  ]);

  return {
    appliedDiscount,
    discountAmount,
    handleApplyDiscount,
    handleRemoveDiscount,
  };
};
