import type { Locale } from './site'

export interface Feature {
  slug: string
  title: string
  eyebrow: string
  headline: string
  description: string
  image: keyof typeof import('./captures.json')
  steps: [string, string][]
  boundary: string
}

const en: Feature[] = [
  { slug: 'pos', title: 'Cashier POS', eyebrow: 'The counter', headline: 'Take the order. Keep the context.', image: 'pos',
    description: 'A menu-led cashier workspace for selecting products and sending confirmed orders into kitchen preparation.',
    steps: [['Select from the menu', 'Browse categories and add the requested products to the cart.'], ['Add the service context', 'Attach a table and guest details where the order needs them.'], ['Send it to preparation', 'Cashier submission confirms the order and creates its kitchen ticket.']],
    boundary: 'The capture shows menu selection with an empty cart. Tax, payment, receipt and offline workflows require separate acceptance; no payment terminal or offline reliability is promised.' },
  { slug: 'kitchen-display', title: 'Kitchen display', eyebrow: 'The kitchen', headline: 'A clear next action for every ticket.', image: 'kitchen',
    description: 'See incoming orders and coordinate preparation through new, in-progress, ready and completed ticket states.',
    steps: [['Receive the confirmed order', 'An open kitchen display receives a ticket for the confirmed order.'], ['Move preparation forward', 'Bump the ticket to in progress, then ready, as work advances.'], ['Review and recall', 'Find completed tickets and recall one when the team needs to revisit it.']],
    boundary: 'Kitchen completion is separate from delivery and payment. The demonstrated handoff is not a guarantee of delivery under every network condition or of multi-station routing.' },
  { slug: 'online-ordering', title: 'Online ordering', eyebrow: 'The guest experience', headline: 'From your guest’s menu to your team’s queue.', image: 'checkout',
    description: 'Guests browse a restaurant storefront, review a cash-on-delivery checkout and receive an order confirmation.',
    steps: [['Choose from the restaurant menu', 'The guest adds a product and reviews the cart.'], ['Send a cash order', 'Checkout creates a pending order and returns its confirmation.'], ['Let staff confirm', 'Staff confirmation sends the order into the kitchen preparation queue.']],
    boundary: 'The demonstrated flow uses the configured demo storefront and cash on delivery. Online card processing and instant availability of newly registered public hostnames are not promised.' },
  { slug: 'website-builder', title: 'Storefront customization', eyebrow: 'Your restaurant’s presentation', headline: 'Your colors. A menu that feels like yours.', image: 'theme',
    description: 'Use the theme editor to adjust colors, save the configuration, preview it and activate a restaurant theme.',
    steps: [['Start with a theme', 'Review a restaurant theme or duplicate it into your own configuration.'], ['Make a considered change', 'Adjust a color and inspect the preview; the editor also surfaces contrast feedback.'], ['Save and activate', 'Saved configuration persists after reload and can be exposed by the restaurant’s public theme endpoint.']],
    boundary: 'This is storefront configuration, not arbitrary page building. Editor previews can use sample content. Public hostname and TLS setup is a separate onboarding requirement; instant publishing is not promised.' },
  { slug: 'inventory', title: 'Inventory', eyebrow: 'Tracked product quantities', headline: 'See the stock effect of a confirmed order.', image: 'inventory',
    description: 'Connect tracked product quantities with order confirmation and eligible unfulfilled cancellations.',
    steps: [['Start with tracked stock', 'The demonstrated product began with 20 available units.'], ['Confirm one item', 'Confirmation reserved one unit and changed the available quantity to 19.'], ['Release an eligible cancellation', 'Cancelling the unfulfilled order restored the quantity to 20.']],
    boundary: 'This example does not establish ingredient recipes, purchasing, historical valuation or automatic physical returns. Refunding money is different from returning stock.' },
  { slug: 'hr', title: 'Staff records', eyebrow: 'The team behind service', headline: 'Keep attendance and leave decisions organized.', image: 'staff',
    description: 'Maintain employee information, attendance records and leave requests within the restaurant workspace.',
    steps: [['Keep employee information', 'Create and review the employee’s record.'], ['Record attendance', 'Record clock-in and maintain clock-out information through attendance editing.'], ['Review a leave request', 'Record approval or rejection and revisit the request’s resulting state.']],
    boundary: 'These are demonstrated record-keeping actions. Statutory payroll, automatic working-hour accuracy, country-specific HR compliance and performance-review workflows are not advertised.' },
  { slug: 'analytics', title: 'Operational overview', eyebrow: 'Daily visibility', headline: 'Know what needs your attention next.', image: 'dashboard',
    description: 'Review recent orders, pending work, basic workspace counts and subscription status in an operational overview.',
    steps: [['Review recent work', 'Use the latest loaded orders and pending-order count as a starting point.'], ['Check the workspace', 'Review available product, team and table information for your current restaurant.'], ['Understand the data context', 'The overview identifies unavailable sources and shows when data was loaded.']],
    boundary: 'Recent paid revenue is based on the latest loaded orders, not a complete accounting period. Financial statements, forecast accuracy and consolidated multi-location reporting are outside this demonstrated scope.' },
]

