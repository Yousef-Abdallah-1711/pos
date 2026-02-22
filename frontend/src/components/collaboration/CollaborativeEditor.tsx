import React, { useEffect, useRef, useState, useCallback } from 'react'
import { useTheme } from '@pos-saas/component-sdk'
import { useCollaboration } from '../../hooks/useCollaboration'
import { RemoteCursor } from './RemoteCursor'
import { CommentPanel } from './CommentPanel'
import styles from './CollaborativeEditor.module.css'

export interface CollaborativeEditorProps {
  resourceId: number
  resourceType: 'component' | 'theme'
  initialContent?: string
  readOnly?: boolean
  showComments?: boolean
  onSave?: (content: string) => Promise<void>
  onError?: (error: Error) => void
}

export const CollaborativeEditor: React.FC<CollaborativeEditorProps> = ({
  resourceId,
  resourceType,
  initialContent = '',
  readOnly = false,
  showComments = true,
  onSave,
  onError,
}) => {
  const { theme } = useTheme()
  const [content, setContent] = useState(initialContent)
  const [cursorPos, setCursorPos] = useState(0)
  const [selectedText, setSelectedText] = useState<{ start: number; end: number } | null>(null)
  const [activeCommentPos, setActiveCommentPos] = useState<number | null>(null)
  const editorRef = useRef<HTMLTextAreaElement>(null)

  const {
    isConnected,
    participants,
    remoteCursors,
    comments,
    error,
    sendEdit,
    updateCursor,
    addComment,
    resolveComment,
  } = useCollaboration(resourceId, resourceType)

  // Update cursor position
  useEffect(() => {
    if (!editorRef.current) return

    const handleSelectionChange = () => {
      const textarea = editorRef.current
      if (textarea) {
        const pos = textarea.selectionStart
        setCursorPos(pos)

        // Calculate line number
        const beforeCursor = content.substring(0, pos)
        const lineNumber = beforeCursor.split('\n').length - 1

        updateCursor(pos, lineNumber, pos % 80)
      }
    }

    const textarea = editorRef.current
    textarea.addEventListener('click', handleSelectionChange)
    textarea.addEventListener('keyup', handleSelectionChange)
    textarea.addEventListener('mouseup', handleSelectionChange)

    return () => {
      textarea.removeEventListener('click', handleSelectionChange)
      textarea.removeEventListener('keyup', handleSelectionChange)
      textarea.removeEventListener('mouseup', handleSelectionChange)
    }
  }, [content, updateCursor])

  // Handle text input
  const handleInput = useCallback(
    (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      const newContent = e.currentTarget.value
      const oldLength = content.length
      const newLength = newContent.length

      // Detect operation type
      if (newLength > oldLength) {
        // Insert
        const insertedText = newContent.substring(cursorPos, cursorPos + (newLength - oldLength))
        sendEdit({
          type: 'insert',
          position: cursorPos,
          operation: 'insert',
          content: insertedText,
          length: insertedText.length,
        })
      } else if (newLength < oldLength) {
        // Delete
        const deletedLength = oldLength - newLength
        sendEdit({
          type: 'delete',
          position: cursorPos,
          operation: 'delete',
          content: '',
          length: deletedLength,
        })
      }

      setContent(newContent)
    },
    [content, cursorPos, sendEdit]
  )

  // Handle comment creation
  const handleCreateComment = useCallback(
    (text: string) => {
      if (selectedText === null) {
        if (onError) {
          onError(new Error('Please select text to comment on'))
        }
        return
      }

      addComment(text, selectedText.start, activeCommentPos || 0)
      setSelectedText(null)
    },
    [selectedText, activeCommentPos, addComment, onError]
  )

  // Handle comment resolve
  const handleResolveComment = useCallback(
    (commentId: number) => {
      resolveComment(commentId)
    },
    [resolveComment]
  )

  const containerBackground = theme?.colors?.background || '#ffffff'
  const saveButtonBackground = theme?.colors?.primary || '#3b82f6'
  const participantRules = participants
    .map(
      (p, index) => `
        .participant-${index} {
          background-color: ${p.color};
        }
      `
    )
    .join('\n')

  return (
    <div className={`${styles.container} collaborative-editor-themed`}>
      <div className={styles.header}>
        <h2 className={styles.title}>
          {resourceType === 'component' ? 'Component' : 'Theme'} Editor
        </h2>

        <div className={styles.statusBar}>
          <span className={`${styles.status} ${isConnected ? styles.connected : styles.disconnected}`}>
            {isConnected ? '●' : '○'} {isConnected ? 'Connected' : 'Disconnected'}
          </span>

          <span className={styles.participants}>
            {participants.length} participant{participants.length !== 1 ? 's' : ''}
          </span>

          {onSave && (
            <button
              className={`${styles.saveButton} themed-save-button`}
              onClick={() => onSave(content)}
            >
              Save
            </button>
          )}
        </div>
      </div>

      <div className={styles.editorContainer}>
        <div className={styles.editorWrapper}>
          <textarea
            ref={editorRef}
            value={content}
            onChange={handleInput}
            onSelect={(e) => {
              const textarea = e.currentTarget
              setSelectedText({
                start: textarea.selectionStart,
                end: textarea.selectionEnd,
              })
            }}
            readOnly={readOnly}
            className={styles.textarea}
            placeholder="Start typing..."
            spellCheck="false"
          />

          {/* Remote cursors */}
          {remoteCursors.map((cursor) => (
            <RemoteCursor
              key={cursor.userId}
              cursor={cursor}
              content={content}
              editorElement={editorRef.current}
            />
          ))}
        </div>

        {/* Comment panel */}
        {showComments && (
          <CommentPanel
            comments={comments}
            selectedRange={selectedText}
            onAddComment={handleCreateComment}
            onResolveComment={handleResolveComment}
            disabled={!isConnected}
          />
        )}
      </div>

      {error && (
        <div className={styles.error}>
          <p>Error: {error.message}</p>
          <button onClick={() => window.location.reload()}>Reconnect</button>
        </div>
      )}

      <div className={styles.footer}>
        <div className={styles.stats}>
          <span>Characters: {content.length}</span>
          <span>Lines: {content.split('\n').length}</span>
          <span>Words: {content.split(/\s+/).filter((w) => w.length > 0).length}</span>
        </div>

        <div className={styles.participants}>
          {participants.map((p, index) => (
            <div
              key={p.userId}
              className={`${styles.participantBadge} participant-${index}`}
              title={p.username}
            >
              {p.username.charAt(0).toUpperCase()}
            </div>
          ))}
        </div>
      </div>
      <style>{`
        .collaborative-editor-themed {
          background-color: ${containerBackground};
        }

        .themed-save-button {
          background-color: ${saveButtonBackground};
        }

        ${participantRules}
      `}</style>
    </div>
  )
}

CollaborativeEditor.displayName = 'CollaborativeEditor'
CollaborativeEditor.description = 'Real-time collaborative editor with live cursors and comments'
CollaborativeEditor.icon = 'editor'
