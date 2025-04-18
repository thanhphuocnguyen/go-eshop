import type { NextConfig } from 'next';

const nextConfig: NextConfig = {
  /* config options here */
  eslint: {
    dirs: ['components', 'app', 'utils', 'libs'],
  },
  reactStrictMode: true,
  typescript: {
    tsconfigPath: 'tsconfig.json',
  },
  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'res.cloudinary.com',
        port: '',
        pathname: '/dboesrelq/image/upload/**',
        search: '',
      },
    ],
  },
};

export default nextConfig;
