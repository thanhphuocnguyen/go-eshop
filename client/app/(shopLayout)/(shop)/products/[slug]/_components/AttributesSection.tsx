'use client';
import { AttributeDetailModel, VariantDetailModel } from '@/app/lib/definitions';
import { Button } from '@headlessui/react';
import clsx from 'clsx';
import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import { useCartCtx } from '@/app/lib/contexts/CartContext';

interface AttributesSectionProps {
  variants: VariantDetailModel[];
}
const AttributesSection: React.FC<AttributesSectionProps> = ({ variants }) => {
  const { updateCartItemQuantity } = useCartCtx();

  const [attributesData, setAttributesData] = useState<AttributeDetailModel[]>(
    []
  );

  const [selectedAttributeValues, setSelectedAttributeValues] = useState<
    Record<string, string>
  >({});

  // This function checks if selecting a specific attribute value would lead to a valid variant
  const isValueAvailable = (attributeId: string, valueId: string): boolean => {
    // Create a potential selection by keeping existing selections and adding/updating this one
    const potentialSelection = {
      ...selectedAttributeValues,
      [attributeId]: valueId,
    };

    // Get all other attributes that have been selected
    const selectedAttributes = Object.keys(potentialSelection).filter(
      (id) => id !== attributeId
    ); // Exclude the current attribute

    // Check if there's at least one variant that would be valid with this selection
    return variants.some((variant) => {
      // First, check if this variant has the attribute value we're checking
      const hasAttributeValue = variant.attributes.some(
        (attr) => attr.id === attributeId && attr.value_object.id === valueId
      );

      if (!hasAttributeValue) return false;

      // Then check if all other selected attributes match this variant
      const otherAttributesMatch = selectedAttributes.every((attrId) => {
        const selectedValueId = potentialSelection[attrId];
        return variant.attributes.some(
          (attr) =>
            attr.id === attrId && attr.value_object.id === selectedValueId
        );
      });

      // This variant is a potential match if it has the attribute value and matches other selections
      return hasAttributeValue && otherAttributesMatch && variant.stock_qty > 0;
    });
  };

  const handleAddToCart = () => {
    const variant = variants.find((variant) =>
      variant.attributes.every((attribute) => {
        return (
          selectedAttributeValues[attribute.id] === attribute.value_object.id
        );
      })
    );

    if (!variant) {
      toast.error('Please select a size and color');
      return;
    }
    // Add to cart logic here
    updateCartItemQuantity(variant.id, 1);
  };

  useEffect(() => {
    if (variants) {
      const attributes = variants.reduce((acc, variant) => {
        variant.attributes.forEach((attribute) => {
          const existingAttribute = acc.find(
            (attr) => attr.id === attribute.id
          );
          if (existingAttribute) {
            // If the attribute already exists, add the value if it doesn't exist
            if (
              !existingAttribute.values.some(
                (value) => value.id === attribute.value_object.id
              )
            ) {
              existingAttribute.values.push(attribute.value_object);
            }
          } else {
            // If the attribute doesn't exist, create a new one
            acc.push({
              id: attribute.id,
              name: attribute.name,
              created_at: attribute.created_at,
              values: [attribute.value_object],
            });
          }
        });
        return acc;
      }, [] as AttributeDetailModel[]);
      setAttributesData([...attributes]);
    }
  }, [variants]);

  return (
    <div className='flex flex-col gap-6 mt-10'>
      {/* Colors */}
      {attributesData.map((attr) => (
        <div key={attr.id}>
          {attr.name.toLowerCase().includes('color') ? (
            <div>
              <h3 className='text-sm text-gray-900 font-medium'>Color</h3>
              <div className='flex items-center space-x-3 mt-2'>
                {attr.values.map((color) => {
                  const isAvailable = isValueAvailable(attr.id, color.id);
                  return (
                    <Button
                      key={color.id}
                      type='button' // Important: prevent form submission
                      disabled={!isAvailable}
                      style={{
                        outlineColor:
                          selectedAttributeValues[attr.id] === color.id
                            ? color.code
                            : '',
                      }}
                      className={clsx(
                        'relative -m-0.5 flex items-center justify-center rounded-full p-0.5 focus:outline-none',
                        !isAvailable
                          ? 'bg-gray-50 text-gray-200 cursor-not-allowed'
                          : '',
                        selectedAttributeValues[attr.id] === color.id
                          ? `ring-2 ring-offset-2`
                          : ''
                      )}
                      onClick={() => {
                        if (
                          !isAvailable ||
                          color.id === selectedAttributeValues[attr.id]
                        ) {
                          return;
                        }
                        setSelectedAttributeValues((prev) => ({
                          ...prev,
                          [attr.id]: color.id,
                        }));
                      }}
                      aria-label={color.code}
                    >
                      <span
                        aria-hidden='true'
                        className={clsx(
                          `h-8 w-8 border border-black border-opacity-10 rounded-full`,
                          !isAvailable ? 'opacity-40' : ''
                        )}
                        style={{
                          backgroundColor: color.code,
                        }}
                      />
                      {!isAvailable && (
                        <span
                          aria-hidden='true'
                          className='absolute inset-0 rounded-full border-2 border-gray-200 pointer-events-none'
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
                {attr.values.map((value) => {
                  const isAvailable = isValueAvailable(attr.id, value.id);
                  return (
                    <button
                      key={value.id}
                      type='button' // Important: prevent form submission
                      disabled={!isAvailable}
                      className={clsx(
                        'group relative border rounded-md py-3 px-4 flex items-center justify-center text-sm font-medium uppercase hover:bg-gray-50 focus:outline-none sm:flex-1',
                        !isAvailable
                          ? 'bg-gray-50 text-gray-200 cursor-not-allowed'
                          : '',
                        selectedAttributeValues[attr.id] === value.id
                          ? 'bg-indigo-600 border-transparent text-white hover:bg-indigo-700' // Selected style
                          : 'bg-white border-gray-200 text-gray-900 hover:bg-gray-50' // Available style,
                      )}
                      onClick={() => {
                        if (
                          !isAvailable ||
                          value.id === selectedAttributeValues[attr.id]
                        ) {
                          return;
                        }
                        setSelectedAttributeValues((prev) => ({
                          ...prev,
                          [attr.id]: value.id,
                        }));
                      }}
                    >
                      {value.code}
                      {!isAvailable && (
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
                  );
                })}
              </div>
            </div>
          )}
        </div>
      ))}
      <Button
        type='button' // Change to type="submit" if this button submits the form
        onClick={handleAddToCart}
        disabled={
          Object.keys(selectedAttributeValues).length !== attributesData?.length
        } // Disable if no size is selected
        className={`mt-10 w-full flex items-center justify-center rounded-md border border-transparent px-8 py-3 text-base font-medium text-white ${
          Object.keys(selectedAttributeValues).length !== attributesData?.length
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
