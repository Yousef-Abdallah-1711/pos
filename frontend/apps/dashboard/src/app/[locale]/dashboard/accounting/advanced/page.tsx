'use client'

import { useMemo, useState } from 'react'
import { Loader, RefreshCw } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { api } from '@/lib/api'
import {
  createBankReconciliationSession,
  createRecurringJournal,
  createSupplierPaymentSchedule,
  downloadCsv,
  getBankReconciliationSessions,
  getBudgets,
  getEInvoiceSubmissions,
  getExchangeRates,
  getRecurringJournals,
  getSupplierPaymentReminders,
  getSupplierPaymentSchedules,
  runRecurringJournal,
  sendSupplierPaymentReminder,
  submitEInvoice,
  upsertBudget,
  upsertExchangeRate,
} from '@/lib/api/accountingAdvancedApi'

function toDateInput(value: Date): string {
  return value.toISOString().slice(0, 10)
}

function downloadBlob(filename: string, blob: Blob) {
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  a.click()
  URL.revokeObjectURL(url)
}

export default function AdvancedAccountingToolsPage() {
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [state, setState] = useState<any>(null)

  const [reconciliationForm, setReconciliationForm] = useState({
    bank_account_id: 1,
    statement_date: toDateInput(new Date()),
    statement_balance: 0,
    book_balance: 0,
  })
  const [budgetForm, setBudgetForm] = useState({
    account_id: 1,
    fiscal_year: new Date().getFullYear(),
    period_month: new Date().getMonth() + 1,
    budget_amount: 0,
    currency: 'EGP',
  })
  const [rateForm, setRateForm] = useState({
    base_currency: 'USD',
    quote_currency: 'EGP',
    rate: 0,
    effective_date: toDateInput(new Date()),
  })
  const [scheduleForm, setScheduleForm] = useState({
    supplier_id: 1,
    scheduled_date: toDateInput(new Date()),
    amount: 0,
    reminder_days: 7,
    currency: 'EGP',
  })
  const [invoiceId, setInvoiceId] = useState(1)

  const currency = useMemo(() => new Intl.NumberFormat('en-US', { style: 'currency', currency: 'EGP', maximumFractionDigits: 2 }), [])

  const loadAll = async () => {
    setLoading(true)
    setError(null)
    try {
      const [
        sessions,
        recurring,
        budgets,
        rates,
        schedules,
        reminders,
        eInvoices,
      ] = await Promise.all([
        getBankReconciliationSessions(),
        getRecurringJournals(),
        getBudgets(),
        getExchangeRates(),
        getSupplierPaymentSchedules(),
        getSupplierPaymentReminders(7),
        getEInvoiceSubmissions(),
      ])
      setState({
        sessions,
        recurring,
        budgets,
        rates,
        schedules,
        reminders,
        eInvoices,
      })
    } catch (err: any) {
      setError(err?.response?.data?.error || err?.message || 'Failed to load advanced accounting tools')
    } finally {
      setLoading(false)
    }
  }

  const onCreateReconciliation = async () => {
    await createBankReconciliationSession(reconciliationForm)
    await loadAll()
  }

  const onCreateBudget = async () => {
    await upsertBudget(budgetForm)
    await loadAll()
  }

  const onCreateRate = async () => {
    await upsertExchangeRate(rateForm)
    await loadAll()
  }

  const onCreateSchedule = async () => {
    await createSupplierPaymentSchedule(scheduleForm)
    await loadAll()
  }

  const onCreateRecurringSample = async () => {
    await createRecurringJournal({
      template_code: `REC-${Date.now()}`,
      template_name: 'Monthly Rent',
      frequency: 'monthly',
      start_date: toDateInput(new Date()),
      auto_post: false,
      lines: [
        { account_id: 1, debit_amount: 5000, credit_amount: 0, line_order: 1, description: 'Rent Expense' },
        { account_id: 2, debit_amount: 0, credit_amount: 5000, line_order: 2, description: 'Cash' },
      ],
    })
    await loadAll()
  }

  const onSubmitEInvoice = async () => {
    await submitEInvoice(invoiceId, 'customer')
    await loadAll()
  }

  const exportCoa = async () => {
    const blob = await downloadCsv('/accounting/excel/chart-of-accounts/export')
    downloadBlob('chart-of-accounts-export.csv', blob)
  }

  const exportJournals = async () => {
    const today = toDateInput(new Date())
    const blob = await downloadCsv('/accounting/excel/journal-entries/export', { date_from: '2024-01-01', date_to: today })
    downloadBlob('journal-entries-export.csv', blob)
  }

  const onImportCsv = async (type: 'coa' | 'journal', file: File | null) => {
    if (!file) return
    const form = new FormData()
    form.append('file', file)
    const url = type === 'coa'
      ? '/accounting/excel/chart-of-accounts/import'
      : '/accounting/excel/journal-entries/import'
    await api.post(url, form, { headers: { 'Content-Type': 'multipart/form-data' } })
    await loadAll()
  }

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 xl:flex-row xl:items-end xl:justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Advanced Accounting Tools</h1>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            Phase 7 operations + Phase 9 integration/import/export utilities.
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" onClick={exportCoa}>Export COA CSV</Button>
          <Button variant="outline" onClick={exportJournals}>Export Journals CSV</Button>
          <Button onClick={loadAll} disabled={loading}>
            {loading ? <Loader className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
            Refresh
          </Button>
        </div>
      </div>

      {error ? (
        <div className="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-800 dark:border-red-800 dark:bg-red-900/20 dark:text-red-300">{error}</div>
      ) : null}

      <Section title="Bank Reconciliation">
        <div className="grid grid-cols-1 gap-2 md:grid-cols-5">
          <input type="number" placeholder="Bank Account ID" value={reconciliationForm.bank_account_id} onChange={(e) => setReconciliationForm((s) => ({ ...s, bank_account_id: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="date" value={reconciliationForm.statement_date} onChange={(e) => setReconciliationForm((s) => ({ ...s, statement_date: e.target.value }))} className="rounded border px-3 py-2" />
          <input type="number" placeholder="Statement Balance" value={reconciliationForm.statement_balance} onChange={(e) => setReconciliationForm((s) => ({ ...s, statement_balance: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="number" placeholder="Book Balance" value={reconciliationForm.book_balance} onChange={(e) => setReconciliationForm((s) => ({ ...s, book_balance: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <Button onClick={onCreateReconciliation}>Create Session</Button>
        </div>
        <div className="space-y-2 text-sm">
          {(state?.sessions || []).slice(0, 8).map((item: any) => (
            <div key={item.id} className="flex items-center justify-between rounded border p-2">
              <span>#{item.id} {item.statement_date} ({item.status})</span>
              <span>{currency.format(item.statement_balance || 0)} / Diff {currency.format(item.difference || 0)}</span>
            </div>
          ))}
        </div>
      </Section>

      <Section title="Recurring Journals">
        <div className="flex flex-wrap gap-2">
          <Button variant="outline" onClick={onCreateRecurringSample}>Create Sample Template</Button>
          {(state?.recurring || []).slice(0, 3).map((item: any) => (
            <Button key={item.id} variant="outline" onClick={() => runRecurringJournal(item.id)}>
              Run {item.template_name}
            </Button>
          ))}
        </div>
      </Section>

      <Section title="Budgets + Exchange Rates">
        <div className="grid grid-cols-1 gap-2 md:grid-cols-6">
          <input type="number" placeholder="Account ID" value={budgetForm.account_id} onChange={(e) => setBudgetForm((s) => ({ ...s, account_id: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="number" placeholder="Year" value={budgetForm.fiscal_year} onChange={(e) => setBudgetForm((s) => ({ ...s, fiscal_year: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="number" placeholder="Month" value={budgetForm.period_month} onChange={(e) => setBudgetForm((s) => ({ ...s, period_month: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="number" placeholder="Budget Amount" value={budgetForm.budget_amount} onChange={(e) => setBudgetForm((s) => ({ ...s, budget_amount: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="text" placeholder="Currency" value={budgetForm.currency} onChange={(e) => setBudgetForm((s) => ({ ...s, currency: e.target.value }))} className="rounded border px-3 py-2" />
          <Button onClick={onCreateBudget}>Save Budget</Button>
        </div>
        <div className="grid grid-cols-1 gap-2 md:grid-cols-5">
          <input type="text" placeholder="Base" value={rateForm.base_currency} onChange={(e) => setRateForm((s) => ({ ...s, base_currency: e.target.value }))} className="rounded border px-3 py-2" />
          <input type="text" placeholder="Quote" value={rateForm.quote_currency} onChange={(e) => setRateForm((s) => ({ ...s, quote_currency: e.target.value }))} className="rounded border px-3 py-2" />
          <input type="number" placeholder="Rate" value={rateForm.rate} onChange={(e) => setRateForm((s) => ({ ...s, rate: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="date" value={rateForm.effective_date} onChange={(e) => setRateForm((s) => ({ ...s, effective_date: e.target.value }))} className="rounded border px-3 py-2" />
          <Button onClick={onCreateRate}>Save FX Rate</Button>
        </div>
      </Section>

      <Section title="Supplier Payment Scheduling + Reminders">
        <div className="grid grid-cols-1 gap-2 md:grid-cols-6">
          <input type="number" placeholder="Supplier ID" value={scheduleForm.supplier_id} onChange={(e) => setScheduleForm((s) => ({ ...s, supplier_id: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="date" value={scheduleForm.scheduled_date} onChange={(e) => setScheduleForm((s) => ({ ...s, scheduled_date: e.target.value }))} className="rounded border px-3 py-2" />
          <input type="number" placeholder="Amount" value={scheduleForm.amount} onChange={(e) => setScheduleForm((s) => ({ ...s, amount: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="number" placeholder="Reminder Days" value={scheduleForm.reminder_days} onChange={(e) => setScheduleForm((s) => ({ ...s, reminder_days: Number(e.target.value) }))} className="rounded border px-3 py-2" />
          <input type="text" placeholder="Currency" value={scheduleForm.currency} onChange={(e) => setScheduleForm((s) => ({ ...s, currency: e.target.value }))} className="rounded border px-3 py-2" />
          <Button onClick={onCreateSchedule}>Create Schedule</Button>
        </div>
        <div className="space-y-2 text-sm">
          {(state?.reminders || []).slice(0, 6).map((item: any) => (
            <div key={item.id} className="flex items-center justify-between rounded border p-2">
              <span>{item.supplier_name || `Supplier #${item.supplier_id}`} {item.scheduled_date}</span>
              <div className="flex items-center gap-2">
                <span>{currency.format(item.amount || 0)}</span>
                <Button size="sm" variant="outline" onClick={() => sendSupplierPaymentReminder(item.id)}>Send Reminder</Button>
              </div>
            </div>
          ))}
        </div>
      </Section>

      <Section title="E-Invoice + Excel-Compatible Import">
        <div className="flex flex-wrap items-center gap-2">
          <input type="number" value={invoiceId} onChange={(e) => setInvoiceId(Number(e.target.value))} className="rounded border px-3 py-2" />
          <Button onClick={onSubmitEInvoice}>Submit E-Invoice (Customer)</Button>
          <label className="rounded border px-3 py-2 text-sm cursor-pointer">
            Import COA CSV
            <input type="file" accept=".csv" className="hidden" onChange={(e) => onImportCsv('coa', e.target.files?.[0] || null)} />
          </label>
          <label className="rounded border px-3 py-2 text-sm cursor-pointer">
            Import Journals CSV
            <input type="file" accept=".csv" className="hidden" onChange={(e) => onImportCsv('journal', e.target.files?.[0] || null)} />
          </label>
        </div>
        <div className="space-y-2 text-sm">
          {(state?.eInvoices || []).slice(0, 8).map((item: any) => (
            <div key={item.id} className="flex items-center justify-between rounded border p-2">
              <span>#{item.id} {item.invoice_type} invoice {item.invoice_id}</span>
              <span>{item.status} {item.eta_uuid ? `(${item.eta_uuid})` : ''}</span>
            </div>
          ))}
        </div>
      </Section>
    </div>
  )
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-gray-200 bg-white p-5 dark:border-gray-700 dark:bg-gray-800">
      <h2 className="mb-3 text-lg font-semibold text-gray-900 dark:text-white">{title}</h2>
      <div className="space-y-3">{children}</div>
    </div>
  )
}
