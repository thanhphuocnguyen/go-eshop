import { Metadata } from 'next';
import RegisterFormComponent from './_components/RegisterForm';

export const metadata: Metadata = {
  title: 'Register simple life shop',
  description: 'An e-commerce site built with Next.js',
};

export default function RegisterPage() {
  return <RegisterFormComponent />;
}
