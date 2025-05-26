import { Checkbox } from '@headlessui/react';
import { CheckIcon } from '@heroicons/react/16/solid';
import React from 'react';

type StyledCheckboxProps = {
  label: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
};
export const StyledCheckbox: React.FC<StyledCheckboxProps> = ({
  checked,
  label,
  onChange,
}) => {
  return (
    <Checkbox
      checked={checked}
      onChange={onChange}
      className='flex items-center cursor-pointer'
    >
      {({ checked }) => (
        <>
          <div className='flex h-5 w-5 items-center justify-center rounded border border-gray-300 bg-white'>
            {checked && <CheckIcon className='h-4 w-4 text-indigo-600' />}
          </div>
          <span className='ml-3 text-sm text-gray-600'>{label}</span>
        </>
      )}
    </Checkbox>
  );
};
