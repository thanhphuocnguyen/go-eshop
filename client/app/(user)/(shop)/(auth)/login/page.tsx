import React from 'react';
import LoginFormComponent from './_components/LoginForm';
import { Metadata } from 'next';

export const metadata: Metadata = {
  title: 'Sign In - eShop',
  description: 'Sign in to your eShop account',
};

export default function LoginPage() {
  return (
    <div className="container mx-auto flex justify-center items-center min-h-[70vh] py-8">
      <div className="w-full max-w-md">
        <LoginFormComponent />
      </div>
    </div>
  );
}
