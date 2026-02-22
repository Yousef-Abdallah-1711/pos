import React, { useEffect, useState, useCallback } from 'react'
import { useTheme } from '@pos-saas/component-sdk'
import { usePersonalization } from '../../hooks/usePersonalization'
import { RecommendationList } from './RecommendationList'
import styles from './PersonalizedContent.module.css'

export interface PersonalizedContentProps {
  userId: number
  restaurantId: number
  maxRecommendations?: number
  onRecommendationClick?: (componentId: number) => void
  showMetrics?: boolean
}

export const PersonalizedContent: React.FC<PersonalizedContentProps> = ({
  userId,
  restaurantId,
  maxRecommendations = 10,
  onRecommendationClick,
  showMetrics = false,
}) => {
  const { theme } = useTheme()
  const {
    recommendations,
    userFeatures,
    metrics,
    loading,
    error,
    refresh,
  } = usePersonalization(userId, restaurantId, maxRecommendations)

  const [isExpanded, setIsExpanded] = useState(false)

  const handleRecommendationClick = useCallback(
    (componentId: number, recommendationId: number) => {
      // Track the click
      onRecommendationClick?.(componentId)

      // Trigger analytics
      if (typeof window !== 'undefined' && (window as any).__personalizationSDK) {
        (window as any).__personalizationSDK.trackRecommendationClick(
          recommendationId,
          userId,
          restaurantId
        )
      }
    },
    [userId, restaurantId, onRecommendationClick]
  )

  if (error) {
    return (
      <div className={styles.error}>
        <p>Failed to load recommendations: {error.message}</p>
        <button onClick={() => refresh()} className={styles.retryButton}>
          Retry
        </button>
      </div>
    )
  }

  const headerBg = theme?.colors?.primary || '#3b82f6'
  const headerText = theme?.colors?.text || '#ffffff'

  return (
    <div className={`personalized-content-themed ${styles.container}`}>
      <div className={styles.header}>
        <h2 className={styles.title}>Personalized For You</h2>
        {showMetrics && userFeatures && (
          <div className={styles.metricsPreview}>
            <span className={styles.metric}>
              Engagement: {(userFeatures.engagementScore * 100).toFixed(0)}%
            </span>
            <span className={styles.metric}>
              Churn Risk: {(userFeatures.churnRisk * 100).toFixed(0)}%
            </span>
          </div>
        )}
      </div>

      <div className={styles.content}>
        {loading && !recommendations.length ? (
          <div className={styles.loading}>
            <div className={styles.spinner} />
            <p>Loading personalized recommendations...</p>
          </div>
        ) : recommendations.length > 0 ? (
          <>
            <RecommendationList
              recommendations={recommendations.slice(0, isExpanded ? undefined : 5)}
              onRecommendationClick={handleRecommendationClick}
              loading={loading}
            />

            {recommendations.length > 5 && (
              <button
                className={styles.expandButton}
                onClick={() => setIsExpanded(!isExpanded)}
              >
                {isExpanded ? 'Show Less' : `Show ${recommendations.length - 5} More`}
              </button>
            )}
          </>
        ) : (
          <div className={styles.empty}>
            <p>No recommendations yet. Keep exploring to see personalized content!</p>
          </div>
        )}
      </div>

      {showMetrics && metrics && (
        <div className={styles.metricsPanel}>
          <h3>Recommendation Performance</h3>
          <div className={styles.metricsGrid}>
            <div className={styles.metricCard}>
              <span className={styles.label}>Click-Through Rate</span>
              <span className={styles.value}>
                {(metrics.clickThroughRate * 100).toFixed(1)}%
              </span>
            </div>
            <div className={styles.metricCard}>
              <span className={styles.label}>Total Shown</span>
              <span className={styles.value}>{metrics.totalRecommendations}</span>
            </div>
            <div className={styles.metricCard}>
              <span className={styles.label}>Avg Relevance</span>
              <span className={styles.value}>
                {(metrics.averageRelevanceScore * 100).toFixed(0)}%
              </span>
            </div>
            <div className={styles.metricCard}>
              <span className={styles.label}>Conversions</span>
              <span className={styles.value}>{metrics.totalConversions}</span>
            </div>
          </div>
        </div>
      )}

      <button className={styles.refreshButton} onClick={() => refresh()}>
        Refresh Recommendations
      </button>

      <style>{`
        .personalized-content-themed .${styles.header} {
          background-color: ${headerBg};
          color: ${headerText};
        }
      `}</style>
    </div>
  )
}

PersonalizedContent.displayName = 'PersonalizedContent'
PersonalizedContent.description = 'Displays personalized content recommendations based on user behavior'