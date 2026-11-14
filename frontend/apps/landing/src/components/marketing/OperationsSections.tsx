import { FeatureActions } from './FeatureActions'
import { getSiteCopy, type Locale } from '@/content/site'
import { ProductCapture } from './ProductCapture'
import { SectionHead } from './SectionHead'

export function OperationsSections({ locale }: { locale: Locale }) {
  const c = getSiteCopy(locale)
  return <>
    <section id="pos" className="sv-section sv-pos"><div className="sv-shell">
      <div className="sv-editorial-row"><SectionHead eyebrow={c.posEyebrow} title={c.posTitle} body={c.posBody} />
        <ol className="sv-annotations">{c.posPoints.map((point, i) => <li key={point}><span>0{i + 1}</span>{point}</li>)}</ol>
      </div>
      <ProductCapture locale={locale} image="pos" alt={c.posTitle} className="sv-wide-capture" />
      <FeatureActions locale={locale} source="home-pos" destination="features/pos" label={c.posLink} />
    </div></section>
    <section id="kitchen" className="sv-section sv-dark sv-kitchen"><div className="sv-shell">
      <div className="sv-kitchen-heading"><SectionHead eyebrow={c.kitchenEyebrow} title={c.kitchenTitle} body={c.kitchenBody} />
        <FeatureActions locale={locale} source="home-kitchen-display" destination="features/kitchen-display" label={c.kitchenLink} />
      </div>
      <div className="sv-kitchen-stage"><ProductCapture locale={locale} image="kitchen" alt={c.kitchenTitle} />
        <div className="sv-state-ledger">{(locale === 'ar' ? ['جديد', 'قيد التجهيز', 'جاهز', 'مكتمل'] : ['New', 'In progress', 'Ready', 'Completed']).map((state, i) => <p key={state}><span>0{i + 1}</span>{state}</p>)}</div>
      </div><p className="sv-note">{c.kitchenNote}</p>
    </div></section>
    <section id="online-ordering" className="sv-section sv-ordering"><div className="sv-shell sv-ordering-grid">
      <div><SectionHead eyebrow={c.orderingEyebrow} title={c.orderingTitle} body={c.orderingBody} />
        <ol className="sv-plain-steps">{c.orderingSteps.map((step, i) => <li key={step}><span>0{i + 1}</span>{step}</li>)}</ol>
        <FeatureActions locale={locale} source="home-online-ordering" destination="features/online-ordering" label={c.orderingLink} />
      </div>
      <div className="sv-guest-captures"><ProductCapture locale={locale} image="checkout" alt={locale === 'ar' ? 'إتمام طلب تجريبي نقدي بقيمة 27.50 جنيه' : 'Synthetic cash checkout totaling 27.50 EGP'} />
        <ProductCapture locale={locale} image="confirmation" alt={locale === 'ar' ? 'تأكيد الطلب التجريبي المرسل من موقع المطعم' : 'Confirmation of the synthetic storefront order'} /></div>
    </div></section>
  </>
}
