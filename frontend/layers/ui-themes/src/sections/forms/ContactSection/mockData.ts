/**
 * Contact Section Mock Data
 * Provides contact information
 */

export const CONTACT_MOCK_DATA = {
  title_en: 'Get In Touch',
  title_ar: 'تواصل معنا',
  description_en: 'We would love to hear from you. Contact us today!',
  description_ar: 'نود أن نسمع منك. تواصل معنا اليوم!',
  address_en: '123 Main Street, City, State 12345',
  address_ar: '123 شارع رئيسي، المدينة، الولاية 12345',
  phone: '+1 (555) 123-4567',
  email: 'info@restaurant.com',
  hours: {
    monday_friday: '9:00 AM - 10:00 PM',
    saturday: '10:00 AM - 11:00 PM',
    sunday: '11:00 AM - 9:00 PM',
  },
  social_links: [
    { platform: 'facebook', url: 'https://facebook.com' },
    { platform: 'instagram', url: 'https://instagram.com' },
    { platform: 'twitter', url: 'https://twitter.com' },
  ],
  show_map: true,
  map_coordinates: {
    latitude: 40.7128,
    longitude: -74.006,
  },
  config: {
    layout: 'side-by-side',
    show_form: true,
  },
}
