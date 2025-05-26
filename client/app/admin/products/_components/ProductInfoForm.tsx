'use client';
import React, { useEffect } from 'react';
import { TextField } from '@/app/components/FormFields';
import { LoadingSpinner } from '@/app/components/Common/Loadings/Loading';
import { StyledComboBoxController } from '@/app/components/FormFields/StyledComboBoxController';
import { StyledMultipleComboBox } from '@/app/components/FormFields/StyledMultipleComboBox';
import { TiptapController } from '@/app/components/Common';
import { Field, Label, Switch } from '@headlessui/react';
import { useCollections } from '../../hooks/useCollections';
import { useBrands } from '../../hooks/useBrands';
import { useAttributes } from '../../hooks/useAttributes';
import { useCategories } from '../../hooks/useCategories';
import { ProductDetailModel, ProductModelForm } from '@/app/lib/definitions';
import { useFormContext } from 'react-hook-form';
import clsx from 'clsx';
import { ProductImagesUploader } from './ProductImagesUploader';

export const ProductInfoForm: React.FC<{
  productDetail?: ProductDetailModel;
}> = ({ productDetail }) => {
  const { categories, isLoading: categoriesLoading } = useCategories();
  const { collections, isLoading: collectionLoading } = useCollections();
  const { brands, isLoading: brandsLoading } = useBrands();
  const { attributes, attributesLoading } = useAttributes();

  const { register, control, watch, formState, setValue } =
    useFormContext<ProductModelForm>();

  useEffect(() => {
    if (!productDetail && categories && categories.length > 0) {
      setValue('productInfo.category', categories[0], {
        shouldDirty: false,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categories]);

  useEffect(() => {
    if (!productDetail && collections && collections.length > 0) {
      setValue('productInfo.collection', collections[0], {
        shouldDirty: false,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [collections]);

  useEffect(() => {
    if (!productDetail && brands && brands.length > 0) {
      setValue('productInfo.brand', brands[0], {
        shouldDirty: false,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [brands]);

  return (
    <>
      <div className='flex gap-4 items-center mb-4'>
        <h2 className='text-xl font-bold text-primary'>Product Information</h2>
        <Field className='flex items-center gap-2'>
          <Switch
            checked={watch('productInfo.isActive')}
            onChange={(value) => setValue('productInfo.isActive', value)}
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
          <Label htmlFor='isActive' className='font-semibold cursor-pointer'>
            Active
          </Label>
        </Field>
      </div>

      {/* Basic Information */}
      <div className='grid grid-cols-4 gap-4 mb-6'>
        <TextField
          label={'Product name'}
          {...register('productInfo.name')}
          error={formState.errors.productInfo?.name?.message}
          placeholder='Enter product name...'
          type='text'
          required
        />
        <TextField
          {...register('productInfo.sku')}
          label={'Sku'}
          placeholder='Enter sku...'
          type='text'
          error={formState.errors.productInfo?.sku?.message}
        />
        <TextField
          {...register('productInfo.price', {
            valueAsNumber: true,
          })}
          label={'Price'}
          placeholder='Enter price...'
          type='number'
          error={formState.errors.productInfo?.price?.message}
        />
        <TextField
          label={'Slug'}
          placeholder='Enter slug...'
          type='text'
          error={formState.errors.productInfo?.slug?.message}
          {...register('productInfo.slug')}
        />
        {attributesLoading ? (
          <div className='flex justify-center items-center'>
            <LoadingSpinner />
          </div>
        ) : attributes ? (
          <StyledMultipleComboBox<{
            id: string;
            name: string;
          }>
            label='Select an attribute'
            setSelected={(values) => {
              setValue(
                'productInfo.attributes',
                values.map((e) => e.id),
                {
                  shouldDirty: false,
                }
              );
            }}
            options={attributes}
            getDisplayValue={(option) => {
              return option?.name || '';
            }}
            selected={watch('productInfo.attributes', []).map((e) => {
              const attribute = attributes.find((a) => a.id === e)!;
              return {
                id: attribute.id,
                name: attribute.name,
              };
            })}
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
            name='productInfo.category'
            label='Category'
            error={formState.errors.productInfo?.category?.message}
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
            name='productInfo.brand'
            nullable
            control={control}
            error={formState.errors.productInfo?.brand?.message}
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
            name='productInfo.collection'
            nullable
            label='Collection'
            error={formState.errors.productInfo?.brand?.message}
            options={
              collections?.map((e) => ({
                id: e.id,
                name: e.name,
              })) ?? []
            }
          />
        )}
      </div>
      {/* Short Description */}
      <Field className='w-full mb-4'>
        <TextField
          label='Short Description'
          {...register('productInfo.shortDescription')}
          placeholder='Enter short product description...'
          type='text'
          error={formState.errors.productInfo?.shortDescription?.message}
          className='w-full'
        />
      </Field>
      {/* Description */}
      <Field className='w-full'>
        <Label className='font-semibold'>Description</Label>
        <TiptapController
          name='productInfo.description'
          control={control}
          error={formState.errors.productInfo?.description?.message}
        />
      </Field>

      <div className='mt-6'>
        <ProductImagesUploader productDetail={productDetail} />
      </div>
    </>
  );
};
