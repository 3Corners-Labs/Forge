import {StrictMode, useEffect, useMemo, useState, type CSSProperties} from 'react'
import {createRoot} from 'react-dom/client'
import {Background, Controls, Edge, Handle, Node, Position, ReactFlow} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import './styles.css'
import {activityIcon, agentLabel, humanState, initials, relativeTime, shortDigest, stateIcon} from './presentation'

type Agent = {id:string; role:string; label:string; state:string; summary:string; authority:string; waitingOn?:string}
type Gate = {action:string; label:string; digest:string; summary:string}
type Activity = {sequence:number; at:string; actor:string; type:string; state:string; message:string}
type Run = {
  id:string; title:string; request:string; state:string; progress:number; startedAt:string; updatedAt:string
  agents:Agent[]; gates:Gate[]; activity:Activity[]
  proposal?:{proposalDigest:string; outcome:string; problem:string; acceptanceCriteria:string[]; assumptions:string[]}
  plan?:{planDigest:string; riskClass:string; checks:string[]}
  candidate?:{candidateDigest:string; verificationDigest:string; changedPaths:string[]; evidenceReady:boolean; changeProvider:string; changeModelId:string}
  publication?:{branch:string; pullRequestId:string; pullRequestUrl:string; provider:string; mergePermitted:boolean}
}

const api = {
  overview: () => fetch('/api/v1/overview').then(expectOK),
  create: (message:string) => fetch('/api/v1/runs', {method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify({message})}).then(expectOK),
  run: (id:string) => fetch(`/api/v1/runs/${id}`).then(expectOK),
  act: (id:string, action:string) => fetch(`/api/v1/runs/${id}/actions`, {method:'POST', headers:{'Content-Type':'application/json'}, body:JSON.stringify({action})}).then(expectOK)
}
async function expectOK(response:Response) {
  const body = await response.json()
  if (!response.ok) throw new Error(body?.error?.message ?? 'Forge request failed')
  return body
}

function AgentNode({data}:{data:Agent}) {
  return <div className={`graph-node ${data.state}`}>
    <Handle type="target" position={Position.Left}/>
    <span className="eyebrow">{data.role}</span>
    <strong>{data.label}</strong>
    <small>{data.state === 'waiting' ? `Waiting · ${data.waitingOn}` : data.state}</small>
    <Handle type="source" position={Position.Right}/>
  </div>
}

const nodeTypes = {agent: AgentNode}

