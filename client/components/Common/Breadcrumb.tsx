"use client";

import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChevronRightIcon } from '@heroicons/react/16/solid';

interface BreadcrumbItem {
  label: string;
  href?: string;
}

interface BreadcrumbProps {
  items: BreadcrumbItem[];
  className?: string;
}

/**
 * Breadcrumb component for navigation
 * @param items - Array of breadcrumb items with labels and optional hrefs
 * @param className - Optional additional CSS classes
 */
export const Breadcrumb: React.FC<BreadcrumbProps> = ({ 
  items,
  className = "",
}) => {
  return (
    <nav aria-label="Breadcrumb" className={`flex items-center text-sm ${className}`}>
      <ol className="flex items-center space-x-2">
        {items.map((item, index) => (
          <li key={index} className="flex items-center">
            {index > 0 && (
              <ChevronRightIcon className="mx-2 size-4 text-gray-400 flex-shrink-0" />
            )}
            {item.href && index !== items.length - 1 ? (
              <Link 
                href={item.href}
                className="text-gray-600 hover:text-gray-900 hover:underline transition-colors"
              >
                {item.label}
              </Link>
            ) : (
              <span className="font-medium text-gray-900">{item.label}</span>
            )}
          </li>
        ))}
      </ol>
    </nav>
  );
};

/**
 * Dynamic breadcrumb component that automatically generates breadcrumb items from the current path
 * @param className - Optional additional CSS classes
 * @param homeLabel - Label for the home page, defaults to "Home"
 * @param excludeSegments - Array of path segments to exclude from breadcrumb
 */
export const DynamicBreadcrumb: React.FC<{
  className?: string;
  homeLabel?: string;
  excludeSegments?: string[];
}> = ({ 
  className = "",
  homeLabel = "Home",
  excludeSegments = []
}) => {
  const pathname = usePathname();
  
  // Skip rendering if we're on the home page
  if (pathname === "/") return null;
  
  const segments = pathname
    .split("/")
    .filter(segment => segment && !excludeSegments.includes(segment));
  
  const breadcrumbItems: BreadcrumbItem[] = [
    { label: homeLabel, href: "/" }
  ];
  
  let path = "";
  segments.forEach((segment) => {
    path += `/${segment}`;
    breadcrumbItems.push({
      label: segment.charAt(0).toUpperCase() + segment.slice(1).replace(/-/g, " "),
      href: path
    });
  });
  
  // Make the last item without href
  if (breadcrumbItems.length > 0) {
    const lastItem = breadcrumbItems[breadcrumbItems.length - 1];
    breadcrumbItems[breadcrumbItems.length - 1] = { label: lastItem.label };
  }
  
  return <Breadcrumb items={breadcrumbItems} className={className} />;
};