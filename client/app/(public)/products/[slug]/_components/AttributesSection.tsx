'use client';
import { AttributeFormModel } from '@/lib/definitions';
import React, { useEffect } from 'react';

interface AttributesSectionProps {
  attributes: AttributeFormModel[];
}
const AttributesSection: React.FC<AttributesSectionProps> = ({
  attributes,
}) => {
  const [attributesFormat, setAttributesFormat] = React.useState<
    { name: string; values: { id: number; value: string }[] }[]
  >([]);

  useEffect(() => {
    const prepareAttribute = attributes.map((attribute) => ({
      name: attribute.name,
      values: [],
    }));
    const valueSet
    prepareAttribute.forEach((attribute) => {
      const attributeFormat = {
        name: attribute.name,
        values: attribute.values.map((value) => ({
          id: value.id,
          value: value.display_value,
        })),
      };
      setAttributesFormat((prev) => [...prev, attributeFormat]);
    });
  }, []);
  return (
    <div>
      {product.sizes.map((size) => (
        <button
          key={size.name}
          type='button' // Important: prevent form submission
          disabled={!size.inStock}
          className={`group relative border rounded-md py-3 px-4 flex items-center justify-center text-sm font-medium uppercase hover:bg-gray-50 focus:outline-none sm:flex-1 ${
            !size.inStock
              ? 'bg-gray-50 text-gray-200 cursor-not-allowed' // Disabled style
              : selectedSize?.name === size.name
                ? 'bg-indigo-600 border-transparent text-white hover:bg-indigo-700' // Selected style
                : 'bg-white border-gray-200 text-gray-900 hover:bg-gray-50' // Available style
          }`}
          onClick={() => size.inStock && setSelectedSize(size)}
        >
          {size.name}
          {!size.inStock && (
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
          )}
        </button>
      ))}
    </div>
  );
};

export default AttributesSection;
