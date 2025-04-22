'use client';
import React, { useEffect } from 'react';
import { TextField } from '@/components/FormFields';
import { LoadingSpinner } from '@/components/Common/Loadings/Loading';
import { StyledComboBoxController } from '@/components/FormFields/StyledComboBoxController';
import { StyledMultipleComboBox } from '@/components/FormFields/StyledMultipleComboBox';
import { TiptapController } from '@/components/Common';
import { Field, Label, Switch } from '@headlessui/react';
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
import { ProductImagesUploader } from './ProductImagesUploader';
import { useProductDetailFormContext } from '../_lib/contexts/ProductFormContext';

export const ProductInfoForm: React.FC<{
  productDetail?: ProductDetailModel;
}> = ({ productDetail }) => {
  const { categories, isLoading: categoriesLoading } = useCategories();
  const { collections, isLoading: collectionLoading } = useCollections();
  const { brands, isLoading: brandsLoading } = useBrands();
  const { attributes, attributesLoading } = useAttributes();

  const { selectedAttributes, setSelectedAttributes } =
    useProductDetailFormContext();

  const { register, control, watch, formState, setValue } =
    useFormContext<ProductModelForm>();

  useEffect(() => {
    if (!productDetail && categories && categories.length > 0) {
      setValue('product_info.category', categories[0], {
        shouldDirty: false,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categories]);

  useEffect(() => {
    if (!productDetail && brands && brands.length > 0) {
      setValue('product_info.brand', brands[0], {
        shouldDirty: false,
      });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [brands]);

  useEffect(() => {
    if (productDetail && attributes) {
      const attrIds = new Set(
        productDetail.variants.map((v) => v.attributes.map((a) => a.id)).flat()
      );
      const selected = attributes.filter((a) => attrIds.has(a.id));
      setSelectedAttributes(selected);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [attributes, productDetail]);

  return (
    <>
      <div className='flex gap-4 items-center mb-4'>
        <h2 className='text-xl font-bold text-primary'>Product Information</h2>
        <Field className='flex items-center gap-2'>
          <Switch
            checked={watch('product_info.is_active')}
            onChange={(value) => setValue('product_info.is_active', value)}
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
          {...register('product_info.name')}
          error={!!formState.errors.product_info?.name}
          message={formState.errors.product_info?.name?.message}
          placeholder='Enter product name...'
          type='text'
          required
        />
        <TextField
          {...register('product_info.sku')}
          label={'Sku'}
          placeholder='Enter sku...'
          type='text'
          error={!!formState.errors.product_info?.sku}
          message={formState.errors.product_info?.sku?.message}
        />
        <TextField
          {...register('product_info.price', {
            valueAsNumber: true,
          })}
          label={'Price'}
          placeholder='Enter price...'
          type='number'
          error={!!formState.errors.product_info?.price}
          message={formState.errors.product_info?.price?.message}
        />
        <TextField
          label={'Slug'}
          placeholder='Enter slug...'
          type='text'
          error={!!formState.errors.product_info?.slug}
          message={formState.errors.product_info?.slug?.message}
          {...register('product_info.slug')}
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
            name='product_info.category'
            label='Category'
            error={!!formState.errors.product_info?.category}
            message={formState.errors.product_info?.category?.message ?? ''}
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
            name='product_info.brand'
            nullable
            control={control}
            error={!!formState.errors.product_info?.brand}
            message={formState.errors.product_info?.brand?.message ?? ''}
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
            name='product_info.collection'
            nullable
            label='Collection'
            error={!!formState.errors.product_info?.brand}
            message={formState.errors.product_info?.brand?.message ?? ''}
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
          name='product_info.description'
          control={control}
          error={!!formState.errors.product_info?.description}
          message={
            formState.errors.product_info?.description?.message as string
          }
        />
      </Field>

      <div className='mt-6'>
        <ProductImagesUploader productDetail={productDetail} />
      </div>
    </>
  );
};
