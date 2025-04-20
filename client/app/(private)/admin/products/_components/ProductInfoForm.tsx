'use client';
import React, { useEffect, useState } from 'react';
import { TextField } from '@/components/FormFields';
import { LoadingSpinner } from '@/components/Common/Loadings/Loading';
import { StyledComboBoxController } from '@/components/FormFields/StyledComboBoxController';
import Image from 'next/image';
import { StyledMultipleComboBox } from '@/components/FormFields/StyledMultipleComboBox';
import { ImageUploader } from '@/components/FormFields/ImageUploader';
import { TiptapController } from '@/components/Common';
import { XMarkIcon } from '@heroicons/react/16/solid';
import { Field, Label, Switch } from '@headlessui/react';
import { ArrowTurnUpLeftIcon } from '@heroicons/react/24/outline';
import { useCollections } from '../../_lib/hooks/useCollections';
import { useBrands } from '../../_lib/hooks/useBrands';
import { useAttributes } from '../../_lib/hooks/useAttributes';
import { useCategories } from '../../_lib/hooks/useCategories';
import {
  AttributeDetailModel,
  ProductDetailModel,
  ProductModelForm,
} from '@/lib/definitions';
import { useFormContext } from 'react-hook-form';
import clsx from 'clsx';
import { useProductDetailFormContext } from '../_lib/contexts/ProductFormContext';