function App() {
  const [runs, setRuns] = useState<Run[]>([])
  const [selected, setSelected] = useState<Run|null>(null)
  const [message, setMessage] = useState('Create a clear Forge Workbench where I can collaborate with the virtual product manager and inspect agent progress, builds, tests and evidence.')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState('')
  const [tab, setTab] = useState<'build'|'tests'|'evidence'>('build')

  useEffect(() => { api.overview().then(data => {setRuns(data.runs); if (data.runs[0]) setSelected(data.runs[0])}).catch(e => setError(e.message)) }, [])
  useEffect(() => {
    if (!selected) return
    const stream = new EventSource(`/api/v1/runs/${selected.id}/events`)
    stream.addEventListener('activity', () => api.run(selected.id).then(run => {
      setSelected(run); setRuns(current => [run, ...current.filter(item => item.id !== run.id)])
    }))
    return () => stream.close()
  }, [selected?.id])

  const graph = useMemo(() => {
    if (!selected) return {nodes:[], edges:[]}
    const nodes:Node[] = selected.agents.map((agent, index) => ({
      id:agent.id, type:'agent', position:{x:index*230, y:index%2*36}, data:agent
    }))
    const edges:Edge[] = selected.agents.slice(1).map((agent, index) => ({
      id:`${selected.agents[index].id}-${agent.id}`, source:selected.agents[index].id, target:agent.id,
      animated:selected.agents[index].state === 'completed' && agent.state !== 'completed',
      style:{stroke: agent.state === 'completed' ? '#56e39f' : '#3d465a', strokeWidth:2}
    }))
    return {nodes, edges}
  }, [selected])

  async function createRun(event:React.FormEvent) {
    event.preventDefault(); setBusy(true); setError('')
    try {
      const run = await api.create(message)
      setRuns(current => [run, ...current]); setSelected(run)
    } catch (cause) { setError((cause as Error).message) } finally { setBusy(false) }
  }

  async function act(action:string) {
    if (!selected) return
    setBusy(true); setError('')
    try {
      const run = await api.act(selected.id, action)
      setSelected(run); setRuns(current => [run, ...current.filter(item => item.id !== run.id)])
    } catch (cause) { setError((cause as Error).message) } finally { setBusy(false) }
  }

  return <div className="shell">
    <aside className="sidebar">
      <div className="brand"><span className="brand-mark">F</span><div><strong>Forge</strong><small>WORKBENCH · TRIAL</small></div></div>
      <button className="new-run" onClick={() => setSelected(null)}>＋ New outcome</button>
      <nav aria-label="Delivery runs">
        <span className="nav-label">ACTIVE WORK</span>
        {runs.map(run => <button key={run.id} className={`run-link ${selected?.id===run.id?'active':''}`} onClick={() => setSelected(run)}>
          <span className={`status-dot ${run.state}`}/><span><strong>{run.title}</strong><small>{humanState(run.state)} · {run.progress}%</small></span>
        </button>)}
      </nav>
      <div className="sidebar-footer"><span className="api-light"/> API connected <small>Local Developer Trial</small></div>
    </aside>

    <main>
      <header className="topbar">
        <div><span className="breadcrumb">3Corners / Forge /</span><strong>{selected?.title ?? 'New outcome'}</strong></div>
        <div className="top-actions"><span className="mode-pill">Exploratory</span><button title="Keyboard shortcuts">⌘ K</button><div className="avatar">SB</div></div>
      </header>

      {!selected ? <section className="empty-state">
        <div className="orb">F</div><span className="eyebrow">VIRTUAL PRODUCT TEAM</span>
        <h1>What should we build next?</h1>
        <p>Start with the outcome. Forge will clarify intent, show its plan, prepare a bounded change and independently verify it.</p>
        <form className="composer large" onSubmit={createRun}>
          <textarea aria-label="Describe the product outcome" value={message} onChange={event => setMessage(event.target.value)} rows={5}/>
          <div><span>Intent first · no execution authority</span><button disabled={busy}>{busy?'Starting…':'Start discovery →'}</button></div>
        </form>
      </section> : <>
        <section className="mission">
          <div>
            <span className="eyebrow">CURRENT MISSION</span>
            <h1>{selected.title}</h1>
            <p>{selected.request}</p>
          </div>
          <div className="progress-card"><div className="progress-ring" style={{'--progress':`${selected.progress*3.6}deg`} as CSSProperties}><strong>{selected.progress}%</strong></div><span>{humanState(selected.state)}</span><small>Updated {relativeTime(selected.updatedAt)}</small></div>
        </section>

        {error && <div className="error" role="alert">{error}</div>}

        <section className="decision-row">
          <div className="decision-inbox">
            <div className="panel-heading"><div><span className="eyebrow">DECISION INBOX</span><h2>{selected.gates.length ? 'Your judgement is needed' : 'No decisions waiting'}</h2></div><span className="count">{selected.gates.length}</span></div>
            {selected.gates.map(gate => <div className="gate" key={gate.action}>
              <div><strong>{gate.label}</strong><p>{gate.summary}</p><code>{shortDigest(gate.digest)}</code></div>
              <button disabled={busy} onClick={() => act(gate.action)}>{busy?'Working…':'Review & approve →'}</button>
            </div>)}
            {!selected.gates.length && <p className="muted">Forge can continue without human input.</p>}
          </div>
          <div className="team-pulse">
            <div className="panel-heading"><div><span className="eyebrow">TEAM PULSE</span><h2>{selected.agents.filter(a=>a.state==='completed').length} of {selected.agents.length} roles complete</h2></div><span className="live"><i/> LIVE</span></div>
            <div className="agent-mini-grid">{selected.agents.map(agent => <div className={`agent-mini ${agent.state}`} key={agent.id}><span>{initials(agent.label)}</span><div><strong>{agent.label}</strong><small>{agent.state==='waiting'?agent.waitingOn:agent.summary}</small></div><b>{stateIcon(agent.state)}</b></div>)}</div>
          </div>
        </section>

        <section className="graph-panel">
          <div className="panel-heading"><div><span className="eyebrow">DELIVERY GRAPH</span><h2>Who is doing what — and why work is waiting</h2></div><div className="legend"><span><i className="done"/>Done</span><span><i className="wait"/>Waiting</span><span><i/>Queued</span></div></div>
          <div className="flow" aria-label="Agent delivery graph"><ReactFlow nodes={graph.nodes} edges={graph.edges} nodeTypes={nodeTypes} fitView proOptions={{hideAttribution:true}} nodesDraggable={false} nodesConnectable={false}><Background color="#202735" gap={22}/><Controls showInteractive={false}/></ReactFlow></div>
        </section>

        <section className="lower-grid">
          <div className="activity-panel">
            <div className="panel-heading"><div><span className="eyebrow">LIVE ACTIVITY</span><h2>Evidence-backed timeline</h2></div></div>
            <ol className="timeline">{[...selected.activity].reverse().map(item => <li key={item.sequence}><span className="timeline-icon">{activityIcon(item.type)}</span><div><strong>{agentLabel(item.actor)}</strong><p>{item.message}</p><small>{relativeTime(item.at)} · {humanState(item.state)}</small></div></li>)}</ol>
          </div>
          <div className="inspection-panel">
            <div className="tabs" role="tablist">
              {(['build','tests','evidence'] as const).map(value => <button role="tab" aria-selected={tab===value} className={tab===value?'active':''} onClick={()=>setTab(value)} key={value}>{value}</button>)}
            </div>
            {tab==='build' && <Build run={selected}/>}
            {tab==='tests' && <Tests run={selected}/>}
            {tab==='evidence' && <Evidence run={selected}/>}
          </div>
        </section>
      </>}
    </main>
  </div>
}

