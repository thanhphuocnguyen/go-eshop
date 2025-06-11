'use client';
import { ADMIN_API_PATHS } from '@/app/lib/constants/api';
import {
  Button,
  Dialog,
  DialogPanel,
  DialogTitle,
  Field,
  Input,
} from '@headlessui/react';
import clsx from 'clsx';
import React, { useState, useRef, useEffect } from 'react';
import { toast } from 'react-toastify';
import {
  AttributeFormModel,
  AttributeDetailModel,
} from '@/app/lib/definitions';
import {
  XCircleIcon,
  PlusIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
} from '@heroicons/react/24/outline';
import { clientSideFetch } from '@/app/lib/api/apiClient';
import { useForm, Controller, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

interface AddNewDialogProps {
  open: boolean;
  onClose: () => void;
  handleSubmitted: (attribute: AttributeDetailModel) => void;
}

// Form schema for React Hook Form
const createAttributeSchema = z.object({
  name: z
    .string()
    .min(1, 'Attribute name is required')
    .max(100, 'Name is too long'),
  values: z
    .array(
      z.object({
        name: z.string().min(1, 'Value name is required'),
        displayOrder: z.number().optional(),
      })
    )
    .min(1, 'At least one value is required'),
});

type CreateAttributeForm = z.infer<typeof createAttributeSchema>;

export const AddNewDialog: React.FC<AddNewDialogProps> = ({
  onClose,
  open,
  handleSubmitted,
}) => {
  const [isLoading, setIsLoading] = useState(false);
  const [currentValue, setCurrentValue] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);
  const nameInputRef = useRef<HTMLInputElement>(null);

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
    watch,
  } = useForm<CreateAttributeForm>({
    resolver: zodResolver(createAttributeSchema),
    defaultValues: {
      name: '',
      values: [],
    },
    mode: 'onChange',
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: 'values',
  });

  const watchedValues = watch('values'); // Focus on name input when dialog opens
  useEffect(() => {
    if (open && nameInputRef.current) {
      setTimeout(() => nameInputRef.current?.focus(), 100);
    }
  }, [open]);

  // Reset form when dialog closes
  useEffect(() => {
    if (!open) {
      reset();
      setCurrentValue('');
    }
  }, [open, reset]);

  const addValue = () => {
    const trimmedValue = currentValue.trim();
    if (!trimmedValue) return;

    // Check for duplicates
    const currentValues = watchedValues || [];
    if (
      currentValues.some(
        (v: any) => v.name?.toLowerCase() === trimmedValue.toLowerCase()
      )
    ) {
      toast.error('This value already exists');
      return;
    }

    append({
      name: trimmedValue,
      displayOrder: currentValues.length + 1,
    });
    setCurrentValue('');
    inputRef.current?.focus();
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      addValue();
    } else if (e.key === 'Escape') {
      setCurrentValue('');
    }
  };

  const onSubmit = async (data: CreateAttributeForm) => {
    setIsLoading(true);

    const body: AttributeFormModel = {
      name: data.name,
      values: data.values.map((value, idx) => ({
        name: value.name,
        displayOrder: idx + 1,
      })),
    };

    try {
      const resp = await clientSideFetch<AttributeDetailModel>(
        ADMIN_API_PATHS.ATTRIBUTES,
        {
          method: 'POST',
          body: body,
        }
      );

      if (resp.error) {
        toast.error('Failed to create attribute');
        setIsLoading(false);
        return;
      }

      toast.success(
        <div className='flex items-center gap-2'>
          <CheckCircleIcon className='h-5 w-5 text-green-500' />
          Attribute created successfully
        </div>
      );
      handleSubmitted(resp.data);
      onClose();
    } catch (error) {
      toast.error('An unexpected error occurred');
      console.error('Error creating attribute:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleClose = () => {
    if (isLoading) return; // Prevent closing while submitting
    onClose();
  };

  return (
    <Dialog
      open={open}
      as='div'
      className='relative z-50 focus:outline-none'
      onClose={handleClose}
    >
      {/* Backdrop with gradient overlay */}
      <div className='fixed inset-0 z-10 w-screen overflow-y-auto'>
        <div className='absolute inset-0 bg-gradient-to-br from-slate-900/80 via-gray-900/80 to-slate-800/80 backdrop-blur-sm' />
        <div className='flex min-h-full items-center justify-center p-4'>
          <DialogPanel
            transition
            className='relative w-full max-w-2xl bg-white rounded-3xl shadow-2xl duration-300 ease-out data-[closed]:transform-[scale(95%)] data-[closed]:opacity-0 overflow-hidden'
          >
            {/* Decorative gradient header */}
            <div className='absolute top-0 left-0 right-0 h-1 bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500' />

            {/* Header */}
            <div className='relative bg-gradient-to-r from-blue-50 to-indigo-50 px-8 py-6 border-b border-blue-100'>
              <div className='flex items-center justify-between'>
                <div className='flex items-center gap-4'>
                  <div className='relative'>
                    <div className='absolute inset-0 bg-blue-500 rounded-2xl blur-lg opacity-30 animate-pulse' />
                    <div className='relative p-3 bg-gradient-to-br from-blue-500 to-blue-600 rounded-2xl shadow-lg'>
                      <PlusIcon className='h-7 w-7 text-white' />
                    </div>
                  </div>
                  <div>
                    <DialogTitle className='text-2xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent'>
                      Create Attribute
                    </DialogTitle>
                    <p className='text-blue-600 font-medium mt-1 text-sm'>
                      ✨ Add a new product attribute with custom values
                    </p>
                  </div>
                </div>
                <Button
                  onClick={handleClose}
                  disabled={isLoading}
                  className='group p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 rounded-xl transition-all duration-200 hover:scale-110'
                >
                  <XCircleIcon className='h-6 w-6 transition-transform group-hover:rotate-90' />
                </Button>
              </div>
            </div>

            {/* Form Content */}
            <div className='px-8 py-8 bg-gradient-to-b from-white to-gray-50/50'>
              <form onSubmit={handleSubmit(onSubmit)} className='space-y-8'>
                {/* Attribute Name */}
                <div className='group'>
                  <div className='flex items-center gap-2 text-sm font-bold text-gray-800 mb-3'>
                    <span className='flex items-center justify-center w-6 h-6 bg-gradient-to-r from-blue-500 to-purple-500 text-white text-xs font-bold rounded-full'>
                      1
                    </span>
                    Attribute Name
                    <span className='text-red-500'>*</span>
                  </div>
                  <div className='relative'>
                    <Controller
                      name='name'
                      control={control}
                      render={({ field }) => (
                        <Input
                          {...field}
                          ref={nameInputRef}
                          disabled={isLoading}
                          id='name'
                          placeholder='e.g., Color, Size, Material, Style...'
                          className={clsx(
                            'block w-full rounded-2xl border-2 py-4 px-6 text-base font-medium transition-all duration-300 bg-white/70 backdrop-blur-sm',
                            errors.name
                              ? 'border-red-300 focus:border-red-500 focus:ring-red-500/30 shadow-red-100'
                              : 'border-gray-200 focus:border-blue-500 focus:ring-blue-500/30 shadow-blue-50 group-hover:border-blue-300',
                            'focus:outline-none focus:ring-4 focus:shadow-xl transform focus:scale-[1.02] placeholder:text-gray-400'
                          )}
                        />
                      )}
                    />
                    <div className='absolute inset-0 rounded-2xl bg-gradient-to-r from-blue-500/10 to-purple-500/10 pointer-events-none opacity-0 group-hover:opacity-100 transition-opacity duration-300' />
                  </div>
                  {errors.name && (
                    <div className='mt-3 flex items-center gap-2 text-red-600 animate-in slide-in-from-top-2 duration-300'>
                      <div className='p-1 bg-red-100 rounded-full'>
                        <ExclamationTriangleIcon className='h-4 w-4' />
                      </div>
                      <span className='text-sm font-medium'>
                        {errors.name.message}
                      </span>
                    </div>
                  )}
                </div>

                {/* Attribute Values */}
                <div className='space-y-4'>
                  <div className='flex items-center gap-2 text-sm font-bold text-gray-800'>
                    <span className='flex items-center justify-center w-6 h-6 bg-gradient-to-r from-purple-500 to-pink-500 text-white text-xs font-bold rounded-full'>
                      2
                    </span>
                    Attribute Values
                    <span className='text-red-500'>*</span>
                  </div>
                  <p className='text-sm text-gray-600 font-medium ml-8 flex items-center gap-2'>
                    <span className='w-2 h-2 bg-blue-400 rounded-full animate-pulse' />
                    Add values for this attribute (e.g., Red, Blue, Green for
                    Color)
                  </p>

                  {/* Values Container */}
                  <div
                    className={clsx(
                      'relative min-h-[140px] p-6 border-2 rounded-3xl transition-all duration-300 bg-gradient-to-br from-white via-gray-50/50 to-white',
                      errors.values
                        ? 'border-red-300 shadow-red-100 bg-red-50/30'
                        : 'border-gray-200 shadow-lg hover:shadow-xl hover:border-blue-300'
                    )}
                  >
                    {/* Decorative corner accents */}
                    <div className='absolute top-3 right-3 w-3 h-3 bg-gradient-to-br from-blue-400 to-purple-400 rounded-full opacity-30' />
                    <div className='absolute bottom-3 left-3 w-2 h-2 bg-gradient-to-br from-pink-400 to-purple-400 rounded-full opacity-30' />

                    {/* Existing Values */}
                    {fields.length > 0 && (
                      <div className='flex flex-wrap gap-3 mb-6'>
                        {fields.map((field, index) => (
                          <div
                            key={field.id}
                            className='group relative animate-in slide-in-from-bottom-2 duration-300'
                            style={{ animationDelay: `${index * 100}ms` }}
                          >
                            <div className='absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-500 rounded-2xl blur-lg opacity-20 group-hover:opacity-40 transition-opacity duration-300' />
                            <span className='relative inline-flex items-center gap-3 px-5 py-3 bg-gradient-to-r from-blue-500 to-blue-600 text-white font-semibold rounded-2xl shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200'>
                              <span className='flex items-center gap-2'>
                                <span className='w-2 h-2 bg-white/70 rounded-full' />
                                {field.name}
                              </span>
                              <button
                                type='button'
                                title='Remove value'
                                onClick={() => remove(index)}
                                disabled={isLoading}
                                className='p-1 hover:bg-red-500 rounded-full transition-all duration-200 disabled:opacity-50 hover:scale-110'
                              >
                                <XCircleIcon className='h-4 w-4' />
                              </button>
                            </span>
                          </div>
                        ))}
                      </div>
                    )}

                    {/* Add New Value Input */}
                    <div className='flex gap-3'>
                      <Field className='flex-1 group'>
                        <div className='relative'>
                          <Input
                            ref={inputRef}
                            disabled={isLoading}
                            value={currentValue}
                            onChange={(e) => setCurrentValue(e.target.value)}
                            onKeyDown={handleKeyDown}
                            placeholder='Type value and press Enter ✨'
                            className='w-full rounded-2xl border-2 border-gray-200 py-3.5 px-5 text-base font-medium bg-white/80 backdrop-blur-sm focus:outline-none focus:ring-4 focus:ring-blue-500/30 focus:border-blue-500 transition-all duration-300 placeholder:text-gray-400 hover:border-blue-300 transform focus:scale-[1.02]'
                          />
                          <div className='absolute right-4 top-1/2 -translate-y-1/2 text-gray-300 group-hover:text-blue-400 transition-colors duration-300'>
                            <span className='text-sm'>⏎</span>
                          </div>
                        </div>
                      </Field>
                      <Button
                        type='button'
                        onClick={addValue}
                        disabled={isLoading || !currentValue.trim()}
                        className={clsx(
                          'relative px-6 py-3.5 rounded-2xl font-bold text-base transition-all duration-300 transform shadow-lg',
                          currentValue.trim() && !isLoading
                            ? 'bg-gradient-to-r from-blue-600 to-purple-600 text-white hover:from-blue-700 hover:to-purple-700 focus:outline-none focus:ring-4 focus:ring-blue-500/50 hover:scale-105 shadow-blue-200'
                            : 'bg-gray-200 text-gray-400 cursor-not-allowed shadow-gray-100'
                        )}
                      >
                        <div className='flex items-center gap-2'>
                          <PlusIcon className='h-5 w-5' />
                          <span>Add</span>
                        </div>
                        {currentValue.trim() && !isLoading && (
                          <div className='absolute inset-0 rounded-2xl bg-gradient-to-r from-blue-600 to-purple-600 opacity-0 hover:opacity-20 transition-opacity duration-300' />
                        )}
                      </Button>
                    </div>

                    {/* Empty State */}
                    {fields.length === 0 && (
                      <div className='text-center py-8 animate-in fade-in-50 duration-500'>
                        <div className='relative mb-4'>
                          <div className='absolute inset-0 flex items-center justify-center'>
                            <div className='w-16 h-16 bg-gradient-to-r from-blue-200 to-purple-200 rounded-full animate-ping opacity-20' />
                          </div>
                          <div className='relative flex items-center justify-center w-16 h-16 bg-gradient-to-r from-blue-100 to-purple-100 rounded-full mx-auto'>
                            <PlusIcon className='h-8 w-8 text-blue-500' />
                          </div>
                        </div>
                        <p className='text-gray-500 font-medium'>
                          No values added yet
                        </p>
                        <p className='text-gray-400 text-sm mt-1'>
                          Start typing to add your first value
                        </p>
                      </div>
                    )}
                  </div>

                  {errors.values && (
                    <div className='mt-4 flex items-center gap-3 text-red-600 animate-in slide-in-from-top-2 duration-300'>
                      <div className='p-2 bg-red-100 rounded-xl'>
                        <ExclamationTriangleIcon className='h-5 w-5' />
                      </div>
                      <div>
                        <span className='text-sm font-bold'>
                          Validation Error
                        </span>
                        <p className='text-sm'>{errors.values.message}</p>
                      </div>
                    </div>
                  )}
                </div>

                {/* Form Actions */}
                <div className='flex items-center justify-between pt-6 border-t border-gradient-to-r from-transparent via-gray-200 to-transparent'>
                  <div className='text-sm text-gray-500 font-medium'>
                    {fields.length > 0 && (
                      <span className='flex items-center gap-2'>
                        <div className='w-2 h-2 bg-green-400 rounded-full animate-pulse' />
                        {fields.length} value
                        {fields.length !== 1 ? 's' : ''} added
                      </span>
                    )}
                  </div>
                  <div className='flex items-center gap-4'>
                    <Button
                      type='button'
                      onClick={handleClose}
                      disabled={isLoading}
                      className='px-8 py-3 text-base font-semibold text-gray-600 bg-white border-2 border-gray-200 rounded-2xl hover:bg-gray-50 hover:border-gray-300 focus:outline-none focus:ring-4 focus:ring-gray-200 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-105 shadow-lg hover:shadow-xl'
                    >
                      Cancel
                    </Button>
                    <Button
                      type='submit'
                      disabled={isLoading || fields.length === 0}
                      className={clsx(
                        'relative px-8 py-3 text-base font-bold rounded-2xl transition-all duration-300 transform shadow-xl',
                        isLoading || fields.length === 0
                          ? 'bg-gray-300 text-gray-500 cursor-not-allowed shadow-gray-200'
                          : 'bg-gradient-to-r from-blue-600 via-purple-600 to-pink-600 text-white hover:from-blue-700 hover:via-purple-700 hover:to-pink-700 focus:outline-none focus:ring-4 focus:ring-blue-500/50 hover:scale-105 shadow-blue-200'
                      )}
                    >
                      <div className='relative flex items-center gap-3'>
                        {isLoading && (
                          <div className='w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin' />
                        )}
                        {!isLoading && <CheckCircleIcon className='h-5 w-5' />}
                        <span>
                          {isLoading ? 'Creating Magic...' : 'Create Attribute'}
                        </span>
                      </div>
                      {!isLoading && fields.length > 0 && (
                        <div className='absolute inset-0 rounded-2xl bg-gradient-to-r from-blue-600 via-purple-600 to-pink-600 opacity-0 hover:opacity-20 transition-opacity duration-300' />
                      )}
                    </Button>
                  </div>
                </div>
              </form>
            </div>
          </DialogPanel>
        </div>
      </div>
    </Dialog>
  );
};
