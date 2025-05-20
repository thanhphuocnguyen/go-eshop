'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useForm, Controller } from 'react-hook-form';
import { 
  ArrowLeftIcon,
  TagIcon,
  ShoppingBagIcon,
  ClipboardDocumentListIcon,
  XMarkIcon,
  PlusIcon,
} from '@heroicons/react/24/outline';
import dayjs from 'dayjs';

// Mock data for products and categories selection
const mockAllProducts = [
  { id: 'prod1', name: 'Summer T-Shirt', price: 29.99 },
  { id: 'prod2', name: 'Beach Shorts', price: 39.99 },
  { id: 'prod3', name: 'Sunglasses', price: 25.0 },
  { id: 'prod4', name: 'Beach Hat', price: 19.99 },
  { id: 'prod5', name: 'Flip Flops', price: 15.99 },
  { id: 'prod6', name: 'Beach Towel', price: 22.99 },
];

const mockAllCategories = [
  { id: 'cat1', name: 'Summer Collection' },
  { id: 'cat2', name: 'Beachwear' },
  { id: 'cat3', name: 'Accessories' },
  { id: 'cat4', name: 'Footwear' },
  { id: 'cat5', name: 'Seasonal' },
];

type DiscountFormData = {
  code: string;
  description: string;
  discountType: string;
  discountValue: number;
  minPurchaseAmount: number | null;
  maxDiscountAmount: number | null;
  usageLimit: number | null;
  isActive: boolean;
  startsAt: string;
  expiresAt: string;
  products: { id: string; name: string; price: number }[];
  categories: { id: string; name: string }[];
};

