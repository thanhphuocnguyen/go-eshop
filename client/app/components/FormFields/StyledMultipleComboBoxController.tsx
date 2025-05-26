/* eslint-disable @typescript-eslint/no-explicit-any */
'use client';

import React, { useEffect, useState } from 'react';
import {
  Combobox,
  ComboboxButton,
  ComboboxInput,
  ComboboxOption,
  ComboboxOptions,
  Field,
  Label,
} from '@headlessui/react';
import clsx from 'clsx';
import { CheckIcon, ChevronDownIcon } from '@heroicons/react/24/outline';
import { Control, Controller, FieldValues, Path } from 'react-hook-form';

interface StyledMultipleComboBoxControllerProps<T extends FieldValues, D> {
  label: string;
  options: D[];
  getDisplayValue: (options: D[]) => string;
  getKey: (option: D) => string;
  error?: boolean;
  message?: string;
  name: Path<T>;
  placeholder?: string;
  disabled?: boolean;
  control: Control<T>;
}

export const StyledMultipleComboBoxController = <T extends FieldValues, D>({
  control,
  label,
  error,
  message,
  disabled,
  placeholder = 'Select options...',
  options,
  getDisplayValue = (option: D[]) =>
    option.map((e) => e.name as unknown as string).join(', '),
  getKey = (option: D) => option as unknown as string,
  name,
}: StyledMultipleComboBoxControllerProps<T, D>) => {
  const [query, setQuery] = useState('');
  const [filteredOptions, setFilteredOptions] = useState(options);

  useEffect(() => {
    setFilteredOptions(
      options.filter((opt) =>
        getDisplayValue(opt).toLowerCase().includes(query.toLowerCase())
      )
    );
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query, options]);

  return (
    <Field className='w-full'>
      <Label className='text-sm/6 text-gray-500 font-semibold'>{label}</Label>
      <Controller
        control={control}
        name={name}
        render={({ field: { onChange, ...rest } }) => (
          <Combobox
            {...rest}
            disabled={disabled}
            as={'div'}
            multiple
            onChange={(selectedValues) => {
              onChange(selectedValues);
            }}
            onClose={() => setQuery('')}
          >
            <div className='relative w-full mt-1'>
              <ComboboxInput
                placeholder={placeholder}
                className={clsx(
                  'w-full rounded-lg border border-gray-300 bg-white py-3 pr-8 pl-3 text-sm/6 text-gray-500 transition-all duration-500 shadow-none ease-in-out',
                  'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 focus:ring-1 focus:ring-sky-400 focus:shadow-lg'
                )}
                displayValue={(options: D) => getDisplayValue(options)}
                onChange={(event) => setQuery(event.target.value)}
              />
              <ComboboxButton className='group absolute inset-y-0 right-0 px-2.5'>
                <ChevronDownIcon className='size-4 fill-white/60 group-data-[hover]:fill-white' />
              </ComboboxButton>
            </div>

            <ComboboxOptions
              anchor='bottom'
              transition
              className={clsx(
                'w-[var(--input-width)] rounded-xl border border-gray-300 bg-white p-1 [--anchor-gap:var(--spacing-1)] empty:invisible',
                'transition duration-100 ease-in data-[leave]:data-[closed]:opacity-0'
              )}
            >
              {filteredOptions.map((opt) => (
                <ComboboxOption
                  key={getKey(opt)}
                  value={opt}
                  className='group flex cursor-default items-center gap-2 rounded-lg py-1.5 px-3 select-none data-[focus]:bg-tertiary'
                >
                  <CheckIcon className='invisible size-4 group-data-[selected]:visible text-green-300' />
                  <div className='text-sm/6 text-gray-500'>
                    {getDisplayValue(opt)}
                  </div>
                </ComboboxOption>
              ))}
            </ComboboxOptions>
          </Combobox>
        )}
      />
      {message && (
        <div
          className={clsx(
            'text-sm mt-2',
            error ? 'text-red-500' : 'text-gray-500'
          )}
        >
          {message}
        </div>
      )}
    </Field>
  );
};
