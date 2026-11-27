import type { Locale } from './site'

const english = {
  sample: 'Illustrative order journey · sample data',
  order: 'Sample order',
  orderStatus: 'Order status',
  kitchenStatus: 'Kitchen ticket',
  annotation: 'Demonstration note',
  screenNote: 'Actual SERVIO product area; the sample order is an explanatory overlay.',
  stages: [
    { title: 'Browse the menu', body: 'A guest chooses items from the restaurant storefront.', order: 'Not submitted', kitchen: 'Not created', note: 'The storefront is public; no staff action has happened yet.', image: 'checkout' },
    { title: 'Submit the order', body: 'The guest sends a cash on delivery request with their details.', order: 'Pending', kitchen: 'Not created', note: 'The request waits for staff review. This is not an online card payment.', image: 'checkout' },
    { title: 'Staff confirms it', body: 'A team member reviews and confirms the online request.', order: 'Confirmed', kitchen: 'New', note: 'Online orders do not enter kitchen preparation before this action.', image: 'dashboard' },
    { title: 'Kitchen receives the ticket', body: 'The confirmed order becomes visible in the kitchen queue.', order: 'Confirmed', kitchen: 'New', note: 'The kitchen display has its own ticket states.', image: 'kitchen' },
    { title: 'Preparation begins', body: 'The kitchen moves the ticket into progress.', order: 'Preparing', kitchen: 'In progress', note: 'The team changes the state; the demonstration does not do it automatically.', image: 'kitchen' },
    { title: 'Ready for fulfillment', body: 'The kitchen marks preparation ready. Fulfillment is still a separate step.', order: 'Ready', kitchen: 'Ready', note: 'Ready describes preparation, not final delivery.', image: 'kitchen' },
    { title: 'Close the loop', body: 'Staff handle delivery and can see the final order state in the workspace.', order: 'Delivered', kitchen: 'Completed', note: 'The delivery path passes through out for delivery before delivered.', image: 'dashboard' },
  ],
} as const

const arabic = {
  sample: 'رحلة طلب توضيحية · بيانات تجريبية',
  order: 'طلب تجريبي',
  orderStatus: 'حالة الطلب',
  kitchenStatus: 'تذكرة المطبخ',
  annotation: 'ملاحظة توضيحية',
  screenNote: 'شاشة فعلية من سيرفيو؛ بيانات الطلب التجريبي موضحة فوقها.',
  stages: [
    { title: 'تصفح المنيو', body: 'الضيف يختار الأصناف من موقع المطعم.', order: 'لم يُرسل', kitchen: 'لم تُنشأ', note: 'الموقع متاح للضيف، والفريق لم يبدأ أي خطوة بعد.', image: 'checkout' },
    { title: 'إرسال الطلب', body: 'الضيف يرسل طلب الدفع عند الاستلام مع بياناته.', order: 'بانتظار التأكيد', kitchen: 'لم تُنشأ', note: 'الطلب ينتظر مراجعة الفريق؛ لا يوجد دفع إلكتروني بالبطاقة هنا.', image: 'checkout' },
    { title: 'تأكيد الفريق', body: 'أحد أفراد الفريق يراجع الطلب الأونلاين ويؤكده.', order: 'مؤكد', kitchen: 'جديدة', note: 'الطلب الأونلاين لا يدخل مرحلة التحضير قبل هذه الخطوة.', image: 'dashboard' },
    { title: 'وصوله للمطبخ', body: 'الطلب المؤكد يظهر في قائمة تذاكر المطبخ.', order: 'مؤكد', kitchen: 'جديدة', note: 'لتذكرة المطبخ حالات مختلفة عن حالات الطلب.', image: 'kitchen' },
    { title: 'بدء التحضير', body: 'الفريق ينقل التذكرة إلى حالة جارٍ التحضير.', order: 'جارٍ التحضير', kitchen: 'جارٍ التحضير', note: 'الفريق يغيّر الحالة بنفسه؛ العرض لا ينفذ الخطوة تلقائياً.', image: 'kitchen' },
    { title: 'جاهز للتسليم', body: 'المطبخ يحدد أن التحضير انتهى، والتسليم خطوة منفصلة.', order: 'جاهز', kitchen: 'جاهزة', note: 'جاهز تعني اكتمال التحضير، وليس اكتمال التوصيل.', image: 'kitchen' },
    { title: 'اكتمال الدورة', body: 'الفريق يتابع التوصيل ويرى الحالة النهائية في مساحة العمل.', order: 'تم التسليم', kitchen: 'مكتملة', note: 'مسار التوصيل يمر بحالة في الطريق قبل تم التسليم.', image: 'dashboard' },
  ],
} as const

export function getOrderStoryCopy(locale: Locale) {
  return locale === 'ar' ? arabic : english
}
