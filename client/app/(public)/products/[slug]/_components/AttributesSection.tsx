'use client';
import {
  AttributeValueDetailModel,
  ProductVariantAttributeModel,
} from '@/lib/definitions';
import { Button } from '@headlessui/react';
import clsx from 'clsx';
import React, { useEffect, useState } from 'react';

interface AttributesSectionProps {
  attributes: ProductVariantAttributeModel[];
}
const AttributesSection: React.FC<AttributesSectionProps> = ({
  attributes,
}) => {
  const [selectedColor, setSelectedColor] = useState<number>();
  const [selectedSize, setSelectedSize] = useState<number>();
  const [attributesFormat, setAttributesFormat] = useState<
    Record<string, AttributeValueDetailModel[]>
  >({});
  const handleAddToCart = () => {
    // Add to cart logic here
    console.log('Adding to cart:', {
      size: selectedSize,
      color: selectedColor,
      quantity: 1,
    });
  };

  useEffect(() => {
    const attributeObj: Record<string, AttributeValueDetailModel[]> = {};
    const attributeValueIds = new Set<number>();
    attributes.reduce((acc, attribute) => {
      const { name, value } = attribute;
      const availableValues: AttributeValueDetailModel[] = [];
      if (!attributeValueIds.has(value.id)) {
        attributeValueIds.add(value.id);
        availableValues.push({ ...value });
      }
      if (!acc[name]) {
        acc[name] = availableValues;
      } else {
        acc[name] = acc[name].concat(availableValues);
      }
      return acc;
    }, attributeObj);
    setAttributesFormat(attributeObj);
  }, [attributes]);
  console.log(attributes);

  return (
    <div className='flex flex-col gap-6 mt-10'>
      {/* Colors */}
      {Object.entries(attributesFormat).map(([key, values]) => (
        <div key={key}>
          {key.toLowerCase().includes('color') ? (
            <div>
              <h3 className='text-sm text-gray-900 font-medium'>Color</h3>
              <div className='flex items-center space-x-3 mt-2'>
                {values.map((color) => {
                  return (
                    <Button
                      key={color.id}
                      type='button' // Important: prevent form submission
                      style={{
                        outlineColor:
                          selectedColor === color.id ? color.value : '',
                      }}
                      className={clsx(
                        'relative -m-0.5 flex items-center justify-center rounded-full p-0.5 focus:outline-none',
                        selectedColor === color.id ? `ring-2 ring-offset-2` : ''
                      )}
                      onClick={() => setSelectedColor(color.id)}
                      aria-label={color.display_value}
                    >
                      <span
                        aria-hidden='true'
                        className={clsx(
                          `h-8 w-8 border border-black border-opacity-10 rounded-full`
                        )}
                        style={{
                          backgroundColor: color.value,
                        }}
                      />
                    </Button>
                  );
                })}
              </div>
            </div>
          ) : (
            <div className='mt-4'>
              {/* Sizes */}
              <div className='flex items-center justify-between'>
                <h3 className='text-sm text-gray-900 font-medium'>Size</h3>
                <a
                  href='#'
                  className='text-sm font-medium text-indigo-600 hover:text-indigo-500'
                >
                  See sizing chart
                </a>
              </div>

              <div className='grid grid-cols-4 gap-4 sm:grid-cols-8 lg:grid-cols-6 mt-4'>
                {values.map((value) => (
                  <button
                    key={value.value}
                    type='button' // Important: prevent form submission
                    // disabled={!value.inStock}
                    className={clsx(
                      'group relative border rounded-md py-3 px-4 flex items-center justify-center text-sm font-medium uppercase hover:bg-gray-50 focus:outline-none sm:flex-1',
                      // !value.inStock
                      //   ? 'bg-gray-50 text-gray-200 cursor-not-allowed'
                      //   : '',
                      selectedSize === value.id
                        ? 'bg-indigo-600 border-transparent text-white hover:bg-indigo-700' // Selected style
                        : 'bg-white border-gray-200 text-gray-900 hover:bg-gray-50' // Available style,
                    )}
                    onClick={() => {
                      // if (value.inStock) {
                      setSelectedSize(value.id);
                      // }
                    }}
                  >
                    {value.value}
                    {/* {!value.inStock && (
                      <span
                        aria-hidden='true'
                        className='absolute -inset-px rounded-md border-2 border-gray-200 pointer-events-none'
                      >
                        <svg
                          className='absolute inset-0 w-full h-full text-gray-200 stroke-2'
                          viewBox='0 0 100 100'
                          preserveAspectRatio='none'
                          stroke='currentColor'
                        >
                          <line
                            x1={0}
                            y1={100}
                            x2={100}
                            y2={0}
                            vectorEffect='non-scaling-stroke'
                          />
                        </svg>
                      </span>
                    )} */}
                  </button>
                ))}
              </div>
            </div>
          )}
        </div>
      ))}
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
    </div>
  );
};

export { AttributesSection };
