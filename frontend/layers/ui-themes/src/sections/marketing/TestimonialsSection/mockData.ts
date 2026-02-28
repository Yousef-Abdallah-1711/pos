/**
 * Testimonials Section Mock Data
 * Provides realistic customer testimonials
 */

export const TESTIMONIALS_MOCK_DATA = {
  title_en: 'What Our Customers Say',
  title_ar: 'ما يقوله عملاؤنا',
  description_en: 'Real reviews from real customers',
  description_ar: 'تقييمات حقيقية من عملائنا الحقيقيين',
  items: [
    {
      author: 'John Smith',
      content_en: 'Best restaurant in town! The food is amazing and the service is outstanding.',
      content_ar: 'أفضل مطعم في المدينة! الطعام رائع والخدمة ممتازة.',
      rating: 5,
      image: 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop',
    },
    {
      author: 'Sarah Johnson',
      content_en: 'Absolutely delicious! I bring my family here every weekend.',
      content_ar: 'طعم لذيذ جداً! أحضر عائلتي هنا كل نهاية أسبوع.',
      rating: 5,
      image: 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=100&h=100&fit=crop',
    },
    {
      author: 'Michael Brown',
      content_en: 'Great ambiance and excellent food quality. Highly recommended!',
      content_ar: 'أجواء رائعة وجودة طعام ممتازة. موصى به بشدة!',
      rating: 5,
      image: 'https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=100&h=100&fit=crop',
    },
    {
      author: 'Emily Davis',
      content_en: 'The staff is very friendly and attentive. Will definitely come back!',
      content_ar: 'الموظفون ودودون جداً واهتمامهم واضح. سأعود بالتأكيد!',
      rating: 5,
      image: 'https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100&h=100&fit=crop',
    },
  ],
  config: {
    layout: 'carousel',
    show_rating: true,
  },
}
