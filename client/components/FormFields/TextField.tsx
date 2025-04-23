import { HookFormProps } from '@/lib/definitions';
import { Field, Input, Label } from '@headlessui/react';
import clsx from 'clsx';
import React from 'react';

interface TextFieldProps extends HookFormProps {
  label: string;
  message?: string;
  className?: string;
  placeholder?: string;
  error?: boolean;
}
export const TextField: React.FC<TextFieldProps> = (props) => {
  const { label, placeholder, error, message, className, ...rest } = props;
  return (
    <Field className={clsx(className, 'w-full ease-in-out')}>
      <Label className={'text-base/6 text-gray-500 font-semibold'}>
        {label}
      </Label>
      <Input
        {...rest}
        className={clsx(
          'border border-gray-300 mt-1 transition-all duration-500 rounded-md p-3 w-full shadow-none',
          'focus:ring-1 focus:ring-sky-400 focus:outline-none focus:shadow-lg',
          rest.required ? 'border-l-4 border-orange-400' : ''
        )}
        placeholder={placeholder}
      />
      {message && (
        <div
          className={clsx(
            'text-sm mt-2',
            error ? 'text-red-500' : 'text-gray-500'
          )}
        >
          {props.message}
        </div>
      )}
    </Field>
  );
};
