import { HookFormProps } from '@/app/lib/definitions';
import { Field, Input, Label } from '@headlessui/react';
import clsx from 'clsx';
import React from 'react';

interface TextFieldProps extends HookFormProps {
  label?: React.ReactNode;
  icon?: React.ReactNode;
  step?: string;
  className?: string;
  placeholder?: string;
  error?: string;
}
export const TextField: React.FC<TextFieldProps> = (props) => {
  const {
    label,
    placeholder,
    error,
    className,
    icon,
    type = 'text',
    ...rest
  } = props;
  return (
    <Field className={clsx(className, 'w-full ease-in-out')}>
      {typeof label === 'string' ? (
        <Label className={'text-base/6 text-gray-500 font-semibold'}>
          {label}
        </Label>
      ) : (
        label
      )}
      <div className='relative'>
        <Input
          {...rest}
          step={props.step}
          type={type}
          className={clsx(
            'border border-gray-300 mt-1 transition-all duration-500 rounded-md p-3 w-full shadow-none',
            'focus:ring-1 focus:ring-sky-400 focus:outline-none focus:shadow-lg',
            rest.required ? 'border-l-4 border-orange-400' : '',
            icon ? 'pl-10' : '' // Add left padding when icon is present
          )}
          placeholder={placeholder}
        />
        {icon && (
          <div className='absolute inset-y-0 left-0 flex items-center pl-3 mt-1 pointer-events-none'>
            {icon}
          </div>
        )}
      </div>
      {error && (
        <div
          className={clsx(
            'text-sm mt-2',
            error ? 'text-red-500' : 'text-gray-500'
          )}
        >
          {props.error}
        </div>
      )}
    </Field>
  );
};
