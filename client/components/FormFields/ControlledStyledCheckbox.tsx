import { Checkbox } from '@headlessui/react';
import { CheckIcon } from '@heroicons/react/16/solid';
import React from 'react';

import { Control, Controller, FieldValues, Path } from 'react-hook-form';
interface ControlledStyledCheckboxProps<T extends FieldValues> {
  label?: React.ReactNode;
  error?: string;
  name: Path<T>;
  getDisplayValue?: (option: any) => string;
  placeholder?: string;
  disabled?: boolean;
  control: Control<T>;
  nullable?: boolean;
}
export const ControlledStyledCheckbox = <T extends FieldValues>(
  props: ControlledStyledCheckboxProps<T>
) => {
  const { control, label, error, disabled, name } = props;
  return (
    <Controller
      name={name}
      control={control}
      disabled={disabled}
      render={({ field: { onChange, ...rest } }) => (
        <div>
          <Checkbox
            {...rest}
            onChange={onChange}
            checked={rest.value}
            className='flex items-center cursor-pointer'
          >
            {({ checked }) => (
              <>
                <div className='flex h-5 w-5 items-center justify-center rounded border border-gray-300 bg-white'>
                  {checked && (
                    <CheckIcon className='h-6 w-6 text-indigo-600 font-bold' />
                  )}
                </div>
                <span className='ml-3 text-sm text-gray-600'>{label}</span>
              </>
            )}
          </Checkbox>
          {error && <span className='text-sm text-red-500'>{error}</span>}
        </div>
      )}
    />
  );
};
