import React, { useEffect, useState } from 'react'
import styles from './RemoteCursor.module.css'

export interface RemoteCursorModel {
  userId: number
  username: string
  position: number
  line: number
  column: number
  color: string
  timestamp: Date
}

export interface RemoteCursorProps {
  cursor: RemoteCursorModel
  content: string
  editorElement: HTMLTextAreaElement | null
}

export const RemoteCursor: React.FC<RemoteCursorProps> = ({ cursor, content, editorElement }) => {
  const [cursorPosition, setCursorPosition] = useState<{ left: number; top: number } | null>(null)
  const [isVisible, setIsVisible] = useState(true)

  useEffect(() => {
    if (!editorElement) return

    try {
      // Calculate character position in pixels
      const textarea = editorElement
      const position = cursor.position

      if (position < 0 || position > content.length) {
        setIsVisible(false)
        return
      }

      // Clone textarea to measure position
      const clone = textarea.cloneNode(true) as HTMLTextAreaElement
      clone.style.position = 'absolute'
      clone.style.visibility = 'hidden'
      clone.value = content.substring(0, position)
      document.body.appendChild(clone)

      // Get the position of the cursor
      const coords = getCaretCoordinates(textarea, position)
      document.body.removeChild(clone)

      setCursorPosition(coords)

      setIsVisible(true)
    } catch {
      setIsVisible(false)
    }
  }, [cursor.position, content, editorElement])

  if (!isVisible || !cursorPosition) {
    return null
  }

  return (
    <div className={`${styles.cursor} remote-cursor`} title={cursor.username}>
      <div className={`${styles.label} remote-cursor-label`}>
        {cursor.username}
      </div>
      <style>{`
        .remote-cursor {
          left: ${cursorPosition.left}px;
          top: ${cursorPosition.top}px;
          background-color: ${cursor.color};
        }

        .remote-cursor-label {
          background-color: ${cursor.color};
        }
      `}</style>
    </div>
  )
}

/**
 * Gets the coordinates (x, y) of a caret position in a textarea
 */
function getCaretCoordinates(
  element: HTMLTextAreaElement,
  position: number
): { left: number; top: number } {
  const div = document.createElement('div')
  const style = window.getComputedStyle(element)

  // Copy relevant styles
  ;([
    'direction', // RTL support
    'boxSizing',
    'width', // on Chrome and IE, exclude the scrollbar, so the mirror div wraps exactly as the textarea does
    'height',
    'overflowX',
    'overflowY', // copy the scrollbar for Firefox
    'borderTopWidth',
    'borderRightWidth',
    'borderBottomWidth',
    'borderLeftWidth',
    'paddingTop',
    'paddingRight',
    'paddingBottom',
    'paddingLeft',
    'fontStyle',
    'fontVariant',
    'fontWeight',
    'fontStretch',
    'fontSize',
    'fontSizeAdjust',
    'lineHeight',
    'fontFamily',
    'textAlign',
    'textTransform',
    'textIndent',
    'textDecoration', // might not make a difference, but shouldn't hurt
    'letterSpacing',
    'wordSpacing',
    'tabSize',
  ] as const).forEach((prop) => {
    div.style[prop] = style[prop]
  })

  div.style.position = 'absolute'
  div.style.visibility = 'hidden'
  div.textContent = element.value.substring(0, position)

  const span = document.createElement('span')
  span.textContent = element.value.substring(position) || '.'
  div.appendChild(span)

  document.body.appendChild(div)

  const { offsetLeft: left, offsetTop: top } = span

  document.body.removeChild(div)

  return { left, top }
}

RemoteCursor.displayName = 'RemoteCursor'