export default function NewDiscountPage() {
  const router = useRouter();
  const [submitting, setSubmitting] = useState(false);
  const [activeTab, setActiveTab] = useState('details');
  const [error, setError] = useState<string | null>(null);
  
  const [availableProducts, setAvailableProducts] = useState<any[]>([]);
  const [availableCategories, setAvailableCategories] = useState<any[]>([]);
  
  // Initialize react-hook-form
  const { 
    register, 
    handleSubmit, 
    control, 
    setValue, 
    watch, 
    formState: { errors } 
  } = useForm<DiscountFormData>({
    defaultValues: {
      code: '',
      description: '',
      discountType: 'percentage',
      discountValue: 0,
      minPurchaseAmount: null,
      maxDiscountAmount: null,
      usageLimit: null,
      isActive: true,
      startsAt: dayjs().format('YYYY-MM-DDTHH:mm'),
      expiresAt: dayjs().add(30, 'day').format('YYYY-MM-DDTHH:mm'),
      products: [],
      categories: []
    }
  });

  // Watch for form values
  const selectedProducts = watch('products') || [];
  const selectedCategories = watch('categories') || [];

  useEffect(() => {
    // In a real implementation, this would fetch available products/categories from an API
    const fetchData = async () => {
      try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 500));
        
        setAvailableProducts(mockAllProducts);
        setAvailableCategories(mockAllCategories);
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };
    
    fetchData();
  }, []);

  // Handler for adding a product to the discount
  const handleAddProduct = (product: any) => {
    const currentProducts = [...selectedProducts];
    currentProducts.push(product);
    setValue('products', currentProducts);
    
    // Remove from available products
    setAvailableProducts(availableProducts.filter(p => p.id !== product.id));
  };

  // Handler for removing a product from the discount
  const handleRemoveProduct = (productId: string) => {
    const removedProduct = selectedProducts.find(p => p.id === productId);
    if (removedProduct) {
      setValue('products', selectedProducts.filter(p => p.id !== productId));
      setAvailableProducts([...availableProducts, removedProduct]);
    }
  };

  // Handler for adding a category to the discount
  const handleAddCategory = (category: any) => {
    const currentCategories = [...selectedCategories];
    currentCategories.push(category);
    setValue('categories', currentCategories);
    
    // Remove from available categories
    setAvailableCategories(availableCategories.filter(c => c.id !== category.id));
  };

  // Handler for removing a category from the discount
  const handleRemoveCategory = (categoryId: string) => {
    const removedCategory = selectedCategories.find(c => c.id === categoryId);
    if (removedCategory) {
      setValue('categories', selectedCategories.filter(c => c.id !== categoryId));
      setAvailableCategories([...availableCategories, removedCategory]);
    }
  };

  // Submit handler
  const onSubmit = async (data: DiscountFormData) => {
    setSubmitting(true);
    setError(null);
    
    try {
      // In a real implementation, this would send the data to an API
      console.log('Submitting discount data:', data);
      
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Redirect to the discounts list page
      router.push('/admin/discounts');
    } catch (err) {
      setError('Failed to create discount. Please try again.');
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="p-4">
      <div className="mb-4 flex items-center">
        <button 
          onClick={() => router.back()}
          className="mr-4 p-2 rounded-full hover:bg-gray-100"
        >
          <ArrowLeftIcon className="h-5 w-5" />
        </button>
        <h1 className="text-2xl font-semibold">Create New Discount</h1>
      </div>
      
      {error && (
        <div className="mb-4 p-3 bg-red-100 text-red-800 rounded-md">
          {error}
        </div>
      )}
      
      <div className="bg-white p-6 rounded-lg shadow">
        <form onSubmit={handleSubmit(onSubmit)}>
          <div className="border-b mb-4 pb-2">
            <nav className="flex -mb-px">
              <button
                type="button"
                onClick={() => setActiveTab('details')}
                className={`px-4 py-2 text-sm font-medium ${
                  activeTab === 'details'
                    ? 'border-b-2 border-primary text-primary'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                <span className="flex items-center">
                  <TagIcon className="h-4 w-4 mr-2" />
                  Details
                </span>
              </button>
              <button
                type="button"
                onClick={() => setActiveTab('products')}
                className={`px-4 py-2 text-sm font-medium ${
                  activeTab === 'products'
                    ? 'border-b-2 border-primary text-primary'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                <span className="flex items-center">
                  <ShoppingBagIcon className="h-4 w-4 mr-2" />
                  Products
                </span>
              </button>
              <button
                type="button"
                onClick={() => setActiveTab('categories')}
                className={`px-4 py-2 text-sm font-medium ${
                  activeTab === 'categories'
                    ? 'border-b-2 border-primary text-primary'
                    : 'text-gray-500 hover:text-gray-700'
                }`}
              >
                <span className="flex items-center">
                  <ClipboardDocumentListIcon className="h-4 w-4 mr-2" />
                  Categories
                </span>
              </button>
            </nav>
          </div>

          {activeTab === 'details' && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="col-span-1">
                <label className="block mb-2 text-sm font-medium">
                  Discount Code <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <TagIcon className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    type="text"
                    {...register("code", { required: "Discount code is required" })}
                    className={`pl-10 block w-full border ${errors.code ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
                    placeholder="e.g. SUMMER25"
                  />
                </div>
                {errors.code && (
                  <p className="mt-1 text-xs text-red-500">{errors.code.message}</p>
                )}
                <p className="mt-1 text-xs text-gray-500">
                  Code used by customers to apply the discount
                </p>
              </div>
              
              <div className="col-span-1">
                <label className="block mb-2 text-sm font-medium">
                  Discount Type <span className="text-red-500">*</span>
                </label>
                <select
                  {...register("discountType", { required: "Discount type is required" })}
                  className={`block w-full border ${errors.discountType ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
                >
                  <option value="percentage">Percentage (%)</option>
                  <option value="fixed_amount">Fixed Amount ($)</option>
                </select>
                {errors.discountType && (
                  <p className="mt-1 text-xs text-red-500">{errors.discountType.message}</p>
                )}
                <p className="mt-1 text-xs text-gray-500">
                  How the discount will be calculated
                </p>
              </div>
              
              <div className="col-span-1">
                <label className="block mb-2 text-sm font-medium">
                  Discount Value <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    {watch('discountType') === 'percentage' ? '%' : '$'}
                  </div>
                  <input
                    type="number"
                    step="0.01"
                    {...register("discountValue", { 
                      required: "Discount value is required",
                      min: { value: 0, message: "Value must be positive" },
                      valueAsNumber: true
                    })}
                  className={`pl-8 block w-full border ${errors.discountValue ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
                  placeholder={watch('discountType') === 'percentage' ? '10' : '10.00'}
                />
              </div>
              <p className="mt-1 text-xs text-gray-500">
                {watch('discountType') === 'percentage' 
                  ? 'Percentage off (0-100)' 
                  : 'Fixed amount to deduct from order'}
              </p>
            </div>
            
            <div className="col-span-1">
              <label className="block mb-2 text-sm font-medium">
                Minimum Purchase Amount
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <span className="text-gray-500">$</span>
                </div>
                <input
                  type="number"
                  step="0.01"
                  {...register("minPurchaseAmount", { 
                    min: { value: 0, message: "Amount must be positive" },
                    valueAsNumber: true
                  })}
                  className={`pl-8 block w-full border ${errors.minPurchaseAmount ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
                  placeholder="0.00"
                />
              </div>
              {errors.minPurchaseAmount && (
                <p className="mt-1 text-xs text-red-500">{errors.minPurchaseAmount.message}</p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                Minimum order amount required (optional)
              </p>
            </div>
            
            <div className="col-span-1">
              <label className="block mb-2 text-sm font-medium">
                Maximum Discount Amount
              </label>
              <div className="relative">
                <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <span className="text-gray-500">$</span>
                </div>
                <input
                  type="number"
                  step="0.01"
                  {...register("maxDiscountAmount", { 
                    min: { value: 0, message: "Amount must be positive" },
                    valueAsNumber: true
                  })}
                  className={`pl-8 block w-full border ${errors.maxDiscountAmount ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
                  placeholder="0.00"
                />
              </div>
              {errors.maxDiscountAmount && (
                <p className="mt-1 text-xs text-red-500">{errors.maxDiscountAmount.message}</p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                Maximum amount to discount (optional)
              </p>
            </div>
            
            <div className="col-span-1">
              <label className="block mb-2 text-sm font-medium">
                Usage Limit
              </label>
              <input
                type="number"
                {...register("usageLimit", { 
                  min: { value: 1, message: "Usage limit must be at least 1" },
                  valueAsNumber: true
                })}
                className={`block w-full border ${errors.usageLimit ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
                placeholder="No limit"
              />
              {errors.usageLimit && (
                <p className="mt-1 text-xs text-red-500">{errors.usageLimit.message}</p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                Maximum number of times this discount can be used (optional)
              </p>
            </div>
            
            <div className="col-span-1">
              <label className="block mb-2 text-sm font-medium">
                Status
              </label>
              <div className="flex items-center">
                <Controller
                  name="isActive"
                  control={control}
                  render={({ field }) => (
                    <input
                      type="checkbox"
                      id="isActive"
                      className="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
                      checked={field.value}
                      onChange={field.onChange}
                    />
                  )}
                />
                <label htmlFor="isActive" className="ml-2 text-sm text-gray-700">
                  Active (can be used by customers)
                </label>
              </div>
            </div>
            
            <div className="col-span-1">
              <label className="block mb-2 text-sm font-medium">
                Start Date <span className="text-red-500">*</span>
              </label>
              <input
                type="datetime-local"
                {...register("startsAt", { required: "Start date is required" })}
                className={`block w-full border ${errors.startsAt ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
              />
              {errors.startsAt && (
                <p className="mt-1 text-xs text-red-500">{errors.startsAt.message}</p>
              )}
            </div>
            
            <div className="col-span-1">
              <label className="block mb-2 text-sm font-medium">
                Expiry Date <span className="text-red-500">*</span>
              </label>
              <input
                type="datetime-local"
                {...register("expiresAt", { required: "Expiry date is required" })}
                className={`block w-full border ${errors.expiresAt ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
              />
              {errors.expiresAt && (
                <p className="mt-1 text-xs text-red-500">{errors.expiresAt.message}</p>
              )}
            </div>
            
            <div className="col-span-2">
              <label className="block mb-2 text-sm font-medium">
                Description
              </label>
              <textarea
                {...register("description")}
                className={`block w-full border ${errors.description ? 'border-red-500' : 'border-gray-300'} rounded-md px-3 py-2`}
                placeholder="Description of this discount"
                rows={3}
              ></textarea>
              {errors.description && (
                <p className="mt-1 text-xs text-red-500">{errors.description.message}</p>
              )}
            </div>
          </div>
          )}

          {/* Products Tab */}
          {activeTab === 'products' && (
            <div>
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-medium">Applied Products</h3>
                <p className="text-sm text-gray-500">
                  If no products are selected, this discount applies to all products.
                </p>
              </div>
              
              {/* Selected Products */}
              <h4 className="font-medium text-sm text-gray-700 mb-2">Selected Products</h4>
              {selectedProducts.length > 0 ? (
                <div className="mb-6 grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                  {selectedProducts.map((product) => (
                    <div
                      key={product.id}
                      className="relative flex items-center border rounded-md p-3"
                    >
                      <div className="flex-1">
                        <p className="font-medium">{product.name}</p>
                        <p className="text-gray-500 text-sm">${product.price.toFixed(2)}</p>
                      </div>
                      <button
                        type="button"
                        onClick={() => handleRemoveProduct(product.id)}
                        className="text-gray-400 hover:text-red-500"
                      >
                        <XMarkIcon className="h-5 w-5" />
                      </button>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="mb-6 bg-gray-50 rounded-md p-6 text-center">
                  <p className="text-gray-500">No products selected. This discount applies to all products.</p>
                </div>
              )}
              
              {/* Available Products */}
              <h4 className="font-medium text-sm text-gray-700 mb-2">Add Products</h4>
              {availableProducts.length > 0 ? (
                <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                  {availableProducts.map((product) => (
                    <div
                      key={product.id}
                      className="flex items-center border rounded-md p-3 cursor-pointer hover:border-primary"
                      onClick={() => handleAddProduct(product)}
                    >
                      <div className="flex-1">
                        <p className="font-medium">{product.name}</p>
                        <p className="text-gray-500 text-sm">${product.price.toFixed(2)}</p>
                      </div>
                      <button
                        type="button"
                        className="text-gray-400 hover:text-green-500"
                      >
                        <PlusIcon className="h-5 w-5" />
                      </button>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="bg-gray-50 rounded-md p-6 text-center">
                  <p className="text-gray-500">All available products have been added to this discount.</p>
                </div>
              )}
            </div>
          )}
          
          {/* Categories Tab */}
          {activeTab === 'categories' && (
            <div>
              <div className="flex justify-between items-center mb-4">
                <h3 className="text-lg font-medium">Applied Categories</h3>
                <p className="text-sm text-gray-500">
                  If no categories are selected, this discount applies to all products.
                </p>
              </div>
              
              {/* Selected Categories */}
              <h4 className="font-medium text-sm text-gray-700 mb-2">Selected Categories</h4>
              {selectedCategories.length > 0 ? (
                <div className="mb-6 grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                  {selectedCategories.map((category) => (
                    <div
                      key={category.id}
                      className="relative flex items-center border rounded-md p-3"
                    >
                      <div className="flex-1">
                        <p className="font-medium">{category.name}</p>
                      </div>
                      <button
                        type="button"
                        onClick={() => handleRemoveCategory(category.id)}
                        className="text-gray-400 hover:text-red-500"
                      >
                        <XMarkIcon className="h-5 w-5" />
                      </button>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="mb-6 bg-gray-50 rounded-md p-6 text-center">
                  <p className="text-gray-500">No categories selected. This discount applies to all products.</p>
                </div>
              )}
              
              {/* Available Categories */}
              <h4 className="font-medium text-sm text-gray-700 mb-2">Add Categories</h4>
              {availableCategories.length > 0 ? (
                <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
                  {availableCategories.map((category) => (
                    <div
                      key={category.id}
                      className="flex items-center border rounded-md p-3 cursor-pointer hover:border-primary"
                      onClick={() => handleAddCategory(category)}
                    >
                      <div className="flex-1">
                        <p className="font-medium">{category.name}</p>
                      </div>
                      <button
                        type="button"
                        className="text-gray-400 hover:text-green-500"
                      >
                        <PlusIcon className="h-5 w-5" />
                      </button>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="bg-gray-50 rounded-md p-6 text-center">
                  <p className="text-gray-500">All available categories have been added to this discount.</p>
                </div>
              )}
            </div>
          )}

          <div className="mt-8 flex justify-end">
            <Link href="/admin/discounts">
              <button type="button" className="mr-3 px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50">
                Cancel
              </button>
            </Link>
            <button
              type="submit"
              disabled={submitting}
              className={`px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/80 ${
                submitting ? 'opacity-75 cursor-not-allowed' : ''
              }`}
            >
              {submitting ? 'Creating...' : 'Create Discount'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
