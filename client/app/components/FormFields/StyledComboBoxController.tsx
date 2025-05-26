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

interface StyledComboBoxControllerProps<T extends FieldValues> {
  label?: React.ReactNode;
  options: any[];
  error?: string;
  name: Path<T>;
  getDisplayValue?: (option: any) => string;
  placeholder?: string;
  disabled?: boolean;
  control: Control<T>;
  nullable?: boolean;
}

export const StyledComboBoxController = <T extends FieldValues>({
  control,
  label,
  error,
  disabled,
  placeholder,
  options,
  name,
  nullable,
  getDisplayValue,
}: StyledComboBoxControllerProps<T>) => {
  const [query, setQuery] = useState('');
  const [filteredOptions, setFilteredOptions] = useState(options);

  useEffect(() => {
    setFilteredOptions(
      options.filter((opt) =>
        (getDisplayValue ? getDisplayValue(opt) : (opt.name ?? ''))
          .toLowerCase()
          .includes(query.toLowerCase())
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
            as={'div'}
            onChange={onChange}
            disabled={disabled}
            onClose={() => setQuery('')}
          >
            <div className='relative w-full mt-1'>
              <ComboboxInput
                placeholder={placeholder}
                className={clsx(
                  'w-full rounded-lg border border-gray-300 bg-white py-3 pr-8 pl-3 text-sm/6 text-gray-500 transition-all duration-500 shadow-none ease-in-out',
                  'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 focus:ring-1 focus:ring-sky-400 focus:shadow-lg'
                )}
                displayValue={
                  getDisplayValue ? getDisplayValue : (opt) => opt?.name
                }
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
              {nullable && (
                <ComboboxOption
                  className='group flex cursor-default items-center gap-2 rounded-lg py-1.5 px-3 select-none data-[focus]:bg-tertiary'
                  value={{ id: null, name: 'None' }}
                >
                  None
                </ComboboxOption>
              )}
              {filteredOptions.map((opt) => (
                <ComboboxOption
                  key={opt.id}
                  value={opt}
                  className='group flex cursor-default items-center gap-2 rounded-lg py-1.5 px-3 select-none data-[focus]:bg-tertiary'
                >
                  <CheckIcon className='invisible size-4 fill-white group-data-[selected]:visible' />
                  <div className='text-sm/6 text-gray-500'>
                    {getDisplayValue ? getDisplayValue(opt) : opt.name}
                  </div>
                </ComboboxOption>
              ))}
            </ComboboxOptions>
          </Combobox>
        )}
      />
      {error && (
        <div
          className={clsx(
            'text-sm mt-2',
            error ? 'text-red-500' : 'text-gray-500'
          )}
        >
          {error}
        </div>
      )}
    </Field>
  );
};
