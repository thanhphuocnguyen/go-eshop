'use client';
import { ProductModelForm } from '@/lib/definitions';
import React from 'react';
import { useFormContext } from 'react-hook-form';
import { TextField } from '@/components/FormFields';
import { ImageUploadForm } from '@/components/FormFields';
import {
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
  Transition,
  Switch,
} from '@headlessui/react';
import { ChevronUpIcon } from '@heroicons/react/24/outline';
import { StyledMultipleComboBoxController } from '@/components/FormFields/StyledMultipleComboBoxController';
import { useProductDetailFormContext } from '../_lib/contexts/ProductFormContext';
import { useAttributes } from '../../_lib/hooks/useAttributes';

// Define the UploadFile type locally
type UploadFile = File & {
  preview: string;
};

interface AttributeFormProps {
  index: number;
  onRemove: (index: number) => void;
}
export const VariantForm: React.FC<AttributeFormProps> = ({
  index,
  onRemove,
}) => {
  const { control, register, getValues, setValue, watch } =
    useFormContext<ProductModelForm>();
  const { productVariantImages, setProductVariantImages } =
    useProductDetailFormContext();
  const { attributes } = useAttributes();

  // Function to convert File to UploadFile
  const convertToUploadFile = (file: File | null): UploadFile | null => {
    if (!file) return null;
    return Object.assign(file, {
      preview: URL.createObjectURL(file),
    }) as UploadFile;
  };

  // Watch the is_active value to use with the Switch component
  const isActive = watch(`variants.${index}.is_active`);

  return (
    <div className='w-full mb-4 border border-gray-200 rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-300'>
      <Disclosure defaultOpen={true}>
        {({ open }) => (
          <div className='w-full'>
            <DisclosureButton
              as='div'
              className={`flex w-full justify-between ${open ? 'bg-blue-50' : 'bg-gray-100 hover:bg-gray-50'} px-4 py-3 text-left text-lg font-semibold focus:outline-none focus-visible:ring focus-visible:ring-primary focus-visible:ring-opacity-75 transition-all duration-300`}
            >
              <span className='transition-colors duration-300'>{`Variant ${index + 1}`}</span>
              <div className='flex items-center'>
                <button
                  type='button'
                  className='btn btn-danger btn-sm mr-3 transition-all duration-300 hover:scale-105'
                  onClick={(e) => {
                    e.stopPropagation();
                    onRemove(index);
                  }}
                >
                  Remove
                </button>
                <ChevronUpIcon
                  className={`${
                    open ? 'rotate-180' : 'rotate-0'
                  } h-5 w-5 text-gray-500 transition-transform duration-300 ease-in-out`}
                />
              </div>
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
              <DisclosurePanel
                className={`${open ? 'disclosure-animate-open' : 'disclosure-animate-close'}`}
              >
                <div className='px-4 pt-4 pb-4 animate-fadeIn'>
                  <div className='grid grid-cols-2 gap-4 mb-4'>
                    <TextField
                      {...register(`variants.${index}.sku`)}
                      label='SKU'
                      placeholder='Enter SKU'
                      type='text'
                      disabled={false}
                    />
                    <TextField
                      {...register(`variants.${index}.price`)}
                      label='Price'
                      placeholder='Enter Price'
                      type='number'
                      disabled={false}
                    />
                    <TextField
                      {...register(`variants.${index}.stock`)}
                      label='Stock'
                      placeholder='Enter Stock'
                      type='number'
                      disabled={false}
                    />
                    <TextField
                      {...register(`variants.${index}.weight`)}
                      label='Weight'
                      placeholder='Enter Weight'
                      type='number'
                      disabled={false}
                    />

                    {getValues(`variants.${index}.attributes`)?.map(
                      (attribute, idx) =>
                        attribute ? (
                          <StyledMultipleComboBoxController
                            control={control}
                            key={attribute.id}
                            getDisplayValue={(attribute) =>
                              attribute.display_value || attribute.value
                            }
                            name={`variants.${index}.attributes.${idx}.values`}
                            label={attribute.name}
                            options={
                              attributes?.find((e) => e.id === attribute.id)
                                ?.values ?? []
                            }
                          />
                        ) : null
                    )}
                  </div>
                  <div className='flex items-center gap-2 mt-2'>
                    <Switch
                      checked={!!isActive}
                      onChange={(checked) =>
                        setValue(`variants.${index}.is_active`, checked)
                      }
                      className={`${
                        isActive ? 'bg-primary' : 'bg-gray-200'
                      } relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2`}
                    >
                      <span className='sr-only'>Active Variant</span>
                      <span
                        className={`${
                          isActive ? 'translate-x-6' : 'translate-x-1'
                        } inline-block h-4 w-4 transform rounded-full bg-white transition-transform`}
                      />
                    </Switch>
                    <span
                      className='font-semibold cursor-pointer'
                      onClick={() =>
                        setValue(`variants.${index}.is_active`, !isActive)
                      }
                    >
                      Active Variant
                    </span>
                  </div>
                  <div className='flex mt-5 justify-center'>
                    <ImageUploadForm
                      label={`Variant ${index + 1} Image`}
                      imageUrl={getValues(`variants.${index}.image_url`)}
                      onFileChange={(file: File | null) => {
                        const newFiles = [...productVariantImages];
                        const variantId = getValues(`variants.${index}.id`);
                        newFiles[index] = {
                          ...newFiles[index],
                          image: convertToUploadFile(file),
                          variantID:
                            typeof variantId === 'number' ? variantId : -1,
                        };
                        setProductVariantImages(newFiles);
                      }}
                    />
                  </div>
                </div>
              </DisclosurePanel>
            </Transition>
          </div>
        )}
      </Disclosure>
    </div>
  );
};
