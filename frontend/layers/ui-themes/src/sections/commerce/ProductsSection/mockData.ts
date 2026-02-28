/**
 * Products Section Mock Data
 * Provides realistic sample products for the Products component
 */

export const PRODUCTS_MOCK_DATA = {
  title_en: 'Our Menu',
  title_ar: 'قائمتنا',
  description_en: 'Explore our carefully selected dishes',
  description_ar: 'استكشف أطباقنا المختارة بعناية',
  items: [
    {
      id: '1',
      name: 'Delicious Burger',
      description: 'Juicy beef burger with fresh toppings',
      price: 12.99,
      image: 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?w=300&h=300&fit=crop',
      category: 'Main Course',
    },
    {
      id: '2',
      name: 'Crispy Pizza',
      description: 'Authentic Italian pizza with mozzarella',
      price: 15.99,
      image: 'https://images.unsplash.com/photo-1604068549290-dea0e4a305ca?w=300&h=300&fit=crop',
      category: 'Main Course',
    },
    {
      id: '3',
      name: 'Fresh Salad',
      description: 'Organic vegetables with house dressing',
      price: 9.99,
      image: 'https://images.unsplash.com/photo-1540189549336-e6e99c3679fe?w=300&h=300&fit=crop',
      category: 'Appetizers',
    },
    {
      id: '4',
      name: 'Chocolate Dessert',
      description: 'Rich chocolate cake with cream',
      price: 7.99,
      image: 'https://images.unsplash.com/photo-1578985545062-69928b1d9587?w=300&h=300&fit=crop',
      category: 'Desserts',
    },
    {
      id: '5',
      name: 'Fresh Juice',
      description: 'Freshly squeezed orange juice',
      price: 5.99,
      image: 'https://images.unsplash.com/photo-1553530666-ba11a7da3888?w=300&h=300&fit=crop',
      category: 'Beverages',
    },
    {
      id: '6',
      name: 'Grilled Chicken',
      description: 'Tender grilled chicken with herbs',
      price: 14.99,
      image: 'https://images.unsplash.com/photo-1598103442097-8b74394b95c6?w=300&h=300&fit=crop',
      category: 'Main Course',
    },
  ],
  config: {
    layout: 'grid',
    columns: 3,
    show_prices: true,
    show_images: true,
  },
}