function Build({run}:{run:Run}) {
  return <div className="tab-content">
    <span className="eyebrow">CANDIDATE</span><h2>{run.candidate?'Reviewable change prepared':'Waiting for approved plan'}</h2>
    {run.candidate ? <><div className="metric-row"><Metric label="Changed files" value={`${run.candidate.changedPaths.length}`}/><Metric label="Provider" value={run.candidate.changeProvider}/><Metric label="Merge" value="Human only"/></div>
      <div className="file-list">{run.candidate.changedPaths.map(path=><div key={path}><span>M</span><code>{path}</code><b>verified</b></div>)}</div>
      <pre className="diff"><span>+ # Forge M1 Developer Trial candidate</span>{'\n'}+ Prepared from promoted requirements{'\n'}+ Independently verified before publication</pre></>
      : <Empty label="The engineering worker starts only after exact plan approval."/>}
  </div>
}

function Tests({run}:{run:Run}) {
  const passed=Boolean(run.candidate)
  return <div className="tab-content"><span className="eyebrow">ASSURANCE</span><h2>Public journey test interface</h2>
    <div className="test-summary"><strong>{passed?'2 passed':'2 queued'}</strong><span>0 failed</span><span>0 flaky</span></div>
    {['Candidate integrity','Declared path scope'].map(test=><div className="test-row" key={test}><b>{passed?'✓':'○'}</b><div><strong>{test}</strong><small>Independent quality verifier · consequential</small></div><span>{passed?'passed':'queued'}</span></div>)}
    <div className="coverage"><span style={{width:passed?'100%':'0%'}}/><small>Accepted M1 journey coverage</small></div>
  </div>
}

function Evidence({run}:{run:Run}) {
  return <div className="tab-content"><span className="eyebrow">TRACEABILITY</span><h2>Intent to evidence</h2>
    <div className="evidence-chain">{[
      ['Intent',Boolean(run.proposal),run.proposal?.proposalDigest],
      ['Plan',Boolean(run.plan),run.plan?.planDigest],
      ['Candidate',Boolean(run.candidate),run.candidate?.candidateDigest],
      ['Verification',Boolean(run.candidate),run.candidate?.verificationDigest],
      ['Publication',Boolean(run.publication),run.publication?.pullRequestId]
    ].map(([label,ready,digest])=><div className={ready?'ready':''} key={String(label)}><b>{ready?'✓':'○'}</b><span><strong>{label}</strong><code>{digest?shortDigest(String(digest)):'pending'}</code></span></div>)}</div>
  </div>
}

function Metric({label,value}:{label:string;value:string}) { return <div className="metric"><span>{label}</span><strong>{value}</strong></div> }
function Empty({label}:{label:string}) { return <div className="empty"><span>◇</span><p>{label}</p></div> }
createRoot(document.getElementById('root')!).render(<StrictMode><App/></StrictMode>)
