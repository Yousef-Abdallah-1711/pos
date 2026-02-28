/**
 * Custom Section Mock Data
 * Provides custom HTML/Markdown content
 */

export const CUSTOM_MOCK_DATA = {
  content_en: `
    <div class="custom-content">
      <h2>Custom Content Section</h2>
      <p>This is a flexible custom section that allows you to add any HTML or Markdown content.</p>
      <ul>
        <li>Feature 1: Easy to customize</li>
        <li>Feature 2: Supports HTML and Markdown</li>
        <li>Feature 3: Great for additional information</li>
      </ul>
    </div>
  `,
  content_ar: `
    <div class="custom-content" dir="rtl">
      <h2>قسم المحتوى المخصص</h2>
      <p>هذا قسم مخصص مرن يسمح لك بإضافة أي محتوى HTML أو Markdown.</p>
      <ul>
        <li>الميزة 1: سهلة التخصيص</li>
        <li>الميزة 2: يدعم HTML و Markdown</li>
        <li>الميزة 3: رائع للمعلومات الإضافية</li>
      </ul>
    </div>
  `,
  config: {
    render_as_html: true,
    allow_markdown: true,
  },
}
