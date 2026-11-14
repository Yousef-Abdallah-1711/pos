import Link from 'next/link'
import { getSiteCopy, type Locale } from '@/content/site'
import { TrialLink } from './TrialLink'
import { SectionHead } from './SectionHead'

export function GettingStarted({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  return <section id="getting-started" className="sv-section sv-setup"><div className="sv-shell">
    <SectionHead eyebrow={c.setupEyebrow} title={c.setupTitle} />
    <ol>{c.setup.map(([title, body], i) => <li key={title}><span>0{i + 1}</span><h3>{title}</h3><p>{body}</p></li>)}</ol>
    <div className="sv-setup-foot"><p className="sv-note">{c.setupNote}</p><TrialLink locale={locale} source="getting-started" /></div>
  </div></section>
}

export function MarketingFAQ({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  const questions = locale === 'ar' ? [
    ['كيف يتصل الكاشير بالمطبخ؟', 'يؤكد إرسال طلب الكاشير الطلب، ثم ينشئ تذكرته في شاشة المطبخ. يتابع الفريق مراحل التجهيز من هناك.'],
    ['هل يمكن للضيف إرسال طلب إلكتروني؟', 'نعم، المسار الموثّق يتيح تصفّح القائمة وإرسال طلب للدفع النقدي عند الاستلام. يؤكد الفريق الطلب قبل انتقاله للمطبخ.'],
    ['ما الذي يمكن تخصيصه في الموقع؟', 'تم التحقق من تعديل ألوان القالب وحفظها ومعاينتها. تجهيز العنوان العام خطوة مستقلة؛ لا نعد بالنشر الفوري أو بإنشاء صفحات حرة.'],
    ['كيف تبدأ التجربة المجانية؟', 'أنشئ حساب المطعم لبدء تجربة لمدة 14 يومًا. يمكنك بعدها إعداد بيانات المطعم وإضافة الأصناف وتجربة مسار الطلب.'],
    ['هل اختيار الفترة السنوية يخصم المبلغ؟', 'لا ينفّذ اختيار فترة الفوترة عملية دفع. الدفع للاشتراك غير متاح هنا حاليًا، ويتطلب التفعيل المدفوع تأكيدًا مستقلًا.'],
    ['هل يتوفر الموقع بالعربية والإنجليزية؟', 'نعم، الموقع التعريفي متاح باللغتين ويدعم الاتجاه من اليمين إلى اليسار. الشاشات التجريبية موسومة بلغتها؛ لا نفترض اكتمال ترجمة كل شاشة داخل المنتج.'],
  ] : [
    ['How do the cashier and kitchen connect?', 'Submitting a cashier order confirms it and creates its kitchen ticket. The kitchen team then moves the ticket through preparation states.'],
    ['Can guests order online?', 'The demonstrated flow lets a guest browse the menu and submit a cash-on-delivery order. Staff confirm it before it enters the kitchen.'],
    ['What can I customize on my storefront?', 'Theme color editing, saving and previewing are demonstrated. Public hostname setup is separate; instant publishing and arbitrary page creation are not promised.'],
    ['How does the free trial start?', 'Create a restaurant account to begin a 14-day trial. Then configure restaurant information, add products and try a supported order workflow.'],
    ['Does selecting annual billing charge me?', 'Selecting a billing preference does not initiate a payment. Paid checkout is not available here; paid activation requires separate confirmation.'],
    ['Are Arabic and English available?', 'The public website supports both languages and right-to-left layout. Product captures are labelled with their language; complete translation of every application screen is not implied.'],
  ]
  return <section id="faq" className="sv-section sv-faq"><div className="sv-shell sv-faq-layout"><SectionHead eyebrow={c.faqEyebrow} title={c.faqTitle} />
    <div>{questions.map(([question, answer]) => <details key={question}><summary>{question}<span aria-hidden="true">+</span></summary><p>{answer}</p></details>)}<Link className="sv-text-link" href={`/${locale}/contact`}>{c.contact}<span aria-hidden="true">↗</span></Link></div>
  </div></section>
}

export function FinalCTA({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  return <section id="start" className="sv-section sv-final"><div className="sv-shell">
    <p className="sv-eyebrow">{c.finalEyebrow}</p><h2>{c.finalTitle}</h2><div className="sv-final-bottom"><p>{c.finalBody}</p><TrialLink locale={locale} source="final" secondary /></div>
  </div></section>
}
