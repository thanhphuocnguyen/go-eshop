'use client';
import { Button } from '@headlessui/react';
import React from 'react';

const AddToCart = () => {
  const [selectedSize, setSelectedSize] = React.useState<string | null>(null);
  const [selectedColor, setSelectedColor] = React.useState<string | null>(null);
  const [quantity, setQuantity] = React.useState<number>(1);
  const handleAddToCart = () => {
    // Add to cart logic here
    console.log('Adding to cart:', {
      size: selectedSize,
      color: selectedColor,
      quantity,
    });
  };

  return (
    <Button
      type='button' // Change to type="submit" if this button submits the form
      onClick={handleAddToCart}
      disabled={!selectedSize} // Disable if no size is selected
      className={`mt-10 w-full flex items-center justify-center rounded-md border border-transparent px-8 py-3 text-base font-medium text-white ${
        !selectedSize
          ? 'bg-gray-400 cursor-not-allowed'
          : 'bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500'
      }`}
    >
      Add to cart
    </Button>
  );
};

export default AddToCart;
