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
import { BaseOption } from '@/lib/definitions';

interface StyledMultipleComboBoxProps<T extends BaseOption> {
  label: string;
  selected: T[] | null[] | undefined;
  options: T[];
  message?: string;
  error?: boolean;
  setSelected: (value: T[]) => void;
}

export const StyledMultipleComboBox = <T extends BaseOption>({
  selected,
  setSelected,
  label,
  error,
  message,
  options,
}: StyledMultipleComboBoxProps<T>) => {
  const [query, setQuery] = useState('');
  const [filteredOptions, setFilteredOptions] = useState(options);

  useEffect(() => {
    setFilteredOptions(
      options.filter((opt) =>
        opt.name.toLowerCase().includes(query.toLowerCase())
      )
    );
  }, [query, options]);

  return (
    <Field className='w-full'>
      <Label className='text-sm/6 text-gray-500 font-semibold'>{label}</Label>
      <Combobox
        value={selected}
        name='category'
        multiple
        onChange={(value) => setSelected(value as T[])}
        onClose={() => setQuery('')}
      >
        <div className='relative w-full mt-1'>
          <ComboboxInput
            placeholder='Select category...'
            className={clsx(
              'w-full rounded-lg border border-gray-300 bg-white py-3 pr-8 pl-3 text-sm/6 text-gray-500 transition-all duration-500 shadow-none ease-in-out',
              'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 focus:ring-1 focus:ring-sky-400 focus:shadow-lg'
            )}
            displayValue={(option: BaseOption[]) =>
              option?.map((opt) => opt.name).join(', ')
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
          {filteredOptions.map((opt) => (
            <ComboboxOption
              key={opt.id}
              value={opt}
              className='group flex cursor-default items-center gap-2 rounded-lg py-1.5 px-3 select-none data-[focus]:bg-gray-100'
            >
              <CheckIcon className='invisible size-4 group-data-[selected]:visible text-green-300' />
              <div className='text-sm/6 text-gray-500'>{opt.name}</div>
            </ComboboxOption>
          ))}
        </ComboboxOptions>
      </Combobox>
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
