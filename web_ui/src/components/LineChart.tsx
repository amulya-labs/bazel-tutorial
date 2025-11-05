import { LineChart as RechartsLineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

interface Observation {
  date: string
  value: number
}

interface LineChartProps {
  data: Observation[]
}

function LineChart({ data }: LineChartProps) {
  // Format data for Recharts
  const chartData = data.map(obs => ({
    date: new Date(obs.date).toLocaleDateString('en-US', { 
      year: 'numeric', 
      month: 'short' 
    }),
    value: obs.value,
  }))

  // Sample data points for better performance on large datasets
  const sampledData = chartData.length > 100 
    ? chartData.filter((_, index) => index % Math.ceil(chartData.length / 100) === 0)
    : chartData

  return (
    <div style={{ width: '100%', height: 400 }}>
      <ResponsiveContainer>
        <RechartsLineChart
          data={sampledData}
          margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
        >
          <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
          <XAxis 
            dataKey="date" 
            tick={{ fontSize: 12 }}
            stroke="#718096"
          />
          <YAxis 
            tick={{ fontSize: 12 }}
            stroke="#718096"
          />
          <Tooltip 
            contentStyle={{
              backgroundColor: 'white',
              border: '1px solid #e2e8f0',
              borderRadius: '6px',
              boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
            }}
          />
          <Line 
            type="monotone" 
            dataKey="value" 
            stroke="#667eea" 
            strokeWidth={2}
            dot={false}
            activeDot={{ r: 6 }}
          />
        </RechartsLineChart>
      </ResponsiveContainer>
    </div>
  )
}

export default LineChart
