import { Button } from '@headlessui/react';
import React from 'react';
import { VariantForm } from './VariantForm';
import { useFieldArray, useFormContext } from 'react-hook-form';
import { ProductModelForm } from '@/app/lib/definitions';
import { PlusIcon } from '@heroicons/react/16/solid';
import clsx from 'clsx';
import { ChevronUpIcon } from '@heroicons/react/24/outline';
import {
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
  Transition,
} from '@headlessui/react';
import { useAttributes } from '../../hooks';

export const VariantInfoForm: React.FC<{
  selectedAttributes: string[];
}> = ({ selectedAttributes }) => {
  const { control, getValues } = useFormContext<ProductModelForm>();
  const { attributes } = useAttributes();

  const { fields, append, remove } = useFieldArray({
    name: 'variants',
    keyName: 'key',
    control,
  });

  return (
    <div>
      <h2 className='text-xl mb-3 font-bold text-primary'>Product Variants</h2>
      {/* Variants Accordion */}
      <Disclosure defaultOpen={true}>
        {({ open }) => (
          <div className='w-full border border-gray-200 rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-300 mb-4'>
            <DisclosureButton
              className={clsx(
                'flex w-full justify-between px-4 py-3 text-left text-lg font-semibold focus:outline-none focus-visible:ring focus-visible:ring-primary focus-visible:ring-opacity-75 transition-all duration-300',
                open ? 'bg-blue-50' : 'bg-gray-100 hover:bg-gray-50'
              )}
            >
              <span>All Variants ({fields.length})</span>
              <ChevronUpIcon
                className={clsx(
                  'h-5 w-5 text-gray-500 transition-transform duration-300 ease-in-out',
                  open ? 'rotate-180' : 'rotate-0'
                )}
              />
            </DisclosureButton>
            <Transition
              show={open}
              enter='transition duration-300 ease-out'
              enterFrom='transform scale-95 opacity-0'
              enterTo='transform scale-100 opacity-100'
              leave='transition duration-200 ease-out'
              leaveFrom='transform scale-100 opacity-100'
              leaveTo='transform scale-95 opacity-0'
            >
              <DisclosurePanel className='p-4'>
                <div className='grid grid-cols-1 md:grid-cols-2 gap-4 mb-4'>
                  {fields.map((item, index) => (
                    <div
                      key={item.key}
                      className='border border-gray-200 rounded-lg p-4 bg-white'
                    >
                      <VariantForm
                        selectedAttributes={selectedAttributes}
                        index={index}
                        item={item}
                        onRemove={() => {
                          remove(index);
                        }}
                      />
                    </div>
                  ))}
                </div>

                {attributes && (
                  <Button
                    onClick={() => {
                      append({
                        attributes: selectedAttributes
                          ? selectedAttributes.map((e) => {
                              const attribute = attributes.find(
                                (attr) => attr.id === e
                              )!;
                              return {
                                id: attribute.id,
                                name: attribute.name,
                                valueObject: {
                                  id: undefined,
                                },
                              };
                            })
                          : [],
                        price: getValues('productInfo.price'),
                        stockQty: 0,
                        weight: undefined,
                        isActive: true,
                      });
                    }}
                    className={clsx('btn btn-primary flex gap-2')}
                  >
                    <PlusIcon className='size-6' />
                    Add Variant
                  </Button>
                )}
              </DisclosurePanel>
            </Transition>
          </div>
        )}
      </Disclosure>
    </div>
  );
};
