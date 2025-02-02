'use client';
import {
  Button,
  Field,
  Fieldset,
  Input,
  Label,
  Legend,
} from '@headlessui/react';
import clsx from 'clsx';
import { useEffect, useState } from 'react';

export default function LoginPage() {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [formErrors, setFormErrors] = useState<{
    username?: string;
    password?: string;
  }>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsSubmitting(true);
    try {
      setIsSubmitting(false);
    } catch (error) {
      setIsSubmitting(false);
    }
  };

  useEffect(() => {
    if (username.length < 3) {
      setFormErrors((prev) => ({
        ...prev,
        username: 'Username must be at least 3 characters long',
      }));
    } else {
      setFormErrors((prev) => ({ ...prev, username: undefined }));
    }
  }, [username]);

  useEffect(() => {
    if (password.length < 6) {
      setFormErrors((prev) => ({
        ...prev,
        password: 'Password must be at least 6 characters long',
      }));
    } else {
      setFormErrors((prev) => ({ ...prev, password: undefined }));
    }
  }, [password]);

  return (
    <div className='w-full h-[90vh] flex items-center justify-center'>
      <div className='w-full max-w-lg px-4'>
        <Fieldset
          as='form'
          onSubmit={handleSubmit}
          className='space-y-6 rounded-xl bg-light-green shadow-md p-6 sm:p-10'
        >
          <Legend className='text-base/7 font-semibold text-white'>
            Login
          </Legend>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Username</Label>
            <Input
              disabled={false}
              id='username'
              name='username'
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {formErrors.username && (
              <Label className='text-red-500 text-sm/6 mt-1'>
                {formErrors.username}
              </Label>
            )}
          </Field>
          <Field>
            <Label className='text-sm/6 font-medium text-white'>Password</Label>
            <Input
              type='password'
              id='password'
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={clsx(
                'mt-3 block w-full rounded-lg border-none bg-white/5 py-1.5 px-3 text-sm/6 text-white',
                'focus:outline-none data-[focus]:outline-2 data-[focus]:-outline-offset-2 data-[focus]:outline-white/25'
              )}
            />
            {
              <Label className='text-red-500 text-sm/6 mt-1'>
                {formErrors.password}
              </Label>
            }
          </Field>
          <Button>{isSubmitting ? 'Logging in' : 'Login'}</Button>
        </Fieldset>
      </div>
    </div>
  );
}
