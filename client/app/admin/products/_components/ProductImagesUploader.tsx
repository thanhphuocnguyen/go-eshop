'use client';
import { useEffect, useState } from 'react';
import { XCircleIcon } from '@heroicons/react/16/solid';
import { ChevronUpIcon } from '@heroicons/react/24/outline';
import { ImageUploader, StyledComboBox } from '@/app/components/FormFields';
import { StyledMultipleComboBox } from '@/app/components/FormFields/StyledMultipleComboBox';
import Image from 'next/image';
import { ProductDetailModel, ProductModelForm } from '@/app/lib/definitions';
import {
  useProductDetailFormContext,
  VariantImage,
} from '../_lib/contexts/ProductFormContext';
import { useFormContext, useWatch } from 'react-hook-form';
import {
  Disclosure,
  DisclosureButton,
  DisclosurePanel,
  Transition,
} from '@headlessui/react';
import clsx from 'clsx';
import { capitalize } from '@/app/lib/utils';

interface ProductImagesUploaderProps {
  productDetail?: ProductDetailModel;
}
type VariantToAssignOption = {
  id: string;
  name: string;
  disabled?: boolean;
};

export const ProductImagesUploader: React.FC<ProductImagesUploaderProps> = (
  props
) => {
  const { productDetail } = props;
  const { setValue, control } = useFormContext<ProductModelForm>();

  const productImages = useWatch({
    control,
    name: 'productInfo.images',
  });

  const { tempProductImages, setTempProductImages } =
    useProductDetailFormContext();
  const [variantOptions, setVariantOptions] = useState<VariantToAssignOption[]>(
    []
  );

  // Handle image upload
  const handleImageUpload = (files: (File & { preview: string })[]) => {
    const newImages: VariantImage[] = files.map((file) => ({
      file,
      preview: file.preview,
      variantIds: [],
    }));

    setTempProductImages([...tempProductImages, ...newImages]);
  };

  // Handle image removal
  const handleRemoveImage = (index: number) => {
    setTempProductImages((prev) => prev.filter((_, idx) => idx !== index));
  };

  // Handle variant assignment for an image
  const handleAssignVariantsForUploadingFile = (
    index: number,
    selectedVariantIds: string[]
  ) => {
    setTempProductImages((prev) => {
      const updated = [...prev];
      updated[index] = {
        ...updated[index],
        variantIds: selectedVariantIds,
      };

      return updated;
    });
  };

  useEffect(() => {
    if (productDetail?.variants) {
      const options = productDetail.variants.reduce((acc, variant) => {
        acc.push({
          id: variant.id,
          name: variant.attributes
            .map(
              (attr) =>
                `${attr.name}: ${attr.valueObject.name || attr.valueObject.code}`
            )
            .join(', '),
        });
        return acc;
      }, [] as VariantToAssignOption[]);

      setVariantOptions(options);
    }
    // Get variant options for the multiple select
  }, [productDetail?.variants, productImages]);

  return (
    <div>
      <Disclosure defaultOpen={true}>
        {({ open }) => (
          <div className='w-full border border-gray-200 rounded-lg overflow-hidden shadow-sm hover:shadow-md transition-shadow duration-300 mb-4'>
            <DisclosureButton
              className={clsx(
                'flex w-full justify-between px-4 py-3 text-left text-lg font-semibold focus:outline-none focus-visible:ring focus-visible:ring-primary focus-visible:ring-opacity-75 transition-all duration-300',
                open ? 'bg-blue-50' : 'bg-gray-100 hover:bg-gray-50'
              )}
            >
              <span className='text-2xl font-semibold text-primary'>
                Product Media
              </span>
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
                <p className='text-sm text-gray-500 mb-4'>
                  Upload images for variants. Each image can be assigned to
                  multiple variants.
                </p>
                <div className='mb-8'>
                  {/* Image uploader */}
                  <div className='mb-2'>
                    <ImageUploader
                      label='Upload Variant Images'
                      multiple={true}
                      onUpload={handleImageUpload}
                    />
                  </div>

                  {/* Uploaded images with variant assignment */}
                  <div className='grid grid-cols-1 gap-6 mt-4'>
                    {!!productImages.length
                      ? productImages.map((image, index) =>
                          image.isRemoved ? null : (
                            <div
                              key={image.id}
                              className='border border-gray-200 rounded-lg p-4 bg-white shadow-sm'
                            >
                              <ImagePreview
                                handleSelectRole={(idx, role) => {
                                  setValue(
                                    `productInfo.images.${idx}.role`,
                                    role,
                                    {
                                      shouldDirty: true,
                                    }
                                  );
                                }}
                                url={image.url}
                                variantIds={image.assignments}
                                index={index}
                                variantOptions={variantOptions}
                                selectedRole={image.role}
                                onRemove={() => {
                                  setValue(
                                    `productInfo.images.${index}.isRemoved`,
                                    true,
                                    {
                                      shouldDirty: true,
                                    }
                                  );
                                }}
                                onAssignVariants={(idx, values) => {
                                  setValue(
                                    `productInfo.images.${idx}.assignments`,
                                    values,
                                    {
                                      shouldDirty: true,
                                    }
                                  );
                                }}
                              />
                            </div>
                          )
                        )
                      : null}
                    {tempProductImages.map((image, index) => (
                      <div
                        key={index}
                        className='border border-gray-200 rounded-lg p-4 bg-white shadow-sm'
                      >
                        <ImagePreview
                          handleSelectRole={(idx, role) => {
                            setTempProductImages((prev) => {
                              const updated = [...prev];
                              updated[idx] = {
                                ...updated[idx],
                                role: role,
                              };

                              return updated;
                            });
                          }}
                          url={image.preview}
                          selectedRole={image.role || 'thumbnail'}
                          variantIds={image.variantIds}
                          index={index}
                          variantOptions={variantOptions}
                          onRemove={handleRemoveImage}
                          onAssignVariants={
                            handleAssignVariantsForUploadingFile
                          }
                        />
                      </div>
                    ))}

                    {!productDetail?.productImages.length &&
                      !tempProductImages.length && (
                        <div className='text-center p-4 border border-dashed border-gray-300 rounded-lg bg-gray-50'>
                          <p className='text-gray-500'>
                            No variant images uploaded yet
                          </p>
                        </div>
                      )}
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

interface ImagePreviewProps {
  url: string;
  index: number;
  variantIds: string[];
  selectedRole?: string | null;
  variantOptions: { id: string; name: string }[];
  handleSelectRole: (idx: number, role: string | undefined | null) => void;
  onRemove: (index: number) => void;
  onAssignVariants: (idx: number, ids: string[]) => void;
}

const ImagePreview: React.FC<ImagePreviewProps> = ({
  url,
  onRemove,
  onAssignVariants,
  handleSelectRole,
  variantOptions,
  variantIds,
  selectedRole,
  index,
}) => {
  return (
    <div className='flex items-start space-x-4'>
      {/* Image preview */}
      <div className='relative h-32 w-32 flex-shrink-0'>
        <Image
          src={url}
          alt={`Variant image ${index + 1}`}
          fill
          className='object-cover rounded-md'
        />
      </div>

      {/* Variant assignment */}
      <div className='flex-1'>
        <div className='flex justify-between mb-2'>
          <h4 className='font-medium text-gray-700'>Image {index + 1}</h4>
          <button
            type='button'
            onClick={() => onRemove(index)}
            className='text-red-500 hover:text-red-700 transition-colors'
          >
            <XCircleIcon className='size-6' />
          </button>
        </div>
        <div className='flex items-center gap-6 mb-2'>
          <div className='w-2/3'>
            <StyledMultipleComboBox
              label='Assign to Variants'
              selected={variantOptions?.filter((opt) =>
                variantIds.includes(opt.id)
              )}
              options={variantOptions}
              getDisplayValue={(opt) => opt.name}
              setSelected={(selected) =>
                onAssignVariants(
                  index,
                  selected.map((s) => s.id)
                )
              }
            />
          </div>
          <div className='w-1/3'>
            <StyledComboBox
              label='Image Role'
              getKey={(opt) => opt || ''}
              selected={selectedRole}
              getDisplayValue={(opt) => (opt ? capitalize(opt) : '')}
              setSelected={(role) => {
                handleSelectRole(index, role);
              }}
              options={['thumbnail', 'gallery', 'additional']}
            />
          </div>
        </div>

        {variantIds.length === 0 && (
          <p className='text-sm text-amber-600 mt-1'>
            This image is not assigned to any variants
          </p>
        )}
      </div>
    </div>
  );
};