const ar: Feature[] = [
  { slug: 'pos', title: 'نقاط البيع', eyebrow: 'عند الكاشير', headline: 'سجّل الطلب مع تفاصيل الخدمة.', image: 'pos', description: 'مساحة عمل تبدأ من القائمة لاختيار الأصناف وإرسال الطلبات المؤكدة إلى المطبخ.', steps: [['اختر الأصناف', 'تصفّح التصنيفات وأضف الأصناف المطلوبة إلى السلة.'], ['أضف سياق الخدمة', 'اربط الطلب بالطاولة وأضف بيانات الضيف عند الحاجة.'], ['أرسله للتجهيز', 'يؤكد إرسال طلب الكاشير الطلب وينشئ تذكرته في المطبخ.']], boundary: 'تعرض الشاشة اختيار القائمة مع سلة فارغة. تتطلب الضرائب والدفع والإيصالات والعمل دون اتصال تحققًا مستقلًا؛ لا نعد بدعم أجهزة دفع أو موثوقية العمل دون اتصال.' },
  { slug: 'kitchen-display', title: 'شاشة المطبخ', eyebrow: 'داخل المطبخ', headline: 'خطوة تالية واضحة لكل تذكرة.', image: 'kitchen', description: 'تابع الطلبات الواردة وتقدّم تجهيزها عبر حالات جديد وقيد التجهيز وجاهز ومكتمل.', steps: [['استقبل الطلب المؤكد', 'تصل تذكرة الطلب المؤكد إلى شاشة المطبخ المفتوحة.'], ['تابع تقدّم التجهيز', 'انقل التذكرة إلى قيد التجهيز ثم إلى جاهز مع تقدّم العمل.'], ['راجع واستدعِ التذكرة', 'ابحث عن التذاكر المكتملة واستدعِ إحداها عند الحاجة.']], boundary: 'اكتمال تذكرة المطبخ مستقل عن التسليم والتحصيل. المسار الموثّق ليس ضمانًا لوصول التذاكر في كل ظروف الشبكة أو لتوجيه الطلبات بين محطات متعددة.' },
  { slug: 'online-ordering', title: 'الطلب الإلكتروني', eyebrow: 'تجربة الضيف', headline: 'من قائمة الضيف إلى قائمة عمل فريقك.', image: 'checkout', description: 'يتصفّح الضيف موقع المطعم ويراجع طلب الدفع النقدي عند الاستلام، ثم يحصل على تأكيد الطلب.', steps: [['اختر من القائمة', 'يضيف الضيف الصنف ويراجع السلة.'], ['أرسل طلبًا نقديًا', 'ينشئ إتمام الطلب طلبًا معلقًا ويعرض تأكيده.'], ['يؤكد فريقك الطلب', 'بعد تأكيد الفريق، ينتقل الطلب إلى قائمة تجهيز المطبخ.']], boundary: 'المسار الموثّق يستخدم الموقع التجريبي المهيّأ والدفع النقدي عند الاستلام. لا نعد هنا بالدفع الإلكتروني بالبطاقات أو بتوفر عنوان موقع جديد فور التسجيل.' },
  { slug: 'website-builder', title: 'تخصيص موقع المطعم', eyebrow: 'هوية مطعمك', headline: 'ألوانك، وقائمة تحمل طابعك.', image: 'theme', description: 'اضبط الألوان في محرر القالب، واحفظ الإعدادات، وعاينها، ثم فعّل قالب المطعم.', steps: [['ابدأ بقالب', 'راجع القالب أو أنشئ نسخة منه لإعداداتك الخاصة.'], ['عدّل بعناية', 'غيّر اللون وراجع المعاينة وملاحظات التباين التي يعرضها المحرر.'], ['احفظ وفعّل', 'تبقى الإعدادات بعد إعادة التحميل ويمكن عرضها عبر واجهة القالب العامة.']], boundary: 'هذه تهيئة لموقع المطعم، وليست إنشاءً حرًا لصفحات عشوائية. قد تستخدم المعاينة محتوى تجريبيًا. تجهيز العنوان العام وTLS متطلب مستقل؛ لا نعد بالنشر الفوري.' },
  { slug: 'inventory', title: 'المخزون', eyebrow: 'كميات الأصناف المتتبَّعة', headline: 'تابع أثر الطلب المؤكد على المخزون.', image: 'inventory', description: 'اربط كميات الأصناف بتأكيد الطلبات وإلغاء الطلبات غير المنفّذة متى كانت مؤهلة.', steps: [['ابدأ بمخزون متتبَّع', 'بدأ الصنف في المثال الموثّق بكمية 20 وحدة.'], ['أكّد طلب وحدة واحدة', 'حجز التأكيد وحدة واحدة فأصبحت الكمية المتاحة 19.'], ['أعد الكمية عند الإلغاء المؤهل', 'أعاد إلغاء الطلب غير المنفّذ الكمية إلى 20.']], boundary: 'لا يثبت هذا المثال تتبّع مكونات الوصفات أو المشتريات أو التقييم التاريخي أو الإرجاع الفعلي الآلي. رد المبلغ يختلف عن إعادة الصنف للمخزون.' },
  { slug: 'hr', title: 'سجلات الفريق', eyebrow: 'الفريق خلف الخدمة', headline: 'نظّم الحضور وقرارات الإجازة.', image: 'staff', description: 'احتفظ ببيانات الموظفين وسجلات الحضور وطلبات الإجازة ضمن مساحة المطعم.', steps: [['نظّم بيانات الموظف', 'أنشئ سجل الموظف وراجعه.'], ['سجّل الحضور', 'سجّل وقت الحضور وحدّث بيانات الانصراف من سجل الحضور.'], ['راجع طلب الإجازة', 'سجّل الموافقة أو الرفض وراجع حالة الطلب بعد القرار.']], boundary: 'هذه إجراءات موثّقة لحفظ السجلات. لا نعلن هنا امتثال الرواتب أو دقة حساب الساعات تلقائيًا أو الامتثال المحلي للموارد البشرية أو اكتمال تقييم الأداء.' },
  { slug: 'analytics', title: 'النظرة التشغيلية', eyebrow: 'متابعة العمل اليومي', headline: 'اعرف ما يحتاج اهتمامك الآن.', image: 'dashboard', description: 'راجع أحدث الطلبات والعمل المعلق والأعداد الأساسية وحالة الاشتراك من لوحة تشغيلية.', steps: [['راجع أحدث العمل', 'ابدأ بأحدث الطلبات المحمّلة وعدد الطلبات المعلقة.'], ['تابع مساحة المطعم', 'راجع بيانات الأصناف والفريق والطاولات المتاحة لمطعمك الحالي.'], ['افهم سياق البيانات', 'توضح اللوحة المصادر غير المتاحة ووقت تحميل البيانات.']], boundary: 'يُحسب الإيراد المدفوع الأخير من أحدث الطلبات المحمّلة، وليس فترة محاسبية كاملة. القوائم المالية ودقة التنبؤ والتقارير المجمعة لعدة مواقع خارج هذا النطاق الموثّق.' },
]

export function getFeatures(locale: Locale): Feature[] {
  return locale === 'ar' ? ar : en
}