export const ProductInfoForm: React.FC<{
  productDetail?: ProductDetailModel;
}> = ({ productDetail }) => {
  const { categories, isLoading: categoriesLoading } = useCategories();
  const { collections, isLoading: collectionLoading } = useCollections();
  const { brands, isLoading: brandsLoading } = useBrands();
  const { attributes, attributesLoading } = useAttributes();

  const [selectedAttributes, setSelectedAttributes] = useState<
    AttributeDetailModel[]
  >([]);

  const { productImages, setProductImages } = useProductDetailFormContext();

  const { register, control, getValues, watch, formState, setValue } =
    useFormContext<ProductModelForm>();

  useEffect(() => {
    if (!productDetail && categories && categories.length > 0) {
      setValue('category', categories[0]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categories]);
  useEffect(() => {
    if (!productDetail && brands && brands.length > 0) {
      setValue('brand', brands[0]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [brands]);

  return (
    <>
      <div className='flex gap-4 items-center mb-4'>
        <h2 className='text-xl font-bold text-primary'>Product Information</h2>
        <Field className='flex items-center gap-2'>
          <Switch
            checked={watch('is_active')}
            onChange={(value) => setValue('is_active', value)}
            className={({ checked }) =>
              clsx(
                'relative inline-flex h-6 w-11 items-center rounded-full',
                checked ? 'bg-primary' : 'bg-gray-200'
              )
            }
          >
            {({ checked }) => (
              <span
                className={clsx(
                  'inline-block h-4 w-4 transform rounded-full bg-white transition',
                  checked ? 'translate-x-6' : 'translate-x-1'
                )}
              />
            )}
          </Switch>
          <Label htmlFor='is_active' className='font-semibold cursor-pointer'>
            Active
          </Label>
        </Field>
      </div>

      {/* Basic Information */}
      <div className='grid grid-cols-4 gap-4 mb-6'>
        <TextField
          label={'Product name'}
          {...register('name')}
          error={!!formState.errors.name}
          message={formState.errors.name?.message}
          placeholder='Enter product name...'
          type='text'
          required
        />
        <TextField
          {...register('sku')}
          label={'Sku'}
          placeholder='Enter sku...'
          type='text'
          error={!!formState.errors.sku}
          message={formState.errors.sku?.message}
        />
        <TextField
          {...register('price', {
            valueAsNumber: true,
          })}
          label={'Price'}
          placeholder='Enter price...'
          type='number'
          error={!!formState.errors.price}
          message={formState.errors.price?.message}
        />
        <TextField
          label={'Slug'}
          placeholder='Enter slug...'
          type='text'
          error={!!formState.errors.slug}
          message={formState.errors.slug?.message}
          {...register('slug')}
        />
        {attributesLoading ? (
          <div className='flex justify-center items-center'>
            <LoadingSpinner />
          </div>
        ) : attributes ? (
          <StyledMultipleComboBox<AttributeDetailModel>
            label='Select an attribute'
            setSelected={(values) => {
              setSelectedAttributes(values);
            }}
            options={attributes}
            getDisplayValue={(option) => {
              return option?.name || '';
            }}
            selected={selectedAttributes}
          />
        ) : null}
        {/* Category, Collections, Brand */}
        {categoriesLoading ? (
          <div className='flex justify-center items-center'>
            <LoadingSpinner />
          </div>
        ) : (
          <StyledComboBoxController
            control={control}
            name='category'
            label='Category'
            error={!!formState.errors.category}
            message={formState.errors.category?.message ?? ''}
            options={
              categories?.map((e) => ({
                id: e.id,
                name: e.name,
              })) ?? []
            }
          />
        )}
        {brandsLoading ? (
          <div className='flex justify-center items-center'>
            <LoadingSpinner />
          </div>
        ) : (
          <StyledComboBoxController
            name='brand'
            nullable
            control={control}
            error={!!formState.errors.brand}
            message={formState.errors.brand?.message ?? ''}
            label='Brand'
            options={
              brands?.map((e) => ({
                id: e.id,
                name: e.name,
              })) ?? []
            }
          />
        )}
        {collectionLoading ? (
          <div className='flex justify-center items-center'>
            <LoadingSpinner />
          </div>
        ) : (
          <StyledComboBoxController
            control={control}
            name='collection'
            nullable
            label='Collection'
            error={!!formState.errors.brand}
            message={formState.errors.brand?.message ?? ''}
            options={
              collections?.map((e) => ({
                id: e.id,
                name: e.name,
              })) ?? []
            }
          />
        )}
      </div>
      {/* Description */}
      <Field className='w-full'>
        <Label className='font-semibold'>Description</Label>
        <TiptapController
          name='description'
          control={control}
          error={!!formState.errors.description}
          message={formState.errors.description?.message as string}
        />
      </Field>

      <div className='mt-6'>
        <ImageUploader
          label='Media'
          multiple={true}
          onUpload={(files) => {
            setProductImages(files);
          }}
        />

        {/* Preview uploaded images */}
        <div className='grid grid-cols-4 gap-4 mt-4'>
          {productDetail?.images.map((image, index) => {
            const isRemoved = watch('removed_images')?.includes(image.id);
            return (
              <div key={index} className='relative rounded-md h-60 w-full'>
                <Image
                  fill
                  src={image.url}
                  objectFit='cover'
                  alt={`Product image ${index + 1}`}
                  className={clsx(isRemoved ? 'opacity-30' : '', 'rounded-md')}
                />
                {isRemoved && (
                  <div className='absolute inset-0 flex items-center justify-center bg-black bg-opacity-20 z-10'>
                    <span className='bg-red-500 text-white px-2 py-1 rounded'>
                      Removed
                    </span>
                  </div>
                )}
                <button
                  type='button'
                  onClick={() => {
                    if (isRemoved) {
                      // Remove from removed_images if it's already there
                      setValue(
                        'removed_images',
                        (getValues('removed_images') || []).filter(
                          (id) => id !== image.id
                        ),
                        { shouldDirty: true }
                      );
                    } else {
                      // Add to removed_images
                      setValue(
                        'removed_images',
                        getValues('removed_images')
                          ? [...getValues('removed_images')!, image.id]
                          : [image.id],
                        { shouldDirty: true }
                      );
                    }
                  }}
                  className={`absolute z-20 top-1 right-1 ${isRemoved ? 'bg-green-500' : 'bg-red-500'} text-white rounded-full p-1 w-6 h-6 flex items-center justify-center`}
                >
                  {isRemoved ? (
                    <ArrowTurnUpLeftIcon className='size-8' />
                  ) : (
                    <XMarkIcon className='size-8' />
                  )}
                </button>
              </div>
            );
          })}
          {productImages.map((file, index) => (
            <div key={index} className='relative rounded-md h-60'>
              <Image
                fill
                src={URL.createObjectURL(file)}
                objectFit='cover'
                alt={`Product image ${index + 1}`}
                className='w-full h-full object-contain'
              />
              <button
                type='button'
                onClick={() => {
                  setProductImages((prev) =>
                    prev.filter((_, idx) => idx !== index)
                  );
                }}
                className='absolute top-1 right-1 bg-red-500 text-white rounded-full p-1 w-6 h-6 flex items-center justify-center'
              >
                <XMarkIcon className='size-8' />
              </button>
            </div>
          ))}
        </div>
      </div>
    </>
  );
};
