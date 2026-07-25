import {describe, expect, it} from 'vitest'
import {agentLabel, humanState, relativeTime, shortDigest, stateIcon} from './presentation'

describe('Workbench projections', () => {
  it('turns machine states into concise human status', () => {
    expect(humanState('awaiting-plan-approval')).toBe('Awaiting Plan Approval')
    expect(agentLabel('quality-verifier')).toBe('Quality Verifier')
    expect(stateIcon('waiting')).toBe('Ⅱ')
  })

  it('keeps exact records recognizable without flooding the interface', () => {
    expect(shortDigest('sha256:abcdefghijklmnopqrstuvwxyz')).toBe('sha256:abcdefg…uvwxyz')
  })

  it('renders deterministic relative time', () => {
    expect(relativeTime('2026-07-25T10:00:00Z', Date.parse('2026-07-25T10:00:42Z'))).toBe('42s ago')
  })
})
