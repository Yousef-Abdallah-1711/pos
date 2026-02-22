'use client';

import React from 'react';

interface ThemePreviewProps {
  primaryColor: string;
  secondaryColor: string;
  accentColor: string;
  logoUrl: string;
  fontFamily: string;
}

export function ThemePreview({
  primaryColor,
  secondaryColor,
  accentColor,
  logoUrl,
  fontFamily,
}: ThemePreviewProps) {
  return (
    <div className="theme-preview">
      {/* Header */}
      <div className="preview-header p-6 text-white flex items-center justify-between">
        {logoUrl ? (
          <img
            src={logoUrl}
            alt="Restaurant Logo"
            className="h-12 w-auto object-contain"
          />
        ) : (
          <div className="text-2xl font-bold">Restaurant</div>
        )}
        <nav className="flex gap-6">
          <a href="#" className="text-sm hover:opacity-80">Menu</a>
          <a href="#" className="text-sm hover:opacity-80">Orders</a>
          <a href="#" className="text-sm hover:opacity-80">Contact</a>
        </nav>
      </div>

      {/* Hero Section */}
      <div className="preview-hero p-12 text-white text-center">
        <h1 className="text-4xl font-bold mb-4">Welcome to Our Restaurant</h1>
        <p className="text-lg mb-6">Experience authentic flavors and exceptional service</p>
        <button className="accent-button px-8 py-3 rounded font-semibold hover:opacity-90 transition-opacity">
          Order Now
        </button>
      </div>

      {/* Featured Section */}
      <div className="p-12 bg-gray-50">
        <h2 className="text-3xl font-bold mb-8 text-gray-900">Featured Items</h2>
        <div className="grid grid-cols-3 gap-6">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-white rounded-lg overflow-hidden shadow">
              <div className="featured-image h-32 flex items-center justify-center">
                <span className="text-white text-sm">Image {i}</span>
              </div>
              <div className="p-4">
                <h3 className="font-semibold text-gray-900 mb-2">Dish Name</h3>
                <p className="text-sm text-gray-600 mb-4">
                  Delicious and authentic dish
                </p>
                <button className="accent-button w-full py-2 rounded text-sm font-medium hover:opacity-90 transition-opacity">
                  View Details
                </button>
              </div>
            </div>
          ))}
        </div>
      </div>

      {/* Why Choose Us Section */}
      <div className="p-12 bg-white">
        <h2 className="text-3xl font-bold mb-8 text-gray-900">Why Choose Us</h2>
        <div className="grid grid-cols-4 gap-6">
          {[
            { title: 'Quality', icon: '✨' },
            { title: 'Fast Delivery', icon: '⚡' },
            { title: 'Wide Variety', icon: '🎯' },
            { title: '24/7 Support', icon: '🛟' },
          ].map((item) => (
            <div key={item.title} className="text-center">
              <div className="text-4xl mb-4">{item.icon}</div>
              <h3 className="primary-text font-semibold text-lg">
                {item.title}
              </h3>
              <p className="text-sm text-gray-600 mt-2">
                Premium quality assurance
              </p>
            </div>
          ))}
        </div>
      </div>

      {/* CTA Section */}
      <div className="preview-cta p-12 text-white text-center">
        <h2 className="text-3xl font-bold mb-4">Ready to Order?</h2>
        <p className="text-lg mb-6">Delicious food is just a click away</p>
        <button className="accent-button px-8 py-3 rounded font-semibold hover:opacity-90 transition-opacity">
          Start Ordering
        </button>
      </div>

      {/* Footer */}
      <div className="bg-gray-900 text-white p-8">
        <div className="grid grid-cols-3 gap-8 mb-8">
          <div>
            <h4 className="font-semibold mb-4">Hours</h4>
            <p className="text-sm text-gray-400">Mon-Fri: 11am - 11pm</p>
          </div>
          <div>
            <h4 className="font-semibold mb-4">Location</h4>
            <p className="text-sm text-gray-400">123 Main Street, City</p>
          </div>
          <div>
            <h4 className="font-semibold mb-4">Contact</h4>
            <p className="text-sm text-gray-400">+1 (555) 123-4567</p>
          </div>
        </div>
        <div className="border-t border-gray-800 pt-8 text-center text-sm text-gray-400">
          <p>© 2024 Our Restaurant. All rights reserved.</p>
        </div>
      </div>

      <style jsx>{`
        .theme-preview {
          font-family: ${fontFamily};
        }

        .preview-header {
          background-color: ${primaryColor};
        }

        .preview-hero {
          background: linear-gradient(135deg, ${primaryColor} 0%, ${secondaryColor} 100%);
        }

        .accent-button {
          background-color: ${accentColor};
          color: #000000;
        }

        .featured-image {
          background-color: ${primaryColor};
        }

        .primary-text {
          color: ${primaryColor};
        }

        .preview-cta {
          background-color: ${secondaryColor};
        }
      `}</style>
    </div>
  );
}