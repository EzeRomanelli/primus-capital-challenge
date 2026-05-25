import { Routes, Route } from 'react-router-dom'
import { Dashboard } from '@/pages/Dashboard'
import { ClienteDetail } from '@/pages/ClienteDetail'

export default function App() {
  return (
    <Routes>
      <Route path="/" element={<Dashboard />} />
      <Route path="/clientes/:id" element={<ClienteDetail />} />
    </Routes>
  )
}
