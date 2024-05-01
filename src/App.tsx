import clsx from 'clsx'
import { Dispatch, SetStateAction, useEffect, useState } from 'react'

export default function App() {
  const [state, setState] = useState(false)
  const [time, setTime] = useState(0)

  return (
    <div>
      GOAT STATE - {time < 10000 && time} -{' '}
      <button
        onClick={() => {
          setState(!state)
          setTime(Date.now())
        }}
      >
        Toggle
      </button>
      {[...new Array(50000)].map((_, i) => (
        <Ez key={i} state={state} setTime={setTime} last={i === 50000 - 1} />
      ))}
    </div>
  )
}

function Ez({ state, last, setTime }: { state: boolean; last: boolean; setTime: Dispatch<SetStateAction<number>> }) {
  useEffect(() => {
    if (!last) return
    setTime(c => Date.now() - c)
  }, [state, last, setTime])
  return <div className={clsx(state ? 'bg-red-500' : 'bg-blue-500')}>Ez</div>
}
