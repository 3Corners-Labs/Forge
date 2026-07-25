export const humanState = (state:string) =>
  state.replaceAll('-', ' ').replace(/\b\w/g, character => character.toUpperCase())

export const shortDigest = (value:string) =>
  value.length > 22 ? `${value.slice(0, 14)}…${value.slice(-6)}` : value

export const relativeTime = (value:string, now = Date.now()) => {
  const seconds = Math.max(0, Math.floor((now - new Date(value).getTime()) / 1000))
  return seconds < 5 ? 'just now' : seconds < 60 ? `${seconds}s ago` : `${Math.floor(seconds / 60)}m ago`
}

export const initials = (value:string) =>
  value.split(' ').map(part => part[0]).join('').slice(0, 2)

export const stateIcon = (state:string) =>
  state === 'completed' ? '✓' : state === 'waiting' ? 'Ⅱ' : '·'

export const activityIcon = (type:string) =>
  type === 'verification' ? '✓' : type === 'publication' ? '↗' : type === 'plan' ? '⌘' : '✦'

export const agentLabel = (id:string) =>
  id.split('-').map(part => part[0].toUpperCase() + part.slice(1)).join(' ')
