import { HookFormProps } from '@/app/lib/definitions';
import { Field, Label, Textarea } from '@headlessui/react';
import clsx from 'clsx';
import React, { useState } from 'react';

interface TextareaFieldProps extends HookFormProps {
  label?: React.ReactNode;
  icon?: React.ReactNode;
  className?: string;
  placeholder?: string;
  error?: string;
  rows?: number;
  maxLength?: number;
  helperText?: string;
}

export const TextAreaField: React.FC<TextareaFieldProps> = (props) => {
  const {
    label,
    placeholder,
    error,
    className,
    icon,
    rows = 4,
    maxLength,
    helperText,
    ...rest
  } = props;
  
  const [charCount, setCharCount] = useState(0);
  
  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    if (rest.onChange) {
      rest.onChange(e);
    }
    setCharCount(e.target.value.length);
  };

  return (
    <Field className={clsx(className, 'w-full')}>
      {typeof label === 'string' ? (
        <Label className={'text-base/6 text-gray-500 font-semibold flex justify-between'}>
          <span>{label}</span>
          {maxLength && (
            <span className="text-sm font-normal text-gray-500">
              {charCount}/{maxLength}
            </span>
          )}
        </Label>
      ) : (
        label
      )}
      <div className='relative mt-1.5'>
        <Textarea
          {...rest}
          rows={rows}
          maxLength={maxLength}
          onChange={handleChange}
          className={clsx(
            'block w-full px-4 py-3 text-gray-900 border border-gray-300 rounded-md shadow-sm',
            'transition-all duration-300 ease-in-out',
            'focus:border-primary focus:ring-2 focus:ring-primary/20 focus:outline-none',
            'placeholder:text-gray-400 resize-y',
            rest.required ? 'border-l-4 border-l-orange-400' : '',
            error ? 'border-red-300 focus:border-red-500 focus:ring-red-200' : '',
            icon ? 'pl-10' : ''
          )}
          placeholder={placeholder}
        />
        {icon && (
          <div className='absolute top-3 left-3 text-gray-400 pointer-events-none'>
            {icon}
          </div>
        )}
      </div>
      
      {helperText && !error && (
        <p className="mt-1.5 text-sm text-gray-500">
          {helperText}
        </p>
      )}
      
      {error && (
        <div className="mt-1.5 text-sm text-red-500 flex items-center">
          <svg 
            xmlns="http://www.w3.org/2000/svg" 
            viewBox="0 0 20 20" 
            fill="currentColor" 
            className="w-4 h-4 mr-1 inline-block"
          >
            <path 
              fillRule="evenodd" 
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" 
              clipRule="evenodd" 
            />
          </svg>
          {props.error}
        </div>
      )}
    </Field>
  );
};
