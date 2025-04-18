'use client';

import { zodResolver } from '@hookform/resolvers/zod';
import React, { useEffect } from 'react';
import { FormProvider, useFieldArray, useForm } from 'react-hook-form';
import { useCategories } from '../../_lib/hooks/useCategories';
import { useCollections } from '../../_lib/hooks/useCollections';
import { useBrands } from '../../_lib/hooks/useBrands';
import {
  Button,
  Field,
  Fieldset,
  Label,
  Legend,
  Switch,
} from '@headlessui/react';
import { TextField } from '@/components/FormFields';
import { LoadingSpinner } from '@/components/Common/Loadings/Loading';
import { StyledComboBoxController } from '@/components/FormFields/StyledComboBoxController';
import clsx from 'clsx';
import {
  ArrowLeftCircleIcon,
  ArrowTurnUpLeftIcon,
  PlusIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';
import Image from 'next/image';
import { useAttributes } from '../../_lib/hooks/useAttributes';
import { StyledMultipleComboBox } from '@/components/FormFields/StyledMultipleComboBox';
import {
  AttributeFormModel,
  ProductDetailModel,
  ProductFormSchema,
  ProductModelForm,
} from '@/lib/definitions';
import { VariantForm } from './VariantForm';
import {
  ProductImageFile,
  useProductDetailFormContext,
} from '../_lib/contexts/ProductFormContext';
import { ImageUploader } from '@/components/FormFields/ImageUploader';
import { TiptapController } from '@/components/Common';
import { XMarkIcon } from '@heroicons/react/16/solid';

interface ProductEditFormProps {
  onSubmit: (
    data: ProductModelForm,
    productImages: ProductImageFile[],
    variantImages: ProductImageFile[]
  ) => Promise<void>;
  productDetail?: ProductDetailModel;
}

export const ProductDetailForm: React.FC<ProductEditFormProps> = ({
  onSubmit,
  productDetail,
}) => {
  const {
    productImages,
    productVariantImages,
    setProductImages,
    setProductVariantImages,
  } = useProductDetailFormContext();

  const productForm = useForm<ProductModelForm>({
    resolver: zodResolver(ProductFormSchema),
    mode: 'onBlur',
    reValidateMode: 'onChange',
    defaultValues: productDetail
      ? {
          id: productDetail.id,
          variants: productDetail.variants.map((variant) => ({
            attributes: variant.attributes.map((attribute) => ({
              id: attribute.id,
              name: attribute.name,
              values: attribute.values,
            })),
            id: variant.id,
            price: variant.price,
            stock: variant.stock_qty,
            sku: variant.sku,
            weight: variant.weight,
            is_active: variant.is_active,
            image_url: variant.image_url,
          })),
          category: {
            id: productDetail.category.id,
            name: productDetail.category.name,
          },
          brand: productDetail.brand
            ? { id: productDetail.brand.id, name: productDetail.brand.name }
            : null,
          collection: productDetail.collection
            ? {
                id: productDetail.collection.id,
                name: productDetail.collection.name,
              }
            : null,
          images: productDetail?.images
            ? productDetail.images.map((image) => ({
                id: image.id,
                image_url: image.url,
              }))
            : [],
          description: productDetail.description,
          name: productDetail.name,
          price: productDetail.price,
          sku: productDetail.sku,
          is_active: productDetail.is_active,
          slug: productDetail.slug,
        }
      : {
          variants: [],
          brand: null,
          category: {
            id: '',
            name: '',
          },
          collection: null,
          description: '',
          name: '',
          sku: '',
          slug: '',
          price: 1,
          is_active: true,
        },
  });
  const { register, control, getValues, watch, formState, setValue } =
    productForm;

  const { fields, append, remove, update } = useFieldArray({
    name: 'variants',
    keyName: 'key',
    control,
  });

  const { categories, isLoading: categoriesLoading } = useCategories();
  const { collections, isLoading: collectionLoading } = useCollections();
  const { brands, isLoading: brandsLoading } = useBrands();
  const { attributes, attributesLoading } = useAttributes();

  async function submitHandler(data: ProductModelForm) {
    await onSubmit(data, productImages, productVariantImages);
    setProductImages([]);
    setProductVariantImages([]);
  }

  useEffect(() => {
    if (!productDetail && categories && categories.length > 0) {
      setValue('category', categories[0]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [categories]);
  return (
    <div className='h-full px-6 py-3 overflow-auto'>
      <FormProvider {...productForm}>
        <Fieldset
          onSubmit={productForm.handleSubmit(submitHandler, console.error)}
          as='form'
        >
          <Link
            href={'/admin/products'}
            className='flex items-center mb-2 space-x-2'
          >
            <ArrowLeftCircleIcon className='size-6 text-primary' />
            <span className='text-primary text-lg hover:underline'>
              Back to Products
            </span>
          </Link>
          <Legend className='text-2xl flex justify-between font-bold text-primary mb-4'>
            {productDetail ? (
              <span>Edit Product: {productDetail.name}</span>
            ) : (
              <span>Create New Product</span>
            )}
            <Button
              disabled={
                formState.isSubmitting ||
                (!formState.isDirty &&
                  !productImages.length &&
                  !productVariantImages.length)
              }
              type='submit'
              className={clsx(
                'btn text-lg btn-primary',
                formState.isSubmitting && 'loading',
                formState.isDirty ? 'btn-primary' : 'btn-disabled'
              )}
            >
              {formState.isSubmitting ? (
                <span>{productDetail ? 'Saving...' : 'Creating...'}</span>
              ) : (
                <span>{productDetail ? 'Save' : 'Create'}</span>
              )}
            </Button>
          </Legend>
          <div className='flex gap-6 mt-4'>
            <div className='w-1/2 flex flex-col gap-4 p-4 bg-white rounded-lg shadow-md'>
              <div className='flex gap-4'>
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
              </div>
              <div className='flex gap-4'>
                <TextField
                  {...register('price', {
                    valueAsNumber: true,
                  })}
                  label={'Price'}
                  placeholder='Enter price...'
                  type='number'
                  error={!!formState.errors.price}
                  message={formState.errors.price?.message ?? 'Invalid price'}
                />
                {attributesLoading ? (
                  <div className='w-full flex justify-center'>
                    <LoadingSpinner />
                  </div>
                ) : attributes ? (
                  <StyledMultipleComboBox<AttributeFormModel>
                    label='Select an attribute'
                    setSelected={(value) => {
                      fields.forEach((field, index) => {
                        update(index, {
                          ...field,
                          attributes: value.map((e) => ({
                            id: e.id,
                            name: e.name,
                            values: [],
                          })),
                        });
                      });
                    }}
                    options={attributes}
                    selected={getValues('variants.0.attributes') ?? []}
                  />
                ) : null}
              </div>

              <div className='flex gap-4'>
                {categoriesLoading ? (
                  <div className='w-full my-auto flex justify-center'>
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
                {collectionLoading ? (
                  <div className='w-full my-auto flex justify-center'>
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
              <div className='flex gap-4'>
                {brandsLoading ? (
                  <div className='w-full my-auto flex justify-center'>
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
                <TextField
                  label={'Slug'}
                  placeholder='Enter slug...'
                  type='text'
                  error={!!formState.errors.slug}
                  message={formState.errors.slug?.message}
                  {...register('slug')}
                />
              </div>

              <div className='flex items-center gap-2'>
                <Field className='flex items-center gap-2'>
                  <Switch
                    checked={getValues('is_active')}
                    onChange={(value) =>
                      productForm.setValue('is_active', value)
                    }
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
                  <Label
                    htmlFor='is_active'
                    className='font-semibold cursor-pointer'
                  >
                    Active Product
                  </Label>
                </Field>
              </div>

              <Field className={clsx('w-full')}>
                <Label className='font-semibold'>Description</Label>
                <TiptapController
                  name='description'
                  control={control}
                  error={!!formState.errors.description}
                  message={formState.errors.description?.message as string}
                />
              </Field>

              <Field className={clsx('w-full')}>
                <Label className='font-semibold'>Media</Label>
                <div className='mt-2'>
                  <ImageUploader
                    label='Product Images'
                    multiple={true}
                    onUpload={(files) => {
                      // Handle base product images
                      const baseProductImages = files.map((file) => ({
                        image: file,
                        variantID: null, // null indicates it's a base product image, not a variant
                      }));

                      setProductImages((prev) => [
                        ...prev.filter((f) => f.variantID !== null), // Keep variant images
                        ...baseProductImages, // Add base product images
                      ]);

                      // Mark form as dirty to enable save button
                      productForm.setValue(
                        'images',
                        files.map((_, index) => ({
                          id: index,
                          image_url: '',
                        })),
                        {
                          shouldDirty: true,
                        }
                      );
                    }}
                  />
                </div>

                {/* Preview uploaded images */}
                <div className='grid grid-cols-4 gap-4 mt-4'>
                  {productDetail?.images.map((image, index) => {
                    const isRemoved = watch('removed_images')?.includes(
                      image.id
                    );
                    return (
                      <div
                        key={index}
                        className='relative rounded-md h-72 w-full'
                        >
                        <Image
                          fill
                          src={image.url}
                          objectFit='cover'
                          alt={`Product image ${index + 1}`}
                          className={`${isRemoved ? 'opacity-30' : ''}`}
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
                  {productImages
                    .filter(
                      (file) => file.variantID === null && file.image !== null
                    )
                    .map((file, index) => (
                      <div
                        key={index}
                        className='relative rounded-md h-72'
                      >
                        <Image
                          fill
                          src={
                            file.image ? URL.createObjectURL(file.image) : ''
                          }
                          objectFit='cover'
                          alt={`Product image ${index + 1}`}
                          className='w-full h-full object-contain'
                        />
                        <button
                          type='button'
                          onClick={() => {
                            setProductImages((prev) =>
                              prev.filter(
                                (_, i) =>
                                  !(prev[i].variantID === null && i === index)
                              )
                            );
                          }}
                          className='absolute top-1 right-1 bg-red-500 text-white rounded-full p-1 w-6 h-6 flex items-center justify-center'
                        >
                          <XMarkIcon className='size-8' />
                        </button>
                      </div>
                    ))}
                </div>
              </Field>
            </div>
            <div className='w-1/2 p-4 bg-white rounded-lg shadow-md'>
              {fields.map((item, index) => (
                <VariantForm
                  key={item.key}
                  index={index}
                  onRemove={() => {
                    remove(index);
                    setProductVariantImages((prev) =>
                      prev.filter((_, idx) => idx !== index)
                    );
                  }}
                />
              ))}
              <div className='flex justify-end'>
                <Button
                  onClick={() => {
                    append({
                      attributes: attributes
                        ? [
                            ...attributes.map((e) => ({
                              id: e.id,
                              name: e.name,
                              values: [],
                            })),
                          ]
                        : [],
                      price: 1,
                      stock: 1,
                      sku: '',
                      weight: undefined,
                      image_url: '',
                      is_active: true,
                    });
                    setProductVariantImages((prev) => [
                      ...prev,
                      {
                        image: null,
                        variantID: -1,
                      },
                    ]);
                  }}
                  className={clsx('btn btn-primary flex gap-2')}
                >
                  <PlusIcon className='size-6' />
                  Add Variant
                </Button>
              </div>
            </div>
          </div>
        </Fieldset>
      </FormProvider>
    </div>
  );
};
