import { createFileRoute } from '@tanstack/react-router'

import { Pie, PieChart, AreaChart, Area, YAxis, XAxis, CartesianGrid, Tooltip, PieLabelRenderProps, Cell, Legend } from 'recharts'

export const Route = createFileRoute('/')({ component: App })

const RADIAN = Math.PI / 180;
const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042'];

const renderCustomizedLabel = ({ cx, cy, midAngle, innerRadius, outerRadius, percent }: PieLabelRenderProps) => {
  if (cx == null || cy == null || innerRadius == null || outerRadius == null) {
    return null
  }
  const radius = innerRadius + (outerRadius - innerRadius) * 0.5
  const ncx = Number(cx)
  const x = ncx + radius * Math.cos(-(midAngle ?? 0) * RADIAN)
  const ncy = Number(cy)
  const y = ncy + radius * Math.sin(-(midAngle ?? 0) * RADIAN)

  return (
    <text x={x} y={y} fill="white" textAnchor={x > ncx ? 'start' : 'end'} dominantBaseline='central'>
      {`${((percent ?? 1) * 100).toFixed(0)}%`}
    </text>
  )
}

type MoneySpent = {
  amount: number
  expenseType: string
}

const data: MoneySpent[] = [
  { expenseType: "Groceries", amount: 300 },
  { expenseType: "Transportation", amount: 300 },
  { expenseType: "Gifts", amount: 200 },
  { expenseType: "Car", amount: 400 },
]

function App() {

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-900 via-slate-800 to-slate-900">
      <PieChart style={{ width: '100%', maxWidth: '700px', maxHeight: '80vh', aspectRatio: 1 }} responsive>
        <Pie
          data={data}
          dataKey="amount"
          nameKey="expenseType"
          cx="50%"
          cy="50%"
          outerRadius="50%"
          label={renderCustomizedLabel}
          isAnimationActive={false}
          labelLine={false}
        >
          {data.map((entry, index) => (
            <Cell key={`cell-{entry.expenseType}`} fill={COLORS[index % COLORS.length]} />
          ))}
        </Pie>
        <Legend layout='vertical' />
      </PieChart>
      <AreaChart
        style={{ width: '100%', maxWidth: '700px', maxHeight: '70vh', aspectRatio: 1.618 }}
        responsive
        data={data}
        margin={{ top: 10, right: 0, left: 0, bottom: 0 }}
      >
        <defs>
          <linearGradient id="colorUv" x1="0" y1="0" x2="0" y2="1">
            <stop offset="5%" stopColor="#8884d8" stopOpacity={0.8} />
            <stop offset="95%" stopColor="#8884d8" stopOpacity={0} />
          </linearGradient>
          <linearGradient id="colorPv" x1="0" y1="0" x2="0" y2="1">
            <stop offset="5%" stopColor="#82ca9d" stopOpacity={0.8} />
            <stop offset="95%" stopColor="#82ca9d" stopOpacity={0} />
          </linearGradient>
        </defs>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="expenseType" />
        <YAxis width="auto" />
        <Tooltip />
        <Area
          type="monotone"
          dataKey="uv"
          stroke="#8884d8"
          fillOpacity={1}
          fill="url(#colorUv)"
        />
        <Area
          type="monotone"
          dataKey="pv"
          stroke="#82ca9d"
          fillOpacity={1}
          fill="url(#colorPv)"
        />
      </AreaChart>
    </div>
  )
}
