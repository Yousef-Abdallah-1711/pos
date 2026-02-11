"use client";

import React, { ReactNode, Component, ErrorInfo } from "react";
import { AlertTriangle, RefreshCw } from "lucide-react";
import { debugError } from "@/lib/debug";

interface Props {
  children: ReactNode;
  locale?: "en" | "ar";
  fallback?: ReactNode;
  onError?: (error: Error, info: ErrorInfo) => void;
  section?: string; // e.g., "products", "categories", "featured"
}

interface State {
  hasError: boolean;
  error: Error | null;
  errorCount: number;
}

/**
 * General Error Boundary component for catching component errors
 * Prevents entire app crash when individual sections fail
 */
export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
      errorCount: 0,
    };
  }

  static getDerivedStateFromError(error: Error): State {
    return {
      hasError: true,
      error,
      errorCount: 0,
    };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    // Log error in development
    if (process.env.NODE_ENV === "development") {
      debugError(`[${this.props.section || "ErrorBoundary"}] Caught error:`, error, info);
    }

    // Call error callback if provided
    this.props.onError?.(error, info);

    // Update error count
    this.setState((prev) => ({
      errorCount: prev.errorCount + 1,
    }));
  }

  handleRetry = () => {
    this.setState({
      hasError: false,
      error: null,
      errorCount: 0,
    });
  };

  render() {
    const { hasError, error } = this.state;
    const { children, locale = "en", fallback, section } = this.props;

    if (!hasError) {
      return children;
    }

    // If custom fallback is provided, use it
    if (fallback) {
      return fallback;
    }

    // Default error UI
    const isRTL = locale === "ar";
    const labels = {
      en: {
        title: "Something went wrong",
        message: `An error occurred while loading ${section || "this section"}.`,
        retry: "Try Again",
      },
      ar: {
        title: "حدث خطأ ما",
        message: `حدث خطأ أثناء تحميل ${section || "هذا القسم"}.`,
        retry: "حاول مرة أخرى",
      },
    };

    const t = labels[locale as keyof typeof labels] || labels.en;

    return (
      <div
        className={`p-4 rounded-lg border ${isRTL ? "text-right" : "text-left"}`}
        style={{
          backgroundColor: "#fef2f2",
          borderColor: "#fecaca",
        }}
      >
        <div className="flex gap-3">
          <AlertTriangle
            className="h-5 w-5 flex-shrink-0 mt-0.5"
            style={{ color: "#dc2626" }}
          />
          <div className="flex-1">
            <h3
              className="font-semibold text-sm"
              style={{ color: "#7f1d1d" }}
            >
              {t.title}
            </h3>
            <p
              className="text-sm mt-1"
              style={{ color: "#991b1b" }}
            >
              {t.message}
            </p>
            {error && (
              <p
                className="text-xs mt-2 font-mono overflow-auto max-h-12"
                style={{ color: "#7f1d1d" }}
              >
                {error.message}
              </p>
            )}
            <button
              onClick={this.handleRetry}
              className="mt-3 inline-flex items-center gap-1 px-3 py-1 rounded text-xs font-medium transition-colors hover:opacity-90"
              style={{
                backgroundColor: "#dc2626",
                color: "#ffffff",
              }}
            >
              <RefreshCw className="h-3 w-3" />
              {t.retry}
            </button>
          </div>
        </div>
      </div>
    );
  }
}

/**
 * HOC to wrap a component with error boundary
 */
export function withErrorBoundary<P extends object>(
  Component: React.ComponentType<P>,
  errorBoundaryProps?: Omit<Props, "children">
) {
  return function WrappedComponent(props: P) {
    return (
      <ErrorBoundary {...errorBoundaryProps}>
        <Component {...props} />
      </ErrorBoundary>
    );
  };
}
