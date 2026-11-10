'use client'

import { useLanguage } from '@/contexts/LanguageContext'

export function FAQ() {
  const { language, dir } = useLanguage()
  const items = language === 'ar'
    ? [
        ['ما هو POOS؟', 'منصة لإدارة المطاعم تربط العمل في واجهة الخدمة بما يجري خلفها.'],
        ['لمن صُمم POOS؟', 'لأصحاب المطاعم المستقلة ومشغلي مجموعات المطاعم متعددة المواقع.'],
        ['ما البيانات المطلوبة لإنشاء الحساب؟', 'الاسم والبريد الإلكتروني وكلمة المرور واسم المطعم ورابطه. رقم الهاتف اختياري.'],
        ['هل يدعم POOS العربية والإنجليزية؟', 'تدعم الواجهة العامة العربية باتجاه من اليمين إلى اليسار والإنجليزية باتجاه من اليسار إلى اليمين.'],
        ['أين أجد الأسعار؟', 'لم تُنشر الأسعار والشروط التجارية على هذه الصفحة.'],
      ]
    : [
        ['What is POOS?', 'POOS is a restaurant management platform connecting work at the front of house with work behind it.'],
        ['Who is POOS for?', 'Independent restaurant owners and operators, and people managing multi-location restaurant groups.'],
        ['What do I need to create an account?', 'Your name, email, password, restaurant name, and restaurant URL. Phone is optional.'],
        ['Does POOS support Arabic and English?', 'The public interface supports Arabic with right-to-left layout and English with left-to-right layout.'],
        ['Where can I find pricing?', 'Pricing and commercial terms are not published on this page.'],
      ]

  return (
    <section id="faq" className="ledger-section" dir={dir}>
      <div className="poos-shell">
        <header className="ledger-section__head">
          <h2>{language === 'ar' ? 'أسئلة عن POOS' : 'Questions about POOS'}</h2>
          <p>{language === 'ar' ? 'إجابات واضحة قبل إنشاء حساب مطعمك.' : 'Clear answers before you create a restaurant account.'}</p>
        </header>
        <div className="ledger-faq">
          {items.map(([question, answer]) => (
            <details key={question}>
              <summary>{question}</summary>
              <p>{answer}</p>
            </details>
          ))}
        </div>
      </div>
    </section>
  )
}
