import clsx from 'clsx'
import { useState } from 'react'

export default function App() {
  const [state, setState] = useState(false)

  return (
    <div className='group'>
      GOAT STATE
      <button onClick={() => setState(c => !c)}>Toggle</button>
      {[...new Array(50000)].map((_, i) => (
        <Ez key={i} state={state} />
      ))}
    </div>
  )
}

function Ez({ state }: { state: boolean }) {
  return <div className={clsx(state ? 'bg-red-500' : 'bg-blue-500')}>Ez</div>
}
