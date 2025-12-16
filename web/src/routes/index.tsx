import { createFileRoute } from '@tanstack/react-router'
import {
  Zap,
  Server,
  Route as RouteIcon,
  Shield,
  Waves,
  Sparkles,
} from 'lucide-react'

import { Pie, PieChart } from 'recharts'

export const Route = createFileRoute('/')({ component: App })

type MoneySpent = {
  amount: number
  expenseType: string
}

const data: MoneySpent[] = [
  { expenseType: "Groceries", amount: 300 }
]

function App() {

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900">
      <PieChart
        style={{ width: '100%', maxWidth: '500px', maxHeight: '80vh', aspectRatio: 1 }}
        responsive
        margin={{ top: 50, right: 50, bottom: 50, left: 50 }}
      >
        <Pie data={data}
          dataKey="amount"
          nameKey="expenseType"
          cx="50%"
          cy="50%"
          outerRadius="50%"
          fill="#8884d8"
          isAnimationActive={true}
        />
      </PieChart>
    </div>
  )
}
